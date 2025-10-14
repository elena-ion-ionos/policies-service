package metrics

import (
	"time"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/metric"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	LabelResult                 = "result"
	LabelResult_success         = "success"
	LabelResult_fail            = "fail"
	LabelLocation               = "location"
	LabelLocation_Repo          = "repo"
	LabelOperation              = "operation"
	LabelOperation_SaveUser     = "SaveUser"
	LabelOperation_FindUserById = "FindUserById"
)

var (
	OpsDurationSeconds = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "ops_duration_seconds",
			Help:       "Duration of the various user operations.",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.03, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 0.999: 0.0001},
			MaxAge:     5 * time.Minute,
		},
		[]string{LabelOperation, LabelLocation},
	)

	OpsNo = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ops_no",
			Help: "Total number of various user operations.",
		},
		[]string{LabelOperation, LabelLocation, LabelResult},
	)

	WorkerLoopBatchSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_loop_ops_batch_size",
			Help: "Number of user operations read by the worker loop in one batch.",
		},
	)

	WorkerLoopContractsPerBatch = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_loop_contracts_per_batch",
			Help: "Number of actual contracts read by the worker loop in one batch.",
		},
	)
)

func init() {
	metric.GlobalRegistry.MustRegister(
		OpsDurationSeconds,
		OpsNo,
		WorkerLoopBatchSize,
		WorkerLoopContractsPerBatch)
}
