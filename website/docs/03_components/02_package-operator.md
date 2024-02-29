# Package Operator

The package operator follows the [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and has two controllers:

## Package Controller

The Package controller manages the `Package` resources of the cluster. 

Whenever a `Package` has been created, changed or deleted these changes will be picked up and applied by the Package controller.

## PackageInfo Controller

The PackageInfo controller syncs the relevant `PackageInfo` resources with the manifests defined in the package repository.

## Handling Package Updates

A Package may have it's `.spec.version` set. 
This instructs the operator to install this exact version of the package. 
We also call this version pinning.
A Package's `.spec.version` may also be empty. 
In that case the operator will try to automatically update the package to it's latest version.

To update a package with a pinned version, run `glasskube update <package>`. 
This will upate the package to the latest version. 

```mermaid
---
title: Package Reconciliation
---
flowchart TB
  start((Start))
  pinned{.spec.version present}
  ensurePinned[Create PackageInfo for .spec.version if missing]
  ensureLatest[Create PackageInfo for latest version if missing]
  ready{PackageInfo is ready}
  ensureManifest[Ensure manifest]
  ensureManifestOk{Success}
  cleanup[Clean up old PackageInfos]
  end1([End])

  start --> pinned
  pinned -- yes --> ensurePinned
  pinned -- no --> ensureLatest
  ensureLatest --> ready
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
  pinned{.spec.version present}
  noVersion{Versionless package?}
  DP[Fetch package.yaml]
  DV[Fetch versions.yaml]
  DVP[Fetch version/package.yaml]
  end1([End])

  start --> pinned
  pinned -- no --> noVersion
  noVersion -- yes --> DP
  noVersion -- no --> DV
  pinned -- yes --> DVP
  DV -- if ok --> DVP
  DP --> end1
  DVP --> end1
```

### FAQ
**How is the latest version determined?**

The package's `versions.yaml` is fetched from the repository. This file contains all available versions.

**How is a specific version of a package fetched?**

Instead of fetching `repository.xyz/package-name/package.yaml`, the operator fetches `repository.xyz/package-name/version/package.yaml`

Check the [package repository documentation](../package-repository#structure) for more information.

