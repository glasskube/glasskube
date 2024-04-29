---
title: Roadmap
---

# Roadmap

Our next [milestones](https://github.com/glasskube/glasskube/milestones) and previous [releases](https://github.com/glasskube/glasskube/releases) are managed on GitHub and are the single source of truth.

## Pre Releases Milestones until our Beta Release ✅ {#pre-release}

Our pre releases are technical proof of concepts that aim to inspire technical folks and will give you a way to try our latest development snapshot where we ship features fast.

|                |                                                                    Features                                                                     |      Timeline      |  Status  |
|----------------|:-----------------------------------------------------------------------------------------------------------------------------------------------:|:------------------:|:--------:|
| v0.0.1 (Alpha) | - first working PackageOperator<br/>- first working client (UI & CLI)<br/>- install first packages (cert-manager, …) <br/>- `bootstrap` command | Released on Jan 31 | Released |
| v0.0.2         |                                   - `open` command<br/>- real-time-updates<br/>- add more supported packages                                    | Released on Feb 09 | Released |
| v0.0.3         |                                  - add package updates and outdated information<br/>- support version pinning                                   | Released on Feb 27 | Released |
| v0.1.0         |                                              - support packages with dependencies<br/>- dark mode                                               | Released on Mar 21 | Released |
| v0.2.0         |                                                             - package configuration                                                             | Released on Apr 18 | Released |
| v0.3.0 (Beta)  |                                                 - markdown support in package long description                                                  | Released on Apr 25 | Released |

## Next Releases until our stable release (v1) 👨🏻‍💻 {#first-release}

Our v1.0.0 milestone is something we will combine with a very special launch on multiple platforms.

Until then multiple features will be needed. Inclusive, but not limited to:

- GitOps integration including a GitHub bot that is able to provide valuable insights
- Package backups and restore functionality
- Custom private package repositories
- Package comments and ratings
- Database support and multi instance support
- Support for the Glasskube Apps Operator (see below)

This should have contained all features we think are necessary for `glasskube` to become the best package manager for Kubernetes and the community can publish their packages via Glasskube.

### Glasskube Apps Operator

The [Glasskube Apps Operator](https://github.com/glasskube/operator/) supports the following Apps and will be part in one of our releases until we launch:

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

## Stable Releases 🥳 {#stable}

Becoming a stable software from a point where interfaces, manifests, and functionality remain unchanged, typically requires significant time, a dedicated community and a big user base.
As we are working towards this goal we will need to iterate fast and things will break. However, already from the beginning we want to provide stable automated upgrade paths so your packages don´t break!

## Packages and Apps schedule 📦

All planned and already supported packages can be found on our [Packages Directory](/packages/).

Can't find a package or want your app included in the list? We are always adding new supported packages & apps,
so just join us on [Discord](https://discord.gg/SxH6KUCGH7) or open up a new issue and let us know what is missing!



