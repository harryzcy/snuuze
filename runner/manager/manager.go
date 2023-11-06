package manager

import (
	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/runner/manager/docker"
	"github.com/harryzcy/snuuze/runner/manager/githubactions"
	"github.com/harryzcy/snuuze/runner/manager/gomod"
	"github.com/harryzcy/snuuze/runner/manager/pip"
	"github.com/harryzcy/snuuze/types"
)

// managers is a list of all supported package managers
var managers = map[types.PackageManager]common.Manager{
	types.PackageManagerDocker:        docker.New(),
	types.PackageManagerGoMod:         gomod.New(),
	types.PackageManagerGitHubActions: githubactions.New(),
	types.PackageManagerPip:           pip.New(),
}
