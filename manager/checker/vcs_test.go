package checker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
	tests := []struct {
		tags       []string
		currentTag string
		want       string
		wantErr    bool
	}{
		{
			tags:       []string{"v1.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1.2.3",
			want:       "v1.2.5",
		},
		{
			tags:       []string{"v1.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1",
			want:       "v1",
		},
		{
			tags:       []string{"v2.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1",
			want:       "v2",
		},
		{
			tags:       []string{"v1.2.4", "v2.2.5", "v1.2.3"},
			currentTag: "v1",
			want:       "v2",
		},
		{
			tags:       []string{"v1.2.4", "v2.2.5-beta", "v1.2.3"},
			currentTag: "v1",
			want:       "v1",
		},
		{
			tags:       []string{"v1.2.4", "v2.2.5-beta", "v1.2.3"},
			currentTag: "v2",
			want:       "v2",
		},
		{
			currentTag: "invalid",
			want:       "",
			wantErr:    true,
		},
		{
			tags:       []string{"invalid"},
			currentTag: "v1",
			want:       "",
			wantErr:    true,
		},
		{
			tags:       []string{"1.2.5"},
			currentTag: "v1",
			want:       "1",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := getLatestTag(test.tags, test.currentTag, true)
			assert.Equal(t, test.want, got)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
