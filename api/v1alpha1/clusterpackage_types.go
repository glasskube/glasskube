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

// ClusterPackageSpec defines the desired state of ClusterPackage
type ClusterPackageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ClusterPackage. Edit clusterpackage_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// ClusterPackageStatus defines the observed state of ClusterPackage
type ClusterPackageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name=Desired version,type=string,JSONPath=".spec.packageInfo.version"
//+kubebuilder:printcolumn:name=Installed version,type=string,JSONPath=".status.version"

// ClusterPackage is the Schema for the clusterpackages API
type ClusterPackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec,omitempty"`
	Status PackageStatus `json:"status,omitempty"`
}

// GetSpec implements PackageCommon.
func (pkg *ClusterPackage) GetSpec() *PackageSpec {
	return &pkg.Spec
}

// GetStatus implements PackageCommon.
func (in *ClusterPackage) GetStatus() *PackageStatus {
	return &in.Status
}

// AutoUpdatesEnabled implements AutoUpdates.
func (in *ClusterPackage) AutoUpdatesEnabled() bool {
	return autoUpdatesEnabled(in.ObjectMeta)
}

// SetAutoUpdatesEnabled implements AutoUpdates.
func (in *ClusterPackage) SetAutoUpdatesEnabled(enabled bool) {
	setAutoUpdatesEnabled(&in.ObjectMeta, enabled)
}

//+kubebuilder:object:root=true

// ClusterPackageList contains a list of ClusterPackage
type ClusterPackageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterPackage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterPackage{}, &ClusterPackageList{})
}
