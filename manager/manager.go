package manager

import (
	"github.com/harryzcy/snuuze/manager/common"
	"github.com/harryzcy/snuuze/manager/githubactions"
	"github.com/harryzcy/snuuze/manager/gomod"
)

// managers is a list of all supported package managers
var managers = []common.Manager{
	gomod.New(),
	githubactions.New(),
}
