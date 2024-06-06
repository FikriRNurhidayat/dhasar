package dhasar

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fikrirnurhidayat/x/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	Echo      *echo.Echo
	Port      uint
	Container *Container
	Logger    logger.Logger
}

type HTTPServerOption struct {
	Container   *Container
	HealthCheck echo.HandlerFunc
	Bootstrap   func(*HTTPServer) error
}

func (s *HTTPServer) HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (server *HTTPServer) HTTPErrorHandler(err error, c echo.Context) {
	if val, ok := err.(*Error); ok {
		c.JSON(val.Code, echo.Map{
			"error": val,
		})

		return
	}

	code := http.StatusInternalServerError
	if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
	}

	if code == http.StatusNotFound {
		c.JSON(code, echo.Map{
			"error": ErrNotFound.Format(c.Request().Method, c.Request().URL),
		})

		return
	}

	c.JSON(code, echo.Map{
		"error": ErrInternalServer,
	})
}

func (s *HTTPServer) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			args := []any{
				logger.String("request-id", v.RequestID),
				logger.String("method", v.Method),
				logger.String("uri", v.URI),
				logger.Int("status", v.Status),
				logger.String("took", fmt.Sprintf("%d ms", v.Latency.Milliseconds())),
			}

			serverLogger := Get[logger.Logger](s.Container, "Logger")

			if v.Error == nil {
				serverLogger.Info("http/OK", args...)
			} else {
				if v.Status == http.StatusNotFound {
					serverLogger.Warn("http/ROUTE_NOT_FOUND", args...)
					return nil
				}
				if val, ok := v.Error.(*Error); ok {
					serverLogger.Warn(fmt.Sprintf("http/%s", val.Reason), args...)
					return nil
				}
				args = append(args, logger.String("error", v.Error.Error()))
				serverLogger.Warn("http/INTERNAL_SERVER_ERROR", args...)
			}
			return nil
		},
		HandleError:      false,
		LogLatency:       true,
		LogProtocol:      false,
		LogRemoteIP:      false,
		LogHost:          false,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       false,
		LogRoutePath:     false,
		LogRequestID:     true,
		LogReferer:       false,
		LogUserAgent:     false,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	})
}

func NewHTTPServer(opt *HTTPServerOption) (*HTTPServer, error) {
	server := &HTTPServer{
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
