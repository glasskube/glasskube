package v1alpha1

import (
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	autoUpdateAnnotation = "packages.glasskube.dev/auto-update"
	dependencyAnnotation = "packages.glasskube.dev/installed-as-dependency"
)

func autoUpdatesEnabled(obj metav1.ObjectMeta) bool {
	if obj.Annotations == nil {
		return false
	} else if enabledStr, ok := obj.Annotations[autoUpdateAnnotation]; !ok {
		return false
	} else {
		enabled, _ := strconv.ParseBool(enabledStr)
		return enabled
	}
}

func setAutoUpdatesEnabled(obj *metav1.ObjectMeta, enabled bool) {
	if obj.Annotations == nil {
		obj.Annotations = make(map[string]string)
	}
	if enabled {
		obj.Annotations[autoUpdateAnnotation] = strconv.FormatBool(true)
	} else {
		delete(obj.Annotations, autoUpdateAnnotation)
	}
}

func installedAsDependency(obj metav1.ObjectMeta) bool {
	if obj.Annotations == nil {
		return false
	} else if enabledStr, ok := obj.Annotations[dependencyAnnotation]; !ok {
		return false
	} else {
		enabled, _ := strconv.ParseBool(enabledStr)
		return enabled
	}
}

func setInstalledAsDependency(obj *metav1.ObjectMeta, value bool) {
	if obj.Annotations == nil {
		obj.Annotations = make(map[string]string)
	}
	if value {
		obj.Annotations[dependencyAnnotation] = strconv.FormatBool(true)
	} else {
		delete(obj.Annotations, dependencyAnnotation)
	}
}

type PackageInfoTemplate struct {
	// Name of the package to install
	Name string `json:"name"`
	// Version of the package to install
	Version string `json:"version"`
	// RepositoryName is the name of the repository to pull the package from (optional)
	RepositoryName string `json:"repositoryName,omitempty"`
}

type ObjectKeyValueSource struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Key       string `json:"key"`
}

type PackageValueSource struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type ValueReference struct {
	ConfigMapRef *ObjectKeyValueSource `json:"configMapRef,omitempty"`
	SecretRef    *ObjectKeyValueSource `json:"secretRef,omitempty"`
	PackageRef   *PackageValueSource   `json:"packageRef,omitempty"`
}

// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:MaxProperties:=1
type ValueConfiguration struct {
	Value     *string         `json:"value,omitempty"`
	ValueFrom *ValueReference `json:"valueFrom,omitempty"`
}

// PackageSpec defines the desired state
type PackageSpec struct {
	PackageInfo PackageInfoTemplate           `json:"packageInfo"`
	Values      map[string]ValueConfiguration `json:"values,omitempty"`
}

// PackageStatus defines the observed state
type PackageStatus struct {
	Version           string             `json:"version,omitempty"`
	Conditions        []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	OwnedResources    []OwnedResourceRef `json:"ownedResources,omitempty"`
	OwnedPackageInfos []OwnedResourceRef `json:"ownedPackageInfos,omitempty"`
	OwnedPackages     []OwnedResourceRef `json:"ownedPackages,omitempty"`
}
