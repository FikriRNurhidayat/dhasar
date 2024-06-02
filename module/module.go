package dhasar_module

import (
	"database/sql"

	sql_database_manager "github.com/fikrirnurhidayat/dhasar/manager/database/sql"
	transaction_manager "github.com/fikrirnurhidayat/dhasar/manager/transaction"
	"github.com/fikrirnurhidayat/dhasar/pkg/logger"
	echo "github.com/labstack/echo/v4"
)

type Module struct {
	Dependency *RootDependency
}

func (m *Module) Wire(dependency *RootDependency) {
	m.Dependency = dependency
}

type RootDependency struct {
	SQLDatabaseManager sql_database_manager.DatabaseManager
	TransactionManager transaction_manager.TransactionManager
	Logger             logger.Logger
}

type HTTPModule interface {
	Wire(*RootDependency)
	WireController(e *echo.Echo) error
}

type HTTPModules []HTTPModule

type CLIModule interface {
	Wire(*RootDependency)
	WireCommand()
}

type CLIModules []CLIModule

func New(db *sql.DB, logger logger.Logger) *RootDependency {
	databaseManager := sql_database_manager.New(logger, db)
	transactionManager := transaction_manager.New(logger, db)

	return &RootDependency{
		SQLDatabaseManager: databaseManager,
		TransactionManager: transactionManager,
		Logger:             logger,
	}
}
