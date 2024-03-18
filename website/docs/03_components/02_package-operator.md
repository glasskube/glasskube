# Package Operator

The package operator follows the [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and has two controllers:

## Package Controller

The Package controller manages the `Package` resources of the cluster. 

Whenever a `Package` has been created, changed or deleted these changes will be picked up and applied by the Package controller.
The Package controller also makes sure that dependencies of a `Package` are met without conflicting installations, by applying the logic described [here](#dependency-management).

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

Check the [package repository documentation](../package-repository#structure) for more information.

## Dependency Management

Dependency Management is a cross-cutting concern that is being handled in all glasskube components (GUI, CLI, Operator). 
The following decision tree states how the Package Operator is handling dependencies. 

### Package Operator – reconciling package P depending on package D (P -> D):

#### Assumptions:
* Each involved referred package has status Ready, i.e. none of the referred packages are currently being deleted or updated, and their installation has not failed.
* Each involved referred package has a `Spec.PackageInfo.Version` set, and it is equal to its `Status.Version`.
* When the result of a situation is a dependency conflict, it might either be resolvable or not. Either way, the operator does not resolve such a conflict directly, but rather
the components interacting with the user (CLI, UI) need to guide them through potential resolution. Consequently, the only time the operator does resolve an unfulfilled
dependency, the "result" is denoted as `install`. 

```
if P requires no version range of D
  if D exists (trivially P -> D is fulfilled anyway)
    if no other package dependent on D
      * P -> D is fulfilled
    if other existing packages X, Y dependent on D
      if X and Y require no version range of D
        * P -> D is fulfilled
      if X requires D to be in version range XDV, or Y requires D to be in version range YDV
        * P -> D is fulfilled
  if D does not exist
    * install D pinned in latest(D)
if P requires D to be in version range PDV
  if D exists (let DV be the version of D)
    if no other existing package dependent on D requires a version range of D
      if DV inside PDV
        * P -> D is fulfilled
      if DV < PDV
        * P -> D not fulfilled – Dependency Conflict
        * resolvable by updating D to max_available(PDV)
      if DV > PDV
        * P -> D not fulfilled – Dependency Conflict
        * not resolvable because P does not support using D in DV yet
    if other existing packages X, Y dependent on D, with X requiring XDV, Y requiring YDV
      if DV inside PDV
        * P -> D is fulfilled
      if DV < PDV
        * P -> D not fulfilled – Dependency Conflict
        * might be resolvable if XDV, YDV and PDV overlap
      if DV > PDV
        * P -> D not fulfilled – Dependency Conflict
        * not resolvable because P does not support using D in DV yet
  if D does not exist
    * install D pinned in max_available(PDV)
```
