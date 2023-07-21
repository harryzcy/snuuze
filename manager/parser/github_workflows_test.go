package parser

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestParseGitHubActions(t *testing.T) {
	tests := []struct {
		path string
		data []byte
		want []types.Dependency
		err  error
	}{
		{
			path: ".github/workflows/lint.yml",
			data: []byte(`
jobs:
  super-linter:
    name: Super Linter
    uses: harryzcy/github-actions/.github/workflows/linter.yml@main
`),
			want: []types.Dependency{
				{
					File:           ".github/workflows/lint.yml",
					Name:           "harryzcy/github-actions/.github/workflows/linter.yml",
					Version:        "main",
					PackageManager: "github-actions",
				},
			},
		},
		{
			path: ".github/workflows/release-drafter.yml",
			data: []byte(`
jobs:
  update_release_draft:
    permissions:
      contents: write
      pull-requests: write
    runs-on: ubuntu-latest
    steps:
      - uses: release-drafter/release-drafter@v5
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`),
			want: []types.Dependency{
				{
					File:           ".github/workflows/release-drafter.yml",
					Name:           "release-drafter/release-drafter",
					Version:        "v5",
					PackageManager: "github-actions",
				},
			},
		},
		{
			data: []byte(``),
			want: []types.Dependency{},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := parseGitHubActions(test.path, test.data)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.want, got)
		})
	}
}
