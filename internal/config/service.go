package config

import (
	"fmt"
	"time"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
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
