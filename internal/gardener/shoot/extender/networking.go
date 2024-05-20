package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

func NetworkingExtender(runtimeShoot imv1.RuntimeShoot, shoot *gardener.Shoot) error {
	runtimeNetworking := runtimeShoot.Networking

	shoot.Spec.Networking = &gardener.Networking{
		Nodes:    &runtimeNetworking.Nodes,
		Pods:     &runtimeNetworking.Pods,
		Services: &runtimeNetworking.Services,
	}

	return nil
}
