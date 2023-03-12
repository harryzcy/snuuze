package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultBranch(t *testing.T) {
	branch := getDefaultBranch(".")
	assert.Equal(t, "main", branch)
}
