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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PackageRepositoryBasicAuthSpec struct {
	Username          *string                   `json:"username,omitempty"`
	UsernameSecretRef *corev1.SecretKeySelector `json:"usernameSecretRef,omitempty"`
	Password          *string                   `json:"password,omitempty"`
	PasswordSecretRef *corev1.SecretKeySelector `json:"passwordSecretRef,omitempty"`
}

type PackageRepositoryBearerAuthSpec struct {
	Token          *string                   `json:"token,omitempty"`
	TokenSecretRef *corev1.SecretKeySelector `json:"tokenSecretRef,omitempty"`
}

type PackageRepositoryAuthSpec struct {
	Basic  *PackageRepositoryBasicAuthSpec  `json:"basic,omitempty"`
	Bearer *PackageRepositoryBearerAuthSpec `json:"bearer,omitempty"`
}

// PackageRepositorySpec defines the desired state of PackageRepository
type PackageRepositorySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Url  string                     `json:"url"`
	Auth *PackageRepositoryAuthSpec `json:"auth,omitempty"`
}

// PackageRepositoryStatus defines the observed state of PackageRepository
type PackageRepositoryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// PackageRepository is the Schema for the packagerepositories API
type PackageRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageRepositorySpec   `json:"spec,omitempty"`
	Status PackageRepositoryStatus `json:"status,omitempty"`
}

var (
	defaultRepositoryAnnotation = "packages.glasskube.dev/defaultRepository"
)

func (repo PackageRepository) IsDefaultRepository() bool {
	return repo.Annotations[defaultRepositoryAnnotation] == "true"
}

func (repo PackageRepository) SetDefaultRepository() {
	repo.SetDefaultRepositoryBool(true)
}

func (repo PackageRepository) SetDefaultRepositoryBool(value bool) {
	if value {
		repo.Annotations[defaultRepositoryAnnotation] = "true"
	} else {
		delete(repo.Annotations, defaultRepositoryAnnotation)
	}
}

//+kubebuilder:object:root=true

// PackageRepositoryList contains a list of PackageRepository
type PackageRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PackageRepository{}, &PackageRepositoryList{})
}
