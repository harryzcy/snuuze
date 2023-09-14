package server

import (
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const Port = "1323"

func Init() {
	authType := config.GetHostingConfig().GitHub.AuthType
	if authType != "github_app" {
		fmt.Fprintln(os.Stderr, "Only GitHub App is supported for running as a server")
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/ping", handler.Ping)

	e.Logger.Fatal(e.Start(":" + Port))
}
