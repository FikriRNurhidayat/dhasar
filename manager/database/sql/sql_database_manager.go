package sql_database_manager

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	manager_values "github.com/fikrirnurhidayat/dhasar/manager/values"
	"github.com/fikrirnurhidayat/dhasar/specification"
	"github.com/fikrirnurhidayat/x/logger"
)

type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type querier struct {
	logger           logger.Logger
	queryContextFunc func(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContextFunc  func(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (q *querier) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	st := time.Now()
	result, err := q.execContextFunc(ctx, query, args...)
	q.logger.Debug("database.sql/QUERY", logger.String("query", query), logger.Any("args", args), logger.String("took", fmt.Sprintf("%d ms", time.Since(st).Milliseconds())))
	return result, err
}

func (q *querier) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	st := time.Now()
	rows, err := q.queryContextFunc(ctx, query, args...)
	q.logger.Debug("database.sql/QUERY", logger.String("query", query), logger.Any("args", args), logger.String("took", fmt.Sprintf("%d ms", time.Since(st).Milliseconds())))
	return rows, err
}

type SQLDatabaseManager interface {
	Querier(ctx context.Context) Querier
	Paginate(builder squirrel.SelectBuilder, specs ...dhasar_specification.Specification) squirrel.SelectBuilder
}

type SQLDatabaseManagerImpl struct {
	db     *sql.DB
	logger logger.Logger
}

func (m *SQLDatabaseManagerImpl) Paginate(builder squirrel.SelectBuilder, specs ...dhasar_specification.Specification) squirrel.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case dhasar_specification.LimitSpecification:
			builder = builder.Limit(uint64(v.Limit))
		case dhasar_specification.OffsetSpecification:
			builder = builder.Offset(uint64(v.Offset))
		}
	}

	return builder
}

func (m *SQLDatabaseManagerImpl) Querier(ctx context.Context) Querier {
	hasExternalTransaction := ctx.Value(manager_values.TxKey{}) != nil
	if !hasExternalTransaction {
		return &querier{
			logger:           m.logger,
			queryContextFunc: m.db.QueryContext,
			execContextFunc:  m.db.ExecContext,
		}
	}

	v := ctx.Value(manager_values.TxKey{})
	tx, ok := v.(*sql.Tx)
	if ok {
		m.logger.Debug("transaction/EXPANDED")
		return &querier{
			logger:           m.logger,
			queryContextFunc: tx.QueryContext,
			execContextFunc:  tx.ExecContext,
		}
	}

	return &querier{
		logger:           m.logger,
		queryContextFunc: m.db.QueryContext,
		execContextFunc:  m.db.ExecContext,
	}
}

func New(logger logger.Logger, db *sql.DB) SQLDatabaseManager {
	return &SQLDatabaseManagerImpl{
		db:     db,
		logger: logger,
	}
}
