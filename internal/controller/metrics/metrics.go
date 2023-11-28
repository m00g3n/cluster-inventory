package metrics

import (
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	shootName       = "shootName"
	globalAccountId = "globalAccountId"
	runtimeId       = "runtimeId"
	state           = "state"
)

var (

	//TODO: test custom metric, remove when done with https://github.com/kyma-project/infrastructure-manager/issues/11
	totalReconciliationLoopsStarted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "im_reconciliation_loops_started_total",
			Help: "Number of times reconciliation loop was started",
		},
	)

	metricGardenerClustersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{ //nolint:gochecknoglobals
			Subsystem: "infrastructure_manager",
			Name:      "im_gardener_clusters_state",
			Help:      "Indicates the Status.state for GardenerCluster CRs",
		}, []string{shootName, state})
)

func init() {
	ctrlMetrics.Registry.MustRegister(totalReconciliationLoopsStarted, metricGardenerClustersState)
}

func IncrementReconciliationLoopsStarted() {
	totalReconciliationLoopsStarted.Inc()
}

func SetGardenerClusterStates(cluster v1.GardenerCluster) {
	// Increase a value using compact (but order-sensitive!) WithLabelValues().
	//metricGardenerClustersState.WithLabelValues(cluster.Spec.Shoot.Name, "READY").Add(1)
	metricGardenerClustersState.WithLabelValues(cluster.Spec.Shoot.Name, cluster.Spec.Shoot.Name).Add(1)
}
