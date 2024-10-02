package clientutils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/glasskube/glasskube/internal/contenttype"
	"github.com/glasskube/glasskube/internal/httperror"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func FetchResourcesFromUrl(url string) ([]unstructured.Unstructured, error) {
	if request, err := NewResourcesRequest(url); err != nil {
		return nil, err
	} else {
		return FetchResources(request)
	}
}

func NewResourcesRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", contenttype.MediaTypeJSON)
	request.Header.Add("Accept", contenttype.MediaTypeYAML)
	return request, nil
}

func FetchResources(request *http.Request) ([]unstructured.Unstructured, error) {
	url := request.URL.Redacted()
	response, err := httperror.CheckResponse(http.DefaultClient.Do(request))
	if err != nil {
		switch {
		case httperror.IsNotFound(err):
			return nil, fmt.Errorf("manifest not found at %v: %v", url, err)
		case httperror.Is(err, http.StatusForbidden):
			return nil, fmt.Errorf("access denied to manifest at %v: %v", url, err)
		case httperror.Is(err, http.StatusUnauthorized):
			return nil, fmt.Errorf("unauthorized to access manifest at %v: %v", url, err)
		default:
			return nil, fmt.Errorf("failed to download manifest from %v: %v", url, err)
		}
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if err := contenttype.IsJsonOrYaml(response); err != nil {
		return nil, fmt.Errorf("could not decode manifest %v: %w", url, err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(response.Body, 4096)
	resources := make([]unstructured.Unstructured, 0)
	for {
		object := unstructured.Unstructured{}
		if err := decoder.Decode(&object); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not decode manifest %v: %w", url, err)
		}
		if len(object.Object) == 0 {
			continue
		}
		resources = append(resources, object)
	}
	return resources, nil
}
