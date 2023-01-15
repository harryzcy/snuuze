package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGoMod(t *testing.T) {
	path := "../go.mod"
	data, err := os.ReadFile(path)
	assert.Nil(t, err)

	dependencies, err := parseGoMod(path, data)
	assert.Nil(t, err)
	assert.NotZero(t, len(dependencies))
}
