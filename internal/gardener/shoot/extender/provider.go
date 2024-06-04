package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/aws"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/azure"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/gcp"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler/openstack"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"slices"
)

func NewProviderExtender(enableIMDSv2 bool) Extend {
	return func(runtimeShoot imv1.RuntimeShoot, shoot *gardener.Shoot) error {
		provider := &shoot.Spec.Provider
		provider.Type = runtimeShoot.Provider.Type
		provider.Workers = runtimeShoot.Provider.Workers

		var err error
		provider.InfrastructureConfig, provider.ControlPlaneConfig, err = getConfig(runtimeShoot)
		if err != nil {
			return err
		}

		if runtimeShoot.Provider.Type == "aws" && enableIMDSv2 {
			provider.Workers[0].ProviderConfig, err = getAWSWorkerConfig()
		}

		return err
	}
}

type InfrastructureProviderFunc func(workersCidr string, zones []string) ([]byte, error)
type ControlPlaneProviderFunc func(zones []string) ([]byte, error)

func getConfig(runtimeShoot imv1.RuntimeShoot) (infrastructureConfig *runtime.RawExtension, controlPlaneConfig *runtime.RawExtension, err error) {
	getConfigForProvider := func(runtimeShoot imv1.RuntimeShoot, infrastructureConfigFunc InfrastructureProviderFunc, controlPlaneConfigFunc ControlPlaneProviderFunc) (*runtime.RawExtension, *runtime.RawExtension, error) {
		zones := getZones(runtimeShoot.Provider.Workers)

		infrastructureConfigBytes, err := infrastructureConfigFunc(runtimeShoot.Networking.Nodes, zones)
		if err != nil {
			return nil, nil, err
		}

		controlPlaneConfigBytes, err := controlPlaneConfigFunc(zones)
		if err != nil {
			return nil, nil, err
		}

		return &runtime.RawExtension{Raw: infrastructureConfigBytes}, &runtime.RawExtension{Raw: controlPlaneConfigBytes}, nil
	}

	switch runtimeShoot.Provider.Type {
	case "aws":
		{
			return getConfigForProvider(runtimeShoot, aws.GetInfrastructureConfig, aws.GetControlPlaneConfig)
		}
	case "azure":
		{
			// Azure shoots are all zoned, put probably it not be validated here.
			return getConfigForProvider(runtimeShoot, azure.GetInfrastructureConfig, azure.GetControlPlaneConfig)
		}
	case "gcp":
		{
			return getConfigForProvider(runtimeShoot, gcp.GetInfrastructureConfig, gcp.GetControlPlaneConfig)
		}
	case "openstack":
		{
			return getConfigForProvider(runtimeShoot, openstack.GetInfrastructureConfig, openstack.GetControlPlaneConfig)
		}
	default:
		return nil, nil, errors.New("provider not supported")
	}
}

func getAWSWorkerConfig() (*runtime.RawExtension, error) {
	workerConfigBytes, err := aws.GetWorkerConfig()
	if err != nil {
		return nil, err
	}

	return &runtime.RawExtension{Raw: workerConfigBytes}, nil
}

func getZones(workers []gardener.Worker) []string {
	var zones []string

	for _, worker := range workers {
		zones = append(zones, worker.Zones...)
	}
	slices.Sort(zones)

	return slices.Compact(zones)
}
