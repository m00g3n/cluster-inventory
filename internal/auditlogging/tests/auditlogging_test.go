package tests

import (
	"context"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/go-logr/logr"
	"github.com/kyma-project/infrastructure-manager/internal/auditlogging"
	"github.com/kyma-project/infrastructure-manager/internal/auditlogging/mocks"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"testing"
)

func TestEnable(t *testing.T) {
	t.Run("Should successfully enable Audit Log for Shoot", func(t *testing.T) {
		// given
		ctx := context.Background()
		configurator := &mocks.AuditLogConfigurator{}
		shoot := shootForTest()
		shoot.Spec.SeedName = ptr.To("seed-name")
		configFromFile := fileConfigData()
		seedKey := types.NamespacedName{Name: "seed-name", Namespace: ""}

		configurator.On("GetLogInstance").Return(logr.Logger{}).Once()
		configurator.On("CanEnableAuditLogsForShoot", "seed-name").Return(true).Once()
		configurator.On("GetConfigFromFile").Return(configFromFile, nil).Once()
		configurator.On("GetPolicyConfigMapName").Return("policyConfigMapName").Once()
		configurator.On("GetSeedObj", ctx, seedKey).Return(seedForTest(), nil).Once()
		configurator.On("UpdateShoot", ctx, shoot).Return(nil).Once()

		// when
		auditLog := &auditlogging.AuditLog{AuditLogConfigurator: configurator}
		enable, err := auditLog.Enable(ctx, shoot)

		// then
		configurator.AssertExpectations(t)
		require.True(t, enable)
		require.NoError(t, err)
	})

	t.Run("Should return false and error if was problem during enabling Audit Logs", func(t *testing.T) {
		// given
		ctx := context.Background()
		configurator := &mocks.AuditLogConfigurator{}
		shoot := shootForTest()
		shoot.Spec.SeedName = ptr.To("seed-name")
		configFromFile := fileConfigData()
		seedKey := types.NamespacedName{Name: "seed-name", Namespace: ""}
		seed := seedForTest()
		// delete shoot region to simulate error
		shoot.Spec.Region = ""

		configurator.On("GetLogInstance").Return(logr.Logger{}).Once()
		configurator.On("CanEnableAuditLogsForShoot", "seed-name").Return(true).Once()
		configurator.On("GetConfigFromFile").Return(configFromFile, nil).Once()
		configurator.On("GetPolicyConfigMapName").Return("policyConfigMapName").Once()
		configurator.On("GetSeedObj", ctx, seedKey).Return(seed, nil).Once()

		// when
		auditLog := &auditlogging.AuditLog{AuditLogConfigurator: configurator}
		enable, err := auditLog.Enable(ctx, shoot)

		// then
		configurator.AssertExpectations(t)
		require.False(t, enable)
		require.Error(t, err)
	})
}

func shootForTest() *gardener.Shoot {
	return &gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-shoot",
			Namespace: "namespace",
		},
		Spec: gardener.ShootSpec{
			Region:   "region",
			Provider: gardener.Provider{Type: "aws"}},
	}
}

func seedForTest() gardener.Seed {
	return gardener.Seed{
		ObjectMeta: metav1.ObjectMeta{
			Name: "seed-name",
		},
		Spec: gardener.SeedSpec{
			Provider: gardener.SeedProvider{
				Type: "aws",
			},
		},
	}
}

func fileConfigData() map[string]map[string]auditlogging.AuditLogData {
	return map[string]map[string]auditlogging.AuditLogData{
		"aws": {
			"region": {
				TenantID:   "tenantID",
				ServiceURL: "serviceURL",
				SecretName: "secretName",
			},
		},
	}
}
