package matcher

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestMatchFile(t *testing.T) {
	tests := []struct {
		path string
		want types.PackageManager
	}{
		{
			path: "go.mod",
			want: types.GoMod,
		},
		{
			path: "server/go.mod",
			want: types.GoMod,
		},
		{
			path: ".github/workflows/go.yml",
			want: types.GitHubActions,
		},
		{
			path: ".github/workflows/go.yaml",
			want: types.GitHubActions,
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
