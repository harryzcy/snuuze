package githubactions

import (
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/harryzcy/snuuze/types"
)

type minimalGitHubActions struct {
	Jobs map[string]job
}

type job struct {
	Uses  string
	Steps []step
}

type step struct {
	Uses string
}

func parseGitHubActions(path string, data []byte) ([]*types.Dependency, error) {
	var content minimalGitHubActions
	err := yaml.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	dependencies := make([]*types.Dependency, 0)
	for _, job := range content.Jobs {
		if job.Uses != "" {
			// reusable workflow
			if dependency, ok := parseWorkflow(path, job.Uses); ok {
				dependencies = append(dependencies, dependency)
				continue
			}
		}

		for _, step := range job.Steps {
			if step.Uses != "" {
				dependency, ok := parseWorkflow(path, step.Uses)
				if ok {
					dependencies = append(dependencies, dependency)
				}
			}
		}
	}

	return dependencies, nil
}

func parseWorkflow(path, uses string) (*types.Dependency, bool) {
	if uses == "" {
		return nil, false
	}

	parts := strings.Split(uses, "@")
	name := parts[0]
	if strings.HasPrefix(name, "./") {
		// local repository, skip
		return nil, false
	}

	version := ""
	if len(parts) > 1 {
		version = parts[1]
	}

	return &types.Dependency{
		File:           path,
		Name:           name,
		Version:        version,
		PackageManager: types.PackageManagerGitHubActions,
	}, true
}
