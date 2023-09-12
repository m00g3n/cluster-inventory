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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type State string

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GardenerCluster is the Schema for the clusters API
type GardenerCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GardenerClusterSpec   `json:"spec,omitempty"`
	Status GardenerClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GardenerClusterList contains a list of GardenerCluster
type GardenerClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GardenerCluster `json:"items"`
}

// GardenerClusterSpec defines the desired state of GardenerCluster
type GardenerClusterSpec struct {
	Shoot      Shoot      `json:"shoot,omitempty"`
	Kubeconfig Kubeconfig `json:"kubeconfig,omitempty"`
}

// Shoot defines the desired state of the Gardener's shoot
type Shoot struct {
	Name string `json:"name,omitempty"`
}

// Kubeconfig defines the desired kubeconfig location
type Kubeconfig struct {
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

// SecretKeyRef defines the location, and structure of the secret containing kubeconfig
type SecretKeyRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Key       string `json:"key,omitempty"`
}

// GardenerClusterStatus defines the observed state of GardenerCluster
type GardenerClusterStatus struct {
	// State signifies current state of Gardener Cluster.
	// Value can be one of ("Ready", "Processing", "Error", "Deleting").
	State State `json:"state,omitempty"`

	// List of status conditions to indicate the status of a ServiceInstance.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&GardenerCluster{}, &GardenerClusterList{})
}
