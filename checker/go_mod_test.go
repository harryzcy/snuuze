package checker

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGoModRepo(t *testing.T) {
	tests := []struct {
		module string
		info   *GoModRepoInfo
		err    error
	}{
		{
			module: "github.com/harryzcy/example",
			info: &GoModRepoInfo{
				Host:  "github.com",
				Owner: "harryzcy",
				Repo:  "example",
			},
		},
		{
			module: "github.com/gofiber/fiber/v2",
			info: &GoModRepoInfo{
				Host:  "github.com",
				Owner: "gofiber",
				Repo:  "fiber",
				Major: "v2",
			},
		},
		{
			module: "github.com/aws/aws-sdk-go-v2/service/s3",
			info: &GoModRepoInfo{
				Host:         "github.com",
				Owner:        "aws",
				Repo:         "aws-sdk-go-v2",
				Subdirectory: "service/s3",
			},
		},
		{
			module: "golang.org/x/sys",
			info: &GoModRepoInfo{
				Host:  "github.com",
				Owner: "golang",
				Repo:  "sys",
			},
		},
		{
			module: "gopkg.in/yaml.v3",
			info: &GoModRepoInfo{
				Host:  "github.com",
				Owner: "go-yaml",
				Repo:  "yaml",
				Major: "v3",
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := parseGoModRepo(test.module)
			assert.Equal(t, test.info, got)
			assert.Equal(t, test.err, err)
		})
	}
}
