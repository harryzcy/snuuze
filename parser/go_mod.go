package parser

import (
	"golang.org/x/mod/modfile"

	"github.com/harryzcy/snuuze/matcher"
	"github.com/harryzcy/snuuze/types"
)

func parseGoMod(path string, data []byte) ([]types.Dependency, error) {
	file, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}

	dependencies := make([]types.Dependency, 0, len(file.Require))
	for _, require := range file.Require {
		dependencies = append(dependencies, types.Dependency{
			Name:           require.Mod.Path,
			Version:        require.Mod.Version,
			Indirect:       require.Indirect,
			PackageManager: matcher.GoMod,
		})
	}

	return dependencies, nil
}
