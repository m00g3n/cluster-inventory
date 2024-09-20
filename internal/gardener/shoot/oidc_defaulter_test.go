package shoot

import (
	"testing"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestOidcDefaulter(t *testing.T) {
	t.Run("Create shoot from Runtime", func(t *testing.T) {
		// given
		runtime := CreateRuntimeStub("runtime")

		// when
		DefaultOidcConfigurationIfNotPresent(runtime)

		// then
		assert.NotNil(t, runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig)
		assert.Equal(t, runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig, (*runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig)[0])
	})
}

func CreateRuntimeStub(resourceName string) *imv1.Runtime {
	resource := &imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: "default",
		},
		Spec: imv1.RuntimeSpec{
			Shoot: imv1.RuntimeShoot{
				Name: resourceName,
				Kubernetes: imv1.Kubernetes{
					KubeAPIServer: imv1.APIServer{
						OidcConfig: gardener.OIDCConfig{
							ClientID:       ptr.To("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
							GroupsClaim:    ptr.To("groups"),
							IssuerURL:      ptr.To("https://example.com"),
							SigningAlgs:    []string{"RSA256"},
							UsernameClaim:  ptr.To("sub"),
							UsernamePrefix: ptr.To("-"),
						},
					},
				},
			},
		},
	}
	return resource
}
