package sqlite_adapter

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func Connect(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath)
}
