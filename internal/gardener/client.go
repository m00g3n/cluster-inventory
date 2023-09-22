package gardener

import (
	"fmt"
	"os"

	gardener_apis "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientFromFile(kubeconfigFilePath string) (*gardener_apis.CoreV1beta1Client, error) {
	rawKubeconfig, err := os.ReadFile(kubeconfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Gardener Kubeconfig from path %s: %s", kubeconfigFilePath, err.Error())
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(rawKubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := gardener_apis.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
