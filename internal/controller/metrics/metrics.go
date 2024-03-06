package metrics

import (
	"fmt"

	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	runtimeIDKeyName               = "runtimeId"
	state                          = "state"
	reason                         = "reason"
	componentName                  = "infrastructure_manager"
	RuntimeIDLabel                 = "kyma-project.io/runtime-id"
	GardenerClusterStateMetricName = "im_gardener_clusters_state"
)

type Metrics struct {
	gardenerClustersStateGaugeVec *prometheus.GaugeVec
}

func NewMetrics() Metrics {
	m := Metrics{
		gardenerClustersStateGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: componentName,
				Name:      GardenerClusterStateMetricName,
				Help:      "Indicates the Status.state for GardenerCluster CRs",
			}, []string{runtimeIDKeyName, state, reason}),
	}
	ctrlMetrics.Registry.MustRegister(m.gardenerClustersStateGaugeVec)
	return m
}

func (m Metrics) SetGardenerClusterStates(cluster v1.GardenerCluster) {
	var runtimeID = cluster.GetLabels()[RuntimeIDLabel]

	if runtimeID != "" {
		if len(cluster.Status.Conditions) != 0 {
			var reason = cluster.Status.Conditions[0].Reason

			// first clean the old metric
			m.cleanUpGardenerClusterGauge(runtimeID)
			m.gardenerClustersStateGaugeVec.WithLabelValues(runtimeID, string(cluster.Status.State), reason).Set(1)
			fmt.Printf("\n++set metrics WithLabelValues(%v, %v, %v)", runtimeID, string(cluster.Status.State), reason)
		}
	}
}

func (m Metrics) UnSetGardenerClusterStates(runtimeID string) {
	m.cleanUpGardenerClusterGauge(runtimeID)
}

func (m Metrics) cleanUpGardenerClusterGauge(runtimeID string) {
	var readyMetric, _ = m.gardenerClustersStateGaugeVec.GetMetricWithLabelValues(runtimeID, "Ready")
	if readyMetric != nil {
		readyMetric.Set(0)
		fmt.Printf("\n--cleanUpGardenerClusterGauge(set value to 0 for %v)", runtimeID)
	}
	var errorMetric, _ = m.gardenerClustersStateGaugeVec.GetMetricWithLabelValues(runtimeID, "Error")
	if errorMetric != nil {
		errorMetric.Set(0)
		fmt.Printf("\n--cleanUpGardenerClusterGauge(set value to 0 for %v)", runtimeID)
	}

	metricsDeleted := m.gardenerClustersStateGaugeVec.DeletePartialMatch(prometheus.Labels{
		runtimeIDKeyName: runtimeID,
	})

	if metricsDeleted > 0 {
		fmt.Printf("\n--cleanUpGardenerClusterGauge(deleted %d metrics for runtimeID %v)", metricsDeleted, runtimeID)
	}
}
