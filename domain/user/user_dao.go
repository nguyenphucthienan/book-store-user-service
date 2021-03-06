package user

import (
	"errors"
	"fmt"
	"github.com/nguyenphucthienan/book-store-user-service/datasource/mysql/db"
	"github.com/nguyenphucthienan/book-store-user-service/logger"
	"github.com/nguyenphucthienan/book-store-user-service/util/mysql_util"
	restErrors "github.com/nguyenphucthienan/book-store-utils-go/errors"
	"strings"
)

const (
	queryGetUser                    = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryInsertUser                 = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryUpdateUser                 = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?;"
	queryDeleteUser                 = "DELETE FROM users WHERE id = ?;"
	queryFindUserByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
	queryFindUserByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email = ? AND password= ?  AND status= ?"
)

func (user *User) Get() restErrors.RestError {
	stmt, err := db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return restErrors.NewInternalServerError("Error when trying to get user",
			errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
		&user.DateCreated, &user.Status); getErr != nil {
		logger.Error("Error when trying to get user", err)
		return restErrors.NewInternalServerError("Error when trying to get user",
			errors.New("database error"))
	}
	return nil
}

func (user *User) Save() restErrors.RestError {
	stmt, err := db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("Error when trying to prepare insert user statement", err)
		return restErrors.NewInternalServerError("Error when trying to save user",
			errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, insertErr := stmt.Exec(user.FirstName, user.LastName, user.Email,
		user.DateCreated, user.Status, user.Password)
	if insertErr != nil {
		logger.Error("Error when trying to insert user", err)
		return restErrors.NewInternalServerError("Error when trying to save user",
			errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("Error when trying to get last insert id after creating a new user", err)
		return restErrors.NewInternalServerError("Error when trying to save user",
			errors.New("database error"))
	}

	user.Id = userId
	return nil
}

func (user *User) Update() restErrors.RestError {
	stmt, err := db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("Error when trying to prepare update user statement", err)
		return restErrors.NewInternalServerError("Error when trying to update user",
			errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Id)
	if updateErr != nil {
		logger.Error("Error when trying to update user", err)
		return restErrors.NewInternalServerError("Error when trying to update user",
			errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() restErrors.RestError {
	stmt, err := db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("Error when trying to prepare delete user statement", err)
		return restErrors.NewInternalServerError("Error when trying to delete user",
			errors.New("database error"))
	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("Error when trying to delete user", err)
		return restErrors.NewInternalServerError("Error when trying to delete user",
			errors.New("database error"))
	}
	return nil
}

func (user *User) SearchByStatus(status string) ([]User, restErrors.RestError) {
	stmt, err := db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("Error when trying to prepare search user by status statement", err)
		return nil, restErrors.NewInternalServerError("Error when trying to search users",
			errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("Error when trying to search user by status", err)
		return nil, restErrors.NewInternalServerError("Error when trying to search users",
			errors.New("database error"))
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
			&user.DateCreated, &user.Status); err != nil {
			logger.Error("Error when trying to search user by status", err)
			return nil, restErrors.NewInternalServerError("Error when trying to search user",
				errors.New("database error"))
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, restErrors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
	}
	return users, nil
}

func (user *User) FindUserByEmailAndPassword() restErrors.RestError {
	stmt, err := db.Client.Prepare(queryFindUserByEmailAndPassword)
	if err != nil {
		logger.Error("Error when trying to prepare get user by email and password statement", err)
		return restErrors.NewInternalServerError("Error when trying to find user",
			errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status);
		getErr != nil {
		if strings.Contains(getErr.Error(), mysql_util.ErrorNoRows) {
			return restErrors.NewNotFoundError("Invalid user credentials")
		}
		logger.Error("Error when trying to get user by email and password", getErr)
		return restErrors.NewInternalServerError("Error when trying to find user",
			errors.New("database error"))
	}
	return nil
}
