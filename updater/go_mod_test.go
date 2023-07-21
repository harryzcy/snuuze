package updater

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGoModuleName(t *testing.T) {
	tests := []struct {
		name        string
		fromVersion string
		toVersion   string
		expect      string
		expectErr   bool
	}{
		{
			name:        "github.com/blevesearch/zapx/v12",
			fromVersion: "v12.0.0",
			toVersion:   "v12.0.1",
			expect:      "github.com/blevesearch/zapx/v12",
			expectErr:   false,
		},
		{
			name:        "github.com/blevesearch/zapx/v12",
			fromVersion: "v12.0.0",
			toVersion:   "v13.0.1",
			expect:      "github.com/blevesearch/zapx/v13",
			expectErr:   false,
		},
		{
			name:        "example.com/m",
			fromVersion: "v4.1.2+incompatible",
			toVersion:   "v4.1.3+incompatible",
			expect:      "example.com/m",
			expectErr:   false,
		},
		{
			name:        "example.com/m",
			fromVersion: "v4.1.2+incompatible",
			toVersion:   "v5.1.3+incompatible",
			expect:      "example.com/m",
			expectErr:   false,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := getGoModuleName(test.name, test.fromVersion, test.toVersion)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expect, actual)
		})
	}
}
