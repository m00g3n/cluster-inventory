# Testing Strategy Infrastructure-Manager

## Introduction
This testing strategy outlines the approach and methodologies to be used for testing a Kubernetes project built using the Kubebuilder framework. It includes various levels of testing such as unit testing, integration testing, and end-to-end testing to ensure the stability, reliability, and correctness of the project.
## Testing Levels
1. **Unit Testing:** Writing and executing tests for individual functions, methods, and components to verify their behavior and correctness in isolation.
2. **Integration Testing:** Validating the integration and interaction between different components, modules, and services in the project.
3. **End-to-End Testing:** Testing the application as a whole in a production-like environment, mimicking real-world scenarios to ensure the entire system functions correctly.
## Testing Tools and Frameworks
Use the following tools and frameworks to implement the above-mentioned testing levels:
- **Ginkgo and Gomega**: For writing and executing unit tests with a BDD-style syntax and assertions.
- **Kubebuilder Test Framework**: For creating and executing integration tests that interact with the Kubernetes cluster.
- **Helm**: For deploying and managing test clusters and environments for end-to-end testing.
- **Kubernetes clients for Go**: Use the official Kubernetes client libraries to interact with the Kubernetes API and validate the behavior of the project.
## Testing Approach
### Unit Testing
1. Identify critical functions, methods, and components that require testing.
2. Write unit tests using Ginkgo and Gomega frameworks.
3. Ensure tests cover various scenarios, edge cases, and possible failure scenarios.
4. Test for both positive and negative inputs to validate the expected behavior.
5. Mock external dependencies and use stubs or fakes to isolate the unit under test.
6. Run unit tests periodically during development and before each commit to prevent regressions.
### Integration Testing
1. Create a separate test suite for integration testing.
2. Use the Kubebuilder Test Framework to create test cases that interact with the Kubernetes cluster.
3. Test the interaction and integration of your custom resources, controllers, and other components with the Kubernetes API.
4. Ensure test cases cover various aspects such as resource creation, updating, deletion, and handling of edge cases.
5. Validate the correctness of event handling, reconciliation, and other control logic.
### End-to-End Testing
1. Use Helm to create, deploy, and manage test clusters and environments that closely resemble the


The following CI/CD jobs are a part of the development cycle:

> **NOTE:** Jobs marked with `pull_request` are triggered with each pull request. Jobs marked with `push` are executed after the merge.

- `golangci-lint / lint (pull_request/push)` - Is responsible for linting and static code analysis. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.golangci.yaml) and [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/golangci-lint.yaml).
- `PR Markdown Link Check / markdown-link-check (pull_request)` - Checks if there are no broken links in the pull request `.md` files. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/mlc.config.json).
- `Run unit tests / validate (pull_request/push)` - Executes basic create/update/delete functional tests of the reconciliation logic. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-tests.yaml).
- `Run vuln check / test (pull_request/push)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-vuln-check.yaml).
- `pull-infrastructure-manager-build` - Triggered with each pull request. It builds the Docker image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).	
- `main-infrastructure-manager-build` - Triggered after the merge. Rebuilds the image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).