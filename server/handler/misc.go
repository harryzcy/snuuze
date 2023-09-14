package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "snuuze")
}

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
