package fsm

import (
	"context"
	"fmt"
	"io"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/yaml"
)

func getWriterForFilesystem(filePath string) (io.Writer, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %w", err)
	}
	return file, nil
}

func persist(path string, s interface{}, saveFunc writerGetter) error {
	writer, err := saveFunc(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	b, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to marshal shoot: %w", err)
	}

	if _, err = writer.Write(b); err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}

func sFnDumpShootSpec(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	paths := createFilesPath(m.PVCPath, s.shoot.Namespace, s.shoot.Name)

	shootCp := s.shoot.DeepCopy()
	runtimeCp := s.instance.DeepCopy()
	shootCp.ManagedFields = nil
	runtimeCp.ManagedFields = nil

	if err := persist(paths["shoot"], shootCp, m.writerProvider); err != nil {
		return updateStatusAndStopWithError(err)
	}

	if err := persist(paths["runtime"], runtimeCp, m.writerProvider); err != nil {
		return updateStatusAndStopWithError(err)
	}
	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}

func createFilesPath(pvcPath, namespace, name string) map[string]string {
	m := make(map[string]string)
	m["shoot"] = fmt.Sprintf("%s/%s-%s-shootCR.yaml", pvcPath, namespace, name)
	m["runtime"] = fmt.Sprintf("%s/%s-%s-runtimeCR.yaml", pvcPath, namespace, name)
	return m
}
