package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	hostName := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, hostName, dbName)

	db, err := sql.Open(dbDriver, dbURI)
	if err != nil {
		panic(err.Error())
	}

	return db
}
