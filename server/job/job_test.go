package job

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/platform"
	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, test.want, got)
		})
	}
}
