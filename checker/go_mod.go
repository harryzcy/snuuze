package checker

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/harryzcy/sailor/types"
)

var (
	DEFAULT_TIMEOUT = 10 * time.Second
	GOPROXY         = os.Getenv("GOPROXY")

	ErrRequestFailed = errors.New("request failed")
)

func init() {
	if GOPROXY == "" {
		GOPROXY = "https://proxy.golang.org"
	} else {
		// use the first proxy
		GOPROXY = strings.Split(GOPROXY, ",")[0]
	}
}

func isUpgradable_GoMod(dep types.Dependency) (UpgradeInfo, error) {
	url := GOPROXY + "/" + dep.Name + "/@latest"
	client := http.Client{
		Timeout: DEFAULT_TIMEOUT,
	}

	req, err := client.Get(url)
	if err != nil {
		return UpgradeInfo{}, err
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return UpgradeInfo{}, ErrRequestFailed
	}
	versionBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return UpgradeInfo{}, err
	}

	version := string(bytes.TrimSpace(versionBytes))
	info := UpgradeInfo{
		Dependency: dep,
	}
	if version != dep.Version {
		info.Upgradable = true
		info.ToVersion = version
		return info, nil
	}
	return info, nil
}
