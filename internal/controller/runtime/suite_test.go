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
	infrastructuremanagerv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_mocks "github.com/kyma-project/infrastructure-manager/internal/gardener/mocks"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	. "github.com/onsi/ginkgo/v2"        //nolint:revive
	. "github.com/onsi/gomega"           //nolint:revive
	. "github.com/stretchr/testify/mock" //nolint:revive
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg            *rest.Config                                              //nolint:gochecknoglobals
	k8sClient      client.Client                                             //nolint:gochecknoglobals
	testEnv        *envtest.Environment                                      //nolint:gochecknoglobals
	suiteCtx       context.Context                                           //nolint:gochecknoglobals
	cancelSuiteCtx context.CancelFunc                                        //nolint:gochecknoglobals
	anyContext     = MatchedBy(func(_ context.Context) bool { return true }) //nolint:gochecknoglobals
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
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = infrastructuremanagerv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())

	mockShootClient := &gardener_mocks.ShootClient{}
	setupShootClientMock(mockShootClient)
	runtimeReconciler := NewruntimeReconciler(mgr, mockShootClient, logger)
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

func setupShootClientMock(shootClientMock *gardener_mocks.ShootClient) {
	runtimeStub := CreateRuntimeStub("test-resource")
	converterConfig := fixConverterConfigForTests()
	converter := gardener_shoot.NewConverter(converterConfig)
	shoot, err := converter.ToShoot(*runtimeStub)
	if err != nil {
		panic(err)
	}

	shootClientMock.On("Create", anyContext, &shoot, Anything).Return(&shoot, nil)
}

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	cancelSuiteCtx()
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

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
