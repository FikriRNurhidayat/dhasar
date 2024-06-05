package http_server

import (
	"fmt"
	"net/http"

	dhasar_container "github.com/fikrirnurhidayat/dhasar/container"
	dhasar_errors "github.com/fikrirnurhidayat/dhasar/errors"
	"github.com/fikrirnurhidayat/x/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (server *Server) HTTPErrorHandler(err error, c echo.Context) {
	if val, ok := err.(*dhasar_errors.Error); ok {
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
			"error": dhasar_errors.ErrNotFound.Format(c.Request().Method, c.Request().URL),
		})

		return
	}

	c.JSON(code, echo.Map{
		"error": dhasar_errors.ErrInternalServer,
	})
}

func (server *Server) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			args := []any{
				logger.String("request-id", v.RequestID),
				logger.String("method", v.Method),
				logger.String("uri", v.URI),
				logger.Int("status", v.Status),
				logger.String("took", fmt.Sprintf("%d ms", v.Latency.Milliseconds())),
			}

			serverLogger, err := dhasar_container.Get[logger.Logger](server.Container, "Logger")
			if err != nil {
				return err
			}

			if v.Error == nil {
				serverLogger.Info("http/OK", args...)
			} else {
				if v.Status == http.StatusNotFound {
					serverLogger.Warn("http/ROUTE_NOT_FOUND", args...)
					return nil
				}
				if val, ok := v.Error.(*dhasar_errors.Error); ok {
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
