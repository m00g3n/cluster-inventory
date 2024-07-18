package runtime_comparer

import (
	"fmt"
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"os"
	"sigs.k8s.io/yaml"
)

func WriteToPV(instance v1.Runtime, path string) error {
	fileName := path + "/" + instance.Name + ".yaml"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	// Remove managed fields from the object
	instance.ManagedFields = nil
	b, err := yaml.Marshal(instance)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to marshal Runtime object: %w", err)
	}

	_, err = file.Write(b)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
