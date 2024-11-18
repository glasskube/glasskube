package util

import (
	"fmt"

	"github.com/glasskube/glasskube/api/v1alpha1"
)

func ComponentName(parentName string, cmp v1alpha1.Component) string {
	if cmp.InstalledName != "" {
		return fmt.Sprintf("%v-%v", parentName, cmp.InstalledName)
	} else {
		return fmt.Sprintf("%v-%v", parentName, cmp.Name)
	}
}
