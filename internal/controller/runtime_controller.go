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

type ShootClient interface {
	Create(ctx context.Context, shoot *gardener.Shoot, opts v1.CreateOptions) (*gardener.Shoot, error)
}

// RuntimeReconciler reconciles a Runtime object
type RuntimeReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	ShootClient ShootClient
	Log         logr.Logger
}

//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=runtimes/finalizers,verbs=update

func (r *RuntimeReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info(request.String())

	var runtime imv1.Runtime

	err := r.Get(ctx, request.NamespacedName, &runtime)

	if err == nil {
		r.Log.Info("Reconciling Runtime", "Name", runtime.Name, "Namespace", runtime.Namespace)
	} else {
		r.Log.Error(err, "unable to fetch Runtime")
		return ctrl.Result{}, err
	}

	shoot := &gardener.Shoot{}

	converterConfig := FixConverterConfig()
	converter := gardener_shoot.NewConverter(converterConfig)

	*shoot, err = converter.ToShoot(runtime)

	if err != nil {
		r.Log.Error(err, "unable to map Runtime to Shoot")
		return ctrl.Result{}, err
	}

	r.Log.Info("Shoot mapped", "Name", shoot.Name, "Namespace", shoot.Namespace, "Shoot", shoot)

	createdShoot, provisioningErr := r.ShootClient.Create(ctx, shoot, v1.CreateOptions{})

	if provisioningErr != nil {
		r.Log.Error(provisioningErr, "unable to create Shoot")
		return ctrl.Result{}, provisioningErr
	}
	r.Log.Info("Shoot created", "Name", createdShoot.Name, "Namespace", createdShoot.Namespace)

	return ctrl.Result{}, nil
}

func FixConverterConfig() gardener_shoot.ConverterConfig {
	return gardener_shoot.ConverterConfig{
		Kubernetes: gardener_shoot.KubernetesConfig{
			DefaultVersion: "1.29", //TODO: Should be parametrised
		},

		DNS: gardener_shoot.DNSConfig{
			SecretName:   "aws-route53-secret-dev",
			DomainPrefix: "dev.kyma.ondemand.com",
			ProviderType: "aws-route53",
		},
		Provider: gardener_shoot.ProviderConfig{
			AWS: gardener_shoot.AWSConfig{
				EnableIMDSv2: true, //TODO: Should be parametrised
			},
		},
		Gardener: gardener_shoot.GardenerConfig{
			ProjectName: "kyma-dev", //TODO: should be parametrised
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
