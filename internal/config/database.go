package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ionos-cloud/go-paaskit/infrastructure/paasql"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/metrics"
	"github.com/ionos-cloud/go-sample-service/internal/migration"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

type Database struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
	SslMode  string

	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration

	NotifyChannelName      string
	MinReconnectIntervalMs int
	MaxReconnectIntervalMs int
}

const (
	defaultMaxOpenConn               = 20
	defaultMaxIdleConn               = 10
	defaultMaxLifetime time.Duration = 0 // unlimited lifetime
)

func (o *Database) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Host, "PGHOST", "localhost", "ServerHost of the database.")
	cmd.Flags().StringVar(&o.Database, "PGDATABASE", "test", "S3 key management database.")
	cmd.Flags().StringVar(&o.Port, "PGPORT", "5432", "Port of the database.")
	cmd.Flags().StringVar(&o.User, "PGUSER", "postgres", "User to connect to the database.")
	cmd.Flags().StringVar(&o.Password, "PGPASSWORD", "root", "Password to connect to the database.")
	cmd.Flags().StringVar(&o.SslMode, "PGSSLMODE", "disable", "SSL mode to connect to the database.")
	cmd.Flags().IntVar(&o.MaxOpenConn, "PG_MAX_OPEN_CONN", defaultMaxOpenConn, "Max number of opened db connections allowed in the db client connection pool.")
	cmd.Flags().IntVar(&o.MaxIdleConn, "PG_MAX_IDLE_CONN", defaultMaxIdleConn, "Max number of idel db connections allowed in the db client connection pool.")
	cmd.Flags().DurationVar(&o.MaxLifetime, "PG_MAX_LIFETIME", defaultMaxLifetime, "Max lifetime of a db connections, 0 means unlimited.")
	cmd.Flags().StringVar(&o.NotifyChannelName, "PG_NOTIFY_CHANNEL_NAME", "notify-worker", "The name of the channel used to notify the worker when there's new work added in the backend table.")
	cmd.Flags().IntVar(&o.MinReconnectIntervalMs, "PG_MIN_RECONNECT_INTERVAL_MS", 100, " The min duration to wait before trying to re-establish the database connection after connection loss.")
	cmd.Flags().IntVar(&o.MaxReconnectIntervalMs, "PG_MAX_RECONNECT_INTERVAL_MS", 10000, " The max duration to wait before trying to re-establish the database connection after connection loss.")
}

func (o *Database) PgConfig() (*paasql.PGConfig, error) {
	port, err := strconv.Atoi(o.Port)
	if err != nil {
		return nil, fmt.Errorf("error converting port to int: %v", err)
	}

	return &paasql.PGConfig{
		Host:     o.Host,
		Port:     port,
		User:     o.User,
		Password: o.Password,
		Name:     o.Database,
		SSLMode:  o.SslMode,
	}, nil
}

func MustNewDB(cfg Database) *sqlx.DB {
	pgConfig, err := cfg.PgConfig()
	if err != nil {
		paaslog.Fatalf("error getting pg config: %v", err)
	}
	// migrations
	if err := migration.Exec(pgConfig); err != nil {
		paaslog.Fatalf("error initializing migrations: %v", err)
	}

	db, err := paasql.PostgresDB(pgConfig)
	if err != nil {
		paaslog.Fatalf("error initializing postgresql: %v", err)
	}

	metrics.MustInitMetrics(
		[]metrics.MonitoredDbConfig{
			{
				DbName: cfg.Database,
				DB:     db.SQLX(),
			},
		})

	configureDBConnectionPool(db.SQLX(), &cfg)

	return db.SQLX()
}

func configureDBConnectionPool(db *sqlx.DB, cfg *Database) {
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
}
