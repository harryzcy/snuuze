package updater

import (
	"testing"

	"github.com/harryzcy/snuuze/types"
)

func TestGenerateBranchName(t *testing.T) {
	tests := []struct {
		group RuleGroup
		want  string
		ok    bool
	}{
		{
			group: RuleGroup{
				Rule: types.Rule{
					Name: "test",
				},
				Infos: []types.UpgradeInfo{},
			},
			want: "snuuze/test",
			ok:   true,
		},
		{
			group: RuleGroup{
				Rule: types.Rule{
					Name: "",
				},
				Infos: []types.UpgradeInfo{
					{
						Dependency: types.Dependency{
							Name: "github.com/aws/aws-sdk-go-v2",
						},
					},
				},
			},
			want: "snuuze/github.com-aws-aws-sdk-go-v2",
			ok:   true,
		},
		{
			group: RuleGroup{
				Rule: types.Rule{
					Name: "",
				},
				Infos: []types.UpgradeInfo{
					{
						Dependency: types.Dependency{
							Name: "github.com/aws/aws-sdk-go-v2",
						},
					},
					{
						Dependency: types.Dependency{
							Name: "github.com/aws/aws-sdk-go-v2/config",
						},
					},
				},
			},
			want: "",
			ok:   false,
		},
	}

	for i, test := range tests {
		got, ok := generateBranchName(test.group)
		if got != test.want || ok != test.ok {
			t.Errorf("test %d: got %s, %t, want %s, %t", i, got, ok, test.want, test.ok)
		}
	}
}
