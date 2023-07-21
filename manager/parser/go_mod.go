package parser

import (
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/harryzcy/snuuze/types"
)

func parseGoMod(path string, data []byte) ([]types.Dependency, error) {
	filename := filepath.Base(path)
	file, err := modfile.Parse(filename, data, nil)
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
				ColStart:  require.Syntax.Start.LineRune,
				ColEnd:    require.Syntax.End.LineRune,
				ByteStart: require.Syntax.Start.Byte,
				ByteEnd:   require.Syntax.End.Byte,
			},
		})
	}

	return dependencies, nil
}
