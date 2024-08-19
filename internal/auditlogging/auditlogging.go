package auditlogging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pkg/errors"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	auditlogSecretReference = "auditlog-credentials"
	auditlogExtensionType   = "shoot-auditlog-service"
)

type AuditLogging interface {
	Enable(ctx context.Context, shoot *gardener.Shoot) (bool, error)
}

type auditLogConfigurator interface {
	canEnableAuditLogsForShoot(seedName string) bool
	getTenantConfigPath() string
	getPolicyConfigMapName() string
	getSeedObj(ctx context.Context, seedKey types.NamespacedName) (gardener.Seed, error)
	getConfigFromFile() (data map[string]map[string]AuditLogData, err error)
}

type AuditLog struct {
	auditLogConfigurator
}

type auditLogConfig struct {
	tenantConfigPath    string
	policyConfigMapName string
	client              client.Client
}

type AuditLogData struct {
	TenantID   string `json:"tenantID"`
	ServiceURL string `json:"serviceURL"`
	SecretName string `json:"secretName"`
}

type AuditlogExtensionConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Type is the type of auditlog service provider.
	Type string `json:"type"`
	// TenantID is the id of the tenant.
	TenantID string `json:"tenantID"`
	// ServiceURL is the URL of the auditlog service.
	ServiceURL string `json:"serviceURL"`
	// SecretReferenceName is the name of the reference for the secret containing the auditlog service credentials.
	SecretReferenceName string `json:"secretReferenceName"`
}

func NewAuditLogging(auditLogTenantConfigPath, auditLogPolicyConfigMapName string, k8s client.Client) AuditLogging {
	return &AuditLog{
		auditLogConfigurator: newAuditLogConfigurator(auditLogTenantConfigPath, auditLogPolicyConfigMapName, k8s),
	}
}

func newAuditLogConfigurator(auditLogTenantConfigPath, auditLogPolicyConfigMapName string, k8s client.Client) auditLogConfigurator {
	return &auditLogConfig{
		tenantConfigPath:    auditLogTenantConfigPath,
		policyConfigMapName: auditLogPolicyConfigMapName,
		client:              k8s,
	}
}

func (a *auditLogConfig) canEnableAuditLogsForShoot(seedName string) bool {
	return seedName != "" && a.tenantConfigPath != ""
}

func (a *auditLogConfig) getTenantConfigPath() string {
	return a.tenantConfigPath
}

func (a *auditLogConfig) getPolicyConfigMapName() string {
	return a.policyConfigMapName
}

func (a *auditLogConfig) getSeedObj(ctx context.Context, seedKey types.NamespacedName) (gardener.Seed, error) {
	var seed gardener.Seed
	if err := a.client.Get(ctx, seedKey, &seed); err != nil {
		return gardener.Seed{}, err
	}
	return seed, nil
}

func (al *AuditLog) Enable(ctx context.Context, shoot *gardener.Shoot) (bool, error) {
	configureAuditPolicy(shoot, al.getPolicyConfigMapName())
	seedName := getSeedName(*shoot)

	if !al.canEnableAuditLogsForShoot(seedName) {
		return false, errors.New("Cannot enable Audit Log for shoot: " + shoot.Name)
	}

	seedKey := types.NamespacedName{Name: seedName, Namespace: ""}
	seed, err := al.getSeedObj(ctx, seedKey)
	if err != nil {
		return false, errors.Wrap(err, "Cannot get Gardener Seed object")
	}

	auditConfigFromFile, err := al.getConfigFromFile()
	if err != nil {
		return false, errors.Wrap(err, "Cannot get Audit Log config from file")
	}

	annotated, err := enableAuditLogs(shoot, auditConfigFromFile, seed.Spec.Provider.Type)

	if err != nil {
		return false, errors.Wrap(err, "Error during enabling Audit Logs on shoot: "+shoot.Name)
	}

	return annotated, nil
}

func enableAuditLogs(shoot *gardener.Shoot, auditConfigFromFile map[string]map[string]AuditLogData, providerType string) (bool, error) {
	providerConfig := auditConfigFromFile[providerType]
	if providerConfig == nil {
		return false, fmt.Errorf("cannot find config for provider %s", providerType)
	}

	auditID := shoot.Spec.Region
	if auditID == "" {
		return false, errors.New("shoot has no region set")
	}

	tenant, ok := providerConfig[auditID]
	if !ok {
		return false, fmt.Errorf("auditlog config for region %s, provider %s is empty", auditID, providerType)
	}

	changedExt, err := configureExtension(shoot, tenant)
	changedSec := configureSecret(shoot, tenant)

	return changedExt || changedSec, err
}

