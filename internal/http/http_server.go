package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func InitServer() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Echo")
	})

	e.Start(":8080")

}
