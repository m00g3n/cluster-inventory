# Testing Strategy

The following CI/CD jobs are a part of the development cycle:

> **NOTE:** Jobs marked with `pull_request` are triggered with each pull request. Jobs marked with `push` are executed after the merge.

- `golangci-lint / lint (pull_request/push)` - Is responsible for linting and static code analysis. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.golangci.yaml) and [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/golangci-lint.yaml).
- `PR Markdown Link Check / markdown-link-check (pull_request)` - Checks if there are no broken links in the pull request `.md` files. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/mlc.config.json).
- `Run unit tests / validate (pull_request/push)` - Executes basic create/update/delete functional tests of the reconciliation logic. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-tests.yaml).
- `Run vuln check / test (pull_request/push)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-vuln-check.yaml).
- `pull-infrastructure-manager-build` - Triggered with each pull request. It builds the Docker image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).	
- `main-infrastructure-manager-build` - Triggered after the merge. Rebuilds the image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).