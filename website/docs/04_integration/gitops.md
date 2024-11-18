# GitOps

All Glasskube features can be used in GitOps powered environments, as all the internal logic is built upon custom resource definitions.
For example, an installation of a cluster package is being represented by a `ClusterPackage` custom resource in your cluster. The same
goes for `Package` and `PackageRepository` resources.

In our [GitOps template](https://github.com/glasskube/gitops-template) we explain how Glasskube can be set up together with ArgoCD,
using the `glasskube bootstrap git` command.

If you prefer a different GitOps tool or need a more customized solution, you can use `glasskube bootstrap --dry-run -o yaml`
in order to generate the Glasskube manifests that you can put into your GitOps repository.
