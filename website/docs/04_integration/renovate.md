# Renovate

[Renovate](https://docs.renovatebot.com/) is a popular application that automates dependency updates in various version control systems.
Developed as an open-source project by Mend and contributors, it integrates deeply into your Git hosting provider,
in order to detect possible version updates and create pull/merge requests using the providers first-party APIs.

Starting with version [38.45.0](https://github.com/renovatebot/renovate/releases/tag/38.45.0), Renovate supports updating Glasskube package versions.
To enable Renovate support for Glasskube in your GitOps repository, add the following configuration to your Renovate config file:

```json title="renovate.json"
"glasskube": {
    "fileMatch": [
      // see below
    ]
}
```

Since there is no canonical naming scheme for files containing Glasskube packages in GitOps repositories,
the manager does not have a default "fileMatch" value, so you have to provider your own.
For example, if you used our GitOps template to set up your repository, you can use `"^packages/.*\\.yaml$"` as a matcher.
This will tell Renovate that all files in the `packages` directory (and its subdirectories) that end with `.yaml` may contain glasskube resources.
The full configuration section should look something like this:

```json title="renovate.json"
"glasskube": {
    "fileMatch": [
      "^packages/.*\\.yaml$"
    ]
}
```

Check out our GitOps template to see the Renovate integration in action: [`glasskube/gitops-template`](https://github.com/glasskube/gitops-template).

The Renovate integration for Glasskube supports scenarios with multiple repositories, but all repositories must be present in the repository.
Furthermore, authenticated repositories are not yet supported.

At Glasskube we use Renovate in almost all of our repositories, and we highly recommend that you install it in your GitOps repository.
Check out the [Renovate documentation](https://docs.renovatebot.com/getting-started/installing-onboarding/) for instructions.
