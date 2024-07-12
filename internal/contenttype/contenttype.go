package contenttype

import (
	"fmt"
	"strings"
)

type ContentType struct {
	MediaType string
	Charset   string
	Boundary  string
}

func ParseContentType(value string) (*ContentType, error) {
	directives := strings.Split(value, ";")
	result := ContentType{
		MediaType: strings.TrimSpace(directives[0]),
	}
	if len(directives) > 1 {
		for _, directive := range directives[1:] {
			parts := strings.SplitN(strings.TrimSpace(directive), "=", 2)
			if len(parts) == 2 {
				switch parts[0] {
				case "charset":
					result.Charset = parts[1]
				case "boundary":
					result.Boundary = parts[1]
				default:
					return nil, fmt.Errorf("content type has unknown directive %v", directive)
				}
			} else {
				return nil, fmt.Errorf("content type has invalid directive %v", directive)
			}
		}
	}
	return &result, nil
}
