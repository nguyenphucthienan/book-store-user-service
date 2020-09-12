package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nguyenphucthienan/book-store-user-service/logger"
	"os"
)

const (
	mysqlHost     = "mysql_user_host"
	mysqlPort     = "mysql_user_port"
	mysqlUsername = "mysql_user_username"
	mysqlPassword = "mysql_user_password"
	mysqlSchema   = "mysql_user_schema"
)

var (
	Client *sql.DB

	host     = os.Getenv(mysqlHost)
	port     = os.Getenv(mysqlPort)
	username = os.Getenv(mysqlUsername)
	password = os.Getenv(mysqlPassword)
	schema   = os.Getenv(mysqlSchema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		username, password, host, port, schema,
	)

	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}

	logger.Info("Database successfully configured")
}
