# Package Operator

The package operator follows the [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and has two controllers:

## Package Controller

The Package controller manages the `Package` resources of the cluster. 

Whenever a `Package` has been created, changed or deleted these changes will be picked up and applied by the Package controller.

## PackageInfo Controller

The PackageInfo controller syncs the relevant `PackageInfo` resources with the manifests defined in the package repository.
