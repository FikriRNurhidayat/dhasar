package dhasar

import (
	"database/sql"
	"fmt"

	"github.com/fikrirnurhidayat/x/logger"
	_ "github.com/lib/pq"
)

type PostgresDatabaseAdapter struct {
	db     *sql.DB
	logger logger.Logger
}

func (p *PostgresDatabaseAdapter) Close() error {
	if err := p.db.Close(); err != nil {
		p.logger.Error("postgres/CLOSE", logger.String("error", err.Error()))
		return err
	}

	p.logger.Debug("postgres/CLOSE", logger.String("status", "OK!"))

	return nil
}

func (p *PostgresDatabaseAdapter) Connect(opt *PostgresDatabaseAdapterOption) (*sql.DB, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", opt.Username, opt.Password, opt.Host, opt.Port, opt.Name, opt.SSLMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		p.logger.Error("postgres/CONNECT", logger.String("error", err.Error()))
		return nil, err
	}

	p.db = db

	p.logger.Debug("postgres/CONNECT", logger.String("status", "OK!"))

	return db, nil
}

type PostgresDatabaseAdapterOption struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func NewPostgresDatabaseAdapter(logger logger.Logger) Adapter[*PostgresDatabaseAdapterOption, *sql.DB] {
	return &PostgresDatabaseAdapter{
		logger: logger,
	}
}
