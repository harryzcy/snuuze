package handler

import (
	"github.com/harryzcy/snuuze/server/job"
	"github.com/labstack/echo/v4"
)

func ListRepos(state *job.State) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(200, state.Repos)
	}
}
