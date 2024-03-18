# Package Repository

The package repository is where `PackageManifest`s are stored, searched for and maintained.
Currently only the glasskube packages repository is supported: [`glasskube/packages`](https://github.com/glasskube/packages)

A `PackageManifest` contains all relevant information needed for identifying and installing a package. 
It can contain either a Helm resource (as used in [cert-manager](https://github.com/glasskube/packages/blob/main/packages/cert-manager/package.yaml)), or a link to a manifest (as used for [cyclops](https://github.com/glasskube/packages/blob/main/packages/cyclops/package.yaml)).

## Structure

A package repository must use the following directory structure to be fully compatible with Glasskube:

```
|-- index.yaml
|-- package-a/
  |-- versions.yaml
  |-- v1.2.3/
  | |-- package.yaml
  |-- v1.3.2/
    |-- package.yaml
```

The root `index.yaml` contains a list of all packages available from this repository. It is used primarily by client software to aid explorability.
All files related to a package reside in a directory that must have the same name as the package. 
Inside a package's directory there must be a `versions.yaml` that contains a list of all versions available for this package.
There must be a subdirectory for each version containing a `package.yaml` file.
A `package.yaml` contains a manifest of that package which holds information such as longer descriptions and included files.

## Package Manifest

### Dependencies

A package can declare dependencies that need to exist in a cluster, before the desired package can be installed. 
Each dependency is a Glasskube package identified by its name. Optionally, a specific version or version range can be defined.
