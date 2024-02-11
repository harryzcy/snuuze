package githubactions

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	err := config.LoadHostingConfig()
	if err != nil {
		panic(err)
	}
	t.Run()
}

func TestIsUpgradable(t *testing.T) {
	m := New()

	info, err := m.IsUpgradable(types.Dependency{
		Name:    "https://gitea.com/actions/checkout@v2",
		Version: "v2",
	})
	assert.NoError(t, err)
	assert.True(t, info.Upgradable)
}

func TestMatch(t *testing.T) {
	tests := []struct {
		path  string
		match bool
	}{
		{path: ".github/workflows/linter.yml", match: true},
		{path: ".github/workflows/linter.yaml", match: true},
		{path: ".github/workflows/linter.json", match: false},
		{path: "linter.yaml", match: false},
		{path: "linter.json", match: false},
	}

	m := New()
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := m.Match(test.path)
			assert.Equal(t, test.match, got)
		})
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		uses   string
		domain string
		owner  string
		repo   string
		err    bool
	}{
		{
			uses:   "actions/checkout@v2",
			domain: "https://github.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "https://github.com/actions/checkout@v4",
			domain: "https://github.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "https://example.com/actions/checkout@v4",
			domain: "https://example.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "harryzcy/github-actions/.github/workflows/linter.yml@main",
			domain: "https://github.com",
			owner:  "harryzcy",
			repo:   "github-actions",
		},
		{
			uses: "actions",
			err:  true,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			domain, owner, repo, err := parseRepo(test.uses)
			assert.Equal(t, test.domain, domain)
			assert.Equal(t, test.owner, owner)
			assert.Equal(t, test.repo, repo)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsSha(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{version: "v1.2.3", want: false},
		{version: "v1.20.3", want: false},
		{version: "a81bbbf8298c0fa03ea29cdc473d45769f953675", want: true},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := isSha(test.version)
			assert.Equal(t, test.want, got)
		})
	}
}
