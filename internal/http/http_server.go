package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func InitServer(handler func(echo.Context) error) *echo.Echo {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Algorithm up")
	})

	e.GET("/recommandation/:username", handler)

	e.Start(":8080")

	return e

}
