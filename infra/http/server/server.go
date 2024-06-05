package http_server

import (
	"io"

	dhasar_container "github.com/fikrirnurhidayat/dhasar/container"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Server struct {
	Echo      *echo.Echo
	Port      uint
	Container *dhasar_container.Container
}

type Option struct {
	Container   *dhasar_container.Container
	HealthCheck echo.HandlerFunc
	Bootstrap   func(*Server) error
}

func New(opt *Option) (*Server, error) {
	server := &Server{
		Port:      viper.GetUint("server.port"),
		Echo:      echo.New(),
		Container: opt.Container,
	}

	if opt.HealthCheck == nil {
		opt.HealthCheck = server.HealthCheck
	}

	server.Echo.Logger.SetOutput(io.Discard)
	server.Echo.Logger.SetLevel(log.OFF)
	server.Echo.HideBanner = true
	server.Echo.HidePort = true
	server.Echo.DisableHTTP2 = true
	server.Echo.Use(middleware.Secure())
	server.Echo.Use(middleware.Timeout())
	server.Echo.Use(middleware.RequestID())
	server.Echo.Use(server.RequestLogger())
	server.Echo.Use(middleware.Recover())
	server.Echo.GET("/health", opt.HealthCheck)
	server.Echo.HTTPErrorHandler = server.HTTPErrorHandler

	if err := opt.Bootstrap(server); err != nil {
		return nil, err
	}

	return server, nil
}
