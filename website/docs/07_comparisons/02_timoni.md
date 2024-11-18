---
slug: timoni
sidebar_label: With Timoni
title: Glasskube vs Timoni Comparison
description: Timoni installs and updates cloud-native applications as bundles (OCI images) with a great support for configuration, but it requires users to create its own bundles.
---

# Glasskube vs Timoni Comparison

[Timoni](https://timoni.sh/) is created by [Stefan Prodman](https://github.com/stefanprodan) to improve bundling cloud native applications and their lifecycle management.

Timoni is closer related to Helm than Glasskube, but brings some features and new approaches to the table Helm doesn't offer.
Its bundles are OCI images, so the only way to install these bundles is to publish them first to an OCI registry,
which makes additional security features like image (co-)signing possible.
Similar to Glasskube it performs garbage collection of orphan resources after an uninstallation operation.

Timoni is generally more mature than Glasskube. However, it doesn't prioritize making standard package installation and upgrades easy.
As of now, there's no plan for GUI support either. Users have to bundle their own applications and provide their registry for images.

Glasskube is not yet compatible with Timoni or its bundles, but in a similar way how Glasskube supports helm charts, Timoni bundles can be supported.
There is currently an [open discussion](https://github.com/glasskube/glasskube/discussions/242) to provide Timoni support. Let us know what you think!
