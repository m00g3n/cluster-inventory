# Testing Strategy

Each pull request to the repository triggers the following CI/CD jobs:

- `golangci-lint / lint (pull_request)` - Is responsible for linting and static code analysis.
- `PR Markdown Link Check / markdown-link-check (pull_request)` - Checks if there are no broken links in the pull request `.md` files.
- `Run unit tests / validate (pull_request)` - Executes basic create/update/delete functional tests of the reconciliation logic.
- `Run vuln check / test (pull_request)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities.
- `pull-infrastructure-manager-build` - builds the Docker image and pushes it to the registry.

After the pull request is merged, the following CI/CD jobs are executed:

 - `Run unit tests / validate (push)` - Executes basic create/update/delete functional tests of the reconciliation logic.
 - `Run vuln check / test (push)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities.
 - `golangci-lint / lint (push)` - Is responsible for linting and static code analysis.
 - `main-infrastructure-manager-build` - Rebuilds the image and pushes it to the registry.