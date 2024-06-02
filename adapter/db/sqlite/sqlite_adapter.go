package postgres

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(connStr string) (*sql.DB, error) {
	return sql.Open("sqlite3", connStr)
}
