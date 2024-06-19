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

//nolint:dupl
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=clpkg
//+kubebuilder:printcolumn:name=Desired version,type=string,JSONPath=".spec.packageInfo.version"
//+kubebuilder:printcolumn:name=Installed version,type=string,JSONPath=".status.version"

// ClusterPackage is the Schema for the clusterpackages API
type ClusterPackage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec,omitempty"`
	Status PackageStatus `json:"status,omitempty"`
}

func (pkg *ClusterPackage) GetSpec() *PackageSpec {
	return &pkg.Spec
}

func (in *ClusterPackage) GetStatus() *PackageStatus {
	return &in.Status
}

func (in *ClusterPackage) AutoUpdatesEnabled() bool {
	return autoUpdatesEnabled(in.ObjectMeta)
}

func (in *ClusterPackage) SetAutoUpdatesEnabled(enabled bool) {
	setAutoUpdatesEnabled(&in.ObjectMeta, enabled)
}

func (pkg *ClusterPackage) InstalledAsDependency() bool {
	return installedAsDependency(pkg.ObjectMeta)
}

func (pkg *ClusterPackage) SetInstalledAsDependency(value bool) {
	setInstalledAsDependency(&pkg.ObjectMeta, value)
}

func (pkg *ClusterPackage) IsNamespaceScoped() bool {
	return false
}

func (pkg *ClusterPackage) IsNil() bool {
	return pkg == nil
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
