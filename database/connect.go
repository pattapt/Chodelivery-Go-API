package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB establishes a connection to the database and returns a *sql.DB object.
func ConnectDB() (*sql.DB, error) {
	username := "appAccount"
	password := "qXS445oMQAPW!fTZ"
	hostname := "localhost"
	port := "3306"
	databaseName := "appdb"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, databaseName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
