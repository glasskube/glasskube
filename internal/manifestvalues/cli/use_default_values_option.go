package cli

import "github.com/glasskube/glasskube/api/v1alpha1"

type UseDefaultValuesOption []string

func (o UseDefaultValuesOption) ShouldUseDefault(name string, def v1alpha1.ValueDefinition) bool {
	for _, v := range o {
		switch v {
		case "all", name:
			return def.DefaultValue != ""
		}
	}
	return false
}
