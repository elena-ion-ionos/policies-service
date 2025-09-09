package port

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *model.User) error
}
