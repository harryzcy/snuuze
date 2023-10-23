package gomod

import (
	"testing"

	"github.com/harryzcy/snuuze/types"
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
	want := []*types.Dependency{
		{
			File:           "go.mod",
			Name:           "github.com/docker/docker",
			Version:        "v20.10.11+incompatible",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      6,
				ColStart:  3,
				ColEnd:    50,
				ByteStart: 64,
				ByteEnd:   111,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/spf13/viper",
			Version:        "v1.14.0",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      7,
				ColStart:  3,
				ColEnd:    33,
				ByteStart: 114,
				ByteEnd:   144,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/stretchr/testify",
			Version:        "v1.8.1",
			Indirect:       false,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      8,
				ColStart:  3,
				ColEnd:    37,
				ByteStart: 147,
				ByteEnd:   181,
			},
		},
		{
			File:           "go.mod",
			Name:           "github.com/davecgh/go-spew",
			Version:        "v1.1.1",
			Indirect:       true,
			PackageManager: "go-mod",
			Position: types.Position{
				Line:      12,
				ColStart:  3,
				ColEnd:    36,
				ByteStart: 200,
				ByteEnd:   233,
			},
		},
	}

	m := New()
	dependencies, err := m.Parse(types.Match{File: path}, data)
	assert.Nil(t, err)
	assert.Equal(t, want, dependencies)
}
