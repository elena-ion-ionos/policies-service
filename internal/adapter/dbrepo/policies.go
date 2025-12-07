package dbrepo

import (
	"context"
	"database/sql"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/policies-service/internal/metrics"
	"github.com/ionos-cloud/policies-service/internal/model"
)

const (
	createLifecycleRuleQuery = "INSERT INTO policies (name, prefix, action, time ) VALUES ($1 , $2, $3, $4)"
)

func txRollback(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		paaslog.ErrorCf(ctx, "Error rolling back transaction: %v", err)
	}
}

func (p PolicyRepoImpl) Save(ctx context.Context, policy *model.Policy) error {
	tx, err := p.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		paaslog.ErrorCf(ctx, "Error starting transaction: %v", err)
		metrics.OpsNo.WithLabelValues(metrics.LabelOperation_SaveUser, metrics.LabelLocation_Repo, metrics.LabelResult_fail).Inc()
		return err
	}

	_, err = createPolicy(ctx, tx, policy)
	if err != nil {
		defer txRollback(ctx, tx)
		paaslog.ErrorCf(ctx, "Error creating policy: %v", err)
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
