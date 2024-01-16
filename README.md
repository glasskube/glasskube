[![GitHub Repo stars](https://img.shields.io/github/stars/glasskube/glasskube)](https://github.com/glasskube/glasskube)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docs](https://img.shields.io/badge/docs-glasskube.dev%2Fdocs-blue)](https://glasskube.dev/docs/)
[![](https://dcbadge.vercel.app/api/server/SxH6KUCGH7?style=flat)](https://discord.gg/SxH6KUCGH7) 
![scarf](https://img.shields.io/static/v1?label=Scarf:%20Downloads&message=8/month&style=flat&color=0572F1&labelColor=374151)
[![twitter](https://img.shields.io/twitter/follow/glasskube?style=social)](https://twitter.com/intent/follow?screen_name=glasskube)

<br>
<div align="center">
  <a href="https://glasskube.dev/">
    <img src="https://raw.githubusercontent.com/glasskube/.github/main/images/glasskube-logo.png" alt="Glasskube Logo" height="160">
  </a>

<h3 align="center">Kubernetes Package Management the easy way </h3>

  <p align="center">
    <a href="https://glasskube.dev/docs/getting-started/install"><strong>Getting started »</strong></a>
    <br> <br>
    <a href="https://glasskube.dev/"><strong>Explore our website »</strong></a>
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
     <a href="https://x.com/glasskube?s=20" target="_blank">Twitter</a>
  </p>
</div>

<hr>

## ⭐️ Why Glasskube?
Using **traditional package managers** or applying manifests can be **super confusing**. Therefore, Glasskube will help you to **install your favorite Kubernetes packages**  using the **Glasskube UI** to reduce complexity and increases transparency. We are also providing a **brew inspired CLI** for advanced users. Our **packages are dependency aware**, as you would expect from a package manager. Designed as a cloud native application, so you can follow your **DevOps approach**.

## ✨ Features
- 💡 **Streamlined UI and CLI Experience**:
<br> We've stripped away unnecessary complexities, providing a simple yet powerful user interface and command-line interface for easy package management.
- 🔄 **Automated Updates**: 
<br> Glasskube ensures your Kubernetes packages and apps are always up-to-date, minimizing the manual effort required for maintenance.
- 🤝 **Dependency Awareness**: 
<br> We understand the interconnected nature of Kubernetes packages. Glasskube intelligently manages dependencies.
- 🛠️ **GitOps Ready** with ArgoCD or Flux: 
<br> Seamlessly integrate Glasskube into your GitOps workflow with support for popular tools like ArgoCD or Flux.
- 📦 **Central Package Register**: 
<br> Keep track of all your packages in one central register, enhancing visibility and control over your Kubernetes environment.
- 🔍 **Cluster Scan** (Version 1.0.0): 
<br> Introducing the Cluster Scan feature in version 1.0.0, which allows you to detect packages in your cluster, providing valuable insights for better management.
- 🔐 **Version Pinning** (Version 1.0.0): 
<br> With version 1.0.0, Glasskube introduces Version Pinning, empowering you to maintain precise control over your package versions for enhanced stability.

## 🗄️ Table Of Contents
- [Quick Start](https://github.com/glasskube/#-quick-start)
- [Supported Tools](https://github.com/glasskube/glasskube#-supported-tools)
- [Screencast](https://github.com/glasskube/glasskube#-screencast)
- [Need help?](https://github.com/glasskube/glasskube#-need-help)
- [Related projects](https://github.com/glasskube/glasskube#-related-projects)
- [How to Contribute](https://github.com/glasskube/glasskube#-how-to-contribute) 
- [Supported by](https://github.com/glasskube/glasskube#-supported-by)

## 🚀 Quick Start - Install your first package in less than 5 minutes.


Install Glasskube via [Homebrew](https://brew.sh/):

```bash
brew tap glasskube/glasskube
brew install glasskube
```

Start the package manager:

```bash
glasskube serve
```

Open [`http://localhost:80805`](http://localhost:80805) and explore available packages.

## 📦 Supported Packages  
- Cert Manager [`cert-manager/cert-manager`](https://github.com/cert-manager/cert-manager)
- Ingress-NGINX Controller [`kubernetes/ingress-nginx`](https://github.com/kubernetes/ingress-nginx) 
- Kubernetes Dashboard [`kubernetes/dashboard`](https://github.com/kubernetes/dashboard)
- Kube-Prometheus-Stack [`kubernetes/dashboard`](https://github.com/kubernetes/dashboard)
- Velero [`vmware-tanzu/velero`](https://github.com/vmware-tanzu/velero)

### Coming Soon
- K8sgpt [`k8sgpt-ai/k8sgpt`](https://github.com/k8sgpt-ai/k8sgpt)
- Keptn [`keptn/lifecycle-toolkit`](https://github.com/keptn/lifecycle-toolkit)
- CCloudNativePG [`cloudnative-pg/cloudnative-pg`](https://github.com/cloudnative-pg/cloudnative-pg])
- MariaDB Operator[`cmariadb-operator/mariadb-operator`](https://github.com/mariadb-operator/mariadb-operator])
- Glasskube Apps Operator [`glasskube/operator`](https://github.com/glasskube/operator/)(with version 1.0.0)
    - Gitea [`go-gitea/gitea`](https://github.com/go-gitea/gitea)
    - GitLab [`gitlab.com/gitlab-org/gitlab`](https://gitlab.com/gitlab-org/gitlab)
    - GlitchTip [`gitlab.com/glitchtip/glitchtip`](https://gitlab.com/glitchtip)
    - Keycloak [`keycloak/keycloak`](https://github.com/keycloak/keycloak)
    - Matomo [`matomo-org/matomo`](https://github.com/matomo-org/matomo)
    - Metabase [`metabase/metabase`](https://github.com/metabase/metabase)
    - Nextcloud [`nextcloud/server`](https://github.com/nextcloud/server)
    - Odoo [`odoo/odoo`](https://github.com/odoo/odoo)
    - Plane [`makeplane/plane`](https://github.com/makeplane/plane)
    - Vault [`hashicorp/vault`](https://github.com/hashicorp/vault)

> Can't find a package or want your app included in the list? We are always adding new supported packages & apps, so just join us on [Discord](https://discord.gg/SxH6KUCGH7) or open up a new issue and let us know what is missing!

## 🎬 How to install you first package

> insert video

## ☝️ Need Help?
If you encounter any problems, we will be happy to support you wherever we can. If you encounter any bugs or issues while working on this project, feel free to contact us on [Discord](https://discord.gg/SxH6KUCGH7). We are happy to assist you with anything related to the project.

## 📎 Related Projects

- Glasskube Apps Operator [`glasskube/operator`](https://github.com/glasskube/operator/)

## 🤝 How to Contribute

See [the contributing guide](CONTRIBUTING.md) for detailed instructions.


## 🤩 Thanks to all our Contributors 

Thanks to everyone, that is supporting this project. We are thankful, for evey contribution, no matter its size! 

<a href="https://github.com/glasskube/glasskube/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=glasskube/glasskube" />
</a>

## 📘 License 

The Glasskube is licensed under the Apache 2.0 license. For more information check the [LICENSE](https://github.com/glasskube/glasskube/blob/main/LICENSE) file for details.