package web

import "strings"

func packagesPathRedactor(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) <= 1 {
		return url
	}
	if parts[1] == "packages" && len(parts) >= 5 && parts[3] != "-" && parts[4] != "-" {
		// when the user opens an installed package, the path is /packages/<manifestName>/<namespace>/<name>
		// so we want to redact namespace and name
		parts[3] = "x"
		namePart := parts[4]
		nameParts := strings.Split(namePart, "?")
		nameParts[0] = "x"
		parts[4] = strings.Join(nameParts, "?")
		return strings.Join(parts, "/")
	}
	return url
}
