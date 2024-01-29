package clientutils

import (
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func FetchResources(url string) (*[]unstructured.Unstructured, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download manifest %v: %w", url, err)
	}
	defer response.Body.Close()
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
		resources = append(resources, object)
	}
	return &resources, nil
}
