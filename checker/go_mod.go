package checker

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if dep.Indirect || strings.HasPrefix(dep.Version, "v0.0.0-") || strings.HasSuffix(dep.Version, "-incompatible") {
		// don't check major version for indirect dependencies and pseudo versions
		// and incompatible versions
		return isUpgradable_GoMod_Minor(dep)
	}

	return isUpgradable_GoMod_Major(dep)
}

type GoModResponse struct {
	Version string `json:"Version"`
}

func isUpgradable_GoMod_Minor(dep types.Dependency) (UpgradeInfo, error) {
	// go proxy requires lowercase module name
	url := GOPROXY + "/" + strings.ToLower(dep.Name) + "/@latest"
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

func isUpgradable_GoMod_Major(dep types.Dependency) (UpgradeInfo, error) {
	modInfo, err := parseGoModRepo(dep.Name)
	if err != nil {
		return UpgradeInfo{}, err
	}

	var source VSCSource
	if modInfo.Host == "github.com" {
		source = &GitHubSource{}
	} else {
		fmt.Println("unsupported host for go modules: " + modInfo.Host)
		// TODO: support other hosts
		return UpgradeInfo{}, nil
	}

	tags, err := source.ListTags(&ListTagsParameters{
		Owner:  modInfo.Owner,
		Repo:   modInfo.Repo,
		Prefix: modInfo.Subdirectory,
	})
	if err != nil {
		return UpgradeInfo{}, err
	}

	latest := getLatestTag(tags, dep.Version)

	info := UpgradeInfo{
		Dependency: dep,
	}
	if latest != dep.Version {
		info.Upgradable = true
		info.ToVersion = latest
		return info, nil
	}
	return info, nil
}

type GoModRepoInfo struct {
	Host         string
	Owner        string
	Repo         string
	Major        string // major version, e.g. v2
	Subdirectory string // subdirectory of a monorepo
}

func parseGoModRepo(module string) (*GoModRepoInfo, error) {
	if info, ok, err := parseGopkgIn(module); ok {
		return info, err
	}
	if info, ok, err := parseGolangX(module); ok {
		return info, err
	}
	if info, ok, err := parseGoogleGolangOrg(module); ok {
		return info, err
	}

	parts := strings.Split(module, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid module: %s", module)
	}
	info := &GoModRepoInfo{
		Host:  parts[0],
		Owner: parts[1],
		Repo:  parts[2],
	}
	lastIdx := len(parts) - 1
	if isMajorVersion(parts[lastIdx]) {
		info.Major = parts[lastIdx]
		lastIdx--
	}
	info.Subdirectory = strings.Join(parts[3:lastIdx+1], "/")
	return info, nil
}

func parseGopkgIn(module string) (*GoModRepoInfo, bool, error) {
	ok := false
	if strings.HasPrefix(module, "gopkg.in/") {
		ok = true
		module = strings.TrimPrefix(module, "gopkg.in/")
		parts := strings.Split(module, ".")
		if len(parts) < 2 {
			return nil, ok, fmt.Errorf("invalid module: %s", module)
		}
		version := parts[1]

		parts = strings.Split(parts[0], "/")
		var owner, repo string
		if len(parts) == 2 {
			owner = parts[0]
			repo = parts[1]
		} else if len(parts) == 1 {
			owner = "go-" + parts[0]
			repo = parts[0]
		} else {
			return nil, ok, fmt.Errorf("invalid module: %s", module)
		}
		info := &GoModRepoInfo{
			Host:  "github.com",
			Owner: owner,
			Repo:  repo,
			Major: version,
		}
		return info, ok, nil
	}
	return nil, ok, nil
}

func parseGolangX(module string) (*GoModRepoInfo, bool, error) {
	if strings.HasPrefix(module, "golang.org/x/") {
		repo := strings.TrimPrefix(module, "golang.org/x/")
		info := &GoModRepoInfo{
			Host:  "github.com", // Repos from go.googlesource.com are mirrored to github.com/golang
			Owner: "golang",
			Repo:  repo,
		}
		return info, true, nil
	}
	return nil, false, nil
}

func parseGoogleGolangOrg(module string) (*GoModRepoInfo, bool, error) {
	if strings.HasPrefix(module, "google.golang.org/") {
		repo := strings.TrimPrefix(module, "google.golang.org/")

		var info *GoModRepoInfo
		switch repo {
		case "appengine":
			info = &GoModRepoInfo{
				Host:  "github.com",
				Owner: "golang",
				Repo:  "appengine",
			}
		case "genproto":
			info = &GoModRepoInfo{
				Host:  "github.com",
				Owner: "googleapis",
				Repo:  "go-genproto",
			}
		case "grpc":
			info = &GoModRepoInfo{
				Host:  "github.com",
				Owner: "grpc",
				Repo:  "grpc-go",
			}
		case "protobuf":
			info = &GoModRepoInfo{
				Host:  "github.com",
				Owner: "protocolbuffers",
				Repo:  "protobuf-go",
			}
		default:
			return nil, false, fmt.Errorf("unsupported google.golang.org module: %s", module)
		}
		return info, true, nil
	}
	return nil, false, nil
}

// isMajorVersion returns true if the version is a major version, e.g. v2
func isMajorVersion(version string) bool {
	if version[0] != 'v' {
		return false
	}
	for _, c := range version[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
