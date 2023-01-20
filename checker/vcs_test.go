package checker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortTags(t *testing.T) {
	tests := []struct {
		tags []string
		want []string
	}{
		{
			tags: []string{"v1.2.3", "v1.20.3", "v1.2.30"},
			want: []string{"v1.20.3", "v1.2.30", "v1.2.3"},
		},
		{
			tags: []string{"v2.2.3", "v1.20.3", "v1.2.30"},
			want: []string{"v2.2.3", "v1.20.3", "v1.2.30"},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := sortTags(test.tags)
			assert.Equal(t, test.want, got)
		})
	}
}
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

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag  string
		want []int
		err  error
	}{
		{tag: "v1.2.3", want: []int{1, 2, 3}},
		{tag: "v1.20.3", want: []int{1, 20, 3}},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := parseTag(test.tag)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.err, err)
		})
	}
}
