package gardener

type KubeconfigProvider struct {
}

func (receiver KubeconfigProvider) Fetch(shootName string) (string, error) {
	return "", nil
}
