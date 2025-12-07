package port

import (
	"context"
	"github.com/ionos-cloud/policies-service/internal/model"
)

// interfata ce defineste operatiile cu baza de date
type PolicyRepo interface {
	Save(ctx context.Context, policy *model.Policy) error
}
