package aws

import (
	"encoding/json"
)

func GetInfrastructureConfig(_ string, _ []string) ([]byte, error) {
	return json.Marshal("{}")
}

func GetControlPlaneConfig() ([]byte, error) {
	return json.Marshal("{}")
}
