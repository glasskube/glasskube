# Package Operator

The package operator follows the [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and has two controllers:

## Package Controller

The Package controller manages the `Package` resources of the cluster. 

Whenever a `Package` has been created, changed or deleted these changes will be picked up and applied by the Package controller.
The Package controller also makes sure that dependencies of a `Package` are met without conflicting installations, by applying the logic described [here](../04_design/dependency-management.md).

## PackageInfo Controller

The PackageInfo controller syncs the relevant `PackageInfo` resources with the manifests defined in the package repository.

## Handling Package Updates

A Package must have it's `.spec.version` set. 
This instructs the operator to install this exact version of the package. 
We also call this version pinning.

To update a package with a pinned version, run `glasskube update <package>`. 
This will upate the package to the latest version. 

```mermaid
---
title: Package Reconciliation
---
flowchart TB
  start((Start))
  ensurePinned[Create PackageInfo for .spec.version if missing]
  ready{PackageInfo is ready}
  ensureManifest[Ensure manifest]
  ensureManifestOk{Success}
  cleanup[Clean up old PackageInfos]
  end1([End])

  start --> ensurePinned
  ensurePinned --> ready
  ready -- yes --> ensureManifest
  ready -- no --> end1
  ensureManifest --> ensureManifestOk
  ensureManifestOk -- yes --> cleanup
  ensureManifestOk -- no --> end1
  cleanup --> end1
```

```mermaid
---
title: PackageInfo Reconciliation
---
flowchart TB
  start((Start))
  DVP[Fetch version/package.yaml]
  end1([End])

  start --> DVP
  DVP --> end1
```

### FAQ
**How is the latest version determined?**

The package's `versions.yaml` is fetched from the repository. This file contains all available versions and a field `latestVersion`. 
Note that `latestVersion` might not be equal to the actual latest available version – think of a package where `latestVersion = v1.1.7`, 
while there might already be a prerelease of a new major version like `v2.0.0-alpha.1`. 

**How is a specific version of a package fetched?**

Instead of fetching `repository.xyz/package-name/package.yaml`, the operator fetches `repository.xyz/package-name/version/package.yaml`

Check the [package repository docs](./03_package-repository.md#structure) and [dependency management docs](../04_design/dependency-management.md) for more information.
