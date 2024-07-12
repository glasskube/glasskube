package contenttype

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	MediaTypeJSON             = "application/json"
	MediaTypeYAML             = "application/yaml"
	MediaTypeOctetStream      = "application/octet-stream"
	MediaTypeTextYAML         = "text/yaml"
	MediaTypeTextPlain        = "text/plain"
	MediaTypeTextXYaml        = "text/x-yaml"
	MediaTypeApplicationXYaml = "application/x-yaml"
)

func IsJsonOrYaml(response *http.Response) error {
	return HasMediaType(response,
		MediaTypeJSON,
		MediaTypeYAML,
		MediaTypeOctetStream,
		MediaTypeTextYAML,
		MediaTypeTextPlain,
		MediaTypeTextXYaml,
		MediaTypeApplicationXYaml)
}

func HasMediaType(response *http.Response, acceptedContentTypes ...string) error {
	contentType, err := ParseContentType(response.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if contentType.MediaType == "" {
		return nil
	}
	for _, t := range acceptedContentTypes {
		if contentType.MediaType == t {
			return nil
		}
	}
	return fmt.Errorf("response has unacceptable media type: %v (acceptable media types are %v)",
		contentType.MediaType, strings.Join(acceptedContentTypes, ", "))
}
