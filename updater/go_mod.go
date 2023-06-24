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

	lines := bytes.Split(file, []byte("\n"))
	lineIdx := info.Dependency.Position.Line - 1 // line number starts from 1

	buffer := bytes.Buffer{}
	for idx, line := range lines {
		if idx == lineIdx {
			buffer.WriteString("\t")
			buffer.WriteString(info.Dependency.Name)
			buffer.WriteString(" ")
			buffer.WriteString(info.ToVersion)
			if info.Dependency.Indirect {
				buffer.WriteString(" // indirect")
			}
		} else {
			buffer.Write(line)
		}
		buffer.WriteByte('\n')
	}
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
