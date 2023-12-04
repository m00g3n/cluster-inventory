package metrics

import (
	"github.com/go-logr/logr"
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

	//TODO: test custom metric, remove when done with https://github.com/kyma-project/infrastructure-manager/issues/11 //nolint:all
	playground_totalReconciliationLoopsStarted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "im_playground_reconciliation_loops_started_total",
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
	ctrlMetrics.Registry.MustRegister(playground_totalReconciliationLoopsStarted, metricGardenerClustersState)
}

func IncrementReconciliationLoopsStarted() {
	playground_totalReconciliationLoopsStarted.Inc()
}

func SetGardenerClusterStates(cluster v1.GardenerCluster, log logr.Logger) {
	metricGardenerClustersState.WithLabelValues(cluster.Spec.Shoot.Name, string(cluster.Status.State)).Set(1)
}
