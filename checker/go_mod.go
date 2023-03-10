package checker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/harryzcy/snuuze/types"
)

var (
	DEFAULT_TIMEOUT = 10 * time.Second
	GOPROXY         = os.Getenv("GOPROXY")

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

func isUpgradable_GoMod(dep types.Dependency) (UpgradeInfo, error) {
	// run `go list -u -m -json <module>` to get the latest version]
	cmd := exec.Command("go", "list", "-u", "-m", "-json", dep.Name)
	cmd.Env = []string{
		"GOPATH=" + os.Getenv("GOPATH"),
	}
	if GOPROXY != "" {
		cmd.Env = append(cmd.Env, "GOPROXY="+GOPROXY)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case <-time.After(DEFAULT_TIMEOUT):
		if err := cmd.Process.Kill(); err != nil {
			return UpgradeInfo{}, err
		}
		return UpgradeInfo{}, ErrRequestFailed
	case err := <-done:
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return UpgradeInfo{}, err
		}
	}

	// parse the output
	// e.g. github.com/shurcooL/githubv4 v0.0.0-20221229060216-a8d4a561cc93 [v0.0.0-20230305132112-efb623903184]
	// e.g. github.com/shurcooL/githubv4 v0.0.0-20221229060216-a8d4a561cc93
	output := out.String()

	info := GoListOutput{}
	err := json.Unmarshal([]byte(output), &info)
	if err != nil {
		return UpgradeInfo{}, err
	}

	if info.Update == nil {
		// no update
		return UpgradeInfo{
			Dependency: dep,
			Upgradable: false,
		}, nil
	}

	return UpgradeInfo{
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
