package manager

import (
	"github.com/harryzcy/snuuze/manager/common"
	"github.com/harryzcy/snuuze/manager/githubactions"
	"github.com/harryzcy/snuuze/manager/gomod"
	"github.com/harryzcy/snuuze/types"
)

// managers is a list of all supported package managers
var managers = map[types.PackageManager]common.Manager{
	types.PackageManagerGoMod:         gomod.New(),
	types.PackageManagerGitHubActions: githubactions.New(),
}
