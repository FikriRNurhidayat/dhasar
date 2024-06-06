package dhasar

import echo "github.com/labstack/echo/v4"

type Controller interface {
	Register(*echo.Echo)
}
