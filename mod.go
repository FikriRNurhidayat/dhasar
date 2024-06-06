package dhasar

import (
	"github.com/fikrirnurhidayat/x/logger"
	echo "github.com/labstack/echo/v4"
)

type Module struct {
	Dependency *RootDependency
}

func (m *Module) Wire(dependency *RootDependency) {
	m.Dependency = dependency
}

type RootDependency struct {
	Logger logger.Logger
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

func NewModule(logger logger.Logger) *RootDependency {
	return &RootDependency{
		Logger: logger,
	}
}
