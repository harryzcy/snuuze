package server

import (
	"errors"
	"fmt"
	"os"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const Port = "1323"

func Run() {
	e, _, err := initialize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = e.Start(":" + Port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initialize() (*echo.Echo, *State, error) {
	authType := config.GetHostingConfig().GitHub.AuthType
	if authType != "github-app" {
		return nil, nil, errors.New("only GitHub App is supported for running as a server")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/ping", handler.Ping)

	state, err := LoadState()
	if err != nil {
		return nil, nil, err
	}
	return e, state, nil
}

func LoadState() (*State, error) {
	client, err := platform.NewClient(platform.NewClientOptions{
		Platform: platform.GitPlatformGitHub,
	})
	if err != nil {
		return nil, err
	}

	repos, err := client.ListRepos()
	if err != nil {
		return nil, err
	}

	return &State{
		Repos: repos,
	}, nil
}
