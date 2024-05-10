package shoot

import (
	gardenerv1beta "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToShoot(runtime imv1.Runtime) gardenerv1beta.Shoot {
	return gardenerv1beta.Shoot{
		ObjectMeta: v1.ObjectMeta{
			Name:        runtime.Spec.Shoot.Name,
			Namespace:   runtime.Namespace,
			Labels:      getLabels(runtime),
			Annotations: getAnnotations(runtime),
		},
		Spec: gardenerv1beta.ShootSpec{},
	}
}

func getLabels(runtime imv1.Runtime) map[string]string {
	return map[string]string{}
}

func getAnnotations(runtime imv1.Runtime) map[string]string {
	return map[string]string{}
}

func getShootSpec(runtimeShoot imv1.RuntimeShoot) gardenerv1beta.ShootSpec {
	return gardenerv1beta.ShootSpec{
		Purpose:           &runtimeShoot.Purpose,
		Region:            runtimeShoot.Region,
		SecretBindingName: &runtimeShoot.SecretBindingName,
		Kubernetes:        getKubernetes(runtimeShoot),
		Networking:        getNetworking(runtimeShoot.Networking),
	}
}

func getKubernetes(runtime imv1.RuntimeShoot) gardenerv1beta.Kubernetes {
	return gardenerv1beta.Kubernetes{
		Version: getKubernetesVersion(runtime),
		KubeAPIServer: &gardenerv1beta.KubeAPIServerConfig{
			OIDCConfig: getOIDCConfig(runtime.Kubernetes.KubeAPIServer.OidcConfig),
		},
	}
}

func getKubernetesVersion(runtime imv1.RuntimeShoot) string {
	if runtime.Kubernetes.Version != nil {
		return *runtime.Kubernetes.Version
	}

	// Determine the default Kubernetes version
	// TODO: it must be read from the configuration (please refer to KEB)
	return ""
}

func getOIDCConfig(oidcConfig gardenerv1beta.OIDCConfig) *gardenerv1beta.OIDCConfig {
	return &gardenerv1beta.OIDCConfig{
		CABundle:       oidcConfig.CABundle,
		ClientID:       oidcConfig.ClientID,
		GroupsClaim:    oidcConfig.GroupsClaim,
		GroupsPrefix:   oidcConfig.GroupsPrefix,
		IssuerURL:      oidcConfig.IssuerURL,
		RequiredClaims: oidcConfig.RequiredClaims,
		SigningAlgs:    oidcConfig.SigningAlgs,
		UsernameClaim:  oidcConfig.UsernameClaim,
		UsernamePrefix: oidcConfig.UsernamePrefix,
	}
}

func getProvider(runtimeProvider imv1.Provider) *gardenerv1beta.Provider {
	return &gardenerv1beta.Provider{
		Type:                 runtimeProvider.Type,
		ControlPlaneConfig:   &runtimeProvider.ControlPlaneConfig,
		InfrastructureConfig: &runtimeProvider.InfrastructureConfig,
	}
}

func getNetworking(runtimeNetworking imv1.Networking) *gardenerv1beta.Networking {
	return &gardenerv1beta.Networking{
		Nodes:    &runtimeNetworking.Nodes,
		Pods:     &runtimeNetworking.Pods,
		Services: &runtimeNetworking.Services,
	}
}

//func getControlPlane(imv1.) gardenerv1beta.ControlPlane {
//  return gardenerv1beta.ControlPlane{
//
//  }
//}
