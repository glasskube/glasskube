# Package Repository

The package repository is where `PackageManifest`s are stored, searched for and maintained.
The default repository is the central Glasskube package repository [`glasskube/packages`](https://github.com/glasskube/packages), however custom repositories can be added. 

A `PackageManifest` contains all relevant information needed for identifying and installing a package. 
It can contain either a Helm resource (as used in [cert-manager](https://github.com/glasskube/packages/blob/main/packages/cert-manager/v1.16.2%2B1/package.yaml)), or a link to a manifest (as used for [cyclops](https://github.com/glasskube/packages/blob/main/packages/cyclops/v0.9.1%2B1/package.yaml)).

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

### Version Numbers

The version number of a package must follow the [semver specification](https://semver.org), with the additional constraint that the build number of a version is only allowed to consist of digits. 
Although the specification states that [build numbers must be ignored](https://semver.org/#spec-item-10) when determining version precedence, 
we see no other way than to do so. Therefore, we only allow digits in the build number, such that we can decide which version is newer. 

This is important for Glasskube to distinguish between two different versions of a package manifest, that might have the same underlying software version of a package.
For example, there could be a package  `kubernetes-dashboard` in version `2.7.0`. 
The preferred way to make this version available in the Glasskube package repository, is to create a version `2.7.0+1`. 
If the package manifest turns out to be faulty and need to be corrected (e.g. typos, wrong entrypoints, wrong dependencies), or if some metadata needs to be changed (e.g. links, description), 
the maintainer can add a version `2.7.0+2` without changing the underlying app version of `2.7.0`. 

## Package Manifest

### Dependencies

A package can declare dependencies that need to exist in a cluster, before the desired package can be installed. 
Each dependency is a Glasskube package identified by its name. Optionally, a specific version or version range can be defined.

#### Version Ranges

Sometimes it is important to pin down what versions of a dependency should be used by a package. 
Often, this will not only be one version, but a range of versions. There is no common specification for semver ranges, 
but there seem to be some common expectations of to what version range specifications look like.
We mostly rely on the [Masterminds/semver](https://github.com/Masterminds/semver) package to do version constraint checks, which itself works with version range syntax close to js/npm and Rust/Cargo. 

Please note that we do not allow build numbers to be part of such version ranges. 
If there is some change in a required package, that would make a dependant package incompatible, this change needs to be reflected in the actual app version anyway. 
