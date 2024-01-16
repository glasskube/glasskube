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