package internal

import (
	"flag"
	"fmt"
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Config struct {
	GardenerKubeconfigPath string
	KcpKubeconfigPath      string
	GardenerProjectName    string
	OutputPath             string
}

func PrintConfig(cfg Config) {
	log.Println("gardener-kubeconfig-path:", cfg.GardenerKubeconfigPath)
	log.Println("kcp-kubeconfig-path:", cfg.KcpKubeconfigPath)
	log.Println("gardener-project-name:", cfg.GardenerProjectName)
	log.Println("output-path:", cfg.OutputPath)
}

// newConfig - creates new application configuration base on passed flags
func NewConfig() Config {
	result := Config{}
	flag.StringVar(&result.KcpKubeconfigPath, "kcp-kubeconfig-path", "", "Kubeconfig file for KCP cluster")
	flag.StringVar(&result.GardenerKubeconfigPath, "gardener-kubeconfig-path", "/gardener/kubeconfig/kubeconfig", "Kubeconfig file for Gardener cluster")
	flag.StringVar(&result.GardenerProjectName, "gardener-project-name", "gardener-project", "Name of the Gardener project")
	flag.StringVar(&result.OutputPath, "output-path", "", "Path where generated yamls will be saved. Directory has to exist")

	flag.Parse()
	return result
}

func addToScheme(s *runtime.Scheme) error {
	for _, add := range []func(s *runtime.Scheme) error{
		corev1.AddToScheme,
		v1.AddToScheme,
	} {
		if err := add(s); err != nil {
			return fmt.Errorf("unable to add scheme: %s", err)
		}

	}
	return nil
}

type GetClient = func() (client.Client, error)

func (cfg *Config) Client() (client.Client, error) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", cfg.KcpKubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch rest config: %w", err)
	}

	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return nil, err
	}

	return client.New(restCfg, client.Options{
		Scheme: scheme,
		Cache: &client.CacheOptions{
			DisableFor: []client.Object{
				&corev1.Secret{},
			},
		},
	})
}
