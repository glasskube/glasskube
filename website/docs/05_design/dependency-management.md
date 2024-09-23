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

<div style={{maxWidth: '100%', overflow: 'auto'}}>
<div id="dep-diagram" style={{width: '1600px'}}>

```mermaid
%%{init: {'themeVariables': {'flowchart': {'useWidth': 10000} }}}%%
flowchart LR

Start --> A("Does P require a<br>version range of D?")

%% Branch when P requires no version range of D

A -->|No| B("Does D exist?")

B -->|Yes| C("Are there other packages<br>dependent on D?")

C -->|No| Fulfilled["P -> D is fulfilled"]

C -->|Yes| E("Do X and Y require<br>no version range of D?")

E ---->|Yes| Fulfilled

E ---->|No| Fulfilled

B ----->|No| F["Install D pinned in latest(D)"]

%% Branch when P requires D to be in version range PDV

A -->|Yes| G("Does D exist?")

G -->|Yes| H("Are there other<br>existing packages<br>dependent on D<br>requiringa version range?")

H -->|No| I("Is DV inside PDV?")

I ---->|Yes| Fulfilled

I ---->|No| K("Is DV < PDV?")

K --->|Yes| L["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Resolvable by updating D to max_available(PDV)"]

K --->|No| M["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Not resolvable because P does not support D in DV yet"]

H -->|Yes| N("Is DV inside PDV?")

N -------->|Yes| Fulfilled

N ---->|No| P("Is DV < PDV?")

P --->|Yes| Q["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Might be resolvable if XDV, YDV, and PDV overlap"]

P --->|No| R["P -> D not fulfilled<br><b>Dependency Conflict</b><br>Not resolvable because P does not support D in DV yet"]

G --------->|No| S["Install D pinned in max_available(PDV)"]

style Fulfilled fill:#006400
style F fill:#B8860B
style S fill:#B8860B
style L fill:#8B0000
style M fill:#8B0000
style Q fill:#8B0000
style R fill:#8B0000
```

</div>
</div>
