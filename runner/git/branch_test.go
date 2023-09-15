package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultBranch(t *testing.T) {
	branch := GetDefaultBranch(".")
	assert.Equal(t, "main", branch)
}
