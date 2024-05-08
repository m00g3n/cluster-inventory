# Introduction
This document defines architecture, and API for provisioning functionality.

# Target architecture

The following picture shows the proposed architecture:
![](./assets/keb-kim-target-arch.drawio.svg)

> Note: at the time of writing the `GardenerCluster` CR is used for generating kubeconfig. The [workplan](https://github.com/kyma-project/infrastructure-manager/issues/112) for delivering provisioning functionality in the Kyma Infrastructure Manager includes renaming the CR to maintain consistency.

The following assumptions were taken:
- KEB is responsible for:
    - Creating `Runtime` CR containing the following data:
      - provider config (type, region, and secret with credentials for hyperscaler)
      - worker pool specification
      - cluster networking settings (nodes, pods, and services API ranges)
      - OIDC settings
      - cluster administrators list
      - Egress network filter settings
      - Control Plane failure tolerance
    - Observing status of the CR to determine whether provisioning succeeded
- Kyma Infrastructure Manager is responsible for:
    - creating shoots based on:
      - corresponding `Runtime` CR properties
      - predefined defaults for the optional properties:
        - Kubernetes version
        - Machine image version
      - predefined configuration for the following extensions:
        - DNS 
        - Certificates
    - upgrading, and deleting shoots for corresponding `Runtime` CRs
    - applying audit log configuration on the shoot resource
    - generating kubeconfig

# API proposal

## CR examples

Please mind that the `Runtime` CR should contain the following labels:
```yaml
 kyma-project.io/instance-id: instance-id
 kyma-project.io/runtime-id: runtime-id
 kyma-project.io/broker-plan-id: plan-id
 kyma-project.io/broker-plan-name: plan-name
 kyma-project.io/global-account-id: global-account-id
 kyma-project.io/subaccount-id: subAccount-id
 kyma-project.io/shoot-name: shoot-name
 kyma-project.io/region: region
 operator.kyma-project.io/kyma-name: kymaName
```

The labels are skipped in the following examples due to clarity.

The example below shows the CR that should be created by the KEB to provision AWS production cluster:
```yaml
apiVersion: infrastructuremanager.kyma-project.io/v1alpha1
kind: Runtime
metadata:
  name: runtime-id
  namespace: kcp-system
spec:
  shoot:
    # spec.shoot.name is set by the KEB, required
    name: shoot-name
    # spec.shoot.purpose is set by the KEB, required
    purpose: production
    kubernetes:
      kubeAPIServer:
        ## spec.shoot.kubernetes.kubeAPIServer.oidcConfig is provided by the KEB, required
        oidcConfig:
          clientID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
          groupsClaim: groups
          issuerURL: https://my.cool.tokens.com
          signingAlgs:
            - RS256
          usernameClaim: sub
    provider:
      type: aws
      region: eu-central-1
      # We must consider whether it makes sense to move HAP into KIM
      secretBindingName: "hypersaler secret"
    # spec.shoot.Networking is Provided by the KEB, required
    networking:
      pods: 100.64.0.0/12
      nodes: 10.250.0.0/16
      services: 100.104.0.0/13
    # spec.shoot.controlPlane is provided by the KEB, optional, default=nil
    controlPlane:
      highAvailability:
        failureTolerance:
          type: node
    workers:
      - machine:
          # spec.shoot.workers.machine.type provided by the KEB, required
          type: m6i.large
        # spec.shoot.workers.zones is provided by the KEB, required
        zones:
          - eu-central-1a
          - eu-central-1b
          - eu-central-1c
        # spec.shoot.workers.minimum is provided by the KEB, required
        minimum: 3
        # spec.shoot.workers.maximum is provided by the KEB, required
        maximum: 20
        # spec.shoot.workers.maxSurge is provided by the KEB, required in the first release.
        # It can be optional in the future, as it equals to zone count
        maxSurge: 3
        # spec.shoot.workers.maxUnavailable is provided by the KEB, required in the first release.
        # It can be optional in the future, as it is always set to 0
        maxUnavailable:  0
  security:
    networking:
      filter:
        # spec.security.networking is provided by the KEB, required
        egress:
          enabled: false
    # spec.security.administrators is provided by the KEB, required
    administrators:
      - admin@myorg.com
```

There are some additional optional fields  that could be specified:
- `spec.shoot.kubernetes.version` ; if not provided default value will be read by KIM from configuration
- `spec.shoot.workers.machine.image` ; if not provided default value will be read by KIM from configuration
- `spec.shoot.kubernetes.kubeAPIServer.additionalOidcConfig` ; if not provided, no addition OIDC provider will be configured
- `spec.shoot.workers.name` ; if not provided, some hardcoded name will be used
- `spec.security.networking.filtering.ingress.enabled` ; if not provided `false` value will be used

The following example shows what `Runtime` CR should be created to provision a cluster with additional OIDC provider, and ingress network filtering enabled:
```yaml
apiVersion: infrastructuremanager.kyma-project.io/v1alpha1
kind: Runtime
metadata:
  name: runtime-id
  namespace: kcp-system
spec:
  shoot:
    # spec.shoot.name is set by the KEB, required
    name: shoot-name
    # spec.shoot.purpose is set by the KEB, required
    purpose: production
    # Will be modified by the SRE
    kubernetes:
      # spec.shoot.kubernetes.version is optional, when not provided default will be used
      version: "1.28.7"
      kubeAPIServer:
        ## spec.shoot.kubernetes.kubeAPIServer.oidcConfig is provided by the KEB, required
        oidcConfig:
          clientID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
          groupsClaim: groups
          issuerURL: https://my.cool.tokens.com
          signingAlgs:
          - RS256
          usernameClaim: sub
        ## spec.shoot.kubernetes.kubeAPIServer.additionalOidcConfig is provided by the KEB, optional, not implemented in the first KIM release
        additionalOidcConfig:
          - clientID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
            groupsClaim: groups
            issuerURL: https://some.others.tokens.com
            signingAlgs:
              - RS256
            usernameClaim: sub
            usernamePrefix: 'someother'
    ## spec.shoot.provider is provided by the KEB, required
    provider:
      type: aws
      region: eu-central-1
      # We must consider whether it makes sense to move HAP into KIM
      secretBindingName: "hypersaler secret"
    # spec.shoot.Networking is Provided by the KEB, required
    networking:
      pods: 100.64.0.0/12
      nodes: 10.250.0.0/16
      services: 100.104.0.0/13
    # spec.shoot.controlPlane is provided by the KEB, optional, default=nil
    controlPlane:
      highAvailability:
        failureTolerance:
          type: node
    workers:
      - machine:
          # spec.shoot.workers.machine.type provided by the KEB, required
          type: m6i.large
          # spec.shoot.workers.machine.image is optional, when not provider default will be used
          # Will be modified by the SRE
          image:
            name: gardenlinux
            version: 1312.3.0
        # spec.shoot.workers.volume is provided by the KEB, required for the first release
        # Probably can be moved into KIM, as it is hardcoded in KEB, and not dependent on plan
        volume:
          type: gp2
          size: 50Gi
        # spec.shoot.workers.zones is provided by the KEB, required
        zones:
          - eu-central-1a
          - eu-central-1b
          - eu-central-1c
        # spec.shoot.workers.name is provided by the KEB. Optional, if not provided default will be used
        name: cpu-worker-0
        # spec.shoot.workers.minimum is provided by the KEB, required
        minimum: 3
        # spec.shoot.workers.maximum is provided by the KEB, required
        maximum: 20
        # spec.shoot.workers.maxSurge is provided by the KEB, required in the first release.
        # It can be optional in the future, as it equals to zone count
        maxSurge: 3
        # spec.shoot.workers.maxUnavailable is provided by the KEB, required in the first release.
        # It can be optional in the future, as it is always set to 0
        maxUnavailable:  0
  security:
    networking:
      filter:
        # spec.security.networking.filter.egress.enabled is provided by the KEB, required
        egress:
          enabled: false
        # spec.security.networking.filter.ingress.enabled will be provided by the KEB, optional (default=false)
        ingress:
          enabled: true
    # spec.security.administrators is provided by the KEB, required
    administrators:
      - admin@myorg.com
```

The following example 

Please, see the following examples to understand what CRs need to be created for particular KEB plans:
- [AWS trial plan](assets/runtime-examples/aws-trial.yaml)
- [Azure](assets/runtime-examples/azure.yaml)
- [Azure lite](assets/runtime-examples/azure-lite.yaml)
- [GCP](assets/runtime-examples/gcp.yaml)
- [SAP Converge Cloud](assets/runtime-examples/sap-converged-cloud.yaml)

## API structures

```go
package v2

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Runtime struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuntimeSpec   `json:"spec"`
	Status RuntimeStatus `json:"status,omitempty"`
}

type RuntimeSpec struct {
	Shoot    Shoot    `json:"spec"`
	Security Security `json:"security"`
}

type Shoot struct {
	Name       string             `json:"name"`
	Purpose    string             `json:"purpose"`
	Kubernetes Kubernetes         `json:"kubernetes"`
	Provider   Provider           `json:"provider"`
	Networking Networking         `json:"networking"`
	Workers    *[]gardener.Worker `json:"workers,omitempty"`
}

type Provider struct {
	Type              string `json:"type"`
	Region            string `json:"region"`
	SecretBindingName string `json:"secretBindingName"`
}

type Networking struct {
	Pods     *string `json:"pods,omitempty"`
	Nodes    *string `json:"nodes,omitempty"`
	Services *string `json:"services,omitempty"`
}

type Kubernetes struct {
	Version       string     `json:"version"`
	KubeAPIServer *APIServer `json:"kubeAPIServer,omitempty"`
}

type APIServer struct {
	oidcConfig           gardener.OIDCConfig    `json:"oidcConfig"`
	additionalOidcConfig *[]gardener.OIDCConfig `json:"additionalOidcConfig""`
}

type Security struct {
	Administrators []string           `json:"administrators"`
	Networking     NetworkingSecurity `json:"networking""`
}

type NetworkingSecurity struct {
	Filter Filter `json:"filter"`
}

type Filter struct {
	Ingress Ingress `json:"ingress"`
	Egress  Egress  `json:"egress"`
}

type Ingress struct {
	Enabled bool `json:"enabled"`
}

type Egress struct {
	Enabled bool `json:"enabled"`
}

type State string

// +kubebuilder:object:root=true
// RuntimeStatus defines the observed state of Runtime
type RuntimeStatus struct {
	// State signifies current state of Runtime.
	// Value can be one of ("Ready", "Processing", "Error", "Deleting").
	State State `json:"state,omitempty"`

	// List of status conditions to indicate the status of a ServiceInstance.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
```
