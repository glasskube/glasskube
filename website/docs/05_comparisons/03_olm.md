---
slug: olm
sidebar_label: With OLM
title: Glasskube vs OLM Comparison
description: The Operator Lifecycle Manager (OLM) was introduced by Red Hat to manage operators. Glasskube supports all kinds of cloud-native packages and a GUI.
---

# Glasskube vs Operator Lifecycle Manager (OLM) Comparison

The operator framework was [introduced by RedHat in 2018](https://www.redhat.com/en/blog/introducing-operator-framework-building-apps-kubernetes) based on a concept from 2016.

The two most popular projects are the Operator SDK and the Operator Lifecycle Management.

The [Operator SDK](https://github.com/operator-framework/operator-sdk) is a wrapper for Kubernetes SIGS project [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder).
An SDK for building Kubernetes APIs using Custom Resource Definitions (CRDs).
The Operator SDK also supports writing Operators with Go, Ansible, Helm and even JVM based languages with the [`java-operator-sdk`](https://github.com/operator-framework/java-operator-sdk)
which Glasskube uses for its [`apps-operator`](https://github.com/glasskube/operator).

The [Operator Lifecycle Manager](https://github.com/operator-framework/operator-lifecycle-manager) is part of the operator framework and a toolkit to provide automatic updates for Kubernetes operators.
Although it is not directly part of OpenShift, it is loosely coupled as the related GUI is not open-source and only part of OpenShift, a commercial product.
Using the OLM is very complex and consists of multiple CRDs and concepts, but you could achieve something similar to `glasskube install cert-manager` with
`kubectl operator install cert-manager -n cert-manager --channel stable --approval Automatic --create-operator-group`, but still lacking a GUI and the simple bootstrap of `glasskube bootstrap`.

OLMs [operatorhub.io](https://operatorhub.io/) lists currently more packages than Glasskube, but [Glasskube packages](https://glasskube.dev/packages/) will rapidly increase over time.

Another difference between Glasskube and OLM is that **OLM only supports operators while Glasskube supports all kinds of cloud-native applications**.

If you are already an OpenShift customer or only want to install operators, OLM might be a good fit for you,
but you will enjoy using Glasskube if you work with Kubernetes clusters and are looking for a package manager that supports all kind of packages.

