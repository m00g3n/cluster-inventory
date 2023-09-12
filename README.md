[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/cluster-inventory)](https://api.reuse.software/info/github.com/kyma-project/cluster-inventory)

> **NOTE:** This is a general template that can be used for a project README.md, example README.md, or any other README.md type in all Kyma repositories in the Kyma organization. Except for the mandatory sections, use only those sections that suit your use case but keep the proposed section order.
>
> Mandatory sections: 
> - `Overview`
> - `Prerequisites`, if there are any requirements regarding hard- or software
> - `Contributing` - do not change this!
> - `Code of Conduct` - do not change this!
> - `Licensing` - do not change this!

# Cluster Inventory

## Overview
<!--- mandatory section --->

This project **will be** responsible for managing [Kyma](https://kyma-project.io/#/) clusters infrastructure. Buil using [kubebuilder framework](https://github.com/kubernetes-sigs/kubebuilder)
It's main responsibilities **will be**:
- Provisioning and deprovisioning Kyma clusters
- Generating dynamic kubeconfigs


## Prerequisites

- Access to a k8s cluster.
- [k3d](https://k3d.io) to get a local cluster for testing, or run against a remote cluster.
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kubebuilder](https://book.kubebuilder.io/)

## Installation

1. Clone the project.

```bash
git clone https://github.com/kyma-project/cluster-inventory.git && cd cluster-inventory/
```

2. Set the `cluster-inventory` image name.

```bash
export IMG=custom-cluster-inventory:0.0.1
export K3D_CLUSTER_NAME=cluster-inventory-demo
```

3. Build the project.

```bash
make build
```

4. Build the image.

```bash
make docker-build
```

5. Push the image to the registry.

<div tabs name="Push image" group="cluster-inventory-installation">
  <details>
  <summary label="k3d">
  k3d
  </summary>

   ```bash
   k3d image import $IMG -c $K3D_CLUSTER_NAME
   ```
  </details>
  <details>
  <summary label="Docker registry">
  Globally available Docker registry
  </summary>

   ```bash
   make docker-push
   ```

  </details>
</div>

6. Deploy.

```bash
make deploy
```

## Usage
TODO:
> Explain how to use the project. You can create multiple subsections (H3). Include the instructions or provide links to the related documentation.

## Development

> Add instructions on how to develop the project or example. It must be clear what to do and, for example, how to trigger the tests so that other contributors know how to make their pull requests acceptable. Include the instructions or provide links to related documentation.

## Troubleshooting

> List potential issues and provide tips on how to avoid or solve them. To structure the content, use the following sections:
>
> - **Symptom**
> - **Cause**
> - **Remedy**

## Contributing
<!--- mandatory section - do not change this! --->

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Code of Conduct
<!--- mandatory section - do not change this! --->

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## Licensing
<!--- mandatory section - do not change this! --->

See the [LICENSE file](./LICENSE)
