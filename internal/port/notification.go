package port

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"
)

type Notifier interface {
	Notify(ctx context.Context, user *model.User, message string) error
}
