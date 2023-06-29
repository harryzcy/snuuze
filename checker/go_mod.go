package checker

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/harryzcy/snuuze/cmdutil"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
)

var (
	ErrRequestFailed = errors.New("request failed")
)

type GoListOutput struct {
	Path    string
	Version string
	Time    string
	Update  *struct {
		Path    string
		Version string
		Time    string
	}
	Dir   string
	GoMod string
}

func isUpgradable_GoMod(dep types.Dependency) (types.UpgradeInfo, error) {
	// run `go list -u -m -json <dep>` to get the latest version
	envs := map[string]string{
		"GOPATH":  config.GoPath(),
		"GOPROXY": config.GoProxy(),
	}

	dir := filepath.Dir(dep.File)
	output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
		Command: []string{"go", "list", "-u", "-m", "-json", dep.Name},
		Env:     envs,
		Dir:     dir,
	})
	if err != nil {
		return types.UpgradeInfo{}, err
	}

	// parse the output
	// e.g. github.com/shurcooL/githubv4 v0.0.0-20221229060216-a8d4a561cc93 [v0.0.0-20230305132112-efb623903184]
	// e.g. github.com/shurcooL/githubv4 v0.0.0-20221229060216-a8d4a561cc93
	info := GoListOutput{}
	err = json.Unmarshal(output.Stdout.Bytes(), &info)
	if err != nil {
		return types.UpgradeInfo{}, err
	}

	if info.Update == nil {
		// no update
		return types.UpgradeInfo{
			Dependency: dep,
			Upgradable: false,
		}, nil
	}

	return types.UpgradeInfo{
		Dependency: dep,
		Upgradable: true,
		ToVersion:  info.Update.Version,
	}, nil
}

type GoModRepoInfo struct {
	Host         string
	Owner        string
	Repo         string
	Major        string // major version, e.g. v2
	Subdirectory string // subdirectory of a monorepo
}
