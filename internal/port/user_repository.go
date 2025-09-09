package port

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}
