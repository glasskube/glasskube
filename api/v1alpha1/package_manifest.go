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

type HelmManifest struct {
	// RepositoryUrl is the remote URL of the helm repository. This is the same URL you would use
	// if you use "helm repo add ...".
	RepositoryUrl string `json:"repositoryUrl" jsonschema:"required"`
	// ChartName is the name of the chart that represents this package.
	ChartName string `json:"chartName" jsonschema:"required"`
	// ChartVersion of the chart that should be installed.
	ChartVersion string `json:"chartVersion" jsonschema:"required"`
	// Values that should be used for the helm release
	Values *JSON `json:"values,omitempty"`
}

type KustomizeManifest struct {
}

// PackageEntrypoint defines a service port a user may use to access the package
type PackageEntrypoint struct {
	// Name of this entrypoint. Used for "glasskube open [package-name] [entypoint-name]" if more
	// than one entrypoint exists. Optional if the package only has one entrypoint.
	Name string `json:"name,omitempty"`
	// ServiceName is the name of a service that is part of
	ServiceName string `json:"serviceName" jsonschema:"required"`
	// Port of the service to bind to
	Port int32 `json:"port" jsonschema:"required"`
}

type PackageManifest struct {
	Name             string `json:"name" jsonschema:"required"`
	ShortDescription string `json:"shortDescription,omitempty"`
	IconUrl          string `json:"iconUrl,omitempty"`
	// Helm instructs the controller to create a helm release when installing this package.
	Helm *HelmManifest `json:"helm,omitempty"`
	// Kustomize instructs the controller to apply a kustomization when installing this package [PLACEHOLDER].
	Kustomize *KustomizeManifest `json:"kustomize,omitempty"`
	// DefaultNamespace to install the package. May be overridden.
	DefaultNamespace string              `json:"defaultNamespace,omitempty" jsonschema:"required"`
	Entrypoints      []PackageEntrypoint `json:"entrypoints,omitempty"`
}
