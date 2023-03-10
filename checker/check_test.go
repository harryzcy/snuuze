package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
