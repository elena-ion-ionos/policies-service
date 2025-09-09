package notification

import (
	"context"
	"testing"

	"github.com/ionos-cloud/go-sample-service/internal/model"
)

func TestEmailNotifier_Notify_Success(t *testing.T) {
	notifier := &EmailNotifier{}
	user := &model.User{} // Use a minimal user; adjust fields as needed

	err := notifier.Notify(context.Background(), user, "test message")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
