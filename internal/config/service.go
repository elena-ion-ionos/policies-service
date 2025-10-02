package config

import (
	"fmt"
	"time"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultMaxAllowedGoroutines int = 10000
	defaultDbPingTimeoutSec         = 3 * time.Second
)

// Service is the common configuration between the services
type Service struct {
	Database
	HttpClient

	MetricsPort     int
	MetricsAddr     string
	HealthCheckPort int
	HealthCheckAddr string

	HealthCheckMaxAllowedGoroutines int           // Liveness check
	HealthCheckDbPingTimeoutSec     time.Duration // Readiness check
}

// MustInitConf returns the configuration for the application read from the environment variables or defaults.
func (co *Service) MustInitConf() {
	cmd := &cobra.Command{}
	co.AddFlags(cmd)

	viper.AutomaticEnv()
	if err := InitViperFlags(cmd, []string{}); err != nil {
		paaslog.Fatalf("failed reading configuration flags, err: %v", err)
	}
}

// AddFlags adds the flags for the worker.
func (co *Service) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&co.MetricsPort, "metrics-port", 8081, "Port to start the metrics API on.")
	co.MetricsAddr = fmt.Sprintf(":%d", co.MetricsPort)
	cmd.Flags().IntVar(&co.HealthCheckPort, "health-check-port", 8082, "Port to start the health check API on.")
	co.HealthCheckAddr = fmt.Sprintf("0.0.0.0:%d", co.HealthCheckPort)

	cmd.Flags().IntVar(&co.HealthCheckMaxAllowedGoroutines, "health-check-max-allowed-goroutines", defaultMaxAllowedGoroutines, "Max allowed no of goroutines per service.")
	cmd.Flags().DurationVar(&co.HealthCheckDbPingTimeoutSec, "health-check-db-ping-timeout-sec", defaultDbPingTimeoutSec, "Max allowed db connection timeout in seconds.")

	co.Database.AddFlags(cmd)
	co.HttpClient.AddFlags(cmd)
}

func (co *Service) ConfigureHealthCheckHandler(health healthcheck.Handler, db *sqlx.DB) error {
	if co.HealthCheckMaxAllowedGoroutines <= 0 {
		return fmt.Errorf("health-check-max-allowed-goroutines must be greater than 0")
	}
	// It's fine to use a higher goroutine number here because v2 worker runs on multiple concurrent goroutines
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(co.HealthCheckMaxAllowedGoroutines))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	// This checks we can get a connection to the db in a reasonable time. 3 sec might look too much but the system is
	// still responsive even with 3 sec db timeout, TODO: consider lowering it once performance is stable.
	if co.HealthCheckDbPingTimeoutSec == 0 {
		return fmt.Errorf("health-check-db-ping-timeout-sec must be greater than 0")
	}
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db.DB, co.HealthCheckDbPingTimeoutSec))

	return nil
}
