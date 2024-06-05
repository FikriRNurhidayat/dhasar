package http_server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HealthCheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
