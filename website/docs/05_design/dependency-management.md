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

### Visualisation:

**Abbreviations:**

- **P** Package that is going to be installed
- **D** Dependency that package *P* requires
- **DV** Installed version of package *D*
- **PDV** Version constraint for package *D* as defined in the dependency relation of package *P*
- **XDV** / **YDV** Version constraint for package *D* as defined in the dependency relation of already installed package *X* / *Y* that also have a dependency on *D*

```mermaid
flowchart TD
Start --> Check_P_Req_D("Does P require a<br>version range of D?")
%% Branch when P requires no version range of D
Check_P_Req_D -->|No| P_NoReq_D__Check_D_Exist("Does D exist?")
P_NoReq_D__Check_D_Exist --->|Yes| State_Fulfilled["P -> D is fulfilled"]
P_NoReq_D__Check_D_Exist -->|No| State_Install_D_latest["Install D pinned in latest(D)"]
%% Branch when P requires D to be in version range PDV
Check_P_Req_D -->|Yes| P_Req_D__Check_D_Exist("Does D exist?")
P_Req_D__Check_D_Exist -->|Yes| Check_DV_Inside_PDV("Is DV inside PDV?")
P_Req_D__Check_D_Exist -->|No| State_Install_D_PDV["Install D pinned in max_available(PDV)"]
Check_DV_Inside_PDV -->|Yes| State_Fulfilled
Check_DV_Inside_PDV -->|No| Check_DV_Less_PDV("Is DV < PDV?")
Check_DV_Less_PDV -->|Yes| Check_OtherPkgs_Req_D("Are there other<br>existing packages<br>dependent on D<br>requiring a version range?")
Check_OtherPkgs_Req_D -->|Yes| State_Conflict_MaybeResolvable["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Might be resolvable if XDV, YDV, and PDV overlap"]
Check_OtherPkgs_Req_D -->|No| State_Conflict_Resolvable["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Resolvable by updating D to max_available(PDV)"]
Check_DV_Less_PDV -->|No| State_Conflict_NotResolvable["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Not resolvable because P does not support D in DV yet"]
%% Styling end nodes
style State_Fulfilled fill:#006400
style State_Install_D_latest fill:#B8860B
style State_Install_D_PDV fill:#B8860B
style State_Conflict_MaybeResolvable fill:#8B0000
style State_Conflict_Resolvable fill:#8B0000
style State_Conflict_NotResolvable fill:#8B0000
```
