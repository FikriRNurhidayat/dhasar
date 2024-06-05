package postgres_repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	sql_database_manager "github.com/fikrirnurhidayat/dhasar/manager/database/sql"
	"github.com/fikrirnurhidayat/dhasar/repository"
	"github.com/fikrirnurhidayat/dhasar/specification"
	"github.com/fikrirnurhidayat/x/logger"
)

type PostgresRepository[Entity any, Specification any, Row any] struct {
	dbm          sql_database_manager.SQLDatabaseManager
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
	TableName          string
	Columns            []string
	Schema             map[string]string
	PrimaryKey         string
	SQLDatabaseManager sql_database_manager.SQLDatabaseManager
	Logger             logger.Logger
	Filter             func(...Specification) sq.Sqlizer
	Scan               func(rows *sql.Rows) (Row, error)
	Entity             func(Row) Entity
	Row                func(Entity) Row
	Values             func(Row) []any
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

func (r *PostgresRepository[Entity, Specification, Row]) Each(ctx context.Context, args dhasar_repository.ListArgs[Specification]) (dhasar_repository.Iterator[Entity], error) {
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
	rows, err := r.query(ctx, dhasar_repository.ListArgs[Specification]{
		Filters: specs,
		Limit:   dhasar_specification.WithLimit(1),
	})

	if err != nil {
		return r.noEntity, err
	}

	defer rows.Close()

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

	defer rows.Close()

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

func (r *PostgresRepository[Entity, Specification, Row]) List(ctx context.Context, args dhasar_repository.ListArgs[Specification]) ([]Entity, error) {
	rows, err := r.query(ctx, args)
	if err != nil {
		return r.noEntities, err
	}

	entities := []Entity{}

	defer rows.Close()

	for rows.Next() {
		row, err := r.scan(rows)
		if err != nil {
			return r.noEntities, err
		}

		entities = append(entities, r.entity(row))
	}

	return entities, nil
}

// Save implements dhasar_repository.Common_repository.
func (r *PostgresRepository[Entity, Specification, Row]) Save(ctx context.Context, entity Entity) error {
	row := r.row(entity)

	query, args, err := sq.
		Insert(r.tableName).
		Columns(r.columns...).
		Values(r.values(row)...).
		PlaceholderFormat(sq.Dollar).
		Suffix(r.upsertSuffix).
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

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func New[Entity any, Specification any, Row any](opt Option[Entity, Specification, Row]) (dhasar_repository.Repository[Entity, Specification], error) {
	r := &PostgresRepository[Entity, Specification, Row]{
		dbm:        opt.SQLDatabaseManager,
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

	return r, nil
}

func (r *PostgresRepository[Entity, Specification, Row]) query(ctx context.Context, args dhasar_repository.ListArgs[Specification]) (*sql.Rows, error) {
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
