package handler

import (
	"github.com/harryzcy/snuuze/server/job"
	"github.com/labstack/echo/v4"
)

func ListRepos(state *job.State) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		state.Lock()
		defer state.Unlock()
		return ctx.JSON(200, state.Repos)
	}
}

func ListDependencies(state *job.State) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		repoURL := ctx.QueryParam("repoURL")
		if repoURL == "" {
			return ctx.JSON(400, "repoURL is required")
		}

		state.Lock()
		defer state.Unlock()
		if _, ok := state.RepoDependencies[repoURL]; !ok {
			return ctx.JSON(404, "repoURL not found")
		}
		return ctx.JSON(200, state.RepoDependencies[repoURL])
	}
}
