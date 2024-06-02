package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(connStr string) (*sql.DB, error) {
	return sql.Open("postgres", connStr)
}
