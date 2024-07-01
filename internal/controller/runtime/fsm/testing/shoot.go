package testing

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ShootNoDNS = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
	}

	ShootNoDNSDomain = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{},
		},
	}

	ShootMissingLastOperation = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{},
	}

	ShootEmptyLastOperation = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{},
		},
	}

	ShootLastOperationProcessing = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeCreate,
				State: gardener.LastOperationStateProcessing,
			},
		},
	}

	ShootLastOperationPending = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeCreate,
				State: gardener.LastOperationStatePending,
			},
		},
	}

	ShootLastOperationSucceeded = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeCreate,
				State: gardener.LastOperationStateSucceeded,
			},
		},
	}

	ShootLastOperationFailed = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeCreate,
				State: gardener.LastOperationStateFailed,
			},
		},
	}

	ShootLastOperationReconcileProcessing = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeReconcile,
				State: gardener.LastOperationStateProcessing,
			},
		},
	}

	ShootLastOperationReconcilePending = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeReconcile,
				State: gardener.LastOperationStatePending,
			},
		},
	}

	ShootLastOperationReconcileFailed = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeReconcile,
				State: gardener.LastOperationStateSucceeded,
			},
		},
	}

	ShootLastOperationReconcileSucceeded = gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
		Spec: gardener.ShootSpec{
			DNS: &gardener.DNS{
				Domain: new(string),
			},
		},
		Status: gardener.ShootStatus{
			LastOperation: &gardener.LastOperation{
				Type:  gardener.LastOperationTypeReconcile,
				State: gardener.LastOperationStateSucceeded,
			},
		},
	}
)
