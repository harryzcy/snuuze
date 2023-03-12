package updater

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/harryzcy/snuuze/cmdutil"
	"github.com/harryzcy/snuuze/types"
)

func upgradeGoMod(cache *Cache, info types.UpgradeInfo) error {
	file, err := cache.Get(info.Dependency.File)
	if err != nil {
		return fmt.Errorf("failed to get file %s from cache: %s", info.Dependency.File, err)
	}

	start := info.Dependency.Position.StartByte
	end := info.Dependency.Position.EndByte

	buffer := bytes.Buffer{}
	buffer.Write(file[:start])
	buffer.WriteString(info.Dependency.Name)
	buffer.WriteString(" ")
	buffer.WriteString(info.ToVersion)
	buffer.Write(file[end:])
	cache.Set(info.Dependency.File, buffer.Bytes())

	return nil
}

func postGoMod(cache *Cache) error {
	goModFiles, err := cache.ListGoMod()
	if err != nil {
		return fmt.Errorf("failed to list go.mod files: %s", err)
	}

	for _, file := range goModFiles {
		dir := filepath.Dir(file)
		output, err := cmdutil.RunCommand(cmdutil.CommandInputs{
			Command: []string{"go", "mod", "tidy"},
			Dir:     dir,
			Env: map[string]string{
				"GOPATH": os.Getenv("GOPATH"),
			},
		})
		if err != nil {
			fmt.Println("failed to run go mod tidy in", dir)
			fmt.Println(output.Stderr.String())
			return fmt.Errorf("failed to run go mod tidy in %s: %s", dir, err)
		}
	}
	return nil
}
