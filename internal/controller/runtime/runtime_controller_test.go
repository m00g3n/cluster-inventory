/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

import (
	"context"
	"time"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("Runtime Controller", func() {
	Context("When reconciling a resource", func() {
		const ResourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      ResourceName,
			Namespace: "default",
		}

		AfterEach(func() {
			resource := &imv1.Runtime{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

		})

		It("Should successfully create new Shoot from provided Runtime and set Ready status on CR", func() {

			By("Setup the fake gardener client for Provisioning")
			setupGardenerTestClientForProvisioning()

			By("Create Runtime CR")
			runtimeStub := CreateRuntimeStub(ResourceName)
			Expect(k8sClient.Create(ctx, runtimeStub)).To(Succeed())

			By("Check if Runtime CR has finalizer")

			Eventually(func() bool {
				runtime := imv1.Runtime{}
				if err := k8sClient.Get(ctx, typeNamespacedName, &runtime); err != nil {
					return false
				}

				return controllerutil.ContainsFinalizer(&runtime, imv1.Finalizer)

			}, time.Second*300, time.Second*3).Should(BeTrue())

			By("Wait for Runtime to process shoot creation process and finish processing in Ready State")

			// should go into Pending Processing state
			Eventually(func() bool {
				runtime := imv1.Runtime{}
				if err := k8sClient.Get(ctx, typeNamespacedName, &runtime); err != nil {
					return false
				}

				// check state
				if runtime.Status.State != imv1.RuntimeStatePending {
					return false
				}

				if !runtime.IsConditionSetWithStatus(imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonShootCreationPending, "Unknown") {
					return false
				}

				return true
			}, time.Second*300, time.Second*3).Should(BeTrue())

			// and end as Ready state with ConfigurationCompleted condition == True
			Eventually(func() bool {
				runtime := imv1.Runtime{}
				if err := k8sClient.Get(ctx, typeNamespacedName, &runtime); err != nil {
					return false
				}

				gardenCluster := imv1.GardenerCluster{}
				if err := k8sClient.Get(ctx, types.NamespacedName{Name: runtime.Labels[imv1.LabelKymaRuntimeID], Namespace: runtime.Namespace}, &gardenCluster); err != nil {
					return false
				}

				if gardenCluster.Status.State != imv1.ReadyState {
					gardenCluster.Status.State = imv1.ReadyState
					k8sClient.Status().Update(ctx, &gardenCluster)
					return false
				}

				// check runtime state
				if runtime.Status.State != imv1.RuntimeStateReady {
					return false
				}

				if !runtime.IsConditionSet(imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConfigurationCompleted) {
					return false
				}

				return true
			}, time.Second*300, time.Second*3).Should(BeTrue())

			Expect(customTracker.IsSequenceFullyUsed()).To(BeTrue())

			By("Wait for Runtime to process shoot update process and finish processing in Ready State")
			setupGardenerTestClientForUpdate()

			runtime := imv1.Runtime{}
			err := k8sClient.Get(ctx, typeNamespacedName, &runtime)

			Expect(err).To(BeNil())

			runtime.Spec.Shoot.Provider.Workers[0].Maximum = 5
			Expect(k8sClient.Update(ctx, &runtime)).To(Succeed())

			// should go into Pending Processing state
			Eventually(func() bool {
				runtime := imv1.Runtime{}
				if err := k8sClient.Get(ctx, typeNamespacedName, &runtime); err != nil {
					return false
				}

				// check state
				if runtime.Status.State != imv1.RuntimeStatePending {
					return false
				}

				if !runtime.IsConditionSetWithStatus(imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonProcessing, "Unknown") {
					return false
				}

				return true
			}, time.Second*300, time.Second*3).Should(BeTrue())

			// and end as Ready state with ConfigurationCompleted condition == True
			Eventually(func() bool {
				runtime := imv1.Runtime{}
				if err := k8sClient.Get(ctx, typeNamespacedName, &runtime); err != nil {
					return false
				}

				// check state
				if runtime.Status.State != imv1.RuntimeStateReady {
					return false
				}

				if !runtime.IsConditionSet(imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConfigurationCompleted) {
					return false
				}

				return true
			}, time.Second*300, time.Second*3).Should(BeTrue())

			Expect(customTracker.IsSequenceFullyUsed()).To(BeTrue())

			// next test will be for runtime deletion
			//
			// By("Delete Runtime CR")
			// Expect(k8sClient.Delete(ctx, &runtime)).To(Succeed())
		})
	})
})

func CreateRuntimeStub(resourceName string) *imv1.Runtime {
	resource := &imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: "default",
			Labels: map[string]string{
				imv1.LabelKymaRuntimeID:       "059dbc39-fd2b-4186-b0e5-8a1bc8ede5b8",
				imv1.LabelKymaInstanceID:      "test-instance",
				imv1.LabelKymaRegion:          "region",
				imv1.LabelKymaName:            "caadafae-1234-1234-1234-123456789abc",
				imv1.LabelKymaBrokerPlanID:    "broker-plan-id",
				imv1.LabelKymaBrokerPlanName:  "aws",
				imv1.LabelKymaGlobalAccountID: "461f6292-8085-41c8-af0c-e185f39b5e18",
				imv1.LabelKymaSubaccountID:    "c5ad84ae-3d1b-4592-bee1-f022661f7b30",
			},
		},
		Spec: imv1.RuntimeSpec{
			Shoot: imv1.RuntimeShoot{
				Name:       resourceName,
				Networking: imv1.Networking{},
				Provider: imv1.Provider{
					Type: "aws",
					Workers: []gardener.Worker{
						{
							Zones:   []string{""},
							Maximum: 1,
						},
					},
				},
			},
			Security: imv1.Security{
				Administrators: []string{},
			},
		},
	}
	return resource
}
