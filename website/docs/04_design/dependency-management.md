# Dependency Management

Dependency Management is a cross-cutting concern that is being handled in all glasskube components (GUI, CLI, Operator).
The following decision tree states how the Package Operator is handling dependencies.

## Package Operator – reconciling package P depending on package D (P -> D):

### Assumptions:
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
