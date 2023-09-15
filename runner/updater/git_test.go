package updater

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitURL(t *testing.T) {
	tests := []struct {
		url         string
		expected    *gitInfo
		expectedErr error
	}{
		{
			url: "https://github.com/harryzcy/snuuze.git",
			expected: &gitInfo{
				owner: "harryzcy",
				repo:  "snuuze",
			},
		},
		{
			url:         "https://github.com",
			expected:    nil,
			expectedErr: errors.New("invalid git url: https://github.com"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			info, err := parseGitURL(test.url)
			assert.Equal(t, test.expected, info)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
