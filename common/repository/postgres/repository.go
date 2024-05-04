package postgres_repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	database_manager "github.com/fikrirnurhidayat/dhasar/common/manager/database"
	"github.com/fikrirnurhidayat/dhasar/common/repository"
	"github.com/fikrirnurhidayat/dhasar/common/specification"
	"github.com/fikrirnurhidayat/dhasar/infra/logger"
)

type PostgresRepository[Entity any, Specification any, Row any] struct {
	dbm          database_manager.DatabaseManager
	logger       logger.Logger
	tableName    string
	columns      []string
	schema       map[string]string
	primaryKey   string
	filter       func(...Specification) sq.Sqlizer
	scan         func(*sql.Rows) (Row, error)
	row          func(Entity) Row
	values       func(Row) []any
	entity       func(Row) Entity
	noEntities   []Entity
	noEntity     Entity
	noRow        Row
	noRows       []Row
	upsertSuffix string
}

type PostgresIterator[Entity any, Row any] struct {
	rows     *sql.Rows
	scan     func(*sql.Rows) (Row, error)
	entity   func(Row) Entity
	noEntity Entity
}

type Option[Entity any, Specification any, Row any] struct {
	TableName       string
	Columns         []string
	Schema          map[string]string
	PrimaryKey      string
	DatabaseManager database_manager.DatabaseManager
	Logger          logger.Logger
	Filter          func(...Specification) sq.Sqlizer
	Scan            func(rows *sql.Rows) (Row, error)
	Entity          func(Row) Entity
	Row             func(Entity) Row
	Values          func(Row) []any
}

func (i *PostgresIterator[Entity, Row]) Current() (Entity, error) {
	row, err := i.scan(i.rows)
	if err != nil {
		return i.noEntity, err
	}

	return i.entity(row), nil
}

func (i *PostgresIterator[Entity, Row]) Next() bool {
	return i.rows.Next()
}

func (r *PostgresRepository[Entity, Specification, Row]) Delete(ctx context.Context, specs ...Specification) error {
	query, args, err := sq.
		Delete(r.tableName).
		Where(r.filter(specs...)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbm.Querier(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository[Entity, Specification, Row]) Each(ctx context.Context, args common_repository.ListArgs[Specification]) (common_repository.Iterator[Entity], error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return nil, err
	}

	return &PostgresIterator[Entity, Row]{
		rows:     rows,
		scan:     r.scan,
		entity:   r.entity,
		noEntity: r.noEntity,
	}, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) Get(ctx context.Context, specs ...Specification) (Entity, error) {
	rows, err := r.query(ctx, common_repository.ListArgs[Specification]{
		Filters: specs,
		Limit:   common_specification.WithLimit(1),
	})
	if err != nil {
		return r.noEntity, err
	}

	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return r.noEntity, err
		}

		return r.entity(row), nil
	}

	return r.noEntity, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) Exist(ctx context.Context, specs ...Specification) (bool, error) {
	builder := sq.
		Select("1").
		From(r.tableName).
		Where(r.filter(specs...)).
		Limit(1)
	queryStr, queryArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
	if err != nil {
		return false, err
	}

	var exist int

	for rows.Next() {
		err := rows.Scan(&exist)
		if err == sql.ErrNoRows {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	return exist == 1, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) List(ctx context.Context, args common_repository.ListArgs[Specification]) ([]Entity, error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return r.noEntities, err
	}

	entities := []Entity{}
	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return r.noEntities, err
		}

		entities = append(entities, r.entity(row))
	}

	return entities, nil
}

// Save implements common_repository.Common_repository.
func (r *PostgresRepository[Entity, Specification, Row]) Save(ctx context.Context, entity Entity) error {
	row := r.row(entity)

	query, args, err := sq.
		Insert(r.tableName).
		Columns(r.columns...).
		Values(r.values(row)...).
		Suffix(r.upsertSuffix).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbm.Querier(ctx).ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository[Entity, Specification, Row]) Size(ctx context.Context, specs ...Specification) (uint32, error) {
	var count uint32
	var err error
	builder := sq.
		Select("COUNT(id)").
		From(r.tableName).
		Where(r.filter(specs...))
	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func New[Entity any, Specification any, Row any](opt Option[Entity, Specification, Row]) (common_repository.Repository[Entity, Specification], error) {
	r := &PostgresRepository[Entity, Specification, Row]{
		dbm:        opt.DatabaseManager,
		logger:     opt.Logger,
		filter:     opt.Filter,
		scan:       opt.Scan,
		entity:     opt.Entity,
		row:        opt.Row,
		values:     opt.Values,
		schema:     opt.Schema,
		columns:    opt.Columns,
		tableName:  opt.TableName,
		primaryKey: opt.PrimaryKey,
	}

	r.upsertSuffix = r.makeUpsertSuffix()
	if err := r.checkSchema(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) query(ctx context.Context, args common_repository.ListArgs[Specification]) (*sql.Rows, error) {
	builder := sq.
		Select(r.columns...).
		From(r.tableName).
		Where(r.filter(args.Filters...))
	builder = r.dbm.Paginate(builder, args.Limit, args.Offset)
	queryStr, queryArgs, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	return r.dbm.Querier(ctx).QueryContext(ctx, queryStr, queryArgs...)
}

func (r *PostgresRepository[Entity, Specification, Row]) makeUpsertSuffix() string {
	parts := make([]string, 0, len(r.columns))
	for _, col := range r.columns {
		parts = append(parts, fmt.Sprintf("%s = excluded.%s", col, col))
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", r.primaryKey, strings.Join(parts, ", "))
}

func (r *PostgresRepository[Entity, Specification, Row]) checkSchema() error {
	ctx := context.Background()
	builder := sq.Select("column_name", "data_type").From("information_schema.columns").Where(sq.Eq{"table_name": r.tableName})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	rows, err := r.dbm.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	var columnName, expectedType string
	for rows.Next() {
		err := rows.Scan(&columnName, &expectedType)
		if err != nil {
			return err
		}
		actualType, ok := r.schema[columnName]
		if !ok || actualType != expectedType {
			r.logger.Error("postgres/INVALID_SCHEMA", "table_name", r.tableName, "column", columnName, "expected_type", expectedType, "actual_type", actualType)
			return ErrInvalidSchema.Format(columnName, expectedType, actualType)
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
