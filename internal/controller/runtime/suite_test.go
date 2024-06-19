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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

import (
	"context"
	gardener_api "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	infrastructuremanagerv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_mocks "github.com/kyma-project/infrastructure-manager/internal/gardener/mocks"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	. "github.com/onsi/ginkgo/v2"        //nolint:revive
	. "github.com/onsi/gomega"           //nolint:revive
	. "github.com/stretchr/testify/mock" //nolint:revive
	"k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg             *rest.Config                                              //nolint:gochecknoglobals
	mockShootClient *gardener_mocks.ShootClient                               //nolint:gochecknoglobals
	k8sClient       client.Client                                             //nolint:gochecknoglobals
	testEnv         *envtest.Environment                                      //nolint:gochecknoglobals
	suiteCtx        context.Context                                           //nolint:gochecknoglobals
	cancelSuiteCtx  context.CancelFunc                                        //nolint:gochecknoglobals
	anyContext      = MatchedBy(func(_ context.Context) bool { return true }) //nolint:gochecknoglobals
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Runtime Controller Suite")
}

var _ = BeforeSuite(func() {
	logger := zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true))
	logf.SetLogger(logger)

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = infrastructuremanagerv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Metrics: metricsserver.Options{
			BindAddress: ":8083",
		},
		Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())

	mockShootClient = &gardener_mocks.ShootClient{}

	runtimeReconciler := NewRuntimeReconciler(mgr, mockShootClient, logger)
	Expect(runtimeReconciler).NotTo(BeNil())
	err = runtimeReconciler.SetupWithManager(mgr)
	Expect(err).To(BeNil())

	//+kubebuilder:scaffold:scheme
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	go func() {
		defer GinkgoRecover()
		suiteCtx, cancelSuiteCtx = context.WithCancel(context.Background())

		err = mgr.Start(suiteCtx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

func clearMockCalls(mock *gardener_mocks.ShootClient) {
	mock.ExpectedCalls = nil
	mock.Calls = nil
}

func setupShootClientMockForProvisioning(shootClientMock *gardener_mocks.ShootClient) {
	runtimeStub := CreateRuntimeStub("test-resource")
	converterConfig := fixConverterConfigForTests()
	converter := gardener_shoot.NewConverter(converterConfig)
	convertedShoot, err := converter.ToShoot(*runtimeStub)
	if err != nil {
		panic(err)
	}

	shootClientMock.On("Create", anyContext, &convertedShoot, Anything).Return(&convertedShoot, nil)

	var shoots []*gardener_api.Shoot = fixGardenerShootsForProvisioning(&convertedShoot)

	for _, shoot := range shoots {

		if shoot != nil {
			shootClientMock.On("Get", anyContext, Anything, Anything).Return(shoot, nil).Once()
			continue
		}

		shootClientMock.On("Get", anyContext, Anything, Anything).Return(nil, errors.NewNotFound(
			schema.GroupResource{Group: "", Resource: "shoots"}, convertedShoot.Name)).Once()
	}
}

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancelSuiteCtx()
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func fixGardenerShootsForProvisioning(shoot *gardener_api.Shoot) []*gardener_api.Shoot {

	var missingShoot *gardener_api.Shoot

	initialisedShoot := shoot.DeepCopy()

	initialisedShoot.Spec.DNS = &gardener_api.DNS{
		Domain: gardener_shoot.PtrTo("test.domain"),
	}

	initialisedShoot.Status = gardener_api.ShootStatus{
		LastOperation: &gardener_api.LastOperation{
			Type:  gardener_api.LastOperationTypeCreate,
			State: gardener_api.LastOperationStatePending,
		},
	}

	processingShoot := initialisedShoot.DeepCopy()

	processingShoot.Status.LastOperation.State = gardener_api.LastOperationStateProcessing

	readyShoot := processingShoot.DeepCopy()

	readyShoot.Status.LastOperation.State = gardener_api.LastOperationStateSucceeded

	//processedShoot := processingShoot.DeepCopy() // will add specific data later

	return []*gardener_api.Shoot{missingShoot, missingShoot, initialisedShoot, processingShoot, readyShoot, readyShoot}
}

func fixConverterConfigForTests() gardener_shoot.ConverterConfig {
	return gardener_shoot.ConverterConfig{
		Kubernetes: gardener_shoot.KubernetesConfig{
			DefaultVersion: "1.29",
		},

		DNS: gardener_shoot.DNSConfig{
			SecretName:   "xxx-secret-dev",
			DomainPrefix: "runtimeprov.dev.kyma.ondemand.com",
			ProviderType: "aws-route53",
		},
		Provider: gardener_shoot.ProviderConfig{
			AWS: gardener_shoot.AWSConfig{
				EnableIMDSv2: true,
			},
		},
	}
}
