# Telemetry

By default, Glasskube collects some client usage data. We use this data to better understand how our users use Glasskube, identify bugs and make Glasskube a better product overall.

## What we collect

### Client (CLI/UI):

- Glasskube version
- Type of Operating System (Linux/Windows/Mac)
- Which Glasskube-clusters the client is interacting with
- Executed commands (without flags' values), whether they succeeded or not, and how long they took to execute
- Browser user agent when using the UI
- Executed operations on the UI (without config values)

### Cluster:

- Glasskube version
- Kubernetes version
- Number of nodes
- Type of cloud provider
- Which Glasskube packages are installed and in which version, and whether the installation is successful or faulty

### Website (glasskube.dev):

- Browser user agent
- Visited Page

## What we **don’t** collect:

- Personal information about you
- Any other information about cluster resources
- Any kind of data about what else is running or stored on your machine or your cluster
- Any sensitive or secret information contained in your cluster
- Any configured values of the Glasskube packages
- Specific data about errors, as they could contain sensitive information about your cluster configuration

## How exactly do we collect the data?

As Glasskube is fully open source, you can find exactly what we track by looking up [occurrences of posthog](https://github.com/search?q=repo%3Aglasskube%2Fglasskube%20posthog&type=code) in our code base.

## Where is the telemetry data stored?

We use [Posthog](https://posthog.com/) to store this kind of information. The data is stored in Posthog’s EU cloud.

## How to opt out?

Use the command `glasskube telemetry status` to check whether telemetry is disabled for this cluster or not. Disable telemetry with `glasskube telemetry disable`. You can always opt in again with glasskube telemetry enable.

At the moment, Glasskube does not keep any client-side configuration, so the information about whether telemetry is enabled, is stored in the cluster.
