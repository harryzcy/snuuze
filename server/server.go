package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/server/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const Port = "1323"

func Run() {
	e, err := initEcho()
	exitOnError(err)

	state, err := loadState()
	exitOnError(err)

	scheduler, err := startCron(state)
	exitOnError(err)

	go func() {
		fmt.Println("Listening on port " + Port)
		if err = e.Start(":" + Port); err != nil && err != http.ErrServerClosed {
			exitOnError(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = e.Shutdown(ctx)
	exitOnError(err)

	stopCron(scheduler)

	fmt.Println("Server gracefully stopped")
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// initEcho initializes the echo server and setup the routes.
func initEcho() (*echo.Echo, error) {
	authType := config.GetHostingConfig().GitHub.AuthType
	if authType != "github-app" {
		return nil, errors.New("only GitHub App is supported for running as a server")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.Index)
	e.GET("/ping", handler.Ping)

	return e, nil
}

// loadState loads the state of the server.
func loadState() (*State, error) {
	client, err := platform.NewClient(platform.NewClientOptions{
		Platform: platform.GitPlatformGitHub,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create platform client: %w", err)
	}

	repos, err := client.ListRepos()
	if err != nil {
		return nil, fmt.Errorf("failed to list repos: %w", err)
	}

	return &State{
		Repos: repos,
	}, nil
}
