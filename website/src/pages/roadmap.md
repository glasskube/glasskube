---
title: Roadmap
---

# Roadmap

Our next [milestones](https://github.com/glasskube/glasskube/milestones) and previous [releases](https://github.com/glasskube/glasskube/releases) are managed on GitHub and are the single source of truth.

## Pre Releases until v1 {#pre-release}

Our pre releases are technical proof of concepts that aim to inspire technical folks and will give you a way to try our latest development snapshot where we ship features fast.

|        |                                                                    Features                                                                     |      Timeline      |     Status     |
|--------|:-----------------------------------------------------------------------------------------------------------------------------------------------:|:------------------:|:--------------:|
| v0.0.1 | - first working PackageOperator<br/>- first working client (UI & CLI)<br/>- install first packages (cert-manager, …) <br/>- `bootstrap` command | Released on Jan 31 |    Released    |
| v0.0.2 |                                   - `open` command<br/>- real-time-updates<br/>- add more supported packages                                    | Released on Feb 09 |    Released    |
| v0.0.3 |                                  - add package updates and outdated information<br/>- support version pinning                                   | Released on Feb 27 |    Released    |
| v0.1.0 |                                              - support packages with dependencies<br/>- dark mode                                               | Released on Mar 21 |    Released    |
| v0.2.0 |                                            - support package suspension<br/>- package configuration                                             | Started on Mar 22  | In development |
| v0.3.0 |                                                     - handle secrets<br/>- support backups                                                      |                    |    Planned     |

## First Releases from v1 {#first-release}

Our v1.0.0 milestone is something we will combine with a launch on ProductHunt and similar platforms.
It will include all features we think are necessary for `glasskube` to become the best package manager for Kubernetes and the community can publish their packages via Glasskube.

|        |                                    Features                                     | Timeline | Status  |
|--------|:-------------------------------------------------------------------------------:|:--------:|:-------:|
| v1.0.0 | - supporting `App`s via the Glasskube Apps Operator<br/>- Your feature requests |          | Planned |

## Stable Releases {#stable}

Becoming a stable software from a point where interfaces, manifests, and functionality remain unchanged, typically requires significant time, a dedicated community and a big user base.
As we are working towards this goal we will need to iterate fast and things will break. However, already from the beginning we want to provide stable automated upgrade paths so your packages don´t break!

## Packages and Apps schedule

Can't find a package or want your app included in the list? We are always adding new supported packages & apps,
so just join us on [Discord](https://discord.gg/SxH6KUCGH7) or open up a new issue and let us know what is missing!

| Version |                                                                                                                                                                                            Package/ App                                                                                                                                                                                            |
|---------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
| v0.0.3  |         Cert Manager[ `cert-manager/cert-manager` ]( https://github.com/cert-manager/cert-manager )<br/>  Ingress-NGINX Controller[ `kubernetes/ingress-nginx` ]( https://github.com/kubernetes/ingress-nginx ) <br/> Kubernetes Dashboard[`kubernetes/dashboard`](https://github.com/kubernetes/dashboard) <br/> Cyclops[ `cyclops-ui/cyclops` ]( https://github.com/cyclops-ui/cyclops )         |
| v0.1.0  |                                                                                                                                                          Keptn[ `keptn/lifecycle-toolkit` ]( https://github.com/keptn/lifecycle-toolkit )                                                                                                                                                          |
| v0.2.0  | Kube-Prometheus-Stack[ `prometheus-community/kube-prometheus-stack` ]( https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack )<br/> MariaDB Operator[ `mariadb-operator/mariadb-operator` ]( https://github.com/mariadb-operator/mariadb-operator )<br/>CloudNativePG[ `cloudnative-pg/cloudnative-pg` ]( https://github.com/cloudnative-pg/cloudnative-pg )  |
| v0.3.0  |                                                                                                                         Velero[ `vmware-tanzu/velero` ]( https://github.com/vmware-tanzu/velero )<br/> K8sGPT[ `k8sgpt-ai/k8sgpt` ]( https://github.com/k8sgpt-ai/k8sgpt )                                                                                                                         |
| v1.0.0  |                                                                                                                                                     Glasskube Apps Operator[ `glasskube/operator` ]( https://github.com/glasskube/operator/ )                                                                                                                                                      |

### Glasskube Apps Operator

The [Glasskube Apps Operator](https://github.com/glasskube/operator/) supports the following Apps:

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
