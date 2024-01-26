# Testing Strategy Infrastructure-Manager

## Introduction
This testing strategy describes how the Framefrog team is testing the product `Kyma Infrastructure Manager`. It outlines the approach and methodologies to be used for testing all layers of this product to ensure the stability, reliability, and correctness.


## Testing Methodology

We investigate the product by separating it into layers:

1. Business Logic

    Includes the technical frameworks (e.g. Kubebuilder) and custom Golang code.

2. Business Features

    Combines the business logic into feature which is consumed by our customers.

3. Product Integration

    Verifies how our product is integarted into the technical landscape, how it interacts with 3rd party systems and how it is accessible by customers or remote systems.
 
For each layer is a dedicated testing approach used:

1. **Unit Testing for Business Logic:** Writing and executing tests for individual functions, methods, and components to verify their behavior and correctness in isolation.
2. **Integration Testing for Business Features:** Validating the integration and interaction between different components, modules, and services in the project.
3. **End-to-End Testing:** Testing the application as a whole in a production-like environment, mimicking real-world scenarios to ensure the entire system functions correctly, is performing well and secure.


## Testing Approach

### Unit Testing
1. Identify critical functions, methods, and components that require testing.
2. Write unit tests using GoUnit tests, Ginkgo and Gomega frameworks.
3. Ensure tests cover various scenarios, edge cases, and possible failure scenarios. We try to verify business relevant logic with at least 65% code coverage.
4. Test for both positive and negative inputs to validate the expected behavior.
5. Mock external dependencies and use stubs or fakes to isolate the unit under test.
6. Run unit tests periodically during development and before each commit to prevent regressions.

### Integration Testing
1. The PO together with the team create an registry of implemented business features and define for each feature a suitable test scenario.
2. Create a separate test suite for integration testing.
3. Each test scenario is implemented in a separte test case. Use the Kubebuilder Test Framework and others to create test cases that interact with the Kubernetes cluster.  
4. Test the interaction and integration of your custom resources, controllers, and other components with the Kubernetes API.
5. Ensure test cases cover various aspects such as resource creation, updating, deletion, and handling of edge cases.
6. Validate the correctness of event handling, reconciliation, and other control logic.

### End-to-End Testing
1. Use Helm to create, deploy, and manage test clusters and environments that closely resemble the productive execution context.
2. Run regularly, but at least once per release, a performance test that measures product KPIs to indicate KPI violations or performance differences between release candidates.

### Testing Tools and Frameworks
Use the following tools and frameworks to implement the above-mentioned testing levels:

- **GoTest**: For unit testing of Golang code.
- **Kubebuilder Test Framework and EnvTest**: For creating and executing integration tests that interact with the Kubernetes API.
- **Ginkgo and Gomega**: For writing and executing unit tests with a BDD-style syntax and assertions.
- **Kubebuilder Test Framework and EnvTest**: For creating and executing integration tests that interact with the Kubernetes API.
- **K3d**: For creating short-living and lightweight Kubernetes clustes running within a Docker context.
- **Helm**: For deploying and managing test clusters and environments for end-to-end testing.
- **K6:**: For performance and stress testing

|Framework|Unit Testing|Integration Testing|End-to-End Testing|
|--|--|--|--|
|GoTest| X |  |  |
|Kubebuilder Test Framework| X | X | |
|EnvTest| X | X |  |
|Ginkgo| X | X |  |
|Gomega| X | X |  |
|K3d|  |  | X |
|Helm|  |  | X |
|K6|  |  | X |


## Test Automation

The following CI/CD jobs are a part of the development cycle and executing quality assurance related steps:

> **NOTE:** Jobs marked with `pull_request` are triggered with each pull request. Jobs marked with `push` are executed after the merge.

- `golangci-lint / lint (pull_request/push)` - Is responsible for linting and static code analysis. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.golangci.yaml) and [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/golangci-lint.yaml).
- `PR Markdown Link Check / markdown-link-check (pull_request)` - Checks if there are no broken links in the pull request `.md` files. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/mlc.config.json).
- `Run unit tests / validate (pull_request/push)` - Executes basic create/update/delete functional tests of the reconciliation logic. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-tests.yaml).
- `Run vuln check / test (pull_request/push)` - Runs [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) on the code to detect known vulnerabilities. It's configured [here](https://github.com/kyma-project/infrastructure-manager/blob/main/.github/workflows/run-vuln-check.yaml).
- `pull-infrastructure-manager-build` - Triggered with each pull request. It builds the Docker image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).	
- `main-infrastructure-manager-build` - Triggered after the merge. Rebuilds the image and pushes it to the registry. It's configured [here](https://github.com/kyma-project/test-infra/blob/a3c2a07da4ba42e468f69cf42f1960d7bfcc3fff/prow/jobs/kyma-project/infrastructure-manager/infrastructure-manager.yaml).