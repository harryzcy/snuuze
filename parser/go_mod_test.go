package parser

import (
	"testing"

	"github.com/harryzcy/sailor/types"
	"github.com/stretchr/testify/assert"
)

func TestParseGoMod(t *testing.T) {
	path := "go.mod"
	data := []byte(`module github.com/harryzcy/test-module

	go 1.19
	
	require (
		github.com/docker/docker v20.10.11+incompatible
		github.com/spf13/viper v1.14.0
		github.com/stretchr/testify v1.8.1
	)
	
	require (
		github.com/davecgh/go-spew v1.1.1 // indirect
	)`)
	want := []types.Dependency{
		{
			Name:           "github.com/docker/docker",
			Version:        "v20.10.11+incompatible",
			Indirect:       false,
			PackageManager: "gomod",
		},
		{
			Name:           "github.com/spf13/viper",
			Version:        "v1.14.0",
			Indirect:       false,
			PackageManager: "gomod",
		},
		{
			Name:           "github.com/stretchr/testify",
			Version:        "v1.8.1",
			Indirect:       false,
			PackageManager: "gomod",
		},
		{
			Name:           "github.com/davecgh/go-spew",
			Version:        "v1.1.1",
			Indirect:       true,
			PackageManager: "gomod",
		},
	}

	dependencies, err := parseGoMod(path, data)
	assert.Nil(t, err)
	assert.Equal(t, want, dependencies)
}
