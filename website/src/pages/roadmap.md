---
title: Roadmap
---

# Roadmap

Our next [milestones](https://github.com/glasskube/glasskube/milestones) and previous [releases](https://github.com/glasskube/glasskube/releases) are managed on GitHub and are the single source of truth.

## Pre Releases until v0.1.0 {#pre-release}

Our pre releases are technical proof of concepts that aim to inspire technical folks and will give you a way to try our latest development snapshot where we ship features fast.

| 	        |                                           Features                                                                                 	                                           |     Timeline        	     |   Status     	   |
|----------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------------------:|:----------------:|
| v0.0.1 	 | - initial project setup<br/>- first version of website<br/>- first working PackageOperator<br/>- first working client (UI & CLI)<br/>- install first packages (cert-manager) 	 | started on january 11th 	 | in development 	 |
| v0.0.2 	 |                  - `glasskube bootstrap` command<br/>- handle package dependencies<br/>- increase supported trivial packages                                	                  |             	             |   planned    	   |
| v0.0.3 	 |                             - package configuration<br/>- support packages with dependencies                                                     	                             |             	             |   planned    	   |
| v0.0.4 	 |                                   - support backups<br/>- handle secrets                                                                  	                                    |             	             |   planned    	   |
| v0.0.5 	 |                              - support version pinning<br/>- support glasskube suspension                                                       	                              |             	             |   planned    	   |

## First Releases from v0.1.0 {#first-release}

Our v0.1.0 milestone is something we will combine with a launch on ProductHunt and similar platforms.
It will include all features we think are necessary for `glasskube` to become the best package manager for Kubernetes and the community can publish their packages via Glasskube.

| 	        |                   Features                                    	                   | Timeline 	 |   Status 	   |
|----------|:---------------------------------------------------------------------------------:|:----------:|:------------:|
| v0.1.0 	 | - supporting `App`s via the Glasskube Apps Operator<br/>- Your feature requests 	 |     	      | planned    	 |

## Stable Releases {#stable}

Becoming a stable software from a point where interfaces, manifests and functionality doesn't change anymore usually takes ages, a dedicated community and a big user base.
As we are working towards this goal we will need to iterate fast and things will break.
Although we are trying to provide automated upgrade paths as we don't want your packages to break!
