package cmd

import (
	"flag"
	"fmt"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/hack/performance/action"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationType int

const (
	Create OperationType = iota
	Delete
	Unknown
)

func Execute() (OperationType, action.Worker) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	loadID := createCmd.String("load-id", "", "the identifier (label) of the created load (required)")
	namePrefix := createCmd.String("name-prefix", "", "the prefix used to generate each runtime name (required)")
	kubeconfig := createCmd.String("kubeconfig", "", "the path to the kubeconfig file (required)")
	rtNumber := createCmd.Int("rt-number", 0, "the number of the runtimes to be created (required)")
	runOnCi := createCmd.String("run-on-ci", "false", "identifies if the load is running on CI")

	loadIDDelete := deleteCmd.String("load-id", "", "the identifier (label) of the created load (required)")
	kubeconfigDelete := deleteCmd.String("kubeconfig", "", "the path to the kubeconfig file (required)")

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or 'delete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *loadID == "" || *namePrefix == "" || *kubeconfig == "" || *rtNumber == 0 {
			fmt.Println("all flags --load-id, --name-prefix, --kubeconfig, and --rt-number are required")
			createCmd.Usage()
			os.Exit(1)
		}
		if *runOnCi == "false" {
			var response string
			fmt.Printf("Do you want to create %d runtimes? [y/n]: ", *rtNumber)
			fmt.Scanln(&response)
			if response != "y" {
				fmt.Println("Operation cancelled.")
				os.Exit(1)
			}
		}
		k8sClient, err := initKubernetesClient(*kubeconfig)
		fmt.Printf("Creating load with ID: %s, Name Prefix: %s, Kubeconfig: %s, Runtime Number: %d\n", *loadID, *namePrefix, *kubeconfig, *rtNumber)
		return Create, action.NewWorker(*loadID, *namePrefix, *rtNumber, k8sClient)
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *loadIDDelete == "" || *kubeconfigDelete == "" {
			fmt.Println("all flags --load-id and --kubeconfig are required")
			deleteCmd.Usage()
			os.Exit(1)
		}
		k8sClient, err := initKubernetesClient(*kubeconfig)
		fmt.Printf("Deleting load with ID: %s, Kubeconfig: %s\n", *loadIDDelete, *kubeconfig)
		return Delete, action.NewWorker(*loadIDDelete, "", 0, k8sClient)
	default:
		fmt.Println("expected 'create' or 'delete' subcommands")
		os.Exit(1)
	}
	return Unknown, nil
}

func initKubernetesClient(kubeconfigPath string) (client.Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err.Error())
	}
	err = imv1.AddToScheme(k8sClient.Scheme())
	return k8sClient, err
}
