package metrics

import (
	"fmt"
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
	"strconv"
	"time"
)

const (
	runtimeIDKeyName               = "runtimeId"
	componentName                  = "infrastructure_manager"
	RuntimeIDLabel                 = "kyma-project.io/runtime-id"
	GardenerClusterStateMetricName = "im_gardener_clusters_state"
	state                          = "state"
	reason                         = "reason"
	KubeconfigExpirationMetricName = "im_kubeconfig_expiration"
	expires                        = "expires"
	lastSyncAnnotation             = "operator.kyma-project.io/last-sync"
)

type Metrics struct {
	gardenerClustersStateGaugeVec *prometheus.GaugeVec
	kubeconfigExpirationGauge     *prometheus.GaugeVec
}

func NewMetrics() Metrics {
	m := Metrics{
		gardenerClustersStateGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      GardenerClusterStateMetricName,
				Help:      "Indicates the Status.state for GardenerCluster CRs",
			}, []string{runtimeIDKeyName, state, reason}),
		kubeconfigExpirationGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      KubeconfigExpirationMetricName,
				Help:      "Exposes current kubeconfig expiration value in epoch timestamp value format",
			}, []string{runtimeIDKeyName, expires}),
	}
	ctrlMetrics.Registry.MustRegister(m.gardenerClustersStateGaugeVec, m.kubeconfigExpirationGauge)
	return m
}

func (m Metrics) SetGardenerClusterStates(cluster v1.GardenerCluster) {
	var runtimeID = cluster.GetLabels()[RuntimeIDLabel]

	if runtimeID != "" {
		if len(cluster.Status.Conditions) != 0 {
			var reason = cluster.Status.Conditions[0].Reason

			// first clean the old metric
			m.CleanUpGardenerClusterGauge(runtimeID)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeID, string(cluster.Status.State), reason).Set(1)
		}
	}
}

func (m Metrics) CleanUpGardenerClusterGauge(runtimeID string) {
	m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})
}

func (m Metrics) CleanUpKubeconfigExpiration(runtimeID string) {
	metricsDeleted := m.kubeconfigExpirationGauge.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})

	fmt.Printf("%v metrics has been deleted\n", metricsDeleted)
}

func (m Metrics) SetKubeconfigExpiration(secret corev1.Secret) {
	var runtimeID = secret.GetLabels()[RuntimeIDLabel]

	// first clean the old metric
	m.CleanUpKubeconfigExpiration(runtimeID)

	if runtimeID != "" {
		var lastSyncTime = secret.GetAnnotations()[lastSyncAnnotation]

		parsedSyncTime, err := time.Parse(time.RFC3339, lastSyncTime)
		if err == nil {
			syncTimeEpoch := parsedSyncTime.Unix()
			syncTimeEpochString := strconv.Itoa(int(syncTimeEpoch))
			m.kubeconfigExpirationGauge.WithLabelValues(runtimeID, syncTimeEpochString)
		}
	}
}
