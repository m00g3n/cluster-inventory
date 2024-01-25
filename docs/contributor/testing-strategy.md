# Testing Strategy

Each pull request to the repository triggers the following CI/CD jobs:

- `golangci-lint / lint (pull_request)` - Is responsible for linting and static code analysis. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.golangci.yaml) and [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/golangci-lint.yaml).
- `PR Markdown Link Check / markdown-link-check (pull_request)` - Checks if there are no broken links in the pull request `.md` files. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/mlc.config.json).
- `Run unit tests / validate (pull_request)` - Executes basic create/update/delete functional tests of the reconciliation logic. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-tests.yaml).
- `Run vuln check / test (pull_request)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-vuln-check.yaml).
- `pull-infrastructure-manager-build` - builds the Docker image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).

After the pull request is merged, the following CI/CD jobs are executed:

 - `Run unit tests / validate (push)` - Executes basic create/update/delete functional tests of the reconciliation logic.
 - `Run vuln check / test (push)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities.
 - `golangci-lint / lint (push)` - Is responsible for linting and static code analysis.
 - `main-infrastructure-manager-build` - Rebuilds the image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).