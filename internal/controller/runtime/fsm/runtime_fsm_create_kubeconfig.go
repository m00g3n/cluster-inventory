package fsm

import (
	"context"
	"fmt"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	labelKymaInstanceId         = "kyma-project.io/instance-id"
	labelKymaRuntimeId          = "kyma-project.io/runtime-id"
	labelKymaShootName          = "kyma-project.io/shootName"
	labelKymaRegion             = "kyma-project.io/region"
	labelKymaName               = "kyma-project.io/kyma-name"
	labelKymaBrokerPlanId       = "kyma-project.io/broker-plan-id"
	labelKymaBrokerPlanName     = "kyma-project.io/broker-plan-name"
	labelKymaGlobalAccountId    = "kyma-project.io/global-account-id"
	labelKymaGlobalSubaccountId = "kyma-project.io/subaccount-id"
	labelKymaManagedBy          = "operator.kyma-project.io/managed-by"
	labelKymaInternal           = "operator.kyma-project.io/internal"
	labelKymaPlatformRegion     = "kyma-project.io/platform-region"
)

func markFailedStatus(m *fsm, s *systemState, condType imv1.RuntimeConditionType, condReason imv1.RuntimeConditionReason, message string) {
	m.log.Error(errors.New(message), message)

	s.instance.UpdateStatePending(
		condType,
		condReason,
		"False",
		message,
	)
}

func sFnCreateKubeconfig(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Create Gardener Cluster CR state - the last one")

	if s.instance.Annotations["labelKymaRuntimeId"] == "" {
		markFailedStatus(m, s, imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonConversionError, "Missing Runtime ID in label \"kyma-project.io/runtime-id\"")
		return updateStatusAndStop()
	}

	runtimeID := s.instance.Labels["labelKymaRuntimeId"]

	var cluster imv1.GardenerCluster
	err := m.Get(ctx, types.NamespacedName{
		Namespace: s.instance.Namespace,
		Name:      runtimeID,
	}, &cluster)

	if err != nil {
		if !k8serrors.IsNotFound(err) {
			markFailedStatus(m, s, imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonKubernetesAPIErr, "GardenerCluster CR read error")
			return updateStatusAndStop()
		}

		m.log.Info("GardenerCluster CR not found, creating a new one", "Name", runtimeID)
		err = m.Create(ctx, makeGardenerClusterForRuntime(s.instance, s.shoot))
		if err != nil {
			markFailedStatus(m, s, imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonKubernetesAPIErr, "GardenerCluster CR create error")
			return updateStatusAndStop()
		}

		m.log.Info("Gardener Cluster CR created successfully")
		s.instance.UpdateStateReady(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonConfigurationCompleted, "Gardener Cluster CR created successfully")
		return updateStatusAndRequeueAfter(controlPlaneRequeueDuration)
	}

	// Gardener cluster exists, check if kubeconfig exists
	kubeConfigExists, err := kubeconfigSecretExists(ctx, m, cluster)
	if err != nil {
		markFailedStatus(m, s, imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonKubernetesAPIErr, err.Error())
		return updateStatusAndStop()
	}

	if !kubeConfigExists {
		// waiting for new kubeconfig to be created
		return requeueAfter(controlPlaneRequeueDuration)
	}
	return switchState(sFnProcessShoot)
}

func kubeconfigSecretExists(ctx context.Context, m *fsm, cluster imv1.GardenerCluster) (bool, error) {
	secret := corev1.Secret{}

	err := m.Get(ctx, types.NamespacedName{
		Namespace: cluster.Spec.Kubeconfig.Secret.Namespace,
		Name:      cluster.Spec.Kubeconfig.Secret.Name,
	}, &secret)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}
		m.log.Error(err, "Error when reading kubeconfig secret for GardenerCluster CR")
		return false, err
	}
	return true, nil
}

func makeGardenerClusterForRuntime(runtime imv1.Runtime, shoot *gardener.Shoot) *imv1.GardenerCluster {
	gardenCluster := &imv1.GardenerCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GardenerCluster",
			APIVersion: imv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      runtime.Spec.Shoot.Name,
			Namespace: runtime.Namespace,
			Annotations: map[string]string{
				"skr-domain": *shoot.Spec.DNS.Domain,
			},
			Labels: map[string]string{
				// values taken from Runtime CR labels
				labelKymaInstanceId:         runtime.Labels[labelKymaInstanceId],
				labelKymaRuntimeId:          runtime.Labels[labelKymaRuntimeId],
				labelKymaBrokerPlanId:       runtime.Labels[labelKymaBrokerPlanId],
				labelKymaBrokerPlanName:     runtime.Labels[labelKymaBrokerPlanName],
				labelKymaGlobalAccountId:    runtime.Labels[labelKymaGlobalAccountId],
				labelKymaGlobalSubaccountId: runtime.Labels[labelKymaGlobalSubaccountId], // most likely this value will be missing
				labelKymaName:               runtime.Labels[labelKymaName],

				// values from Runtime CR fields
				labelKymaPlatformRegion: runtime.Spec.Shoot.PlatformRegion,
				labelKymaRegion:         runtime.Spec.Shoot.Region,
				labelKymaShootName:      shoot.Name,

				// hardcoded values
				labelKymaManagedBy: "lifecycle-manager",
				labelKymaInternal:  "true",
			},
		},
		Spec: imv1.GardenerClusterSpec{
			Shoot: imv1.Shoot{
				Name: shoot.Name,
			},
			Kubeconfig: imv1.Kubeconfig{
				Secret: imv1.Secret{
					Name:      fmt.Sprintf("kubeconfig-%s ", runtime.Labels[labelKymaRuntimeId]),
					Namespace: runtime.Namespace,
					Key:       "config",
				},
			},
		},
	}

	return gardenCluster
}
