//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPackage) DeepCopyInto(out *ClusterPackage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPackage.
func (in *ClusterPackage) DeepCopy() *ClusterPackage {
	if in == nil {
		return nil
	}
	out := new(ClusterPackage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPackage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPackageList) DeepCopyInto(out *ClusterPackageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterPackage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPackageList.
func (in *ClusterPackageList) DeepCopy() *ClusterPackageList {
	if in == nil {
		return nil
	}
	out := new(ClusterPackageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterPackageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Component) DeepCopyInto(out *Component) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Component.
func (in *Component) DeepCopy() *Component {
	if in == nil {
		return nil
	}
	out := new(Component)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dependency) DeepCopyInto(out *Dependency) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dependency.
func (in *Dependency) DeepCopy() *Dependency {
	if in == nil {
		return nil
	}
	out := new(Dependency)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmManifest) DeepCopyInto(out *HelmManifest) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = new(JSON)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmManifest.
func (in *HelmManifest) DeepCopy() *HelmManifest {
	if in == nil {
		return nil
	}
	out := new(HelmManifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSON) DeepCopyInto(out *JSON) {
	*out = *in
	if in.Raw != nil {
		in, out := &in.Raw, &out.Raw
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSON.
func (in *JSON) DeepCopy() *JSON {
	if in == nil {
		return nil
	}
	out := new(JSON)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KustomizeManifest) DeepCopyInto(out *KustomizeManifest) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KustomizeManifest.
func (in *KustomizeManifest) DeepCopy() *KustomizeManifest {
	if in == nil {
		return nil
	}
	out := new(KustomizeManifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectKeyValueSource) DeepCopyInto(out *ObjectKeyValueSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectKeyValueSource.
func (in *ObjectKeyValueSource) DeepCopy() *ObjectKeyValueSource {
	if in == nil {
		return nil
	}
	out := new(ObjectKeyValueSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnedResourceRef) DeepCopyInto(out *OwnedResourceRef) {
	*out = *in
	out.GroupVersionKind = in.GroupVersionKind
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnedResourceRef.
func (in *OwnedResourceRef) DeepCopy() *OwnedResourceRef {
	if in == nil {
		return nil
	}
	out := new(OwnedResourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Package) DeepCopyInto(out *Package) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Package.
func (in *Package) DeepCopy() *Package {
	if in == nil {
		return nil
	}
	out := new(Package)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Package) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageEntrypoint) DeepCopyInto(out *PackageEntrypoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageEntrypoint.
func (in *PackageEntrypoint) DeepCopy() *PackageEntrypoint {
	if in == nil {
		return nil
	}
	out := new(PackageEntrypoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfo) DeepCopyInto(out *PackageInfo) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfo.
func (in *PackageInfo) DeepCopy() *PackageInfo {
	if in == nil {
		return nil
	}
	out := new(PackageInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageInfo) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfoList) DeepCopyInto(out *PackageInfoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfoList.
func (in *PackageInfoList) DeepCopy() *PackageInfoList {
	if in == nil {
		return nil
	}
	out := new(PackageInfoList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageInfoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfoSpec) DeepCopyInto(out *PackageInfoSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfoSpec.
func (in *PackageInfoSpec) DeepCopy() *PackageInfoSpec {
	if in == nil {
		return nil
	}
	out := new(PackageInfoSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfoStatus) DeepCopyInto(out *PackageInfoStatus) {
	*out = *in
	if in.Manifest != nil {
		in, out := &in.Manifest, &out.Manifest
		*out = new(PackageManifest)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LastUpdateTimestamp != nil {
		in, out := &in.LastUpdateTimestamp, &out.LastUpdateTimestamp
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfoStatus.
func (in *PackageInfoStatus) DeepCopy() *PackageInfoStatus {
	if in == nil {
		return nil
	}
	out := new(PackageInfoStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfoTemplate) DeepCopyInto(out *PackageInfoTemplate) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfoTemplate.
func (in *PackageInfoTemplate) DeepCopy() *PackageInfoTemplate {
	if in == nil {
		return nil
	}
	out := new(PackageInfoTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageList) DeepCopyInto(out *PackageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Package, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageList.
func (in *PackageList) DeepCopy() *PackageList {
	if in == nil {
		return nil
	}
	out := new(PackageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageManifest) DeepCopyInto(out *PackageManifest) {
	*out = *in
	if in.Scope != nil {
		in, out := &in.Scope, &out.Scope
		*out = new(PackageScope)
		**out = **in
	}
	if in.References != nil {
		in, out := &in.References, &out.References
		*out = make([]PackageReference, len(*in))
		copy(*out, *in)
	}
	if in.Helm != nil {
		in, out := &in.Helm, &out.Helm
		*out = new(HelmManifest)
		(*in).DeepCopyInto(*out)
	}
	if in.Kustomize != nil {
		in, out := &in.Kustomize, &out.Kustomize
		*out = new(KustomizeManifest)
		**out = **in
	}
	if in.Manifests != nil {
		in, out := &in.Manifests, &out.Manifests
		*out = make([]PlainManifest, len(*in))
		copy(*out, *in)
	}
	if in.ValueDefinitions != nil {
		in, out := &in.ValueDefinitions, &out.ValueDefinitions
		*out = make(map[string]ValueDefinition, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.TransitiveResources != nil {
		in, out := &in.TransitiveResources, &out.TransitiveResources
		*out = make([]corev1.TypedLocalObjectReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Entrypoints != nil {
		in, out := &in.Entrypoints, &out.Entrypoints
		*out = make([]PackageEntrypoint, len(*in))
		copy(*out, *in)
	}
	if in.Dependencies != nil {
		in, out := &in.Dependencies, &out.Dependencies
		*out = make([]Dependency, len(*in))
		copy(*out, *in)
	}
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = make([]Component, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageManifest.
func (in *PackageManifest) DeepCopy() *PackageManifest {
	if in == nil {
		return nil
	}
	out := new(PackageManifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageReference) DeepCopyInto(out *PackageReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageReference.
func (in *PackageReference) DeepCopy() *PackageReference {
	if in == nil {
		return nil
	}
	out := new(PackageReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepository) DeepCopyInto(out *PackageRepository) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepository.
func (in *PackageRepository) DeepCopy() *PackageRepository {
	if in == nil {
		return nil
	}
	out := new(PackageRepository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageRepository) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryAuthSpec) DeepCopyInto(out *PackageRepositoryAuthSpec) {
	*out = *in
	if in.Basic != nil {
		in, out := &in.Basic, &out.Basic
		*out = new(PackageRepositoryBasicAuthSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Bearer != nil {
		in, out := &in.Bearer, &out.Bearer
		*out = new(PackageRepositoryBearerAuthSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryAuthSpec.
func (in *PackageRepositoryAuthSpec) DeepCopy() *PackageRepositoryAuthSpec {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryAuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryBasicAuthSpec) DeepCopyInto(out *PackageRepositoryBasicAuthSpec) {
	*out = *in
	if in.Username != nil {
		in, out := &in.Username, &out.Username
		*out = new(string)
		**out = **in
	}
	if in.UsernameSecretRef != nil {
		in, out := &in.UsernameSecretRef, &out.UsernameSecretRef
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Password != nil {
		in, out := &in.Password, &out.Password
		*out = new(string)
		**out = **in
	}
	if in.PasswordSecretRef != nil {
		in, out := &in.PasswordSecretRef, &out.PasswordSecretRef
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryBasicAuthSpec.
func (in *PackageRepositoryBasicAuthSpec) DeepCopy() *PackageRepositoryBasicAuthSpec {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryBasicAuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryBearerAuthSpec) DeepCopyInto(out *PackageRepositoryBearerAuthSpec) {
	*out = *in
	if in.Token != nil {
		in, out := &in.Token, &out.Token
		*out = new(string)
		**out = **in
	}
	if in.TokenSecretRef != nil {
		in, out := &in.TokenSecretRef, &out.TokenSecretRef
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryBearerAuthSpec.
func (in *PackageRepositoryBearerAuthSpec) DeepCopy() *PackageRepositoryBearerAuthSpec {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryBearerAuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryList) DeepCopyInto(out *PackageRepositoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageRepository, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryList.
func (in *PackageRepositoryList) DeepCopy() *PackageRepositoryList {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageRepositoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositorySpec) DeepCopyInto(out *PackageRepositorySpec) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(PackageRepositoryAuthSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositorySpec.
func (in *PackageRepositorySpec) DeepCopy() *PackageRepositorySpec {
	if in == nil {
		return nil
	}
	out := new(PackageRepositorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageRepositoryStatus) DeepCopyInto(out *PackageRepositoryStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageRepositoryStatus.
func (in *PackageRepositoryStatus) DeepCopy() *PackageRepositoryStatus {
	if in == nil {
		return nil
	}
	out := new(PackageRepositoryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSpec) DeepCopyInto(out *PackageSpec) {
	*out = *in
	out.PackageInfo = in.PackageInfo
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make(map[string]ValueConfiguration, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSpec.
func (in *PackageSpec) DeepCopy() *PackageSpec {
	if in == nil {
		return nil
	}
	out := new(PackageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageStatus) DeepCopyInto(out *PackageStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.OwnedResources != nil {
		in, out := &in.OwnedResources, &out.OwnedResources
		*out = make([]OwnedResourceRef, len(*in))
		copy(*out, *in)
	}
	if in.OwnedPackageInfos != nil {
		in, out := &in.OwnedPackageInfos, &out.OwnedPackageInfos
		*out = make([]OwnedResourceRef, len(*in))
		copy(*out, *in)
	}
	if in.OwnedPackages != nil {
		in, out := &in.OwnedPackages, &out.OwnedPackages
		*out = make([]OwnedResourceRef, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageStatus.
func (in *PackageStatus) DeepCopy() *PackageStatus {
	if in == nil {
		return nil
	}
	out := new(PackageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageValueSource) DeepCopyInto(out *PackageValueSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageValueSource.
func (in *PackageValueSource) DeepCopy() *PackageValueSource {
	if in == nil {
		return nil
	}
	out := new(PackageValueSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PartialJsonPatch) DeepCopyInto(out *PartialJsonPatch) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PartialJsonPatch.
func (in *PartialJsonPatch) DeepCopy() *PartialJsonPatch {
	if in == nil {
		return nil
	}
	out := new(PartialJsonPatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlainManifest) DeepCopyInto(out *PlainManifest) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlainManifest.
func (in *PlainManifest) DeepCopy() *PlainManifest {
	if in == nil {
		return nil
	}
	out := new(PlainManifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueConfiguration) DeepCopyInto(out *ValueConfiguration) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(ValueReference)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueConfiguration.
func (in *ValueConfiguration) DeepCopy() *ValueConfiguration {
	if in == nil {
		return nil
	}
	out := new(ValueConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueDefinition) DeepCopyInto(out *ValueDefinition) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Constraints.DeepCopyInto(&out.Constraints)
	if in.Targets != nil {
		in, out := &in.Targets, &out.Targets
		*out = make([]ValueDefinitionTarget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueDefinition.
func (in *ValueDefinition) DeepCopy() *ValueDefinition {
	if in == nil {
		return nil
	}
	out := new(ValueDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueDefinitionConstraints) DeepCopyInto(out *ValueDefinitionConstraints) {
	*out = *in
	if in.Min != nil {
		in, out := &in.Min, &out.Min
		*out = new(int)
		**out = **in
	}
	if in.Max != nil {
		in, out := &in.Max, &out.Max
		*out = new(int)
		**out = **in
	}
	if in.MinLength != nil {
		in, out := &in.MinLength, &out.MinLength
		*out = new(int)
		**out = **in
	}
	if in.MaxLength != nil {
		in, out := &in.MaxLength, &out.MaxLength
		*out = new(int)
		**out = **in
	}
	if in.Pattern != nil {
		in, out := &in.Pattern, &out.Pattern
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueDefinitionConstraints.
func (in *ValueDefinitionConstraints) DeepCopy() *ValueDefinitionConstraints {
	if in == nil {
		return nil
	}
	out := new(ValueDefinitionConstraints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueDefinitionMetadata) DeepCopyInto(out *ValueDefinitionMetadata) {
	*out = *in
	if in.Hints != nil {
		in, out := &in.Hints, &out.Hints
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueDefinitionMetadata.
func (in *ValueDefinitionMetadata) DeepCopy() *ValueDefinitionMetadata {
	if in == nil {
		return nil
	}
	out := new(ValueDefinitionMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueDefinitionTarget) DeepCopyInto(out *ValueDefinitionTarget) {
	*out = *in
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		*out = new(corev1.TypedObjectReference)
		(*in).DeepCopyInto(*out)
	}
	if in.ChartName != nil {
		in, out := &in.ChartName, &out.ChartName
		*out = new(string)
		**out = **in
	}
	out.Patch = in.Patch
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueDefinitionTarget.
func (in *ValueDefinitionTarget) DeepCopy() *ValueDefinitionTarget {
	if in == nil {
		return nil
	}
	out := new(ValueDefinitionTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueReference) DeepCopyInto(out *ValueReference) {
	*out = *in
	if in.ConfigMapRef != nil {
		in, out := &in.ConfigMapRef, &out.ConfigMapRef
		*out = new(ObjectKeyValueSource)
		**out = **in
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(ObjectKeyValueSource)
		**out = **in
	}
	if in.PackageRef != nil {
		in, out := &in.PackageRef, &out.PackageRef
		*out = new(PackageValueSource)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueReference.
func (in *ValueReference) DeepCopy() *ValueReference {
	if in == nil {
		return nil
	}
	out := new(ValueReference)
	in.DeepCopyInto(out)
	return out
}
