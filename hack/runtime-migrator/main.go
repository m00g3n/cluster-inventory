package main

import (
	"context"
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardener_types "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	migrator "github.com/kyma-project/infrastructure-manager/hack/runtime-migrator/internal"
	"github.com/kyma-project/infrastructure-manager/internal/gardener"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/kubeconfig"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
	"time"
)

const (
	migratorLabel = "operator.kyma-project.io/created-by-migrator"
)

func main() {
	cfg := migrator.NewConfig()
	migrator.PrintConfig(cfg)

	gardenerNamespace := fmt.Sprintf("garden-%s", cfg.GardenerProjectName)

	gardenerShootClient := setupGardenerShootClient(cfg.GardenerKubeconfigPath, gardenerNamespace)
	list, err := gardenerShootClient.List(context.Background(), metav1.ListOptions{})
	//k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})

	if err != nil {
		log.Fatal(err)
	}

	provider, err := setupKubernetesKubeconfigProvider(cfg.GardenerKubeconfigPath, gardenerNamespace, 60*time.Minute)
	if err != nil {
		log.Fatal("Failed to create kubeconfig provider")
	}

	for _, shoot := range list.Items {
		var subjects = getAdministratorsList(provider, shoot.Name)

		var licenceType = shoot.Annotations["kcp.provisioner.kyma-project.io/licence-type"]
		var nginxIngressEnabled = isNginxIngressEnabled(shoot)
		var hAFailureToleranceType = getFailureToleranceType(shoot)

		var runtime = v1.Runtime{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Runtime",
				APIVersion: "infrastructuremanager.kyma-project.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:                       shoot.Name,
				GenerateName:               shoot.GenerateName,
				Namespace:                  "kcp-system",
				DeletionTimestamp:          shoot.DeletionTimestamp,
				DeletionGracePeriodSeconds: shoot.DeletionGracePeriodSeconds,
				Labels:                     getAllRuntimeLabels(shoot, cfg.Client),
				Annotations:                shoot.Annotations,
				OwnerReferences:            shoot.OwnerReferences,
				Finalizers:                 shoot.Finalizers,
				ManagedFields:              nil, // deliberately left empty "This is mostly for migrator housekeeping, and users typically shouldn't need to set or understand this field."
			},
			Spec: v1.RuntimeSpec{
				Shoot: v1.RuntimeShoot{
					Name:              shoot.Name,
					Purpose:           *shoot.Spec.Purpose,
					Region:            shoot.Spec.Region,
					LicenceType:       &licenceType, //TODO: consult if this is a valid approach
					SecretBindingName: *shoot.Spec.SecretBindingName,
					Kubernetes: v1.Kubernetes{
						Version: &shoot.Spec.Kubernetes.Version,
						KubeAPIServer: v1.APIServer{
							OidcConfig: v1beta1.OIDCConfig{
								CABundle:             nil, //deliberately left empty
								ClientAuthentication: shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.ClientAuthentication,
								ClientID:             shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.ClientID,
								GroupsClaim:          shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.GroupsClaim,
								GroupsPrefix:         shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.GroupsPrefix,
								IssuerURL:            shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.IssuerURL,
								RequiredClaims:       shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.RequiredClaims,
								SigningAlgs:          shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.SigningAlgs,
								UsernameClaim:        shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.UsernameClaim,
								UsernamePrefix:       shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig.UsernamePrefix,
							},
							AdditionalOidcConfig: nil, //deliberately left empty for now
						},
					},
					Provider: v1.Provider{
						Type: shoot.Spec.Provider.Type,
						ControlPlaneConfig: runtime.RawExtension{
							Raw:    shoot.Spec.Provider.ControlPlaneConfig.Raw,
							Object: shoot.Spec.Provider.ControlPlaneConfig.Object,
						},
						InfrastructureConfig: runtime.RawExtension{
							Raw:    shoot.Spec.Provider.InfrastructureConfig.Raw,
							Object: shoot.Spec.Provider.InfrastructureConfig.Object,
						},
						Workers: shoot.Spec.Provider.Workers,
					},
					Networking: v1.Networking{
						Pods:     *shoot.Spec.Networking.Pods,
						Nodes:    *shoot.Spec.Networking.Nodes,
						Services: *shoot.Spec.Networking.Services,
					},
					ControlPlane: v1beta1.ControlPlane{
						HighAvailability: &v1beta1.HighAvailability{
							FailureTolerance: v1beta1.FailureTolerance{
								Type: hAFailureToleranceType, //TODO: verify if needed/present shoot.Spec.ControlPlane.HighAvailability.FailureTolerance.Type
								//TODO: check on prod
							},
						},
					},
				},
				Security: v1.Security{
					Administrators: subjects,
					Networking: v1.NetworkingSecurity{
						Filter: v1.Filter{
							Ingress: &v1.Ingress{
								Enabled: nginxIngressEnabled, //TODO: consult if this is a valid approach
							},
							Egress: v1.Egress{
								Enabled: false, //TODO: fix me
							},
						},
					},
				},
			},
			Status: v1.RuntimeStatus{
				State:      "",  //deliberately left empty by our migrator to show that controller has not picked it yet
				Conditions: nil, //deliberately left nil by our migrator to show that controller has not picked it yet
			},
		}

		shootAsYaml, err := getYamlSpec(runtime)
		writeSpecToFile(cfg.OutputPath, shoot, err, shootAsYaml)
	}
}

