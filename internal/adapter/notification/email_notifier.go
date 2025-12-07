package notification

import (
	"context"

	"github.com/ionos-cloud/policies-service/internal/model"
)

type EmailNotifier struct{}

func (n *EmailNotifier) Notify(ctx context.Context, user *model.User, message string) error {
	// Send email (omitted)
	return nil
}
