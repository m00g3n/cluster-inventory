/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/client-go/tools/record"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/go-logr/logr"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// runtime reconciler specific configuration
type RCCfg struct {
	Finalizer string
}

//go:generate mockery --name=ShootClient
type ShootClient interface {
	Create(ctx context.Context, shoot *gardener.Shoot, opts v1.CreateOptions) (*gardener.Shoot, error)
	Get(ctx context.Context, name string, opts v1.GetOptions) (*gardener.Shoot, error)
	//Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	//List(ctx context.Context, opts v1.ListOptions) (*gardener.ShootList, error)
}

// RuntimeReconciler reconciles a Runtime object
type RuntimeReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	ShootClient   ShootClient
	Log           logr.Logger
	Cfg           RCCfg
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/finalizers,verbs=update

var requCounter = 0

func (r *RuntimeReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info(request.String())

	var runtime imv1.Runtime

	if err := r.Get(ctx, request.NamespacedName, &runtime); err != nil {
		return ctrl.Result{
			Requeue: true,
		}, client.IgnoreNotFound(err)
	}

	r.Log.Info("Reconciling Runtime", "Name", runtime.Name, "Namespace", runtime.Namespace)

	stateFSM := NewFsm(
		r.Log.WithName(fmt.Sprintf("reqID %d", requCounter)),
		r.Cfg,
		K8s{
			Client:        r.Client,
			shootClient:   r.ShootClient,
			EventRecorder: r.EventRecorder,
		})
	return stateFSM.Run(ctx, runtime)

}

func FixConverterConfig() gardener_shoot.ConverterConfig {
	return gardener_shoot.ConverterConfig{
		Kubernetes: gardener_shoot.KubernetesConfig{
			DefaultVersion: "1.29",
		},

		DNS: gardener_shoot.DNSConfig{
			SecretName:   "xxx-secret-dev",
			DomainPrefix: "runtimeprov.dev.kyma.ondemand.com",
			ProviderType: "aws-route53",
		},
		Provider: gardener_shoot.ProviderConfig{
			AWS: gardener_shoot.AWSConfig{
				EnableIMDSv2: true,
			},
		},
	}
}

func NewruntimeReconciler(mgr ctrl.Manager, shootClient ShootClient, logger logr.Logger) *RuntimeReconciler {
	return &RuntimeReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		ShootClient: shootClient,
		Log:         logger,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *RuntimeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&imv1.Runtime{}).
		WithEventFilter(predicate.And(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
