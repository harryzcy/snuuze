package updater

import (
	"strconv"
	"testing"

	"github.com/harryzcy/snuuze/types"
	"github.com/stretchr/testify/assert"
)

func TestPrepareCommit(t *testing.T) {
	tests := []struct {
		group RuleGroup
		want  *commitInfo
		ok    bool
	}{
		{
			group: RuleGroup{
				Rule: types.Rule{
					Name: "test",
				},
				Infos: []types.UpgradeInfo{},
			},
			want: &commitInfo{
				branchName: "snuuze/test",
				message:    "Update test",
			},
			ok: true,
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
						ToVersion: "v1.0.0",
					},
				},
			},
			want: &commitInfo{
				branchName: "snuuze/github.com-aws-aws-sdk-go-v2",
				message:    "Update github.com/aws/aws-sdk-go-v2 to v1.0.0",
			},
			ok: true,
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
			want: nil,
			ok:   false,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, ok := prepareCommit(test.group)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.want, got)
		})
	}
}
