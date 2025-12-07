package notification

import (
	"context"

	"github.com/ionos-cloud/policies-service/internal/model"
)

type SMSNotifier struct{}

func (n *SMSNotifier) Notify(ctx context.Context, user *model.User, message string) error {
	// Send SMS (omitted)
	return nil
}
