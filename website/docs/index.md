---
sidebar_position: 1
---

# Welcome

📦️ Glasskube will help you to **install your favorite Kubernetes packages**.

🤯 Using **traditional package managers** or applying manifests can be **super confusing**.

😍 Using the **Glasskube UI** for reduces complexity and increases transparency.

🧑‍💻 Still providing a **brew inspired CLI** for advanced users.

🏗️ Glasskube **packages are dependency aware**, as you would expect from a package manager.

🤖  Designed as a cloud native application, so you can follow your **DevOps approach**.



> Glasskube is an open-source Kubernetes package manager the simplifies package management for Kubernetes.


## Fast Track ⏱️ {#fast-track}

Install your first package in 5 minutes.


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


For more installation options 

## Architecture 📏 {#architecture}

Glasskube uses multiple components, most notably a client and a package operator.

## Components

### Client

The client is an executable written in Go. It accepts user inputs via a UI and CLI.

They client creates packages in the form of Kubernetes Resources via the Kubernetes API.

### Package Operator

The package operator has two controllers:

1. The Package Controller (managing packages)
2. The PackageInfo Controller (syncing package infos from a repository)

### Packages Repository

A place where `PackageManifest`s are stored, searched and maintained.
Currently only the glasskube packages repository is supported: [`glasskube/packages`](https://github.com/glasskube/packages)

## Commands

### `glasskube`

Will start the UI server and opens a local browser on localhost:80805.

### `glasskube bootstrap`

### `glasskube install`

```

                                          glasskube install <package>



                               1. validate package
                                                 ┌────────────────┐
                                                 │                │     6. pull package info from repo (and keep up to date)
                               ┌───────────────► │    Repo        │ ◄─────────────────────────────────────────────────────────────┐
                               │                 │                │                                                               │
                               │                 └────────────────┘                                                               │
                               │                                                                                       ┌──────────┴──────────────┐
                               │                                                    5. reconcile `PackageInfo`         │                         │
                               │                                               ┌──────────────────────────────────────►│  PackageInfoController  │
                               │                                               │                                       │                         │
                               │                                               │                                       └─────────────────────────┘
┌────────┐                     │                                               │
│        │                     │                                               │
│  UI    │                     │                                               │             3. reconcile `Pacakge`
│        │                     │                                               │            ─────────────────────────►
└────────┘                     │                                               │             4. create `PackageInfo`
    ▲                          │                                               │                if not present
    │  via local server        │                                               │            ◄──────────────────────────
    │  localhost:80805   ┌─────┴────┐  2. create `Package` CR          ┌───────┴──────────┐                             ┌─────────────────────┐
    └────────────────────┤          │ ────────────────────────────────►│                  │  7. reconcile `PacakgeInfo` │                     │
                         │  client  │  10. pull latest `Package` status│  Kubernetes API  │ ─────────────────────────►  │  PackageController  │
    ┌────────────────────┤          │      and finish inststall cmd    │                  │  9. update `Package` status │                     │
    │ cobra cli          └──────────┘ ────────────────────────────────►└──────────────────┘ ◄────────────────────────── └─────────────────────┘
    │                                                                                                                     8. create `Release`
    ▼
┌────────┐
│        │
│  CLI   │
│        │
└────────┘
```

To edit use: https://asciiflow.com/#/

### `glasskube uninstall <package>`

### `glasskube search <name>`

### `glasskube list`

Lists all installed packages

### `glasskube --version`

Prints the client version.

### `glasskube version`

Prints the version of all components.
