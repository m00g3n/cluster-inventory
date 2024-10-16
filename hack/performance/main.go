package main

import (
	"fmt"
	"github.com/kyma-project/infrastructure-manager/hack/performance/cmd"
)

/*

rt-load - generates performance test load

Usage:
  rt-load create --load-id <LOAD-ID> --name-prefix <STRING> --kubeconfig <FILE> --rt-number <RT-NUMBER> [--rt-template] <RT-TEMPLATE> [--run-on-ci] <BOOL>
  rt-load delete --kubeconfig <FILE> <LOAD-ID>

Options:
  --load-id     - the identifier (label) of the created load (wymagany o podania, w CI przekazujemy na sztywno wartosc)
  --name-prefix - the prefix used to generate each runtime name
  --kubeconfig  - the path to the kubeconfig file tam gdzie chcemy stworzyc CRy
  --rt-number   - the number of the runtimes to be created
  --rt-template - sciezka do pliku
  --run-on-ci   - wiadomo o co chodzi (opcjonalne domylsnie false)
*/
import (
	"log/slog"
)

func main() {

	log := slog.Default()
	operation, worker, err := cmd.Execute()

	if err != nil {
		log.Error("error during initialization", err)
		return
	}

	switch operation {
	case cmd.Create:
		fmt.Println("create")
		err = worker.Create()
		if err != nil {
			log.Error("error in creation", err)
		}
	case cmd.Delete:
		fmt.Println("delete")
		worker.Delete()
	default:
		fmt.Println("unknown")
	}
}
