package port

import (
	"context"

	"github.com/ionos-cloud/policies-service/internal/model"
)

type Notifier interface {
	Notify(ctx context.Context, policy *model.Policy, message string) error
}
