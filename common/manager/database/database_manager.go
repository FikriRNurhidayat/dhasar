package database_manager

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	manager_values "github.com/fikrirnurhidayat/dhasar/common/manager/values"
	"github.com/fikrirnurhidayat/dhasar/common/specification"
	"github.com/fikrirnurhidayat/dhasar/infra/logger"
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

type DatabaseManager interface {
	Querier(ctx context.Context) Querier
	Paginate(builder squirrel.SelectBuilder, specs ...common_specification.Specification) squirrel.SelectBuilder
}

type DatabaseManagerImpl struct {
	db     *sql.DB
	logger logger.Logger
}

func (m *DatabaseManagerImpl) Paginate(builder squirrel.SelectBuilder, specs ...common_specification.Specification) squirrel.SelectBuilder {
	for _, spec := range specs {
		switch v := spec.(type) {
		case common_specification.LimitSpecification:
			builder = builder.Limit(uint64(v.Limit))
		case common_specification.OffsetSpecification:
			builder = builder.Offset(uint64(v.Offset))
		}
	}

	return builder
}

func (m *DatabaseManagerImpl) Querier(ctx context.Context) Querier {
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

func New(logger logger.Logger, db *sql.DB) DatabaseManager {
	return &DatabaseManagerImpl{
		db:     db,
		logger: logger,
	}
}
