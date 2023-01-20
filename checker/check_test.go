package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/harryzcy/latte/matcher"
	"github.com/harryzcy/latte/types"
	"github.com/stretchr/testify/assert"
)

func TestIsUpgradable(t *testing.T) {
	tests := []struct {
		dep  types.Dependency
		want UpgradeInfo
		err  error
	}{
		{
			dep: types.Dependency{
				Name:           "github.com/gofiber/fiber/v2",
				Version:        "v2.40.0", // use an old version
				PackageManager: matcher.GoMod,
			},
			want: UpgradeInfo{
				Dependency: types.Dependency{
					Name:           "github.com/gofiber/fiber/v2",
					Version:        "v2.40.0",
					PackageManager: matcher.GoMod,
				},
				Upgradable: true,
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			info, err := isUpgradable(test.dep)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.want.Upgradable, info.Upgradable)
			assert.Equal(t, test.want.Dependency, info.Dependency)
		})
	}
}

func TestMain(m *testing.M) {
	envFile := filepath.Join("..", ".env")
	content, err := os.ReadFile(envFile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		fmt.Println(parts)
		if len(parts) != 2 {
			panic(fmt.Sprintf("invalid line: %s", line))
		}

		if parts[0] == "GITHUB_TOKEN" {
			originalToken := GITHUB_TOKEN
			GITHUB_TOKEN = parts[1]
			defer func() {
				GITHUB_TOKEN = originalToken
			}()
		}
	}

	m.Run()
}
