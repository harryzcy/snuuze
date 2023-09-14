package server

import (
	"github.com/harryzcy/snuuze/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const Port = "1323"

func Init() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/ping", handler.Ping)

	e.Logger.Fatal(e.Start(":" + Port))
}
