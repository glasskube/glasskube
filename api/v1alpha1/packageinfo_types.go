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

// PackageInfoSpec defines the desired state of PackageInfo
type PackageInfoSpec struct {
	// +kubebuilder:validation:Required
	Name           string `json:"name"`
	Version        string `json:"version,omitempty"`
	RepositoryName string `json:"repositoryUrl,omitempty"`
}

// PackageInfoStatus defines the observed state of PackageInfo
type PackageInfoStatus struct {
	Manifest            *PackageManifest   `json:"manifest,omitempty"`
	ResolvedUrl         string             `json:"resolvedUrl,omitempty"`
	Conditions          []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
	LastUpdateTimestamp *metav1.Time       `json:"lastUpdateTimestamp,omitempty"`
	Version             string             `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=packageinfos,scope=Cluster,shortName=pkgi
//+kubebuilder:printcolumn:name=Desired version,type=string,JSONPath=".spec.version"
//+kubebuilder:printcolumn:name=Current version,type=string,JSONPath=".status.version"
//+kubebuilder:printcolumn:name="Last Updated",type=date,JSONPath=".status.lastUpdateTimestamp"

// PackageInfo is the Schema for the packageinfos API
type PackageInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageInfoSpec   `json:"spec,omitempty"`
	Status PackageInfoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PackageInfoList contains a list of PackageInfo
type PackageInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PackageInfo{}, &PackageInfoList{})
}
