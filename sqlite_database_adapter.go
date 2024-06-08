package dhasar

import (
	"database/sql"

	"github.com/fikrirnurhidayat/x/logger"
	_ "modernc.org/sqlite"
)

type SQLiteDatabaseAdapter struct {
	db     *sql.DB
	logger logger.Logger
}

type SQLiteDatabaseAdapterOption struct {
	FilePath string
}

func (s *SQLiteDatabaseAdapter) Close() error {
	if err := s.db.Close(); err != nil {
		s.logger.Error("sqlite/CLOSE", logger.String("error", err.Error()))
		return err
	}

	s.logger.Debug("sqlite/CLOSE", logger.String("status", "OK!"))

	return nil
}

func (s *SQLiteDatabaseAdapter) Connect(opt *SQLiteDatabaseAdapterOption) (*sql.DB, error) {
	db, err := sql.Open("sqlite", opt.FilePath)

	if err != nil {
		s.logger.Error("sqlite/CONNECT", logger.String("error", err.Error()))
		return nil, err
	}

	s.db = db

	s.logger.Debug("sqlite/CONNECT", logger.String("status", "OK!"))

	return db, nil
}

func NewSQLiteDatabaseAdapter(logger logger.Logger) Adapter[*SQLiteDatabaseAdapterOption, *sql.DB] {
	return &SQLiteDatabaseAdapter{
		logger: logger,
	}
}
