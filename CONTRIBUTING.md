# Contributing guide for Glasskube 🧊

Welcome, and thank you for deciding to invest some of your time in contributing to the Glasskube project!
The goal of this document is to define some guidelines to streamline our contribution workflow.

## Before you get started ✋

---

There are many types of issues you can take on when contributing to the Glasskube project. We try our best to provide a wide array of open issues that vary in levels of complexity. From beginners to seasoned developers, everyone should be able to find something to work on.

### Let's find the perfect open issue for you!

- If you are new to the project, please check out the [good first issue](https://github.com/glasskube/glasskube/labels/good%20first%20issue) label.
- If you are ready to make a big impact on the project, or are already a seasoned contributor, check out our [unassigned `"help wanted"` issues](https://github.com/glasskube/glasskube/issues?q=is%3Aopen+label%3A%22help+wanted%22+no%3Aassignee+-label%3A%22good+first+issue%22).
- If you are looking for something specific to work on, check out our [open issues](https://github.com/glasskube/glasskube/issues?q=is%3Aissue+is%3Aopen+no%3Aassignee+-label%3Aneeds-triage) and filter against the available tags such as `component: cli`, `component: ui` `component: repo`, `documentation`.
- If you have an idea for a new feature, please open an issue, and we can discuss it.
- We are also happy to help you find something to work on. Just reach out to us.

### Get in touch with the community

- Join our [Glasskube Discord Server](https://discord.gg/SxH6KUCGH7).
- Introduce yourself in the [`🎤|Intros`](https://discord.com/channels/1176558649250951330/1184508688757694634) channel or open an issue to let us know that you are interested in contributing.

### Discuss issues

- If you have a way of approaching an issue that is outside the scope of the issues description, propose and discuss your solution in the issue itself.
- If you are unsure about something, don't hesitate to ask the community.

## 🚨 Contributing best practices

> - Please `only work on one` issue at a time.
> - If you're unable to continue with an assigned task, inform us promptly.
> - Ensure to `TEST` your feature contributions locally before requesting reviews.
> - Need assistance? Utilize the issue or `help-forum` on [Discord](https://discord.gg/SxH6KUCGH7)
> - While Generative AI can be useful, minimize its use for `direct team communication`. We value concise, genuine exchanges over scripted messages.

## How to contribute? 🤷

---

Following these steps will ensure that your contributions are well-received, reviewed, and integrated effectively into Komiser's codebase.

### Issue assigning

1. Assign yourself to the issue, if you are working on it (if you are not a member of the organization, please leave a comment on the issue and we will assign you to it.)

### Fork and Pull Request Flow 🪜

1. Head over to the [Glasskube GitHub repo](https://github.com/glasskube/glasskube) and "fork it" into your own GitHub account.
2. Clone your fork to your local machine, using the following command:

```shell
# replace USERNAME with your GitHub username
git clone git@github.com:USERNAME/glasskube.git
```

3. Please use a feature branch based on **\`main\`** for your changes. This allows easier synchronization with the main repository:

```shell
git switch main
git switch -c your-awesome-new-feature
```

4. Implement the changes or additions you intend to contribute. Whether it's **bug fixes**, **new features**, or **enhancements**, this is where you put your coding skills to use.

5. Once your changes are ready, you may then commit and push the changes from your working branch:

```shell
git commit -m "fix: nice commit description"
git push origin your-awesome-new-feature
```

6. Create a Pull Request following our [pull request template](.github/PULL_REQUEST_TEMPLATE.md) to request a code-review.

## Format for Commit Message and Pull Request Titles 💬

Glasskube uses a workflow based on GitHubs "Squash & Merge" feature.
We therefore require all pull request titles to adher to the syntax specified by [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
We do not restrict the format of your commit messages, however we do encourage using the Conventional Commits syntax as well.

In case you've never heard of Conventional Commits, here's a brief summary:

1. Every message consists of a header and optional body and footer (for PR titles, there is no body or footer).
2. The header consists of a type, an optional scope in parentheses and a description.

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

The following `<type>`s are available:

- `fix` (bug fix)
- `feat` (includes new feature)
- `docs` (update to our documentation)
- `build` (update to the build config)
- `perf` (performance improvement)
- `style` (code style change without any other changes)
- `refactor` (code refactoring)
- `chore` (misc. routine tasks; e.g. dependency updates)

For more details, please refer to the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

## Keeping your Fork Up-to-Date 🆕

Glasskube has an active community of contributors, with new PRs being created and merged almost every day.
This means that the upstream repository might change during the time between you creating your fork and your PR being accepted.
To do this without any special tooling, first, add the upstream repository as a remote, then merge the main branch into your feature branch:

```sh
git remote add upstream git@github.com:glasskube/glasskube.git
git fetch upstream
git merge upstream/main
```

For more information, check out the [official documentation](https://docs.github.com/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork).

**Now, your feature branch is up-to-date with everything modified upstream!**

Please avoid rebasing or force-pushing your branch, because this prevents our code-review team from tracking changes since their last review.

## Opening PRs 📩

- As long as you are working on your PR, please mark it as a draft
- Please make sure that your PR is up-to-date with the latest changes in `main`
- Fill out the PR template
- Mention the issue that your PR is addressing (closes: #<id>)
- Make sure that your PR passes all checks
- Keep pull requests small and focused, if you have multiple changes, please open multiple PRs
- Make sure to test your changes

## Reviewing PRs 🕵️

- Be respectful and constructive
- Assign yourself to the PR
- Check if all checks are passing
- Suggest changes instead of simply commenting on found issues
- If you are unsure about something, ask the author
- If you are not sure if the changes work, try them out
- Reach out to other reviewers if you are unsure about something
- If you are happy with the changes, approve the PR
- Merge the PR once it has all approvals and the checks are passing

## Development Guide 👨‍💻

Glasskube is developed using the [Go](https://golang.org/) programming language. The current version of Go being used is [v1.22](https://go.dev/doc/go1.22). It uses go modules for dependency management.

We use [GNU Make](https://www.gnu.org/software/make/) and do not support other make flavors.

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

#### glasskube

In order to execute the `glasskube` binary locally, you can do this manually by creating a copy of it to your project directory.

However, there's an easy and preferred way for doing this by creating an `alias` using the following command:

```
alias <alias-name> = /path/to/glasskube/binary
```

This will make sure the `alias-name` is in sync with your glasskube binary. However, this is a temporary alias. If you'd like to create a permanent alias, you can read more about it [here](https://www.freecodecamp.org/news/how-to-create-your-own-command-in-linux/).

**Note:** Don't use `alias-name` as _glasskube_ since the actual glasskube CLI tool installed locally will get in conflict with executable `glasskube` binary.

#### dependencies

Install dependencies using the following command:

```
kubectl apply -k config/dependencies
```

#### package-operator

**tl;dr:** Use [minikube](https://minikube.sigs.k8s.io/docs/) during development and run: `make install webhook cert run`

For development, we provide the following `make` targets:

- `make install` installs the package-operator CRDs in the current cluster.
- `make dependencies` installs the package-operator dependencies (Flux source-controller and helm-controller) in the current cluster.
- `make webhook` installs the package-operator webhook CRDs in the current cluster, including a patch to allow using the package-operator running on the local machine for the validating admission webhook (only works on minikube).
- `make cert` runs the package-operator cert-manager locally to generate a self signed TLS certificate and patch the ValidatingWebhookConfiguration with the CA bundle. The TLS certificate is valid for 1 year, but is saved in a temporary directory, so it is recommended to run this task at least once everytime the machine is restarted.
- `make setup` is like `make install dependencies webhook cert` but also creates the default Glasskube package repository in the cluster. This is probably what you want to use to get started quickly.
- `make run` runs the package-operator locally.
- `make docker-build` builds a docker image for the package-operator.
- `make deploy` applies the full package-operator manifest (excluding dependencies) in the current cluster.

The package-operator ships with a ValidatingAdmissionWebhook. While it is not mandatory to use it during development, we do recommend that you do. Just follow these steps:

1. `make install` creates the package-operator CRDs in your cluster.
2. `make webhook` creates an "ExternalName" service in your cluster that points to your host machine. This only works if you use [minikube](https://minikube.sigs.k8s.io/docs/), if you want to use [kind](https://kind.sigs.k8s.io/) instead, take a look at [this issue](https://github.com/kubernetes-sigs/kind/issues/1200).
3. With the webhook configuration in place, you have to generate a TLS certificate locally and patch the webhook configuration with the CA bundle by running `make cert`.
4. Alternatively, you can just run `make setup`, which includes all of the above.
5. Run the operator using your preferred environment, or `make run`.

When changing the manifests, it is recommended to deploy the package-operator in a minikube cluster. To achieve this, you will have to do three things:

1. Point your local docker CLI to the minikube docker daemon:
   `minikube docker-env` for more info
2. Build the docker image:
   `make docker-build`
3. Deploy the operator using the locally built image:
   `make deploy`

#### Web Development

We have a minimal set of dependencies that need to be installed to work on the web UI locally. Install them with `make web`.
This will download and install the [glasskube theme](https://github.com/glasskube/theme), [Bootstrap](https://getbootstrap.com/) and [htmx](https://htmx.org).

After this you are ready to go by running the `serve` command: `go run cmd/glasskube/main.go serve`.

We are aware that the developer experience for the web part could be improved, e.g. by [introducing hot reload](https://github.com/glasskube/glasskube/issues/170).

#### Custom Package Repository

Sometimes it's necessary to develop and test new features and their different edge cases, and the official package repository does not include these cases yet.

In this case, you can host your own repository locally and add it to the list of repositories in your cluster.

1. Clone the [packages repository](https://github.com/glasskube/packages).
2. Make your changes locally and host it, e.g. with [caddy](https://caddyserver.com/docs/command-line): `caddy file-server --root . --listen :9684` from the root directory of the packages project.
3. Run `glasskube repo add local http://localhost:9684` to add this repository with the name `local`.
4. Make sure to restart your applications (operator, CLI, UI), such that the local repository is being used everywhere.

Also note that some of the information in the repository is redundant by design, to reduce the amount of queries against the repo.
For example, the `index.yaml` contains a `latestVersion` for each package, but the `latestVersion` is also defined in each package index file.
Please make sure to have consistent and valid state in your local repo.

Also please be aware of the package repo cache: When changing something in the repo, you might want to restart the applications again (otherwise you might have to wait up to 5 minutes).
There is no option yet to override the cache time, but you could locally change it in `internal/repo/client/clientset.go:NewClientset`.

## Testing

> It's crucial to acknowledge the significance of various types of testing. Alongside conducting unit tests for your contributed code, it's imperative to locally build Glasskube and `test it within a Kubernetes cluster`. ☸️

### Set up a local Minikube cluster for testing locally

In case you don't have access to a remote Kubernetes cluster, set up a local testing environment using Minikube. [This guide](https://minikube.sigs.k8s.io/docs/tutorials/kubernetes_101/module1/) will help you set up a single node cluster in no time, which will be more than enough for you Glasskube testing needs.

### Test locally

Install dependencies and build

```shell
npm ci
make all
```

### Unit tests

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

## Contributor Guidelines Video

[![Contributor Guidelines Video](https://github.com/glasskube/glasskube/assets/38757612/3c5141e8-f541-4064-8a29-fcfe674da5a3)](https://www.youtube.com/watch?v=1V5fBjSU7EI)
