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
//+kubebuilder:resource:scope=Namespaced,shortName=pkg
//+kubebuilder:printcolumn:name=Desired version,type=string,JSONPath=".spec.packageInfo.version"
//+kubebuilder:printcolumn:name=Installed version,type=string,JSONPath=".status.version"

// Package is the Schema for the packages API
type Package struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageSpec   `json:"spec,omitempty"`
	Status PackageStatus `json:"status,omitempty"`
}

func (pkg *Package) GetSpec() *PackageSpec {
	return &pkg.Spec
}

func (pkg *Package) GetStatus() *PackageStatus {
	return &pkg.Status
}

func (pkg *Package) AutoUpdatesEnabled() bool {
	return autoUpdatesEnabled(pkg.ObjectMeta)
}

func (pkg *Package) SetAutoUpdatesEnabled(enabled bool) {
	setAutoUpdatesEnabled(&pkg.ObjectMeta, enabled)
}

func (pkg *Package) InstalledAsDependency() bool {
	return false
}

func (pkg *Package) SetInstalledAsDependency(value bool) {
	panic("illegal operation: package can not be installed as dependency")
}

func (pkg *Package) IsNamespaceScoped() bool {
	return true
}

func (pkg *Package) IsNil() bool {
	return pkg == nil
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
