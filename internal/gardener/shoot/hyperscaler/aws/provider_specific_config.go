package aws

import (
	"encoding/json"
)

func GetInfrastructureConfig(workersCidr string, zones []string) ([]byte, error) {
	return json.Marshal("{}")
}

func GetControlPlaneConfig() ([]byte, error) {
	return json.Marshal("{}")
}
