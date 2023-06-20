package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dbDriver := "mysql"
	dbUser := "user"
	dbPass := "secret"
	dbName := "ah_bonus"
	hostName := "fullstack-mysql"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+hostName+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
