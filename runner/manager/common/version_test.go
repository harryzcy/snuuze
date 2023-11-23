package common

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestTag(t *testing.T) {
	tests := []struct {
		tags       []string
		currentTag string
		delimiter  string
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
			want:       "v1",
		},
		{
			tags:       []string{"1.2.5"},
			currentTag: "v1",
			want:       "1",
		},
		{
			// alpine
			tags:       []string{"3.8.4", "20230901", "3.8.3", "3.8.2"},
			currentTag: "3.8.3",
			want:       "3.8.4",
		},
		{
			tags:       []string{"1.20-alpine", "1.20", "1.21", "1.21.1"},
			currentTag: "1.20",
			want:       "1.21",
		},
		{
			tags:       []string{"1.20-alpine", "1.20", "1.21-alpine"},
			currentTag: "1.20-alpine",
			want:       "1.21-alpine",
			delimiter:  "-",
		},
		// {
		// 	tags: []string{
		// 		"1.20-alpine", "1.20-alpine3.17", "1.20-alpine3.18",
		// 		"1.21-alpine", "1.21-alpine3.18", "1.21-alpine3.17",
		// 	},
		// 	currentTag: "1.20-alpine3.17",
		// 	want:       "1.21-alpine3.18",
		// 	delimiter:  "-",
		// },
		// {
		// 	tags:       []string{"3.17.0", "3.17", "3.18", "3.18.0"},
		// 	currentTag: "3.17",
		// 	want:       "3.18",
		// 	delimiter:  "-",
		// },
		// {
		// 	tags:       []string{"3.17.0", "3.17", "3.18", "3.18.0"},
		// 	currentTag: "3.17.0",
		// 	want:       "3.18.0",
		// 	delimiter:  "-",
		// },
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := GetLatestTag(&GetLatestTagInput{
				DepName:    "test",
				Tags:       test.tags,
				CurrentTag: test.currentTag,
				AllowMajor: true,
				Delimiter:  test.delimiter,
			})
			assert.Equal(t, test.want, got)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
