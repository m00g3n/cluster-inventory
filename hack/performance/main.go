package main

import (
	"fmt"
	"github.com/kyma-project/infrastructure-manager/hack/performance/cmd"
)

func main() {
	operation, worker := cmd.Execute()

	switch operation {
	case cmd.Create:
		fmt.Println("Operation to be performed: create")
		worker.Create()
	case cmd.Delete:
		fmt.Println("Operation to be performed: delete")
		worker.Delete()
	default:
		fmt.Println("Operation to be performed: unknown")
	}
}
