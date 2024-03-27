---
slug: helm
sidebar_label: With Helm
title: Glasskube vs Helm Comparison
description: Discover Glasskube, a fresh Kubernetes package manager challenging Helm's shortcomings. Explore automatic updates, dependency management and much more.
---

# Glasskube vs Helm Comparison

In late 2023, we published a [blog post](https://glasskube.eu/en/r/knowledge/5-helm-shortcomings/) about, what we thought, were the biggest shortcomings of [Helm](https://helm.sh/)—undoubtedly the most popular package manager for Kubernetes.
After receiving *TONS* of feedback from the community we realized that we were not alone in being unsatisfied with the status quo.
That is why we started thinking about and prototyping Glasskube—a *new* vision of what a package manager for Kubernetes could look like.

We take a new approach on multiple features, but want to highlight the 5 most important differences that got either just released, are in development right now or are on our roadmap and will be started soon:

1. **Package updates**<br/>Glasskube provides the possibility to install packages in `@latest` version which leads to automatic updates or install a package in specific version and explore new versions via the `outdated` command and upgrade them with the `upgrade` command.
2. **Dependency management**<br/>Glasskube offers—as you would expect—dependency management out of the box. Multiple packages can require a specific package (e.g. `cert-manager`) being installed and updated in the Kubernetes cluster and its preferred namespace.
3. **Custom Resources Definition (CRD) changes**<br/>Upgrading CRDs will be taken care of by Glasskube to ensure CRs and its operators don't get out-of-sync.
4. **Cloud-Native architecture**<br/>As Helm's architecture is purely client side, it renders the templates and applies them via the Kubernetes API. Although releases are stored in Kubernetes Secrets there is no first party server-side component for helm making it harder to install packages via the GitOps approach.
5. **Kubernetes version upgrade compatibility**<br/>Glasskube tries to make Kubernetes version upgrades as smooth as possible by automatically testing all package (combinations) across multiple Kubernetes versions.

We acknowledge that Helm, with its flexibility and extensibility, has its place in a seasoned DevOps engineer's tool belt, and it's status as one of the most popular methods to deploy applications in Kubernetes is not without merit.
However, Helm's extensive flexibility comes, at least in part, at the cost of the user, which is especially true for junior and novice Kubernetes administrators.
That is why Glasskube is laser focused on delivering a tool for administrators who need to only manage a couple of applications, but who also need to make sure that a multitude of infrastructure components are kept up-to-date and secure throughout multiple Kubernetes version upgrades while also adapting to inevitable breaking changes.

With that being said, **Glasskube is not a full replacement of Helm**, neither do we aspire for it to become one in the future.
Rather, Glasskube is designed to integrate with established workflows and work in synergy with existing tools, such as Helm, Flux, ArgoCD and many more.
