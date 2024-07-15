package fsm

import (
	"context"
	"fmt"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateKubeconfig(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Create Gardener Cluster CR state")

	if err := s.instance.ValidateRequiredLabels(); err != nil {
		m.log.Error(err, "Failed to validate required Runtime labels", "name", s.instance.Name)
		s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonConversionError, "False", err.Error())
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
			m.log.Error(err, "GardenerCluster CR read error", "name", runtimeID)
			s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonKubernetesAPIErr, "False", err.Error())
			return updateStatusAndStop()
		}

		m.log.Info("GardenerCluster CR not found, creating a new one", "Name", runtimeID)
		err = m.Create(ctx, makeGardenerClusterForRuntime(s.instance, s.shoot))
		if err != nil {
			m.log.Error(err, "GardenerCluster CR read error", "name", runtimeID)
			s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonKubernetesAPIErr, "False", err.Error())
			return updateStatusAndStop()
		}

		m.log.Info("Gardener Cluster CR created, waiting for readiness", "Name", runtimeID)
		s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonGardenerCRCreated, "Unknown", "Gardener Cluster CR created, waiting for readiness")
		return updateStatusAndRequeueAfter(controlPlaneRequeueDuration)
	}

	if cluster.Status.State != imv1.ReadyState {
		m.log.Info("GardenerCluster CR not ready yet, requeue", "Name", runtimeID, "State", cluster.Status.State)
		return requeueAfter(controlPlaneRequeueDuration)
	}

	m.log.Info("GardenerCluster CR is ready", "Name", runtimeID)

	return ensureStatusConditionIsDoneAndContinue(&s.instance,
		imv1.ConditionTypeRuntimeKubeconfigReady,
		imv1.ConditionReasonGardenerCRReady,
		"Gardener Cluster CR is ready.",
		sFnProcessShoot)
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
				imv1.LabelKymaInstanceId:         runtime.Labels[imv1.LabelKymaInstanceId],
				imv1.LabelKymaRuntimeId:          runtime.Labels[imv1.LabelKymaRuntimeId],
				imv1.LabelKymaBrokerPlanId:       runtime.Labels[imv1.LabelKymaBrokerPlanId],
				imv1.LabelKymaBrokerPlanName:     runtime.Labels[imv1.LabelKymaBrokerPlanName],
				imv1.LabelKymaGlobalAccountId:    runtime.Labels[imv1.LabelKymaGlobalAccountId],
				imv1.LabelKymaGlobalSubaccountId: runtime.Labels[imv1.LabelKymaGlobalSubaccountId], // BTW most likely this value will be missing
				imv1.LabelKymaName:               runtime.Labels[imv1.LabelKymaName],

				// values from Runtime CR fields
				imv1.LabelKymaPlatformRegion: runtime.Spec.Shoot.PlatformRegion,
				imv1.LabelKymaRegion:         runtime.Spec.Shoot.Region,
				imv1.LabelKymaShootName:      shoot.Name,

				// hardcoded values
				imv1.LabelKymaManagedBy: "lifecycle-manager",
				imv1.LabelKymaInternal:  "true",
			},
		},
		Spec: imv1.GardenerClusterSpec{
			Shoot: imv1.Shoot{
				Name: shoot.Name,
			},
			Kubeconfig: imv1.Kubeconfig{
				Secret: imv1.Secret{
					Name:      fmt.Sprintf("kubeconfig-%s ", runtime.Labels[imv1.LabelKymaRuntimeId]),
					Namespace: runtime.Namespace,
					Key:       "config",
				},
			},
		},
	}

	return gardenCluster
}
