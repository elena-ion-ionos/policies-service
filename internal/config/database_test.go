package config

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkDatabase(t *testing.T, opts *Database, cmd *cobra.Command) {
	// Check all flags
	assert.NotNil(t, cmd.Flag("PGHOST"))
	assert.NotNil(t, cmd.Flag("PGDATABASE"))
	assert.NotNil(t, cmd.Flag("PGPORT"))
	assert.NotNil(t, cmd.Flag("PGUSER"))
	assert.NotNil(t, cmd.Flag("PGPASSWORD"))
	assert.NotNil(t, cmd.Flag("PGSSLMODE"))
	assert.NotNil(t, cmd.Flag("PG_MAX_OPEN_CONN"))
	assert.NotNil(t, cmd.Flag("PG_MAX_IDLE_CONN"))
	assert.NotNil(t, cmd.Flag("PG_MAX_LIFETIME"))
	assert.NotNil(t, cmd.Flag("PG_NOTIFY_CHANNEL_NAME"))
	assert.NotNil(t, cmd.Flag("PG_MIN_RECONNECT_INTERVAL_MS"))
	assert.NotNil(t, cmd.Flag("PG_MAX_RECONNECT_INTERVAL_MS"))

	// Check default values
	assert.Equal(t, "localhost", opts.Host)
	assert.Equal(t, "test", opts.Database)
	assert.Equal(t, "5432", opts.Port)
	assert.Equal(t, "postgres", opts.User)
	assert.Equal(t, "root", opts.Password)
	assert.Equal(t, "disable", opts.SslMode)
	assert.Equal(t, 20, opts.MaxOpenConn)
	assert.Equal(t, 10, opts.MaxIdleConn)
	assert.Equal(t, time.Duration(0), opts.MaxLifetime)
	assert.Equal(t, "notify-worker", opts.NotifyChannelName)
	assert.Equal(t, 100, opts.MinReconnectIntervalMs)
	assert.Equal(t, 10000, opts.MaxReconnectIntervalMs)
}

func TestDatabaseAddFlagsDefault(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &Database{}

	opts.AddFlags(cmd)

	checkDatabase(t, opts, cmd)
}

func TestDatabase_PgConfig(t *testing.T) {
	// Valid port
	opts := &Database{
		Host:     "h",
		Database: "d",
		User:     "u",
		Password: "p",
		Port:     "1234",
		SslMode:  "ssl",
	}
	cfg, err := opts.PgConfig()
	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "h", cfg.Host)
	assert.Equal(t, 1234, cfg.Port)
	assert.Equal(t, "u", cfg.User)
	assert.Equal(t, "p", cfg.Password)
	assert.Equal(t, "d", cfg.Name)
	assert.Equal(t, "ssl", cfg.SSLMode)
}

// To test error handling, use a recover since PgConfig() calls Fatalf
func TestDatabase_PgConfig_InvalidPort(t *testing.T) {
	opts := &Database{Port: "notanint"}
	cfg, err := opts.PgConfig()
	assert.Nil(t, cfg)
	assert.NotNil(t, err)
}
