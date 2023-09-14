package updater

import (
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/types"
)

var getRules = config.GetRules

type RuleGroup struct {
	Rule  types.Rule
	Infos []*types.UpgradeInfo
}

func groupUpdates(infos []*types.UpgradeInfo) []RuleGroup {
	rules := getRules()

	infoStatus := make([]bool, len(infos))

	groups := make([]RuleGroup, 0)
	for _, rule := range rules {
		group := RuleGroup{
			Rule:  rule,
			Infos: make([]*types.UpgradeInfo, 0),
		}

		for i, info := range infos {
			if matchRule(rule, info) {
				group.Infos = append(group.Infos, info)
				infoStatus[i] = true
			}
		}
		if len(group.Infos) > 0 {
			groups = append(groups, group)
		}
	}

	// add remaining infos separately
	for i, info := range infos {
		if !infoStatus[i] {
			groups = append(groups, RuleGroup{
				Rule: types.Rule{
					PackageManagers: []types.PackageManager{info.Dependency.PackageManager},
					PackageNames:    []string{info.Dependency.Name},
				},
				Infos: []*types.UpgradeInfo{info},
			})
		}
	}

	return groups
}

func matchRule(rule types.Rule, info *types.UpgradeInfo) bool {
	if len(rule.PackageManagers) > 0 && !contains(rule.PackageManagers, info.Dependency.PackageManager) {
		return false
	}
	if len(rule.PackageNames) > 0 && !contains(rule.PackageNames, info.Dependency.Name) {
		return false
	}
	if len(rule.PackageTypes) > 0 && !contains(rule.PackageTypes, types.PackageTypeAll) {
		if info.Dependency.Indirect && !contains(rule.PackageTypes, types.PackageTypeIndirect) {
			return false
		}
		if !info.Dependency.Indirect && !contains(rule.PackageTypes, types.PackageTypeDirect) {
			return false
		}
	}

	return true
}

func contains[elem comparable](list []elem, value elem) bool {
	for _, e := range list {
		if e == value {
			return true
		}
	}
	return false
}
