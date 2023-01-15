package parser

import (
	"strconv"
	"testing"

	"github.com/harryzcy/sailor/types"
	"github.com/stretchr/testify/assert"
)

func TestParseGitHubActions(t *testing.T) {
	tests := []struct {
		data []byte
		want []types.Dependency
		err  error
	}{
		{
			data: []byte(`
jobs:
  super-linter:
    name: Super Linter
    uses: harryzcy/github-actions/.github/workflows/linter.yml@main

  go-lint:
    name: Go Lint
    uses: harryzcy/github-actions/.github/workflows/golangci-lint.yml@main
`),
			want: []types.Dependency{
				{
					Name:           "harryzcy/github-actions/.github/workflows/linter.yml",
					Version:        "main",
					PackageManager: "github-actions",
				},
				{
					Name:           "harryzcy/github-actions/.github/workflows/golangci-lint.yml",
					Version:        "main",
					PackageManager: "github-actions",
				},
			},
		},
		{
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
			got, err := parseGitHubActions(test.data)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.want, got)
		})
	}
}
