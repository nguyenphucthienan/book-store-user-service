package mysql_util

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	restErrors "github.com/nguyenphucthienan/book-store-utils-go/errors"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) restErrors.RestError {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return restErrors.NewNotFoundError("No record matching with given id")
		}
		return restErrors.NewInternalServerError("Error parsing database response",
			errors.New("database error"))
	}

	switch sqlErr.Number {
	case 1062:
		return restErrors.NewBadRequestError("Invalid data")
	default:
		return restErrors.NewInternalServerError("Error processing request",
			errors.New("database error"))
	}
}
