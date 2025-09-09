package migration

import (
	"embed"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ionos-cloud/go-paaskit/infrastructure/paasql"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
)

//go:embed postgresql/*.sql
var fs embed.FS

func migrations() (source.Driver, error) {
	return iofs.New(fs, "postgresql")
}

// Exec executes all DB migrations based on provided PG configuration
func Exec(cfg *paasql.PGConfig) error {
	m, err := migrations()
	if err != nil {
		paaslog.Infof("error during migrations: %v", err)
		return err
	}
	paaslog.Infof("migrating database")

	err = paasql.MigratePostgresDB(cfg, m)
	if err != nil {
		paaslog.Infof("error during migrations: %v", err)
		return err
	}
	paaslog.Infof("migrated database")
	return err
}
