---
sidebar_position: 1
---

# Welcome

Glasskube is an open-source Kubernetes package manager the simplifies package management for Kuberenetes.

# Architecture

Glasskube uses multiple components, most notably a client and a package operator.

## Components

### Client

The client is an executable written in Go. It accepts user inputs via a UI and CLI.

They client creates packages in the form of Kubernetes Resources via the Kubernetes API.

### Package Operator

The package operator has two controllers:

1. The Package Controller (managing packages)
2. The PackageInfo Controller (syncing package infos from a registry)

### Package Repo

A place where `PackageInfo`s are stored, searched and maintained.

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
