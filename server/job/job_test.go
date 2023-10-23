package job

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateDependencyIndex(t *testing.T) {
	state := &State{
		RepoDependencies: make(map[string]map[types.PackageManager][]*types.Dependency),
		ReverseDependencyIndex: make(map[string]struct {
			Dependency *types.Dependency
			Repos      []platform.Repo
		}),
	}

	repo := platform.Repo{
		URL: "https://github.com/harryzcy/snuuze",
	}
	deps1 := map[types.PackageManager][]*types.Dependency{
		types.PackageManagerGoMod: {
			{
				Name:    "github.com/spf13/viper",
				Version: "v1.16.0",
			},
			{
				Name:    "github.com/stretchr/testify",
				Version: "v1.8.4",
			},
		},
	}
	updateDependencyIndex(state, repo, deps1)
	assert.Len(t, state.RepoDependencies, 1)
	assert.Equal(t, deps1, state.RepoDependencies[repo.URL])
	assert.Len(t, state.ReverseDependencyIndex, 2)
	assert.Contains(t, state.ReverseDependencyIndex, deps1[types.PackageManagerGoMod][0].Hash())
	assert.Contains(t, state.ReverseDependencyIndex, deps1[types.PackageManagerGoMod][1].Hash())

	deps2 := map[types.PackageManager][]*types.Dependency{
		types.PackageManagerGoMod: {
			{
				Name:    "github.com/labstack/echo/v4",
				Version: "v4.11.1",
			},
			{
				Name:    "github.com/stretchr/testify",
				Version: "v1.8.4",
			},
		},
	}
	updateDependencyIndex(state, repo, deps2)
	assert.Len(t, state.RepoDependencies, 1)
	assert.Equal(t, deps2, state.RepoDependencies[repo.URL])
	assert.Len(t, state.ReverseDependencyIndex, 2)
	assert.Contains(t, state.ReverseDependencyIndex, deps2[types.PackageManagerGoMod][0].Hash())
	assert.Contains(t, state.ReverseDependencyIndex, deps2[types.PackageManagerGoMod][1].Hash())
	assert.Len(t, state.ReverseDependencyIndex[deps2[types.PackageManagerGoMod][0].Hash()].Repos, 1)
	assert.Len(t, state.ReverseDependencyIndex[deps2[types.PackageManagerGoMod][1].Hash()].Repos, 1)
	assert.Equal(t, repo, state.ReverseDependencyIndex[deps2[types.PackageManagerGoMod][0].Hash()].Repos[0])
	assert.Equal(t, repo, state.ReverseDependencyIndex[deps2[types.PackageManagerGoMod][1].Hash()].Repos[0])
}

func TestContainRepo(t *testing.T) {
	tests := []struct {
		repos []platform.Repo
		repo  platform.Repo
		want  bool
	}{
		{
			repos: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/snuuze",
			},
			want: true,
		},
		{
			repos: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/mailbox",
				},
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/snuuze",
			},
			want: true,
		},
		{
			repos: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/random",
			},
			want: false,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := containRepo(test.repos, test.repo)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestRemoveRepo(t *testing.T) {
	tests := []struct {
		before []platform.Repo
		repo   platform.Repo
		after  []platform.Repo
	}{
		{
			before: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/snuuze",
			},
			after: []platform.Repo{},
		},
		{
			before: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/mailbox",
				},
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/snuuze",
			},
			after: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/mailbox",
				},
			},
		},
		{
			before: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
			repo: platform.Repo{
				URL: "https://github.com/harryzcy/random",
			},
			after: []platform.Repo{
				{
					URL: "https://github.com/harryzcy/snuuze",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := removeRepo(test.before, test.repo)
			assert.Equal(t, test.after, got)
		})
	}
}

func TestFlattenDependencies(t *testing.T) {
	tests := []struct {
		deps map[types.PackageManager][]*types.Dependency
		want []*types.Dependency
	}{
		{
			deps: map[types.PackageManager][]*types.Dependency{
				types.PackageManagerGoMod: {
					{
						Name:    "snuuze",
						Version: "0.0.1",
					},
				},
			},
			want: []*types.Dependency{
				{
					Name:    "snuuze",
					Version: "0.0.1",
				},
			},
		},
		{
			deps: map[types.PackageManager][]*types.Dependency{
				types.PackageManagerGoMod: {
					{
						Name:    "snuuze",
						Version: "0.0.1",
					},
					{
						Name:    "mailbox",
						Version: "1.0.0",
					},
				},
				types.PackageManagerGitHubActions: {
					{
						Name:    "actions/checkout",
						Version: "v4",
					},
				},
			},
			want: []*types.Dependency{
				{
					Name:    "snuuze",
					Version: "0.0.1",
				},
				{
					Name:    "mailbox",
					Version: "1.0.0",
				},
				{
					Name:    "actions/checkout",
					Version: "v4",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := flattenDependencies(test.deps)
			assert.Len(t, got, len(test.want))
			for i := range got {
				assert.Equal(t, *test.want[i], *got[i])
			}
		})
	}
}
