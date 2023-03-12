package parser

import (
	"golang.org/x/mod/modfile"

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
			File:           path,
			Name:           require.Mod.Path,
			Version:        require.Mod.Version,
			Indirect:       require.Indirect,
			PackageManager: types.PackageManagerGoMod,
			Position: types.Position{
				Line:      require.Syntax.Start.Line,
				StartByte: require.Syntax.Start.Byte,
				EndByte:   require.Syntax.End.Byte,
			},
		})
	}

	return dependencies, nil
}
