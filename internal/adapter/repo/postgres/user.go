package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-paaskit/service/contract"
	"github.com/ionos-cloud/go-sample-service/internal/model"

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
	createUserQuery     = "INSERT INTO users (id, contract_number, phone, email ) VALUES ($1 , $2, $3, $4)"
	getUserByIdQueryStr = "SELECT * FROM users WHERE id=$1"
)

type userRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepoImpl {
	return &userRepoImpl{db: db}
}

func txRollback(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		paaslog.ErrorCf(ctx, "Error rolling back transaction: %v", err)
	}
}

func (r *userRepoImpl) Save(ctx context.Context, user *model.User) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		paaslog.ErrorCf(ctx, "Error starting transaction: %v", err)
		return err
	}
	defer txRollback(ctx, tx)

	_, err = createKeyWTx(ctx, tx, user)
	if err != nil {
		paaslog.ErrorCf(ctx, "Error creating user: %v", err)
		return err
	}
	return tx.Commit()
}

func createKeyWTx(ctx context.Context, tx *sql.Tx, key *model.User) (*model.User, error) {
	_, err := tx.ExecContext(ctx, createUserQuery,
		key.UserID, key.ContractNumber, key.Phone, key.Email)
	return nil, err
}

func (r *userRepoImpl) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	// Query user from DB (omitted for brevity)
	var dboUser UserDBO

	err := r.db.GetContext(ctx, &dboUser, getUserByIdQueryStr, userID)
	if err != nil {
		paaslog.ErrorCf(ctx, "Error querying user by ID: %v", err)
		return nil, err
	}

	return &model.User{
		UserID:         dboUser.ID,
		Phone:          dboUser.Phone,
		Email:          dboUser.Email,
		ContractNumber: contract.MustParse(dboUser.ContractNumber),
	}, nil
}
