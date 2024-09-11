# Package Manifest

## Properties

| Name                | Type                                                                                                                                | Required / Default | Description                 |
| ------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ------------------ | --------------------------- |
| scope               | string                                                                                                                              | `"Namespaced"`     | One of: Cluster, Namespaced |
| name                | string                                                                                                                              | required           | Name of the package         |
| shortDescription    | string                                                                                                                              |                    |                             |
| longDescription     | string                                                                                                                              |                    |                             |
| defaultNamespace    | string                                                                                                                              | required           |
| references          | [][PackageReference](#packagereference)                                                                                             |                    |
| iconUrl             | string                                                                                                                              |                    |
| helm                | [HelmManifest](#helmmanifest)                                                                                                       |                    |
| manifests           | [][PlainManifest](#plainmanifest)                                                                                                   |                    |
| valueDefinitions    | map[string][ValueDefinition](#valuedefinition)                                                                                      |                    |
| transformations     | [][TransformationDefinition](#transformationdefinition)                                                                             |                    |                             |
| transitiveResources | [][TypedLocalObjectReference](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/typed-local-object-reference/) |                    |
| entrypoints         | [][PackageEntrypoint](#packageentrypoint)                                                                                           |                    |
| dependencies        | [][Dependency](#dependency)                                                                                                         |                    |
| components          | [][Component](#component)                                                                                                           |                    |

## Subresources

### PackageReference

| Name  | Type   | Required / Default | Description |
| ----- | ------ | ------------------ | ----------- |
| label | string | required           |             |
| url   | string | required           |             |

### HelmManifest

| Name          | Type       | Required / Default | Description |
| ------------- | ---------- | ------------------ | ----------- |
| repositoryUrl | string     | required           |             |
| chartName     | string     | required           |             |
| chartVersion  | string     | required           |             |
| values        | _any JSON_ |                    |             |

### PlainManifest

| Name             | Type   | Required / Default | Description                                  |
| ---------------- | ------ | ------------------ | -------------------------------------------- |
| url              | string | required           |                                              |
| defaultNamespace | string |                    | overrides the package-level defaultNamespace |

### ValueDefinition

| Name         | Type                                                      | Required / Default | Description                            |
| ------------ | --------------------------------------------------------- | ------------------ | -------------------------------------- |
| type         | string                                                    | required           | One of: boolean, text, number, options |
| metadata     | [ValueDefinitionMetadata](#valuedefinitionmetadata)       |                    |                                        |
| defaultValue | string                                                    |                    |                                        |
| options      | []string                                                  |                    |                                        |
| constrains   | [ValueDefinitionConstraints](#valuedefinitionconstraints) |                    |                                        |
| targets      | [ValueDefinitionTarget](#valuedefinitiontarget)           |                    |                                        |

### TransformationDefinition

| Name    | Type                                              | Required / Default | Description |
| ------- | ------------------------------------------------- | ------------------ | ----------- |
| source  | TransformationSource                              | required           |             |
| targets | [][ValueDefinitionTarget](#valuedefinitiontarget) | required           |             |

### ValueDefinitionMetadata

| Name        | Type     | Required / Default | Description                                    |
| ----------- | -------- | ------------------ | ---------------------------------------------- |
| label       |          |                    | form label to show on the UI                   |
| description |          |                    | longer description to show along side the form |
| hints       | []string |                    | currently unused                               |

### ValueDefinitionConstraints

| Name      | Type   | Required / Default | Description                                                     |
| --------- | ------ | ------------------ | --------------------------------------------------------------- |
| required  | bool   | required           | whether this value **must** be specified. It may still be empty |
| min       | int    |                    | minimum value for values with type number                       |
| max       | int    |                    | maximum number for values with type number                      |
| minLength | int    |                    | minimum length for values with type text                        |
| maxLength | int    |                    | maximum lenght for values with type text                        |
| pattern   | string |                    | regex pattern for validation                                    |

### ValueDefinitionTarget

| Name          | Type                                  | Required / Default | Description                                   |
| ------------- | ------------------------------------- | ------------------ | --------------------------------------------- |
| resource      | TypedObjectReference                  |                    | reference of a resource owned by this package |
| chartName     | string                                |                    | name of a helm chart managed by this package  |
| patch         | [PartialJsonPatch](#partialjsonpatch) | required           |                                               |
| valueTemplate | string                                |                    |                                               |

Either `resource` or `chartName` must be specified.

### PartialJsonPatch

| Name | Type   | Required / Default | Description |
| ---- | ------ | ------------------ | ----------- |
| op   | string | required           |             |
| path | string | required           |             |

The `value` to create a complete JSON Patch is supplied by the controller.
See https://jsonpatch.com/ for a complete reference.

### TransformationSource

| Name     | Type                                                                                                                              | Required / Default | Description                                                                                    |
| -------- | --------------------------------------------------------------------------------------------------------------------------------- | ------------------ | ---------------------------------------------------------------------------------------------- |
| resource | [TypedLocalObjectReference](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/typed-local-object-reference/) |                    | leave empty to reference the current package                                                   |
| path     | string                                                                                                                            | required           | JSON path to a property of the resource ([reference](https://goessner.net/articles/JsonPath/)) |

### PackageEntrypoint

| Name        | Type   | Required / Default | Description |
| ----------- | ------ | ------------------ | ----------- |
| name        | string |                    |             |
| serviceName | string | required           |             |
| port        | int32  | required           |             |
| localPort   | int32  |                    |             |
| scheme      | string |                    |             |

### Dependency

| Name    | Type   | Required / Default | Description                             |
| ------- | ------ | ------------------ | --------------------------------------- |
| name    | string | required           |                                         |
| version | string |                    | a semver constraint for this dependency |

### Component

| Name          | Type                                                             | Required / Default | Description                            |
| ------------- | ---------------------------------------------------------------- | ------------------ | -------------------------------------- |
| name          | string                                                           | required           |                                        |
| installedName | string                                                           |                    | name suffix for the created `Package`  |
| version       | string                                                           |                    | a semver constraint for this component |
| values        | map[string][InlineValueConfiguration](#inlinevalueconfiguration) |                    | specify values for this component      |

### InlineValueConfiguration

A stripped down variant of a package's value configuration that only supports directly specified values and no reference values.

| Name  | Type   | Required / Default | Description |
| ----- | ------ | ------------------ | ----------- |
| value | string | required           |             |

## Complete Example

```yaml title="package.yaml"
name: example
scope: Namespaced
defaultNamespace: default
iconUrl: https://example.com/logo.jpeg
shortDescription: An example package for this documentation
longDescription: |
  An extended description of this package. 

  Markdown is **supported**.
references:
  - label: GitHub
    url: https://github.com/example/example
entrypoints:
  - name: ui
    serviceName: example-ui
    port: 443
    localPort: 8443
    scheme: https
helm:
  repositoryUrl: https://charts.example.com/
  chartName: example
  chartVersion: v1.0.0
  values:
    config:
      database:
        driver: postgresql
manifests:
  - url: https://github.com/example/example/releases/v1.0.0/example.yaml
dependencies:
  - name: cloudnative-pg
    version: '1.x.x'
components:
  - name: postgresql
    installedName: db
    version: '>=1.0.0'
    values:
      enableSuperuserAccess:
        value: 'true'
transitiveResources:
  - apiVersion: v1
    kind: Secret
    name: db-app
valueDefinitions:
  someConfigAttribute:
    type: text
    metadata:
      label: Some text
      description: Longer text
    constraints:
      required: true
      minLength: 5
      maxLength: 24
    targets:
      - chartName: example
        patch:
          op: add
          path: /config/database/dbName
      - resource:
          apiVersion: v1
          kind: ConfigMap
          name: example
        patch:
          op: add
          path: /data/DB_NAME
transformations:
  - source:
      path: '{ $.metadata.name }'
    targets:
      - chartName: example
        patch:
          op: add
          path: /config/database/dbHost
          valueTemplate: '"{{.}}-db-rw"'
```

## JSON Schema

An up to date JSON schema file is available at https://glasskube.dev/schemas/v1/package-manifest.json.
