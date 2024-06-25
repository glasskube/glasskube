[![GitHub Repo stars](https://img.shields.io/github/stars/glasskube/glasskube?style=flat)](https://github.com/glasskube/glasskube)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docs](https://img.shields.io/badge/docs-glasskube.dev%2Fdocs-blue)](https://glasskube.dev/docs/?utm_source=github)
[![PRs](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![](https://dcbadge.vercel.app/api/server/SxH6KUCGH7?style=flat)](https://discord.gg/SxH6KUCGH7)
[![Downloads](https://img.shields.io/github/downloads/glasskube/glasskube/total)](https://github.com/glasskube/glasskube/releases)
[![CNCF Landscape](https://img.shields.io/badge/CNCF%20Landscape-5699C6)](https://landscape.cncf.io/?item=app-definition-and-development--application-definition-image-build--glasskube)
[![Go Reference](https://pkg.go.dev/badge/github.com/glasskube/glasskube)](https://pkg.go.dev/github.com/glasskube/glasskube)
[![Go Report Card](https://goreportcard.com/badge/github.com/glasskube/glasskube)](https://goreportcard.com/report/github.com/glasskube/glasskube)

<br>
<div align="center">
  <a href="https://glasskube.dev?utm_source=github">
    <img src="https://raw.githubusercontent.com/glasskube/.github/main/images/glasskube-logo.png" alt="Glasskube Logo" height="160">
  </a>
  <img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=899d5aee-625c-4345-bad0-713d29caf929" />

<h3 align="center">🧊 The next generation Package Manager for Kubernetes 📦 (Beta Version)</h3>

  <p align="center">
    <a href="https://glasskube.dev/docs/getting-started/install?utm_source=github"><strong>Getting started »</strong></a>
    <br> <br>
    <a href="https://glasskube.dev?utm_source=github"><strong>Explore our website »</strong></a>
    <br>
    <br>
    <a href="https://github.com/glasskube" target="_blank">GitHub</a>
    .
    <a href="https://hub.docker.com/u/glasskube" target="_blank">Docker Hub</a>
    .
    <a href="https://artifacthub.io/packages/search?org=glasskube" target="_blank">Artifact Hub</a>
    .
    <a href="https://www.linkedin.com/company/glasskube/" target="_blank">LinkedIn</a>
    . 
     <a href="https://x.com/intent/follow?screen_name=glasskube" target="_blank">Twitter / X</a>
  </p>
</div>

<hr>

![Glasskube GUI](https://github.com/glasskube/glasskube/assets/3041752/54b20ffe-1daf-4905-abc5-37e99e056b02)


## 📦 What is Glasskube?

Glasskube is an **Open Source package manager for Kubernetes**.
It makes deploying, updating, and configuring packages on Kubernetes **20 times faster** than tools like **Helm or Kustomize**.
Inspired by the simplicity of Homebrew and npm. You can decide if you want to use the Glasskube UI, CLI, or directly deploy packages via GitOps.

## ⭐️ Why Glasskube?

We have been working in the Kubernetes ecosystem for over five years.
During this time, we have consistently struggled with package management, configuration, and distribution.
We've spent countless hours templating and writing documentation for commands and concepts that were difficult to grasp.

In contrast, tools like Homebrew, apt, and dnf felt easy to use and rarely caused problems.
While we worked on other cloud-native projects, our users consistently highlighted several common pain points.
This realization prompted us to tackle the larger issue of package management in Kubernetes, leading to the development of Glasskube.

## 🗄️ Table Of Contents

- [Features](https://github.com/glasskube/#-features)
- [Quick Start](https://github.com/glasskube/#-quick-start)
- [How to install your first package](https://github.com/glasskube/glasskube#-how-to-install-you-first-package)
- [Supported Packages](https://github.com/glasskube/glasskube#-supported-packages)
- [Architecture Diagram](https://github.com/glasskube/glasskube#architecture-diagram)
- [Need help?](https://github.com/glasskube/glasskube#-need-help)
- [Related projects](https://github.com/glasskube/glasskube#-related-projects)
- [How to Contribute](https://github.com/glasskube/glasskube#-how-to-contribute)
- [Supported by](https://github.com/glasskube/glasskube#-thanks-to-all-our-contributors)
- [Activity](https://github.com/glasskube/glasskube#-activity)
- [License](https://github.com/glasskube/glasskube#-license)

## ✨ Features
|                                                                                                                                                                                                                                                                  |                                                                                                              |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| **Focusing on simplicity and reliability with our CLI and UI** <br> Easily install packages in your cluster via the Glasskube UI, where all packages are conveniently located, eliminating the need to search for a Helm repository.                             | ![Glasskube GUI](https://github.com/glasskube/glasskube/assets/3041752/323994d6-6b08-4dca-ac59-d29ae6b37f94) |
| **Package configurations** <br> Configure packages with typesafe input values via the UI or interactive CLI questionnaire. Inject values from other packages, ConfigMaps, and Secrets easily. No more untyped and undocumented `values.yaml` files.              | ![Configuration](https://github.com/glasskube/glasskube/assets/3041752/df6bd7d4-7cac-435b-b3a0-31c3cab6069b) |
| **Dependency Management** <br> Glasskube packages are dependency aware, so they can be used and referenced by multiple other packages. They will also get installed in the correct namespace. This is how umbrella charts should have worked from the beginning. | ![Dependency](https://github.com/glasskube/glasskube/assets/3041752/9588b3fc-2a87-454e-97ff-b0f7558717bc)    |
| **Safe Package Updates** <br> Preview and perform pending updates to your desired version with a single click (or CLI command). All updates are pre-tested by the Glasskube test suite.                                                                          | ![Updates](https://github.com/glasskube/glasskube/assets/3041752/a6e6dc72-9919-4d15-addf-bc709ec76d9d)       |
| **Reactions and comments** <br> Discuss and upvote your favorit Kubernetes package on [GitHub](https://github.com/glasskube/glasskube/discussions/categories/packages) or right inside the Glasskube UI.                                                         | ![Reactions](https://github.com/glasskube/glasskube/assets/3041752/56f08373-fbbe-46fd-820e-fb637114336b)     |
| **GitOps Integration** <br> All Glasskube packages are custom resources, manageable via GitOps. We're also integrating with [renovate](https://github.com/renovatebot/renovate/issues/29322)                                                                     | ![GitOps](https://github.com/glasskube/glasskube/assets/3041752/8c359e61-9eec-4413-9c13-bca5cd8710d1)        |
| **Multiple Repositories and private packages** <br> Use multiple repositories and publish your own private packages. This could be your companies internal services packages, so all developers will have up-to-date and easily configured internal services.    | ![Repo](https://github.com/glasskube/glasskube/assets/130456438/e2f4472b-5b80-4043-9c78-9ccabd8f3337)        |


## 🚀 Quick Start - Install the Beta Version.

You can install Glasskube via [Homebrew](https://brew.sh/):

```bash
brew install glasskube/tap/glasskube
```

For other installation options check out our [installation guide](https://glasskube.dev/docs/getting-started/install).

Once the CLI is installed, the first step is to install the necessary components in your cluster. To do that, run
```sh
glasskube bootstrap
```

After successfully bootstrapping your cluster, you are ready to start the package manager UI:

```bash
glasskube serve
```

This command will open [`http://localhost:8580`](http://localhost:8580) in your default browser.
Congratulations, you can now explore and install all our available packages! 🎉

## 🎬 Glasskube Demo Video


[![Glasskube Demo Video](https://i.ytimg.com/vi/aIeTHGWsG2c/hq720.jpg)](https://www.youtube.com/watch?v=aIeTHGWsG2c)

## 📦 Supported Packages

Glasskube already supports a wide range of packages, including, but not limited to:

- Kubernetes Dashboard [`kubernetes/dashboard`](https://github.com/kubernetes/dashboard)
- Cert Manager [`cert-manager/cert-manager`](https://github.com/cert-manager/cert-manager)
- Ingress-NGINX Controller [`kubernetes/ingress-nginx`](https://github.com/kubernetes/ingress-nginx)
- Kube Prometheus Stack [`prometheus-operator/kube-prometheus`](https://github.com/prometheus-operator/kube-prometheus)
- Cloud Native PG [`cloudnative-pg/cloudnative-pg`](https://github.com/cloudnative-pg/cloudnative-pg)

You can find all supported and planned packages on [glasskube.dev/packages](https://glasskube.dev/packages/).

## Architecture Diagram

```mermaid
---
title: glasskube install [package]
---
flowchart BT
  UI([UI])-- via local server<br>http://localhost:8580 ---Client(Client)
  CLI([CLI])-- cobra cli ---Client
  Client-- 1. validate package -->Repo[(Package Repo)]
  Client-- 2. create<br>`Package` CR -->Kubernetes(((Kubernetes API)))
  subgraph Cluster
    Kubernetes-- 3. reconcile<br>`Package` -->PackageController
    PackageController-- 4. create `PackageInfo`<br>if not present-->Kubernetes
    Kubernetes-- 5. reconcile<br>`PackageInfo`-->PackageInfoController
    end
  PackageInfoController<-- 6. update package manifest -->Repo
  subgraph Cluster
    PackageInfoController-- 7. update manifest<br>in `PackageInfo` -->Kubernetes
    Kubernetes-- 8. reconcile<br>`PackageInfo` -->PackageController
    PackageController-- 9. deploy package -->Kubernetes
  end

  Kubernetes-- 10. package status -->Client 
```

## ☝️ Need Help or Want to Provide Feedback?

If you encounter any problems, we will be happy to support you wherever we can on our [Discord](https://discord.gg/SxH6KUCGH7).
For bugs, issues or feature requests fee free to [open an issue](https://github.com/glasskube/glasskube/issues/new/choose).
We are happy to assist you with anything related to the project.

## 📎 Related Projects

- Glasskube Apps Operator [`glasskube/operator`](https://github.com/glasskube/operator/)

## 🤝 How to Contribute to Glasskube

Your feedback is invaluable to us as we continue to improve Glasskube. If you'd like to contribute, consider trying out the beta version, reporting any issues, and sharing your suggestions. See [the contributing guide](CONTRIBUTING.md) for detailed instructions on how you can contribute.

## 🤩 Thanks to all our Contributors

Thanks to everyone, that is supporting this project. We are thankful, for every contribution, no matter its size!

<a href="https://github.com/glasskube/glasskube/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=glasskube/glasskube" />
</a>

## 👾 Activity

![Glasskube Activity](https://repobeats.axiom.co/api/embed/c5aac6f5d22bd6b83a21ae51353dd7bcb43f9517.svg "Glasskube activity image")

## 📘 License

The Glasskube is licensed under the Apache 2.0 license. For more information check the [LICENSE](https://github.com/glasskube/glasskube/blob/main/LICENSE) file for details.
