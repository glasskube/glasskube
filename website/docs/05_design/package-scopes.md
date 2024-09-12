# Package Scopes

At Glasskube, we understand that granular management of Kubernetes resources is crucial. That's why we offer two distinct resource types that determine how widely a package can be installed:

-   **Cluster-scoped packages (ClusterPackages):** Installed once per cluster and accessible from all namespaces. Ideal for cluster-wide functionality and shared resources.  
-   **Namespace-scoped packages (Packages):** Installed in specific namespaces, offering isolation and flexibility. Perfect for application-specific tools or multiple instances of the same package.
    
Let's explore when you'd opt for each scope.

### ClusterPackages: Managing Resources at the Cluster Level

ClusterPackages are the prefered when your package:

-   Needs access to all namespaces in your cluster.
-   Utilizes `ClusterRoles` and `ClusterRoleBindings` for authorization.
-   Provides functionality or manages resources that are relevant cluster-wide, such as networking, monitoring, or security tools.
    
### Packages: Isolation at the Namespace Level

Packages are ideal when you require:

-   Isolation: Your package's functionality is contained within a specific namespace.  
-   Flexibility: You need to install the same package multiple times, each instance operating independently in its own namespace.  
-   When the package uses `Roles` and `RoleBindings` for access control.  
-   Application-specific tools or components that don't need cluster-wide access.  

### Dependencies and Package Scopes

The distinction between ClusterPackages and Packages also impacts dependencies.

-   **ClusterPackages:** If multiple ClusterPackages/Packages depend on another ClusterPackage, they'll share that dependency. This is also the expected behaviour for dependency on operators. 
-   **Packages:** When a ClusterPackage/Package depends on a Package, that Package is instantiated as a `component` alongside the original package. Each instance of the original package gets its own, isolated instance of the dependent Package.
    
### Defining Package Scope: A Note for Package Authors

While Glasskube packages default to cluster-scoped, we strongly recommend explicitly defining the scope in your package definition. This ensures clarity and avoids any potential ambiguity about how your package is intended to be used.

To designate a package scope, include the `scope`: `Namespaced` or `Cluster` key-value pair.

```
name: <packageName>
shortDescription: <description>
scope: Namespaced or Cluster
defaultNamespace: <defaultNamespaceName>
iconUrl: <iconURL>
```
For a complete overview of all available options in the package manifest, refer to the [full package manifest reference page](https://glasskube.dev/docs/reference/package-manifest/).

### Evolving Package Scopes

As Glasskube continues to evolve, we recognize that the distinction between cluster-scoped and namespace-scoped packages might need further refinement to address a broader range of use cases, at the moment each package can only be assigned a single scope that does not change, we're actively exploring ways to enhance flexibility and support scenarios where a package could benefit from being installed both cluster-wide and within specific namespaces.

We value community input and encourage you to join the [ongoing discussion](https://github.com/glasskube/glasskube/discussions/1220) about the future of package scopes. Your insights and experiences will help us shape the evolution of this feature to better serve the diverse needs of Glasskube users.

> Note: The specific mechanisms and implementation details of any future changes to package scopes are still under consideration.