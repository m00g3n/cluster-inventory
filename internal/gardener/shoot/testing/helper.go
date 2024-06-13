package testing

import (
	"fmt"
	"os"

	gardener_api "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type ShootExample string

var (
	decoderBufferSize = 2048
	examplesPath      = "data"

	ErrInvalidValue = fmt.Errorf("invalid value")
)

func (s *ShootExample) Get() (gardener_api.Shoot, error) {
	var result gardener_api.Shoot
	if s == nil {
		return result, fmt.Errorf("%w: nil", ErrInvalidValue)
	}
	file, err := os.Open(string(*s))
	if err != nil {
		return result, err
	}

	err = yaml.NewYAMLOrJSONDecoder(file, decoderBufferSize).Decode(&result)
	return result, err
}
