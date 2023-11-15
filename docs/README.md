# Docs

## Overview

This folder contains documents that relate to the project.

## Development

run `make test` to see if all tests are passing. 

## Configuration

It's possible to configure Infrastructure Manager deployment with following arguments:
1. `gardener-kubeconfig-path` - defines path to the gardener project kubeconfig used during API calls
2. `gardener-project` - name of the gardener project where the infrastructure operations are performed
3. `kubeconfig-expiration-time` - maximum time after which kubeconfig is rotated. The rotation will happen sometime between `0.6 * kubeconfig-expiration-time` and `kubeconfig-expiration-time`. 

See [manager_gardener_secret_patch.yaml](../config/default/manager_gardener_secret_patch.yaml) for default values.

## Troubleshooting

> TBD: List potential issues and provide tips on how to avoid or solve them. To structure the content, use the following sections:
>
> - **Symptom**
> - **Cause**
> - **Remedy**
