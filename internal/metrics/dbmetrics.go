package metrics

import (
	"fmt"
	"time"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/metric"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/jmoiron/sqlx"
)

const (
	metricNamespace = "sqlx"
	metricSubsystem = "connection"
)

// supportedMetrics exposes a subset of sql.DBStats
type supportedMetrics struct {
	// max num of open connections to the database
	MaxOpenConnections int

	// Pool status
	OpenConnections int
	InUse           int
	Idle            int

	// Total number of connections waited for
	WaitCount int64
	// Total time blocked waiting for a new connection
	WaitDuration time.Duration
}

// promMetric contains data needed to register a Prometheus metric
type promMetric struct {
	Desc    *prometheus.Desc
	Eval    func(supportedMetrics) float64
	ValType prometheus.ValueType
}

// metricsCollector implements a Prometheus collector for client side db metrics
type metricsCollector struct {
	dbs map[string]*sqlx.DB

	promMetrics []promMetric
}

// MonitoredDbConfig contains necessary data to register a db for monitoring
type MonitoredDbConfig struct {
	DbName string
	DB     *sqlx.DB
}

var metricsCollector_ *metricsCollector

// init will do collector registration one time only
func init() {
	metricsCollector_ = &metricsCollector{
		dbs: map[string]*sqlx.DB{},
	}
	metricsCollector_.promMetrics = metricsCollector_.mustNewPromMetrics()
	metric.GlobalRegistry.MustRegister(metricsCollector_)
}

// MustInitMetrics will initialize collection of Prometheus metrics for all provided database configurations.
// It can monitor multiple databases at the same time.
func MustInitMetrics(cfgs []MonitoredDbConfig) {
	for _, cfg := range cfgs {
		paaslog.Infof("initializing database metrics for database `%s`", cfg.DbName)

		if _, found := metricsCollector_.dbs[cfg.DbName]; found {
			err := fmt.Errorf("database `%s` has already been registered for metrics", cfg.DbName)
			paaslog.Warnf("%v", err)
		}

		metricsCollector_.dbs[cfg.DbName] = cfg.DB
	}
}

// mustNewPromMetrics contains definition of all collected db metrics
func (m *metricsCollector) mustNewPromMetrics() []promMetric {
	// db_name: configured database name
	defaultLabels := []string{"db_name"}

	return []promMetric{
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "max_open"),
				"Maximum number of open connections to the database.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.MaxOpenConnections)
			},
			ValType: prometheus.GaugeValue,
		},
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "open"),
				"Number of established connections to the database both in use and idle.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.OpenConnections)
			},
			ValType: prometheus.GaugeValue,
		},
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "in_use"),
				"Number of connections currently in use.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.InUse)
			},
			ValType: prometheus.GaugeValue,
		},
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "idle"),
				"Number of idle connections.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.Idle)
			},
			ValType: prometheus.GaugeValue,
		},
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "wait_total"),
				"Total number of connections waited for.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.WaitCount)
			},
			ValType: prometheus.CounterValue,
		},
		{
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(metricNamespace, metricSubsystem, "wait_duration_total"),
				"Total time blocked waiting for a new connection in seconds.",
				defaultLabels, nil),
			Eval: func(metric supportedMetrics) float64 {
				return float64(metric.WaitDuration)
			},
			ValType: prometheus.CounterValue,
		},
	}
}

// Describe is required by Prometheus Collector interface
func (m *metricsCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range m.promMetrics {
		ch <- metric.Desc
	}
}

// Collect is required by Prometheus Collector interface
// The method must be concurrent safe as per Prometheus requirements
func (m *metricsCollector) Collect(ch chan<- prometheus.Metric) {
	for dbname, db := range m.dbs {
		stats := db.DB.Stats()

		collectedMetric := supportedMetrics{
			MaxOpenConnections: stats.MaxOpenConnections,
			OpenConnections:    stats.OpenConnections,
			InUse:              stats.InUse,
			Idle:               stats.Idle,
			WaitCount:          stats.WaitCount,
			WaitDuration:       stats.WaitDuration,
		}

		for _, i := range m.promMetrics {
			ch <- prometheus.MustNewConstMetric(i.Desc, i.ValType, i.Eval(collectedMetric), dbname)
		}
	}
}
