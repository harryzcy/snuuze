package manager

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFile(t *testing.T) {
	tests := []struct {
		path string
		want PackageManager
	}{
		{
			path: "go.mod",
			want: GoMod,
		},
		{
			path: "server/go.mod",
			want: GoMod,
		},
		{
			path: ".github/workflows/go.yml",
			want: GitHubActions,
		},
		{
			path: ".github/workflows/go.yaml",
			want: GitHubActions,
		},
		{
			path: "unknown",
			want: "",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, ok := matchFile(test.path)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.want != "", ok)
		})
	}
}
