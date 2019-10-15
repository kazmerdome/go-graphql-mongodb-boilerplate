package server

import (
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New Init new Echo Server
func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: loggerFormat()}))

	GetRoutes(e)

	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + "9090"))
	return e
}

func loggerFormat() string {
	blue := color.New(color.FgBlue).SprintFunc()
	return blue(" latency: ") + "${latency_human}" +
		// blue(" uri: ") + "${uri}" +
		// blue(" method: ") + "${method}" +
		blue(" status: ") + "${status}" +
		blue(" timeRequest: ") + "${time_rfc3339_nano}" + "\n"
}
