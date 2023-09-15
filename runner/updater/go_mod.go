package updater

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/runner/command"
	"github.com/harryzcy/snuuze/types"
	"github.com/hashicorp/go-version"
)

type ReplaceItem struct {
	Dir     string
	OldName string
	NewName string
}

func upgradeGoMod(cache *Cache, info *types.UpgradeInfo) (*ReplaceItem, error) {
	file, err := cache.Get(info.Dependency.File)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %s from cache: %s", info.Dependency.File, err)
	}

	lines := bytes.Split(file, []byte("\n"))
	lineIdx := info.Dependency.Position.Line - 1 // line number starts from 1

	moduleName, err := getGoModuleName(info.Dependency.Name, info.Dependency.Version, info.ToVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %s", err)
	}

	var replace *ReplaceItem
	if moduleName != info.Dependency.Name {
		replace = &ReplaceItem{
			Dir:     filepath.Dir(info.Dependency.File),
			OldName: info.Dependency.Name,
			NewName: moduleName,
		}
	}

	buffer := bytes.Buffer{}
	for idx, line := range lines {
		if idx == lineIdx {
			buffer.WriteString("\t")
			buffer.WriteString(moduleName)
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

	return replace, nil
}

// getGoModuleName returns the module name of the current major version.
func getGoModuleName(name, fromVersion, toVersion string) (string, error) {
	if strings.HasSuffix(toVersion, "+incompatible") {
		return name, nil
	}

	toV, err := version.NewVersion(toVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse version %s: %s", toVersion, err)
	}

	fromV, err := version.NewVersion(fromVersion)
	if err != nil {
		return "", fmt.Errorf("failed to parse version %s: %s", fromVersion, err)
	}

	shouldUpdate := toV.Segments()[0] > fromV.Segments()[0]
	if !shouldUpdate {
		return name, nil
	}

	oldSuffix := "/v" + strconv.Itoa(fromV.Segments()[0])
	newSuffix := "/v" + strconv.Itoa(toV.Segments()[0])

	baseName := strings.TrimSuffix(name, oldSuffix)
	return baseName + newSuffix, nil
}

func postGoMod(cache *Cache, goReplaceItems []*ReplaceItem) error {
	goModFiles, err := cache.ListGoMod()
	if err != nil {
		return fmt.Errorf("postGoMod: failed to list go.mod files: %s", err)
	}

	for _, file := range goModFiles {
		dir := filepath.Dir(file)

		for _, replace := range goReplaceItems {
			if replace.Dir != dir {
				continue
			}

			err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if !strings.HasSuffix(path, ".go") {
					return nil
				}

				data, err := cache.Get(path)
				if err != nil {
					return fmt.Errorf("failed to get file %s from cache: %s", path, err)
				}
				data = bytes.ReplaceAll(data, []byte(replace.OldName), []byte(replace.NewName))
				cache.Set(path, data)
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to walk dir %s: %s", dir, err)
			}
		}
		err = cache.Commit()
		if err != nil {
			return fmt.Errorf("postGoMod: failed to commit cache: %s", err)
		}

		output, err := command.RunCommand(command.CommandInputs{
			Command: []string{"go", "mod", "tidy"},
			Dir:     dir,
			Env: map[string]string{
				"GOPATH": config.GoPath(),
			},
		})
		if err != nil {
			fmt.Println("postGoMod: failed to run go mod tidy in", dir)
			fmt.Println(output.Stderr.String())
			return fmt.Errorf("postGoMod: failed to run go mod tidy in %s: %s", dir, err)
		}
	}
	return nil
}
