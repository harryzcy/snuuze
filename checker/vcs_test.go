package checker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
	tests := []struct {
		sortedTags []string
		currentTag string
		want       string
	}{
		{
			sortedTags: []string{"v1.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1.2.3",
			want:       "v1.2.5",
		},
		{
			sortedTags: []string{"v1.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1.2",
			want:       "v1.2",
		},
		{
			sortedTags: []string{"v1.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1",
			want:       "v1",
		},
		{
			sortedTags: []string{"v2.2.5", "v1.2.4", "v1.2.3"},
			currentTag: "v1",
			want:       "v2",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := getLatestTag(test.sortedTags, test.currentTag)
			assert.Equal(t, test.want, got)
		})
	}
}
