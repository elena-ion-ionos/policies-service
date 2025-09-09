package postgres

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepoImpl {
	return &userRepoImpl{db: db}
}

func (r *userRepoImpl) Save(ctx context.Context, user *model.User) error {
	// Insert user into DB (omitted for brevity)
	return nil
}

func (r *userRepoImpl) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	// Query user from DB (omitted for brevity)
	return &model.User{
		UserID: userID,
	}, nil
}
