package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRest(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "<3")
	})
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
