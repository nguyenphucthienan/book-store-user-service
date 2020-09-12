package user

import (
	"fmt"
	"github.com/nguyenphucthienan/book-store-user-service/datasource/mysql/db"
	"github.com/nguyenphucthienan/book-store-user-service/util/errors"
	"github.com/nguyenphucthienan/book-store-user-service/util/mysql_util"
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
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
		&user.DateCreated, &user.Status); getErr != nil {
		return mysql_util.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestError {
	stmt, err := db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, insertErr := stmt.Exec(user.FirstName, user.LastName, user.Email,
		user.DateCreated, user.Status, user.Password)
	if insertErr != nil {
		return mysql_util.ParseError(insertErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_util.ParseError(err)
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestError {
	stmt, err := db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Id)
	if updateErr != nil {
		return mysql_util.ParseError(updateErr)
	}
	return nil
}

func (user *User) Delete() *errors.RestError {
	stmt, err := db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		return mysql_util.ParseError(deleteErr)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestError) {
	stmt, err := db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email,
			&user.DateCreated, &user.Status); err != nil {
			return nil, mysql_util.ParseError(err)
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
	}
	return users, nil
}
