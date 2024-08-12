package clientutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/controller/ctrlpkg"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

type OutputFormat string

const (
	OutputFormatJSON OutputFormat = "json"
	OutputFormatYAML OutputFormat = "yaml"
)

func Format(outputFormat OutputFormat, showAll bool, pkgs ...ctrlpkg.Package) (string, error) {
	if !showAll {
		for _, pkg := range pkgs {
			switch p := pkg.(type) {
			case *v1alpha1.ClusterPackage:
				p.ObjectMeta = pruneExtraFields(p.ObjectMeta)
			case *v1alpha1.Package:
				p.ObjectMeta = pruneExtraFields(p.ObjectMeta)
			default:
				panic("unexpected package type")
			}
		}
	}

	var outputData []byte
	for i := range pkgs {
		if gvks, _, err := scheme.Scheme.ObjectKinds(pkgs[i]); err == nil && len(gvks) == 1 {
			pkgs[i].SetGroupVersionKind(gvks[0])
		} else {
			return "", fmt.Errorf("failed to set GVK for package: %w\n", err)
		}
	}
	switch outputFormat {
	case OutputFormatJSON:
		var err error
		res := make([]map[string]any, 0)
		for _, pkg := range pkgs {
			if pkgAsMap, pruneErr := packageAsMap(pkg, showAll); pruneErr != nil {
				return "", fmt.Errorf("failed to prune status during marshalling: %w", pruneErr)
			} else {
				res = append(res, pkgAsMap)
			}
		}
		if outputData, err = json.MarshalIndent(res, "", "  "); err != nil {
			return "", fmt.Errorf("failed to marshal output: %w\n", err)
		}
	case OutputFormatYAML:
		var buffer bytes.Buffer
		l := len(pkgs)
		for _, pkg := range pkgs {
			if p, pruneErr := packageAsMap(pkg, showAll); pruneErr != nil {
				return "", fmt.Errorf("failed to prune status during marshalling: %w", pruneErr)
			} else if data, err := yaml.Marshal(p); err != nil {
				return "", fmt.Errorf("failed to marshal output: %w\n", err)
			} else {
				if l > 1 {
					buffer.WriteString("---\n")
				}
				buffer.Write(data)
			}
		}
		outputData = buffer.Bytes()
	default:
		return "", fmt.Errorf("unsupported output format: %v", outputFormat)
	}

	return string(outputData), nil
}

func pruneExtraFields(original v1.ObjectMeta) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:        original.Name,
		Namespace:   original.Namespace,
		Labels:      pruneNonGlasskubeKeys(original.Labels),
		Annotations: pruneNonGlasskubeKeys(original.Annotations),
	}
}

func pruneNonGlasskubeKeys(original map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range original {
		if strings.Contains(k, "packages.glasskube.dev") {
			res[k] = v
		}
	}
	return res
}

func packageAsMap(pkg ctrlpkg.Package, showAll bool) (map[string]any, error) {
	var pkgAsMap map[string]any
	if data, err := json.Marshal(pkg); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &pkgAsMap); err != nil {
		return nil, err
	} else {
		if !showAll {
			delete(pkgAsMap, "status")
		}
		return pkgAsMap, nil
	}
}
