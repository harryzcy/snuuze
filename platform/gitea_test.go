package platform

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestGetGiteaToken(t *testing.T) {
	originalGetGiteaConfigs := getGiteaConfigs
	defer func() {
		getGiteaConfigs = originalGetGiteaConfigs
	}()
	getGiteaConfigs = func() []types.GiteaConfig {
		return []types.GiteaConfig{
			{
				Host:  "https://gitea.com",
				Token: "token1",
			},
			{
				Host:  "https://git.example.com",
				Token: "token2",
			},
		}
	}

	tests := []struct {
		host  string
		token string
	}{
		{
			host:  "https://gitea.com",
			token: "token1",
		},
		{
			host:  "https://git.example.com",
			token: "token2",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := getGiteaToken(test.host)
			assert.Equal(t, test.token, actual)
		})
	}
}
