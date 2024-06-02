package transaction_manager

import (
	"context"
	"database/sql"

	manager_values "github.com/fikrirnurhidayat/dhasar/manager/values"
	"github.com/fikrirnurhidayat/x/logger"
)

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

	if err := fn(context.WithValue(ctx, manager_values.TxKey{}, tx)); err != nil {
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

func New(logger logger.Logger, db *sql.DB) TransactionManager {
	return &TransactionManagerImpl{
		db:     db,
		logger: logger,
	}
}
