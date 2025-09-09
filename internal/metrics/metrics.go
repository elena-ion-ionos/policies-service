package metrics

import (
	"time"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/metric"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	LabelResult         = "result"
	LabelResult_success = "success"
	LabelResult_fail    = "fail"

	LabelLocation                     = "location"
	LabelLocation_Worker              = "worker"
	LabelLocation_WorkerBackendSplit  = "workerBackendSplit"
	LabelLocation_WorkerFrontendSplit = "workerFrontendSplit"
	LabelLocation_CephRepo            = "cephRepo"
	LabelLocation_CloudianRepo        = "cloudianRepo"
	LabelLocation_Controller          = "controller"
	LabelLocation_DbRepo              = "v2DbRepo"

	LabelOperation = "operation"
	// dbrepo ops
	LabelOperation_CreateKey         = "CreateKey"
	LabelOperation_DeleteKey         = "DeleteKey"
	LabelOperation_GetKey            = "GetKey"
	LabelOperation_ListContractKeys  = "ListContractKeys"
	LabelOperation_ListKeys          = "ListKeys"
	LabelOperation_RenewSecret       = "RenewSecret"
	LabelOperation_UpdateKey         = "UpdateKey"
	LabelOperation_GetKeyByAccessKey = "GetKeyByAccessKey"
	LabelOperation_ChangeKeyStatus   = "ChangeKeyStatus"

	// backend queue ops
	LabelOperation_GetKeyOperations = "GetKeyOperations"

	// Worker
	// how long it takes the worker to process a whole batch
	LabelOperation_WorkerProcessBatch = "WorkerProcessBatch"
	// how long it takes one of the worker goroutines to process all ops of a contract from a batch (includes processing time + time waiting for goroutine to start)
	LabelOperation_WorkerProcessContractOps = "WorkerProcessContractOps"
	// how much time it takes one group of key ops to be processed since they were read from db by worker
	LabelOperation_KeyOpThreadTimeout = "KeyOpThreadTimeout"
	LabelOperation_KeyOpHandleTimeout = "KeyOpHandleTimeout"
	LabelOperation_KeyOpTotalTimeout  = "KeyOpTotalTimeout"
)

var (
	OpsDurationSeconds = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "s3_ops_duration_seconds",
			Help:       "Duration of the various key operations.",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.03, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 0.999: 0.0001},
			MaxAge:     5 * time.Minute,
		},
		[]string{LabelOperation, LabelLocation},
	)

	OpsNo = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "s3_ops_no",
			Help: "Total number of various key operations.",
		},
		[]string{LabelOperation, LabelLocation, LabelResult},
	)

	WorkerLoopBatchSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "s3_worker_loop_ops_batch_size",
			Help: "Number of key operations read by the worker loop in one batch.",
		},
	)

	WorkerLoopContractsPerBatch = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "s3_worker_loop_contracts_per_batch",
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
