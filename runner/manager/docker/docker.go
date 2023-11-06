package docker

import (
	"bytes"
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
	result, err := parser.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	dependencies := make([]*types.Dependency, 0)
	for _, child := range result.AST.Children {
		if child.Value != "FROM" {
			continue
		}

		image, version, versionType := parseDockerfileFromDirective(child.Original)
		dependencies = append(dependencies, &types.Dependency{
			File:           match.File,
			Name:           image,
			Version:        version,
			Indirect:       false,
			PackageManager: types.PackageManagerDocker,
			Position: types.Position{
				Line: child.StartLine,
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
// FROM [--platform=<platform>] <image>[:<tag>] [AS <name>] or
// FROM [--platform=<platform>] <image>[@<digest>] [AS <name>]
//
// Returns the image, version, and version type (tag or digest).
func parseDockerfileFromDirective(value string) (image, version, versionType string) {
	parts := strings.Split(value, " ")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "FROM" || strings.HasPrefix(part, "--") {
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

	tags, err := getDockerImageTags(dep.Name)
	if err != nil {
		return nil, err
	}

	latest, err := common.GetLatestTag(dep.Name, tags, dep.Version, true)
	if err != nil {
		return nil, err
	}
	if latest != dep.Version {
		info.Upgradable = true
		info.ToVersion = latest
	}
	return info, nil
}

func getDockerImageTags(name string) ([]string, error) {
	endpoints, image := parseImageName(name)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://" + endpoints + "/v2/" + image + "/tags/list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if endpoints == "index.docker.io" {
		token, err := getDockerHubToken(client, image)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
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

	return tagsResponse.Tags, nil
}

type dockerHubTokenResponse struct {
	Token string `json:"token"`
}

// getDockerHubToken gets a token from auth.docker.io/token.
func getDockerHubToken(client *http.Client, image string) (token string, err error) {
	url := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + image + ":pull"
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", types.ErrRequestFailed
	}

	var tokenResponse dockerHubTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Token, nil
}

func parseImageName(name string) (endpoint, image string) {
	imageParts := strings.Split(name, "/")
	if len(imageParts) == 1 {
		endpoint = "index.docker.io"
		image = "library/" + imageParts[0]
	} else {
		if strings.Contains(imageParts[0], ".") {
			// first part is a domain
			endpoint = imageParts[0]
			image = strings.Join(imageParts[1:], "/")
		} else {
			// first part is a user
			endpoint = "index.docker.io"
			image = strings.Join(imageParts, "/")
		}
	}

	return endpoint, image
}
