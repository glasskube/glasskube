# Contributing guide for Glasskube 🧊

Welcome, and thank you for deciding to invest some of your time in contributing to the Glasskube project!
The goal of this document is to define some guidelines to streamline our contribution workflow.

## Before you get started ✋
---
There are many types of issues you can take on when contributing to the Glasskube project. We try our best to provide a wide array of open issues that vary in levels of complexity. From beginners to seasoned developers, everyone should be able to find something to work on.

### Let's find the perfect open issue for you!

- If you are new to the project, please check out the [good first issue](https://github.com/glasskube/glasskube/labels/good%20first%20issue) label.
- If you are ready to make a big impact on the project, check out the [current milestone](https://github.com/glasskube/glasskube/milestones) that is being worked on and filter the issues by `"help-wanted"`, these issues are the ones that will make it into the next official release. 
- If you are looking for something specific to work on, check out our [open issues](hhttps://github.com/glasskube/glasskube/issues) and filter against the available tags such as `component: cli`, `component: ui` `component: repo`, `documentation`.
- If you have an idea for a new feature, please open an issue, and we can discuss it.
- We are also happy to help you find something to work on. Just reach out to us.

### Getting in touch with the community

- Join our [Glasskube Discord Server](https://discord.gg/SxH6KUCGH7).
- Introduce yourself in the `🎤|Intros` channel or open an issue to let us know that you are interested in contributing.

### Discuss issues

- If you have a way of approaching an issue that is outside of the scope of the issues description, propose and discuss your solution in the issue itself.
- If you are unsure about something, don't hesitate to ask the community.

## 🚨 Contributing best practices
>  - Please `only work on one` issue at a time.
>  - If you're unable to continue with an assigned task, inform us promptly. 
>  - Ensure to `TEST` your feature contributions locally before requesting reviews. 
>  - Need assistance? Utilize the issue or `help-forum` on [Discord](https://discord.gg/SxH6KUCGH7)
>  - While Generative AI can be useful, minimize its use for `direct team communication`. We value concise, genuine exchanges over scripted messages.


## How to contribute? 🤷
---

Following these steps will ensure that your contributions are well-received, reviewed, and integrated effectively into Komiser's codebase.

### Issue assigning 
1. Assign yourself to the issue, if you are working on it (if you are not a member of the organization, please leave a comment on the issue and we will assign you to it.)

### Fork and Pull Request Flow 🪜

1. Head over to the [Glasskube GitHub repo](https://github.com/glasskube/glasskube) and "fork it" into your own GitHub account.
2. Clone your fork to your local machine, using the following command:
```bash
git clone git@github.com:USERNAME/FORKED-PROJECT.git
```

3. Create a new branch based-off **\`main\`** branch:
```bash
git checkout main
git checkout -b fix/XXX-something
```

4. Implement the changes or additions you intend to contribute. Whether it's **bug fixes**, **new features**, or **enhancements**, this is where you put your coding skills to use.

5. Once your changes are ready, you may then commit and push the changes from your working branch:
```bash
git commit -m "fix(xxxx-name_of_bug): nice commit description"
git push origin feature/add-new-package
```

## Commit Message Format 💬

We require all commits in this repository to adhere to the following commit message format.

```
<type>: <description> (#<issue number>)

[optional body]
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

This format is based on [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
Please refer to the Conventional Commits specification for more details.

## Keeping your Fork Up-to-Date 🆕
If you plan on doing anything more than just a tiny quick fix, you’ll want to **make sure you keep your fork up to date** by tracking the original [“upstream” repo](https://github.com/glasskube/glasskube) that you forked.

Follow the steps given below to do so:

1. Add the 'upstream' repo to list of remotes:
```bash
git remote add upstream https://github.com/glasskube/glasskube.git
```

2. Fetch upstream repo’s branches and latest commits:
```bash
git fetch upstream
```

3. Checkout to the **\`feature/add-new-package\`** branch and merge the upstream:
```bash
git checkout feature/add-new-package
git merge upstream/feature/add-new-package
```

**Now, your local 'feature/add-new-package' branch is up-to-date with everything modified upstream!**

- Now it's time to create a pull request back to the upstream repository and follow the [pull request template](.github/pull_request_template.md) guidelines.
- Wait for a review and address any comments.

## Opening PRs 📩

- As long as you are working on your PR, please mark it as a draft
- Please make sure that your PR is up-to-date with the latest changes in `main`
- Fill out the PR template
- Mention the issue that your PR is addressing (closes: #<id>)
- Make sure that your PR passes all checks
- Keep pull requests small and focused, if you have multiple changes, please open multiple PRs
- Make sure to test your changes
- If you have multiple commits in your PR, that solve the same problem, please squash the commits

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

Glasskube is developed using the [Go](https://golang.org/) programming language. The current version of Go being used is [v1.21](https://go.dev/doc/go1.21). It uses go modules for dependency management.

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

#### package-operator

**tl;dr:** Use [minikube](https://minikube.sigs.k8s.io/docs/) during development and run: `make install webhook cert run`

For development, we provide the following `make` targets:

- `make install` installs the package-operator CRDs in the current cluster.
- `make webhook` installs the package-operator webhook CRDs in the current cluster, including a patch to allow using the package-operator running on the local machine for the validating admission webhook (only works on minikube).
- `make cert` runs the package-operator cert-manager locally to generate a self signed TLS certificate and patch the ValidatingWebhookConfiguration with the CA bundle. The TLS certificate is valid for 1 year, but is saved in a temporary directory, so it is recommended to run this task at least once everytime the machine is restarted.
- `make run` runs the package-operator locally.
- `make docker-build` builds a docker image for the package-operator.
- `make deploy` applies the full package-operator manifest (excluding dependencies) in the current cluster.

The package-operator ships with a ValidatingAdmissionWebhook. While it is not mandatory to use it during development, we do recommend that you do. Just follow these steps:

1. `make install` creates the package-operator CRDs in your cluster.
2. `make webhook` creates an "ExternalName" service in your cluster that points to your host machine. This only works if you use [minikube](https://minikube.sigs.k8s.io/docs/), if you want to use [kind](https://kind.sigs.k8s.io/) instead, take a look at [this issue](https://github.com/kubernetes-sigs/kind/issues/1200).
3. With the webhook configuration in place, you have to generate a TLS certificate locally and patch the webhook configuration with the CA bundle by running `make cert`.
4. Run the operator using your preferred environment, or `make run`.

When changing the manifests, it is recommended to deploy the package-operator in a minikube cluster. To achieve this, you will have to do three things:

1. Point your local docker CLI to the minikube docker daemon:
   `minikube docker-env` for more info
2. Build the docker image:
   `make docker-build`
3. Deploy the operator using the locally built image:
   `make deploy`

#### Web Development

We have a minimal set of dependencies that need to be installed to work on the web UI locally. Install them with `make web`.
This will download and install the [glasskube theme](https://github.com/glasskube/theme),  [Bootstrap](https://getbootstrap.com/) and [htmx](https://htmx.org). 

After this you are ready to go by running the `serve` command: `go run cmd/glasskube/main.go serve`. 

We are aware that the developer experience for the web part could be improved, e.g. by [introducing hot reload](https://github.com/glasskube/glasskube/issues/170).

#### Custom Package Repository

Sometimes it's necessary to develop and test new features and their different edge cases, and the official package repository does not include these cases yet.

In this case, you can host your own repository locally and change the package repository URL in the code.

1. Clone the [packages repository](https://github.com/glasskube/packages).
2. Make your changes locally and host it, e.g. with [caddy](https://caddyserver.com/docs/command-line): `caddy file-server --root . --listen :9684` from the root directory of the packages project.
3. In `internal/repo/interface.go`, change the repository URL to `http://localhost:9684/packages`. 
4. Make sure to restart your applications (operator, CLI, UI), such that the local repository is being used everywhere.

We do not have a command line option yet to change the repository URL, so for now the code change is necessary.

Also note that some of the information in the repository is redundant by design, to reduce the amount of queries against the repo.
For example, the `index.yaml` contains a `latestVersion` for each package, but the `latestVersion` is also defined in each package index file. 
Please make sure to have consistent and valid state in your local repo. 

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
