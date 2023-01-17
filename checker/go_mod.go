package checker

import (
	"encoding/json"
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

type GoModResponse struct {
	Version string `json:"Version"`
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
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return UpgradeInfo{}, err
	}

	response := GoModResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UpgradeInfo{}, err
	}

	info := UpgradeInfo{
		Dependency: dep,
	}
	if response.Version != dep.Version {
		info.Upgradable = true
		info.ToVersion = response.Version
		return info, nil
	}
	return info, nil
}
