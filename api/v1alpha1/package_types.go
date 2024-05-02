/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

// PackageSpec defines the desired state of Package
type PackageSpec struct {
	PackageInfo PackageInfoTemplate           `json:"packageInfo"`
	Values      map[string]ValueConfiguration `json:"values,omitempty"`
}

// PackageStatus defines the observed state of Package
type PackageStatus struct {
	Version           string             `json:"version,omitempty"`
	Conditions        []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	OwnedResources    []OwnedResourceRef `json:"ownedResources,omitempty"`
	OwnedPackageInfos []OwnedResourceRef `json:"ownedPackageInfos,omitempty"`
	OwnedPackages     []OwnedResourceRef `json:"ownedPackages,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=pkg
//+kubebuilder:printcolumn:name=Desired version,type=string,JSONPath=".spec.packageInfo.version"
//+kubebuilder:printcolumn:name=Installed version,type=string,JSONPath=".status.version"

// Package is the Schema for the packages API
type Package struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec,omitempty"`
	Status PackageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PackageList contains a list of Package
type PackageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Package `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Package{}, &PackageList{})
}
