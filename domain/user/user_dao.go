package user

import (
	"fmt"
	"github.com/nguyenphucthienan/book-store-user-service/datasource/mysql/db"
	"github.com/nguyenphucthienan/book-store-user-service/logger"
	"github.com/nguyenphucthienan/book-store-user-service/util/errors"
)

const (
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryUpdateUser       = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?;"
	queryDeleteUser       = "DELETE FROM users WHERE id = ?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
)

func (user *User) Get() *errors.RestError {
	stmt, err := db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("Database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
		&user.DateCreated, &user.Status); getErr != nil {
		logger.Error("Error when trying to get user", err)
		return errors.NewInternalServerError("Database error")
	}
	return nil
}

func (user *User) Save() *errors.RestError {
	stmt, err := db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("Error when trying to prepare insert user statement", err)
		return errors.NewInternalServerError("Database error")
	}
	defer stmt.Close()

	insertResult, insertErr := stmt.Exec(user.FirstName, user.LastName, user.Email,
		user.DateCreated, user.Status, user.Password)
	if insertErr != nil {
		logger.Error("Error when trying to insert user", err)
		return errors.NewInternalServerError("Database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("Error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("Database error")
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestError {
	stmt, err := db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("Error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("Database error")
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Id)
	if updateErr != nil {
		logger.Error("Error when trying to update user", err)
		return errors.NewInternalServerError("Database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestError {
	stmt, err := db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("Error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("Database error")
	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("Error when trying to delete user", err)
		return errors.NewInternalServerError("Database error")
	}
	return nil
}

func (user *User) SearchByStatus(status string) ([]User, *errors.RestError) {
	stmt, err := db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("Error when trying to prepare search user by status statement", err)
		return nil, errors.NewInternalServerError("Database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("Error when trying to search user by status", err)
		return nil, errors.NewInternalServerError("Database error")
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
			&user.DateCreated, &user.Status); err != nil {
			logger.Error("Error when trying to search user by status", err)
			return nil, errors.NewInternalServerError("Database error")
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
	}
	return users, nil
}
