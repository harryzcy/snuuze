package server

import (
	"errors"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const Port = "1323"

func Init() (*State, error) {
	authType := config.GetHostingConfig().GitHub.AuthType
	if authType != "github-app" {
		return nil, errors.New("only GitHub App is supported for running as a server")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/ping", handler.Ping)

	repos, err := LoadRepos()
	if err != nil {
		return nil, err
	}

	err = e.Start(":" + Port)
	if err != nil {
		return nil, err
	}

	return &State{
		Repos: repos,
	}, nil
}

func LoadRepos() ([]platform.Repo, error) {
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

	return repos, nil
}
