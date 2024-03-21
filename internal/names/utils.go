package names

import (
	"regexp"
	"strings"
)

var (
	resourceNameRegex = regexp.MustCompile(`[^\w.-]`)
)

func escapeResourceName(name string) string {
	return strings.ToLower(resourceNameRegex.ReplaceAllString(name, "--"))
}
