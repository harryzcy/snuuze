package githubactions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			domain: "github.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "https://github.com/actions/checkout@v4",
			domain: "github.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "https://example.com/actions/checkout@v4",
			domain: "example.com",
			owner:  "actions",
			repo:   "checkout",
		},
		{
			uses:   "harryzcy/github-actions/.github/workflows/linter.yml@main",
			domain: "github.com",
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