func isNginxIngressEnabled(shoot v1beta1.Shoot) bool {
	return shoot.Spec.Addons.NginxIngress != nil && shoot.Spec.Addons.NginxIngress.Enabled
}

func getFailureToleranceType(shoot v1beta1.Shoot) v1beta1.FailureToleranceType {
	if shoot.Spec.ControlPlane != nil {
		if shoot.Spec.ControlPlane.HighAvailability != nil {
			return shoot.Spec.ControlPlane.HighAvailability.FailureTolerance.Type
		}
	}
	return ""
}

func getAdministratorsList(provider kubeconfig.Provider, shootName string) []string {
	var kubeconfig, _ = provider.Fetch(context.Background(), shootName)
	if kubeconfig == "" {
		log.Fatal("failed to get dynamic kubeconfig")
	}

	restClientConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		log.Fatal("failed to create REST client from kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		log.Fatal("failed to create clientset from restconfig")
	}

	var clusterRoleBindings, _ = clientset.RbacV1().ClusterRoleBindings().List(context.Background(), metav1.ListOptions{
		LabelSelector: "reconciler.kyma-project.io/managed-by=reconciler,app=kyma",
	})

	var subjects = []string{}
	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		for _, subject := range clusterRoleBinding.Subjects {
			subjects = append(subjects, subject.Name)
		}
	}

	return subjects
}

func getYamlSpec(shoot v1.Runtime) ([]byte, error) {
	shootAsYaml, err := yaml.Marshal(shoot)
	return shootAsYaml, err
}

func writeSpecToFile(outputPath string, shoot v1beta1.Shoot, err error, shootAsYaml []byte) {
	var fileName = fmt.Sprintf("%sshoot-%s.yaml", outputPath, shoot.Name)

	err = os.WriteFile(fileName, shootAsYaml, 0644)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s created\n", fileName)
}

func setupGardenerShootClient(kubeconfigPath, gardenerNamespace string) gardener_types.ShootInterface {
	//TODO: use the same client factory as in the main code
	restConfig, err := gardener.NewRestConfigFromFile(kubeconfigPath)
	if err != nil {
		return nil
	}

	gardenerClientSet, err := gardener_types.NewForConfig(restConfig)
	if err != nil {
		return nil
	}

	shootClient := gardenerClientSet.Shoots(gardenerNamespace)

	return shootClient
}

func setupKubernetesKubeconfigProvider(kubeconfigPath string, namespace string, expirationTime time.Duration) (kubeconfig.Provider, error) {
	restConfig, err := gardener.NewRestConfigFromFile(kubeconfigPath)
	if err != nil {
		return kubeconfig.Provider{}, err
	}

	gardenerClientSet, err := gardener_types.NewForConfig(restConfig)
	if err != nil {
		return kubeconfig.Provider{}, err
	}

	gardenerClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return kubeconfig.Provider{}, err
	}

	shootClient := gardenerClientSet.Shoots(namespace)
	dynamicKubeconfigAPI := gardenerClient.SubResource("adminkubeconfig")

	err = v1beta1.AddToScheme(gardenerClient.Scheme())
	if err != nil {
		return kubeconfig.Provider{}, errors.Wrap(err, "failed to register Gardener schema")
	}

	return kubeconfig.NewKubeconfigProvider(shootClient,
		dynamicKubeconfigAPI,
		namespace,
		int64(expirationTime.Seconds())), nil
}

func getAllRuntimeLabels(shoot v1beta1.Shoot, getClient migrator.GetClient) map[string]string {
	enrichedRuntimeLabels := map[string]string{}

	// add all labels from the shoot
	for labelKey, labelValue := range shoot.Labels {
		enrichedRuntimeLabels[labelKey] = labelValue
	}

	// add agreed labels from the GardenerCluster CR
	k8sClient, err := getClient()
	if err != nil {
		log.Fatal(err)
	}
	gardenerCluster := v1.GardenerCluster{}
	shootKey := types.NamespacedName{Name: "runtime-id", Namespace: "kcp-system"}
	k8sClient.Get(context.Background(), shootKey, &gardenerCluster)

	enrichedRuntimeLabels["kyma-project.io/broker-plan-id"] = gardenerCluster.Labels["kyma-project.io/broker-plan-id"]
	enrichedRuntimeLabels["kyma-project.io/runtime-id"] = gardenerCluster.Labels["kyma-project.io/runtime-id"]
	enrichedRuntimeLabels["kyma-project.io/subaccount-id"] = gardenerCluster.Labels["kyma-project.io/subaccount-id"]
	enrichedRuntimeLabels["kyma-project.io/broker-plan-name"] = gardenerCluster.Labels["kyma-project.io/broker-plan-name"]
	enrichedRuntimeLabels["kyma-project.io/global-account-id"] = gardenerCluster.Labels["kyma-project.io/global-account-id"]
	enrichedRuntimeLabels["kyma-project.io/instance-id"] = gardenerCluster.Labels["kyma-project.io/instance-id"]
	enrichedRuntimeLabels["kyma-project.io/region"] = gardenerCluster.Labels["kyma-project.io/region"]
	enrichedRuntimeLabels["kyma-project.io/shoot-name"] = gardenerCluster.Labels["kyma-project.io/shoot-name"]
	enrichedRuntimeLabels["operator.kyma-project.io/kyma-name"] = gardenerCluster.Labels["operator.kyma-project.io/kyma-name"]

	// add custom label for the migrator
	enrichedRuntimeLabels[migratorLabel] = "true"

	return enrichedRuntimeLabels
}
