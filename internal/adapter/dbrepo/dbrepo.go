package dbrepo

import (
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/policies-service/internal/config"
	"github.com/ionos-cloud/policies-service/internal/port"
	"github.com/jmoiron/sqlx"
)

//asta e adapterul ce va implementa portul pentru a face legatura cu baza de date

type PolicyRepoImpl struct {
	DB *sqlx.DB
	//DBKeys         port.DBKeys
	//DBKeysOpsQueue port.DBKeysOpsQueue
}

func NewPolicyRepo(db *sqlx.DB) *PolicyRepoImpl {
	return &PolicyRepoImpl{DB: db}
}

func CreateFromConn(db *sqlx.DB) (*sqlx.DB, port.PolicyRepo) {
	paaslog.Infof("creating db repos ")
	policiesRepo := NewPolicyRepo(db)

	return db, policiesRepo
}

func MustCreateFromConfig(cfg config.Database) (*sqlx.DB, port.PolicyRepo) {
	paaslog.Infof("creating db repos ")

	db := config.MustNewDB(cfg)

	return CreateFromConn(db)
}
