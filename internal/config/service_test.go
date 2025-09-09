package config

import (
	"net/http"
	"testing"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkService(t *testing.T, opts *Service, cmd *cobra.Command) {
	// Check all Service-specific flags
	assert.NotNil(t, cmd.Flag("metrics-port"))
	assert.NotNil(t, cmd.Flag("health-check-port"))
	assert.NotNil(t, cmd.Flag("health-check-max-allowed-goroutines"))
	assert.NotNil(t, cmd.Flag("health-check-db-ping-timeout-sec"))

	// Check default values
	assert.Equal(t, 8081, opts.MetricsPort)
	assert.Equal(t, ":8081", opts.MetricsAddr)
	assert.Equal(t, 8082, opts.HealthCheckPort)
	assert.Equal(t, "0.0.0.0:8082", opts.HealthCheckAddr)
	assert.Equal(t, 10000, opts.HealthCheckMaxAllowedGoroutines)
	assert.Equal(t, 3*time.Second, opts.HealthCheckDbPingTimeoutSec)

	checkDatabase(t, &opts.Database, cmd)
	checkHttpClient(t, &opts.HttpClient, cmd)
}

func TestCommonAddFlagsDefault(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &Service{}

	opts.AddFlags(cmd)

	checkService(t, opts, cmd)
}

// Mock healthcheck.Handler
type mockHealthcheck struct {
	livenessChecks  map[string]healthcheck.Check
	readinessChecks map[string]healthcheck.Check
}

func newMockHealthcheck() *mockHealthcheck {
	return &mockHealthcheck{
		livenessChecks:  make(map[string]healthcheck.Check),
		readinessChecks: make(map[string]healthcheck.Check),
	}
}

func (m *mockHealthcheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func (m *mockHealthcheck) AddLivenessCheck(name string, check healthcheck.Check) {
	m.livenessChecks[name] = check
}

func (m *mockHealthcheck) AddReadinessCheck(name string, check healthcheck.Check) {
	m.readinessChecks[name] = check
}
func (m *mockHealthcheck) LiveEndpoint(http.ResponseWriter, *http.Request)  {}
func (m *mockHealthcheck) ReadyEndpoint(http.ResponseWriter, *http.Request) {}

func TestConfigureHealthCheckHandler_Valid(t *testing.T) {
	opts := &Service{
		HealthCheckMaxAllowedGoroutines: 10,
		HealthCheckDbPingTimeoutSec:     2 * time.Second,
	}
	health := newMockHealthcheck()
	db := &sqlx.DB{}

	err := opts.ConfigureHealthCheckHandler(health, db)
	assert.NoError(t, err)
	assert.Contains(t, health.livenessChecks, "goroutine-threshold")
	assert.NotNil(t, health.livenessChecks["goroutine-threshold"])
	assert.Contains(t, health.readinessChecks, "database")
	assert.NotNil(t, health.readinessChecks["database"])
}

func TestConfigureHealthCheckHandler_InvalidGoroutines(t *testing.T) {
	opts := &Service{
		HealthCheckMaxAllowedGoroutines: 0,
		HealthCheckDbPingTimeoutSec:     2 * time.Second,
	}
	health := newMockHealthcheck()
	db := &sqlx.DB{}

	err := opts.ConfigureHealthCheckHandler(health, db)
	assert.Error(t, err)
	assert.EqualError(t, err, "health-check-max-allowed-goroutines must be greater than 0")
}

func TestConfigureHealthCheckHandler_InvalidDbTimeout(t *testing.T) {
	opts := &Service{
		HealthCheckMaxAllowedGoroutines: 10,
		HealthCheckDbPingTimeoutSec:     0,
	}
	health := newMockHealthcheck()
	db := &sqlx.DB{}

	err := opts.ConfigureHealthCheckHandler(health, db)
	assert.Error(t, err)
	assert.EqualError(t, err, "health-check-db-ping-timeout-sec must be greater than 0")
}
