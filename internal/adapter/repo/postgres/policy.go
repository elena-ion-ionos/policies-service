package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	_ "github.com/ionos-cloud/go-paaskit/service/contract"
	"github.com/ionos-cloud/policies-service/internal/metrics"
	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserDBO struct {
	Phone          string     `db:"phone"`
	ContractNumber string     `db:"contract_number"`
	ID             uuid.UUID  `db:"id"`
	Email          string     `db:"email"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

const (
	createLifecycleRuleQuery = "INSERT INTO policies (name, prefix, action, time ) VALUES ($1 , $2, $3, $4, $5)"
)

type lifecycleRuleRepo struct {
	db *sqlx.DB
}

func NewLifeCycleRule(db *sqlx.DB) *lifecycleRuleRepo {
	return &lifecycleRuleRepo{db: db}
}

func txRollback(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		paaslog.ErrorCf(ctx, "Error rolling back transaction: %v", err)
	}
}

func (r *lifecycleRuleRepo) Save(ctx context.Context, user *model.Policy) error {
	timer := prometheus.NewTimer(metrics.OpsDurationSeconds.WithLabelValues(metrics.LabelOperation_SaveUser, metrics.LabelLocation_Repo))
	defer timer.ObserveDuration()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		paaslog.ErrorCf(ctx, "Error starting transaction: %v", err)
		metrics.OpsNo.WithLabelValues(metrics.LabelOperation_SaveUser, metrics.LabelLocation_Repo, metrics.LabelResult_fail).Inc()
		return err
	}
	defer txRollback(ctx, tx)

	_, err = createPolicy(ctx, tx, user)
	if err != nil {
		paaslog.ErrorCf(ctx, "Error creating user: %v", err)
		metrics.OpsNo.WithLabelValues(metrics.LabelOperation_SaveUser, metrics.LabelLocation_Repo, metrics.LabelResult_fail).Inc()
		return err
	}
	metrics.OpsNo.WithLabelValues(metrics.LabelOperation_SaveUser, metrics.LabelLocation_Repo, metrics.LabelResult_success).Inc()
	return tx.Commit()
}

func createPolicy(ctx context.Context, tx *sql.Tx, policy *model.Policy) (*model.Policy, error) {
	_, err := tx.ExecContext(ctx, createLifecycleRuleQuery,
		policy.Name, policy.Prefix, policy.Action, policy.Time)
	return nil, err
}
