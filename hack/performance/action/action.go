package action

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Worker interface {
	Create()
	Delete()
}

type WorkerData struct {
	LoadID     string
	NamePrefix string
	Client     client.Client
	RtNumber   int
}

func NewWorker(loadID, namePrefix string, rtNumber int, client client.Client) Worker {
	return &WorkerData{
		LoadID:     loadID,
		NamePrefix: namePrefix,
		RtNumber:   rtNumber,
		Client:     client,
	}
}

func (w WorkerData) Create() {
	client, err := w.initClient()
	if err != nil {
		return err
	}
	return nil
}

func (w WorkerData) Delete() {
}
