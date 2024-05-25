package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func New() (*sql.DB, error) {
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbName := viper.GetString("database.name")
	dbSSLMode := viper.GetString("database.sslmode")

	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	return sql.Open("postgres", dbUrl)
}
