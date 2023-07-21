package checker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRepo(t *testing.T) {
	tests := []struct {
		uses  string
		owner string
		repo  string
		err   error
	}{
		{
			uses:  "actions/checkout@v2",
			owner: "actions",
			repo:  "checkout",
		},
		{
			uses:  "harryzcy/github-actions/.github/workflows/linter.yml@main",
			owner: "harryzcy",
			repo:  "github-actions",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			owner, repo, err := parseRepo(test.uses)
			assert.Equal(t, test.owner, owner)
			assert.Equal(t, test.repo, repo)
			assert.Equal(t, test.err, err)
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
