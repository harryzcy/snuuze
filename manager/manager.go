package manager

import (
	"github.com/harryzcy/snuuze/manager/common"
	"github.com/harryzcy/snuuze/manager/githubactions"
	"github.com/harryzcy/snuuze/manager/gomod"
)

var managers = []common.Manager{
	gomod.New(),
	githubactions.New(),
}