func configureExtension(shoot *gardener.Shoot, config AuditLogData) (changed bool, err error) {
	changed = false
	const (
		extensionKind    = "AuditlogConfig"
		extensionVersion = "service.auditlog.extensions.gardener.cloud/v1alpha1"
		extensionType    = "standard"
	)

	ext := findExtension(shoot)
	if ext != nil {
		cfg := AuditlogExtensionConfig{}
		err = json.Unmarshal(ext.ProviderConfig.Raw, &cfg)
		if err != nil {
			return
		}

		if cfg.Kind == extensionKind &&
			cfg.Type == extensionType &&
			cfg.TenantID == config.TenantID &&
			cfg.ServiceURL == config.ServiceURL &&
			cfg.SecretReferenceName == auditlogSecretReference {
			return false, nil
		}
	} else {
		shoot.Spec.Extensions = append(shoot.Spec.Extensions, gardener.Extension{
			Type: auditlogExtensionType,
		})
		ext = &shoot.Spec.Extensions[len(shoot.Spec.Extensions)-1]
	}

	changed = true

	cfg := AuditlogExtensionConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       extensionKind,
			APIVersion: extensionVersion,
		},
		Type:                extensionType,
		TenantID:            config.TenantID,
		ServiceURL:          config.ServiceURL,
		SecretReferenceName: auditlogSecretReference,
	}

	ext.ProviderConfig = &runtime.RawExtension{}
	ext.ProviderConfig.Raw, err = json.Marshal(cfg)

	return
}

func findExtension(shoot *gardener.Shoot) *gardener.Extension {
	for i, e := range shoot.Spec.Extensions {
		if e.Type == auditlogExtensionType {
			return &shoot.Spec.Extensions[i]
		}
	}

	return nil
}

func findSecret(shoot *gardener.Shoot) *gardener.NamedResourceReference {
	for i, e := range shoot.Spec.Resources {
		if e.Name == auditlogSecretReference {
			return &shoot.Spec.Resources[i]
		}
	}

	return nil
}

func configureSecret(shoot *gardener.Shoot, config AuditLogData) (changed bool) {
	changed = false

	sec := findSecret(shoot)
	if sec != nil {
		if sec.Name == auditlogSecretReference &&
			sec.ResourceRef.APIVersion == "v1" &&
			sec.ResourceRef.Kind == "Secret" &&
			sec.ResourceRef.Name == config.SecretName {
			return
		}
	} else {
		shoot.Spec.Resources = append(shoot.Spec.Resources, gardener.NamedResourceReference{})
		sec = &shoot.Spec.Resources[len(shoot.Spec.Resources)-1]
	}

	changed = true

	sec.Name = auditlogSecretReference
	sec.ResourceRef.APIVersion = "v1"
	sec.ResourceRef.Kind = "Secret"
	sec.ResourceRef.Name = config.SecretName

	return
}

func (a *auditLogConfig) getConfigFromFile() (data map[string]map[string]AuditLogData, err error) {
	file, err := os.Open(a.tenantConfigPath)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func getSeedName(shoot gardener.Shoot) string {
	if shoot.Spec.SeedName != nil {
		return *shoot.Spec.SeedName
	}
	return ""
}

func configureAuditPolicy(shoot *gardener.Shoot, policyConfigMapName string) {
	if shoot.Spec.Kubernetes.KubeAPIServer == nil {
		shoot.Spec.Kubernetes.KubeAPIServer = &gardener.KubeAPIServerConfig{}
	}

	shoot.Spec.Kubernetes.KubeAPIServer.AuditConfig = newAuditPolicyConfig(policyConfigMapName)
}

func newAuditPolicyConfig(policyConfigMapName string) *gardener.AuditConfig {
	return &gardener.AuditConfig{
		AuditPolicy: &gardener.AuditPolicy{
			ConfigMapRef: &v12.ObjectReference{Name: policyConfigMapName},
		},
	}
}
