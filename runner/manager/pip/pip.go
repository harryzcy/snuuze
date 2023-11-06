package pip

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/types"
)

type PipManager struct{}

func New() common.Manager {
	return &PipManager{}
}

func (m *PipManager) Name() types.PackageManager {
	return types.PackageManagerPip
}

func (m *PipManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "requirements.txt"
}

func (m *PipManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	text := string(data)

	lines := strings.Split(text, "\n")

	dependencies := make([]*types.Dependency, 0)
	for lineNum, line := range lines {
		parser := NewParser(match.File, line, lineNum)
		dep, err := parser.Parse()
		if err != nil {
			continue
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func (m *PipManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *PipManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

type pipJson struct {
	Releases map[string]interface{} `json:"releases"`
}

func (m *PipManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	info := &types.UpgradeInfo{
		Dependency: dep,
	}

	if dep.Version == "" {
		return info, nil
	}

	if constraints, ok := dep.Extra["constraints"].([][2]string); !ok || len(constraints) > 1 {
		// Dependencies with multiple constraints are not supported
		return info, nil
	}

	if len(dep.Version) > 2 && (dep.Version[:2] != "==" || dep.Version[:2] == ">=" || dep.Version[:1] == ">") {
		// Only support exact version match or greater than match
		return info, nil
	}

	currentVersion := strings.TrimSpace(dep.Version[2:])
	versions, err := getPipPackageVersions(dep.Name)
	if err != nil {
		return nil, err
	}

	latest, err := common.GetLatestTag(dep.Name, versions, currentVersion, true)
	if err != nil {
		return nil, err
	}
	if latest != dep.Version {
		info.Upgradable = true
		info.ToVersion = latest
	}
	return info, nil
}

func getPipPackageVersions(name string) ([]string, error) {
	jsonURL := "https://pypi.org/pypi/" + name + "/json"
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(jsonURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, types.ErrRequestFailed
	}

	var data pipJson
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	for version := range data.Releases {
		versions = append(versions, version)
	}
	return versions, nil
}
