package metrics

import (
	"strconv"
	"time"

	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	runtimeIDKeyName               = "runtimeId"
	shootNameIDKeyName             = "shootName"
	rotationDuration               = "rotationDuration"
	expirationDuration             = "expirationDuration"
	componentName                  = "infrastructure_manager"
	RuntimeIDLabel                 = "kyma-project.io/runtime-id"
	ShootNameLabel                 = "kyma-project.io/shoot-name"
	GardenerClusterStateMetricName = "im_gardener_clusters_state"
	RuntimeStateMetricName         = "im_runtime_state"
	provider                       = "provider"
	state                          = "state"
	reason                         = "reason"
	message                        = "message"
	KubeconfigExpirationMetricName = "im_kubeconfig_expiration"
	expires                        = "expires"
	lastSyncAnnotation             = "operator.kyma-project.io/last-sync"
)

type Metrics struct {
	gardenerClustersStateGaugeVec *prometheus.GaugeVec
	kubeconfigExpirationGauge     *prometheus.GaugeVec
	runtimeStateGauge             *prometheus.GaugeVec
}

func NewMetrics() Metrics {
	m := Metrics{
		gardenerClustersStateGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      GardenerClusterStateMetricName,
				Help:      "Indicates the Status.state for GardenerCluster CRs",
			}, []string{runtimeIDKeyName, shootNameIDKeyName, state, reason}),
		kubeconfigExpirationGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      KubeconfigExpirationMetricName,
				Help:      "Exposes current kubeconfig expiration value in epoch timestamp value format",
			}, []string{runtimeIDKeyName, shootNameIDKeyName, expires, rotationDuration, expirationDuration}),
		runtimeStateGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      RuntimeStateMetricName,
				Help:      "Exposes current Status.state for Runtime CRs",
			}, []string{runtimeIDKeyName, shootNameIDKeyName, provider, state, reason, message}),
	}
	ctrlMetrics.Registry.MustRegister(m.gardenerClustersStateGaugeVec, m.kubeconfigExpirationGauge, m.runtimeStateGauge)
	return m
}

func (m Metrics) SetGardenerClusterStates(cluster v1.GardenerCluster) {
	var runtimeID = cluster.GetLabels()[RuntimeIDLabel]
	var shootName = cluster.GetLabels()[ShootNameLabel]

	if runtimeID != "" {
		if len(cluster.Status.Conditions) != 0 {
			var reason = cluster.Status.Conditions[0].Reason

			// first clean the old metric
			m.CleanUpGardenerClusterGauge(runtimeID)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeID, shootName, string(cluster.Status.State), reason).Set(1)
		}
	}
}

func (m Metrics) CleanUpGardenerClusterGauge(runtimeID string) {
	m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})
}

func (m Metrics) SetRuntimeStates(runtime v1.Runtime) {
	runtimeID := runtime.GetLabels()[RuntimeIDLabel]

	if runtimeID != "" {
		if len(runtime.Status.Conditions) != 0 {
			var reason = runtime.Status.Conditions[0].Reason // will change it

			// first clean the old metric
			m.CleanUpRuntimeGauge(runtimeID)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeID, runtime.Spec.Shoot.Name, runtime.Spec.Shoot.Provider.Type, string(runtime.Status.State), reason).Set(1)
		}
	}
}

func (m Metrics) CleanUpRuntimeGauge(runtimeID string) {
	m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})
}

func (m Metrics) CleanUpKubeconfigExpiration(runtimeID string) {
	m.kubeconfigExpirationGauge.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})
}

func computeExpirationInSeconds(rotationPeriod time.Duration, minimalRotationTimeRatio float64) float64 {
	return rotationPeriod.Seconds() / minimalRotationTimeRatio
}

func (m Metrics) SetKubeconfigExpiration(secret corev1.Secret, rotationPeriod time.Duration, minimalRotationTimeRatio float64) {
	var runtimeID = secret.GetLabels()[RuntimeIDLabel]
	var shootName = secret.GetLabels()[ShootNameLabel]

	// first clean the old metric
	m.CleanUpKubeconfigExpiration(runtimeID)

	if runtimeID != "" {
		var lastSyncTime = secret.GetAnnotations()[lastSyncAnnotation]

		parsedSyncTime, err := time.Parse(time.RFC3339, lastSyncTime)
		if err == nil {
			expirationTimeInSeconds := computeExpirationInSeconds(rotationPeriod, minimalRotationTimeRatio)
			expirationTime := parsedSyncTime.Add(time.Duration(expirationTimeInSeconds * float64(time.Second)))

			expirationTimeEpoch := expirationTime.Unix()
			expirationTimeEpochString := strconv.Itoa(int(expirationTimeEpoch))
			rotationPeriodString := strconv.FormatFloat(rotationPeriod.Seconds(), 'G', -1, 64)

			expirationTimeInSecondsString := strconv.FormatFloat(expirationTimeInSeconds, 'G', -1, 64)

			m.kubeconfigExpirationGauge.WithLabelValues(
				runtimeID,
				shootName,
				expirationTimeEpochString,
				rotationPeriodString,
				expirationTimeInSecondsString,
			).Set(float64(expirationTimeEpoch))
		}
	}
}
