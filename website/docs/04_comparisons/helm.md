---
slug: helm
title: Glasskube vs Helm Comparison
sidebar_label: With Helm
---

# Glasskube vs Helm Comparison

In late 2023, we published a [blog post](https://glasskube.eu/en/r/knowledge/5-helm-shortcomings/) about, what we thought, were the biggest shortcomings of Helm—undoubtedly the most popular package manager for Kubernetes.
After receiving TONS of feedback from the community we realized that we were not alone in being unsatisfied with the status quo.
That is why we started thinking about and prototyping Glasskube—a *new* vision of what a package manager for Kubernetes could look like.

We acknowledge that Helm, with it's flexibility and extensibility, has it's place in a seasoned DevOps engineer's tool belt and it's status as one of the most popular methods to deploy applications in Kubernetes is not without merit.
However, Helm's extensive flexibility comes, at least in part, at the cost of the user, which is especially true for junior and novice Kubernetes administrators.
That is why Glasskube is laser focused on delivering a tool for administrators who need to only manage a couple of applications, but who also need to make sure that a multitude of infrastructure components are kept up-to-date and secure throughout multiple Kubernetes version upgrades while also adapting to inevitable breaking changes.

With that being said, **Glasskube is not a full replacement of Helm**, neither do we aspire for it to become one in the future.
Rather, Glasskube is designed to integrate with established workflows and work in synergy with existing tools, such as Helm, Flux, ArgoCD and many more.
