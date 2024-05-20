package shoot

import (
	"encoding/json"
	"fmt"

	gardenerv1beta "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
)

type DNSExtensionProviderConfig struct {
	// APIVersion is gardener extension api version
	APIVersion string `json:"apiVersion"`
	// DnsProviderReplication indicates whether dnsProvider replication is on
	DNSProviderReplication *DNSProviderReplication `json:"dnsProviderReplication,omitempty"`
	// Kind is extension type
	Kind string `json:"kind"`
}

type DNSProviderReplication struct {
	// Enabled indicates whether replication is on
	Enabled bool `json:"enabled"`
}

func newDNSExtensionConfig() *DNSExtensionProviderConfig {
	return &DNSExtensionProviderConfig{
		APIVersion:             "service.dns.extensions.gardener.cloud/v1alpha1",
		DNSProviderReplication: &DNSProviderReplication{Enabled: true},
		Kind:                   "DNSConfig",
	}
}

func newDNSExtender(secretName, domainPrefix, dnsProviderType string) extender {
	return func(runtime imv1.RuntimeShoot, shoot *gardenerv1beta.Shoot) error {
		domain := fmt.Sprintf("%s.%s", runtime.Name, domainPrefix)
		isPrimary := true

		shoot.Spec.DNS = &gardenerv1beta.DNS{
			Domain: &domain,
			Providers: []gardenerv1beta.DNSProvider{
				{
					Domains: &gardenerv1beta.DNSIncludeExclude{
						Include: []string{
							domain,
						},
					},
					Primary:    &isPrimary,
					SecretName: &secretName,
					Type:       &dnsProviderType,
				},
			},
		}

		extensionJSON, err := json.Marshal(newDNSExtensionConfig())
		if err != nil {
			return err
		}

		dnsExtension := gardenerv1beta.Extension{
			Type: "shoot-dns-service",
			ProviderConfig: &apimachineryruntime.RawExtension{
				Raw: extensionJSON,
			},
		}

		shoot.Spec.Extensions = append(shoot.Spec.Extensions, dnsExtension)

		return nil
	}
}
