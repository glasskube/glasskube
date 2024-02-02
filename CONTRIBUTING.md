# Contributing guide for Glasskube

Welcome, and thank you for deciding to invest some of your time in contributing to the Glasskube project!
The goal of this document is to define some guidelines to streamline our contribution workflow.

## Getting started

### Where should I start?
- If you are new to the project, please check out the [good first issue](https://github.com/glasskube/glasskube/labels/good%20first%20issue) label.
- If you are looking for something to work on, check out our [open issues](hhttps://github.com/glasskube/glasskube/issues).
- If you have an idea for a new feature, please open an issue, and we can discuss it.
- We are also happy to help you find something to work on. Just reach out to us.

### Getting in touch with the community
- Join our [Glasskube Discord Channel](https://discord.gg/SxH6KUCGH7)
- Introduce yourself on the intros channel or open an issue to let us know that you are interested in contributing

### Discuss issues
- Before you start working on something, propose and discuss your solution on the issue
- If you are unsure about something, ask the community

### How do I contribute?
- Fork the repository and clone it locally
- Create a new branch and follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) guidelines for work undertaken
- Assign yourself to the issue, if you are working on it (if you are not a member of the organization, please leave a comment on the issue)
- Make your changes
- Keep pull requests small and focused, if you have multiple changes, please open multiple PRs
- Create a pull request back to the upstream repository and follow follow the [pull request template](.github/pull_request_template.md) guidelines.
- Wait for a review and address any comments

### Opening PRs
- As long as you are working on your PR, please mark it as a draft
- Please make sure that your PR is up-to-date with the latest changes in `main`
- Fill out the PR template
- Mention the issue that your PR is addressing (closes: #<id>)
- Make sure that your PR passes all checks

### Reviewing PRs
- Be respectful and constructive
- Assign yourself to the PR
- Check if all checks are passing
- Suggest changes instead of simply commenting on found issues
- If you are unsure about something, ask the author
- If you are not sure if the changes work, try them out
- Reach out to other reviewers if you are unsure about something
- If you are happy with the changes, approve the PR
- Merge the PR once it has all approvals and the checks are passing

## Commit Message Format

We require all commits in this repository to adhere to the following commit message format.

```
<type>: <description> (#<issue number>)

[optional body]
```

The following `<type>`s are available:

* `fix` (bug fix)
* `feat` (includes new feature)
* `docs` (update to our documentation)
* `build` (update to the build config)
* `perf` (performance improvement)
* `style` (code style change without any other changes)
* `refactor` (code refactoring)
* `chore` (misc. routine tasks; e.g. dependency updates)

This format is based on [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
Please refer to the Conventional Commits specification for more details.

## Development Guide

Glasskube is developed using the [Go](https://golang.org/) programming language. The current version of Go being used is [v1.21](https://go.dev/doc/go1.21). It uses go modules for dependency management.

### Building

Once you've made your changes, you might want to build a binary of the glasskube CLI containing your changes to test them out. This can be done by running the following command at the root of the project:

```
make
```

This will create the `glasskube` and `package-operator` binary in the `dist` folder. You can execute the binary by running the following:

```
dist/glasskube
```

After you make more changes, simply run `make` again to recompile your changes.

### Executing

In order to execute the `glasskube` binary locally, you can do this manually by creating a copy of it to your project directory.

However, there's an easy and preferred way for doing this by creating an `alias` using the following command:

```
alias <alias-name> = /path/to/glasskube/binary
```

This will make sure the `alias-name` is in sync with your glasskube binary. However, this is a temporary alias. If you'd like to create a permanent alias, you can read more about it [here](https://www.freecodecamp.org/news/how-to-create-your-own-command-in-linux/).

**Note:** Don't use `alias-name` as _glasskube_ since the actual glasskube CLI tool installed locally will get in conflict with executable `glasskube` binary.

### Testing

Unit tests for the project can be executed by running:

```
make test
```

This command will run all the unit tests, will try to detect race conditions, and will generate a test coverage report.

### Linting

Before making a PR, we recommend contributors to run a lint check on their code by running:

```
make lint
```

Some linting errors can be automatically fixed by running:

```
make lint-fix
```
