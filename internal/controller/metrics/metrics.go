package metrics

import (
	"fmt"

	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	runtimeIDKeyName = "runtimeId"
	state            = "state"
	reason           = "reason"
	runtimeIDLabel   = "kyma-project.io/runtime-id"
	componentName    = "infrastructure_manager"
)

type Metrics struct {
	gardenerClustersStateGaugeVec *prometheus.GaugeVec
}

func NewMetrics() Metrics {
	m := Metrics{
		gardenerClustersStateGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      "im_gardener_clusters_state",
				Help:      "Indicates the Status.state for GardenerCluster CRs",
			}, []string{runtimeIDKeyName, state, reason}),
	}
	ctrlMetrics.Registry.MustRegister(m.gardenerClustersStateGaugeVec)
	return m
}

func (m Metrics) SetGardenerClusterStates(cluster v1.GardenerCluster) {
	var runtimeId = cluster.GetLabels()[runtimeIDLabel]

	if runtimeId != "" {
		if len(cluster.Status.Conditions) != 0 {
			var reason = cluster.Status.Conditions[0].Reason

			// first clean the old metric
			m.cleanUpGardenerClusterGauge(runtimeId)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeId, string(cluster.Status.State), reason).Set(1)
		}
	}
}

func (m Metrics) UnSetGardenerClusterStates(runtimeId string) {
	m.cleanUpGardenerClusterGauge(runtimeId)
}

func (m Metrics) cleanUpGardenerClusterGauge(runtimeId string) {
	var readyMetric, _ = m.gardenerClustersStateGaugeVec.GetMetricWithLabelValues(runtimeId, "Ready")
	if readyMetric != nil {
		readyMetric.Set(0)
	}
	var errorMetric, _ = m.gardenerClustersStateGaugeVec.GetMetricWithLabelValues(runtimeId, "Error")
	if errorMetric != nil {
		errorMetric.Set(0)
	}
	fmt.Printf("GardenerClusterStates set value to 0 for %v", runtimeId)

	metricsDeleted := m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeId,
	})

	if metricsDeleted > 0 {
		fmt.Printf("gardenerClusterStateGauge deleted %d metrics for runtimeId %v", metricsDeleted, runtimeId)
	}
}
