# Package Configuration

## Problem analysis

At Glasskube, what we mean when we talk about package configuration, is a controlled alteration of the resources that are part of a package depending on values specified by the user.
This can be thought of as three sub-problems:

### What values are available for a given package

The author of a package may include a declaration of the possible configuration items in the package manifest, where each package may have multiple **value definitions**.
A value definition contains information to help clients display an appropriate form item for data entry and constraints to be validated by the client as well as the package operator.

### How does a value affect the deployed package resources

Additionally, a value definition describes some alterations to the deployed package.
We call those alterations **targets**.
Each value definition may contain a number of targets.
A target can be one of two things:

1. A change to the default values of a helm release contained in the package
2. A JSON patch that should be applied to a resource contained in a `manifests` entry of the package

### How does a value affect the deployed package resources

A user, when installing a package using Glasskube, may declare a **value configuration** for each value definition of that package.
A value configuration can hold either a **literal value** or a **reference value**.
Literal values represent a simple value (e.g. a string entered via a text field).
Reference values represent references to values in other resources in the same Kubernetes cluster.
Such references can be secrets, configmaps or other packages.

Values are non-mandatory by default, however, a package author may opt to make any of their packages values required by
specifying a constraint on that value definition.
If a package has no value configuration for a given value definition that is non-mandatory, that values targets will not be
applied and it is the package authors responsibility to ensure that their package also works in this case.

## Design proposal

The `PackageManifest` has a property `Values` of type `map[string]ValueDefinition`.
The key in this map is referred to as that values **name**
`ValueDefinition` is a struct with the following properties:

- **`Type`** (`string` enum):
  Every value must have a type, so that we know what kind of input field to show for this value.
  Initially, this can be one of `boolean`, `options`, `text`, `number` but it is possible to add more types in future releases.
- **`Metadata`**:
  A colletion of (mostly) UI-related metadata with the following (optional) properties:
  - **`Label`** (`string`):
    The label is used to denote an input field related to this value in a UI.
    By default the name of the value should be used.
  - **`Description`** (`string`):
    The description can be used to give more context to a value.
  - **`Hints`** (`string` enum):
    Hints offer package maintainers the ability to make some elements of the UI more prominent.
    For example, every value can be set to reference the value of a secret key, but if a value has the
    "SuggestSecretRef" hint, this option can be highlighted by the UI or enabled by default.
    _Available hints and whether they will be included in the initial release is TBD_
- **`DefaultValue`** (`string`):
  The default value is pre-selected/pre-filled in the form field of this value for new packages.
- **`Options`** (`[]string`):
  Available choices for values of type options.
  Should be ignored for other types.
- **`Constraints`**:
  Specifying a number of constraints is possible.
  Available constraints are
  `Required` (`bool`), `Min` (`int`), `Max` (`int`), `MinLength` (`int`), `MaxLength` (`int`), `Pattern` (`string`).
  These should be checked by the UI, as well as by the validating webhook.
  Not all constrains apply to all types of value. Non-applicable constraints are ignored.
  For example a "text" value with constraints.max = 3 is the same as a "text" value with no constraints.
- **`Targets`**:
  Where to apply this value.
  Either the name of a helm chart, or a `TypedObjectReference` combined with patch information.
  Initially, the idea is to use RFC 6902 JSON patches but this is still TBD.
  We use Unstructured for plain resources, which already supports setting values via a kind of JSON path,
  but it does not support setting values in lists.
  So, for example, it would not be possible to change something in the container of a deployment, since the containers are a list.
  Maybe [evanphx/json-patch](https://github.com/evanphx/json-patch) can be a useful alternative, but it only works on byte slices.

The `PackageSpec` has a property `Values` of type `map[string]ValueConfiguration`.
The key in this map used to identify the corresponding `ValueDefinition` with the same name in the `PackageManifest`s `Values` map.
A `ValueConfiguration` must have exactly one of the following properties:

- **`Value`**:
  A literal value.
- **`ValueFrom`**:
  To represent a reference value. `ValueFrom` must have exactly one of the following properties:
  - **`ConfigMapRef`**:
    To reference a `Key` of a ConfigMap with `Name` in `Namespace`
  - **`SecretRef`**:
    To reference a `Key` of a Secret with `Name` in `Namespace`
  - **`PackageRef`**:
    To reference the value of the `ValueConfiguration` with name `Value` of a package with `Name`.

## Examples

```yaml title="PackageManifest with a simple value specification"
name: foo
helm:
  repositoryUrl: 'https://charts.example.com'
  chartName: 'foo'
  chartVersion: 'v1.0.0'
  values: {}
valueDefinitions:
  ingress:
    type: 'boolean'
    label: 'Enable Ingress'
    description: 'Whether an ingress resource should be created for this Package'
    defaultValue: 'true'
    targets:
      - chartName: 'foo'
        patch:
          - op: 'add'
            path: 'ingress/enabled'
```

```yaml title="PackageManifest with a value specification that has multiple targets"
name: foo
valueDefinitions:
  host:
    type: 'text'
    constraints:
      required: true
    targets:
      - resource:
          kind: 'Ingress'
          apiVersion: 'networking.k8s.io/v1'
          name: 'foo'
          namespace: 'foo'
        patch:
          - op: 'add'
            path: '/spec/rules/0/host'
          - op: 'add'
            path: '/spec/tls/0/hosts/-'
      - resource:
          apiVersion: 'apps/v1'
          kind: 'Deployment'
          name: 'foo'
          namespace: 'foo'
        patch:
          - op: 'add'
            path: '/spec/template/spec/containers/0/env/-'
        valueTemplate: |
          { "name": "APP_HOST", "value": "{{ .value }}" }
```

```yaml title="Package with a variety of value configurations"
apiVersion: 'packages.glasskube.dev/v1alpha1'
kind: 'Package'
metadata:
  name: 'foo'
spec:
  packageInfo:
    name: 'foo'
    version: 'v1.0.0'
  values:
    ingress:
      value: 'false'
    host:
      valueFrom:
        configMapRef:
          name: 'foo-prod-config'
          value: 'host'
    apiKey:
      valueFrom:
        secretRef:
          name: 'api-key-secret'
          key: 'apiKey'
```

## Known Limitations/caveats

- Value configurations can not have list types
- More possibilities for deadlocks with required values that reference other packages
