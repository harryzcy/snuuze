package pip

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		line string
		want *types.Dependency
	}{
		{
			line: "django==1.11.0",
			want: &types.Dependency{
				Name:           "django",
				Version:        "== 1.11.0",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"constraints": [][2]string{{"==", "1.11.0"}},
				},
			},
		},
		{
			line: "SomeProject",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "",
				PackageManager: types.PackageManagerPip,
				Extra:          map[string]interface{}{},
			},
		},
		{
			line: "SomeProject == 1.3",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "== 1.3",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"constraints": [][2]string{{"==", "1.3"}},
				},
			},
		},
		{
			line: "SomeProject >= 1.2, < 2.0",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        ">= 1.2, < 2.0",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"constraints": [][2]string{{">=", "1.2"}, {"<", "2.0"}},
				},
			},
		},
		{
			line: "SomeProject[foo, bar]",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"extras": "foo, bar",
				},
			},
		},
		{
			line: "SomeProject[foo, bar] == 1.3",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "== 1.3",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"extras":      "foo, bar",
					"constraints": [][2]string{{"==", "1.3"}},
				},
			},
		},
		{
			line: "SomeProject ~= 1.4.2",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "~= 1.4.2",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"constraints": [][2]string{{"~=", "1.4.2"}},
				},
			},
		},
		{
			line: "SomeProject == 5.4 ; python_version < '3.8'",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "== 5.4",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"constraints": [][2]string{{"==", "5.4"}},
				},
			},
		},
		{
			line: "SomeProject ; sys_platform == 'win32'",
			want: &types.Dependency{
				Name:           "SomeProject",
				Version:        "",
				PackageManager: types.PackageManagerPip,
				Extra:          map[string]interface{}{},
			},
		},
		{
			line: "requests [security] >= 2.8.1, == 2.8.* ; python_version < \"2.7\"",
			want: &types.Dependency{
				Name:           "requests",
				Version:        ">= 2.8.1, == 2.8.*",
				PackageManager: types.PackageManagerPip,
				Extra: map[string]interface{}{
					"extras":      "security",
					"constraints": [][2]string{{">=", "2.8.1"}, {"==", "2.8.*"}},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			parser := NewParser("", test.line)
			dep, err := parser.Parse()

			assert.Nil(t, err)
			assert.Equal(t, test.want, dep)
		})
	}
}
