package dhasar

import (
	"context"
	"database/sql"

	"github.com/fikrirnurhidayat/x/logger"
)

type TxKey struct{}

type TransactionManager interface {
	Execute(ctx context.Context, fn func(context.Context) error) error
}

type TransactionManagerImpl struct {
	db     *sql.DB
	logger logger.Logger
}

func (m *TransactionManagerImpl) Execute(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	m.logger.Debug("transaction/STARTED")

	if err := fn(context.WithValue(ctx, TxKey{}, tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			m.logger.Debug("transaction/ABORTED")
			return err
		}

		m.logger.Debug("transaction/ABORTED")
		return err
	}

	if err := tx.Commit(); err != nil {
		m.logger.Debug("transaction/ABORTED")
		return err
	}

	m.logger.Debug("transaction/COMMITED")
	return nil
}

func NewTransactionManager(logger logger.Logger, db *sql.DB) TransactionManager {
	return &TransactionManagerImpl{
		db:     db,
		logger: logger,
	}
}
