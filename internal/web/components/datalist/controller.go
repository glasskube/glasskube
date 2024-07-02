package datalist

import (
	"fmt"
)

func ForDatalist(valueName string, postfix string, options []string) map[string]any {
	return map[string]any{
		"Id":      fmt.Sprintf("%s-%s", valueName, postfix),
		"Options": options,
	}
}
