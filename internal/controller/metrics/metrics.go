package metrics

import (
	v1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/prometheus/client_golang/prometheus"
	ctrlMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	shootName = "shootName"
	state     = "state"
)

var (

	//nolint:godox //TODO: test custom metric, remove when done with https://github.com/kyma-project/infrastructure-manager/issues/11
	playgroundTotalReconciliationLoopsStarted = prometheus.NewCounter( //nolint:gochecknoglobals
		prometheus.CounterOpts{
			Name: "im_playground_reconciliation_loops_started_total",
			Help: "Number of times reconciliation loop was started",
		},
	)

	metricGardenerClustersState = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{ //nolint:gochecknoglobals
			Subsystem: "infrastructure_manager",
			Name:      "im_gardener_clusters_state",
			Help:      "Indicates the Status.state for GardenerCluster CRs",
		}, []string{shootName, state})
)

func init() {
	ctrlMetrics.Registry.MustRegister(playgroundTotalReconciliationLoopsStarted, metricGardenerClustersState)
}

func IncrementReconciliationLoopsStarted() {
	playgroundTotalReconciliationLoopsStarted.Inc()
}

func SetGardenerClusterStates(cluster v1.GardenerCluster) {
	metricGardenerClustersState.WithLabelValues(cluster.Spec.Shoot.Name, string(cluster.Status.State)).Set(1)
}
