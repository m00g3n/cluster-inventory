package controller

import (
	"context"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Gardener Cluster controller", func() {
	Context("Secret with kubeconfig doesn't exist", func() {
		kymaName := "kymaname1"
		secretName := "secret-name1"
		shootName := "shootName1"
		namespace := "default"

		It("Create secret", func() {
			By("Create GardenerCluster CR")

			gardenerClusterCR := fixGardenerClusterCR(kymaName, namespace, shootName, secretName)
			Expect(k8sClient.Create(context.Background(), &gardenerClusterCR)).To(Succeed())

			By("Wait for secret creation")
			var kubeconfigSecret corev1.Secret
			key := types.NamespacedName{Name: secretName, Namespace: namespace}

			Eventually(func() bool {
				return k8sClient.Get(context.Background(), key, &kubeconfigSecret) == nil
			}, time.Second*30, time.Second*3).Should(BeTrue())

			err := k8sClient.Get(context.Background(), key, &kubeconfigSecret)
			Expect(err).To(BeNil())
			expectedSecret := fixNewSecret(secretName, namespace, kymaName, shootName, "kubeconfig1")
			Expect(kubeconfigSecret.Labels).To(Equal(expectedSecret.Labels))
			Expect(kubeconfigSecret.Data).To(Equal(expectedSecret.Data))
			Expect(kubeconfigSecret.Annotations[lastKubeconfigSyncAnnotation]).To(Not(BeEmpty()))
		})
	})
})

func fixNewSecret(name, namespace, kymaName, shootName, data string) corev1.Secret {
	labels := fixSecretLabels(kymaName, shootName)

	builder := newTestSecret(name, namespace)
	return builder.WithLabels(labels).WithData(data).ToSecret()
}

func (sb *TestSecret) WithLabels(labels map[string]string) *TestSecret {
	sb.secret.Labels = labels

	return sb
}

func (sb *TestSecret) WithData(data string) *TestSecret {
	sb.secret.Data = map[string][]byte{"config": []byte(data)}

	return sb
}

func (sb *TestSecret) ToSecret() corev1.Secret {
	return sb.secret
}

func newTestSecret(name, namespace string) *TestSecret {
	return &TestSecret{
		secret: corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

type TestSecret struct {
	secret corev1.Secret
}

func fixSecretLabels(kymaName, shootName string) map[string]string {
	labels := fixClusterInventoryLabels(kymaName, shootName)
	labels["operator.kyma-project.io/managed-by"] = "infrastructure-manager"
	labels["operator.kyma-project.io/cluster-name"] = kymaName
	return labels
}

func fixGardenerClusterCR(kymaName, namespace, shootName, secretName string) imv1.GardenerCluster {
	return newTestInfrastructureManagerCR(kymaName, namespace, shootName, secretName).
		WithLabels(fixClusterInventoryLabels(kymaName, shootName)).ToCluster()
}

func newTestInfrastructureManagerCR(name, namespace, shootName, secretName string) *TestGardenerClusterCR {
	return &TestGardenerClusterCR{
		gardenerCluster: imv1.GardenerCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: imv1.GardenerClusterSpec{
				Shoot: imv1.Shoot{
					Name: shootName,
				},
				Kubeconfig: imv1.Kubeconfig{
					Secret: imv1.Secret{
						Name:      secretName,
						Namespace: namespace,
						Key:       "config", //nolint:all TODO: fill it up with the actual data
					},
				},
			},
		},
	}
}

func (sb *TestGardenerClusterCR) WithLabels(labels map[string]string) *TestGardenerClusterCR {
	sb.gardenerCluster.Labels = labels

	return sb
}

func (sb *TestGardenerClusterCR) ToCluster() imv1.GardenerCluster {
	return sb.gardenerCluster
}

type TestGardenerClusterCR struct {
	gardenerCluster imv1.GardenerCluster
}

func fixClusterInventoryLabels(kymaName, shootName string) map[string]string {
	labels := map[string]string{}

	labels["kyma-project.io/instance-id"] = "instanceID"
	labels["kyma-project.io/runtime-id"] = "runtimeID"
	labels["kyma-project.io/broker-plan-id"] = "planID"
	labels["kyma-project.io/broker-plan-name"] = "planName"
	labels["kyma-project.io/global-account-id"] = "globalAccountID"
	labels["kyma-project.io/subaccount-id"] = "subAccountID"
	labels["kyma-project.io/shoot-name"] = shootName
	labels["kyma-project.io/region"] = "region"
	labels["operator.kyma-project.io/kyma-name"] = kymaName

	return labels
}
