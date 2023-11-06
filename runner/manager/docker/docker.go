package docker

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/harryzcy/snuuze/runner/manager/common"
	"github.com/harryzcy/snuuze/types"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type DockerManager struct{}

func New() common.Manager {
	return &DockerManager{}
}

func (m *DockerManager) Name() types.PackageManager {
	return types.PackageManagerDocker
}

func (m *DockerManager) Match(path string) bool {
	filename := filepath.Base(path)
	return filename == "Dockerfile"
}

func (m *DockerManager) Parse(match types.Match, data []byte) ([]*types.Dependency, error) {
	p := &parser.DirectiveParser{}
	directives, err := p.ParseAll(data)
	if err != nil {
		return nil, err
	}

	dependencies := make([]*types.Dependency, 0)
	for _, directive := range directives {
		if directive.Name != "FROM" {
			continue
		}

		image, version, versionType := parseDockerfileFromDirective(directive.Value)
		dependencies = append(dependencies, &types.Dependency{
			File:           match.File,
			Name:           image,
			Version:        version,
			Indirect:       false,
			PackageManager: types.PackageManagerDocker,
			Position: types.Position{
				Line:     directive.Location[0].Start.Line,
				ColStart: directive.Location[0].Start.Character,
				ColEnd:   directive.Location[0].End.Character,
			},
			Extra: map[string]interface{}{
				"versionType": versionType,
			},
		})
	}

	return dependencies, nil
}

// parseDockerfileFromDirective parses a FROM directive in a Dockerfile.
//
// The value of the directive is expected to be in the format of:
// [--platform=<platform>] <image>[:<tag>] [AS <name>] or
// [--platform=<platform>] <image>[@<digest>] [AS <name>]
//
// Returns the image, version, and version type (tag or digest).
func parseDockerfileFromDirective(value string) (image, version, versionType string) {
	parts := strings.Split(value, " ")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.HasPrefix(part, "--") {
			continue
		}

		if strings.Contains(part, ":") {
			image = strings.Split(part, ":")[0]
			version = strings.Split(part, ":")[1]
			versionType = "tag"
		} else if strings.Contains(part, "@") {
			image = strings.Split(part, "@")[0]
			version = strings.Split(part, "@")[1]
			versionType = "digest"
		} else {
			image = part
		}
	}

	return image, version, versionType
}

func (m *DockerManager) FindDependencies(matches []types.Match) ([]*types.Dependency, error) {
	return common.FindDependencies(m, matches)
}

func (m *DockerManager) ListUpgrades(matches []types.Match) ([]*types.UpgradeInfo, error) {
	return common.ListUpgrades(m, matches)
}

type dockerTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (m *DockerManager) IsUpgradable(dep types.Dependency) (*types.UpgradeInfo, error) {
	info := &types.UpgradeInfo{
		Dependency: dep,
	}

	endpoints, image := parseImageName(dep.Name)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := endpoints + "/v2/" + image + "/tags/list"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, types.ErrRequestFailed
	}

	var tagsResponse dockerTagsResponse
	err = json.NewDecoder(resp.Body).Decode(&tagsResponse)
	if err != nil {
		return nil, err
	}

	latest, err := common.GetLatestTag(dep.Name, tagsResponse.Tags, dep.Version, true)
	if err != nil {
		return nil, err
	}
	if latest != dep.Version {
		info.Upgradable = true
		info.ToVersion = latest
	}
	return info, nil
}

func parseImageName(name string) (endpoint, image string) {
	imageParts := strings.Split(name, "/")
	if len(imageParts) == 1 {
		endpoint = "https://index.docker.io"
		image = imageParts[0]
	} else {
		if strings.Contains(imageParts[0], ".") {
			// first part is a domain
			endpoint = "https://" + imageParts[0]
			image = strings.Join(imageParts[1:], "/")
		} else {
			// first part is a user
			endpoint = "https://index.docker.io"
			image = strings.Join(imageParts, "/")
		}
	}

	return endpoint, image
}
