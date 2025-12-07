package notification

import (
	"context"
	"testing"

	"github.com/ionos-cloud/policies-service/internal/model"
)

func TestSMSNotifier_Notify_Success(t *testing.T) {
	notifier := &SMSNotifier{}
	user := &model.User{} // Minimal user; adjust fields as needed

	err := notifier.Notify(context.Background(), user, "test message")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
