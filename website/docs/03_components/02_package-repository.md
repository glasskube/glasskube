# Package Repository

The package repository is where `PackageManifest`s are stored, searched for and maintained.
Currently only the glasskube packages repository is supported: [`glasskube/packages`](https://github.com/glasskube/packages)

A `PackageManifest` contains all relevant information needed for identifying and installing a package. 
It can contain either a Helm resource (as used in [cert-manager](https://github.com/glasskube/packages/blob/main/packages/cert-manager/package.yaml)), or a link to a manifest (as used for [cyclops](https://github.com/glasskube/packages/blob/main/packages/cyclops/package.yaml)).
