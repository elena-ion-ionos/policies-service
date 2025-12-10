package controller

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/ionos-cloud/policies-service/internal/port"
)

//de sters
//type RegisterRule interface {
//	// RegisterRule registers a new user and notifies them.
//	RegisterPolicy(ctx context.Context, user *model.Policy) error
//}

type CreatePolicy struct {
	repo      port.PolicyRepo
	notifiers []port.Notifier
}

func NewCreatePolicyCtrl(repo port.PolicyRepo, notifiers ...port.Notifier) (*CreatePolicy, error) {
	return &CreatePolicy{repo: repo, notifiers: notifiers}, nil
}

func (s *CreatePolicy) RegisterPolicy(ctx context.Context, policy *model.Policy) error {
	policy.ID = uuid.New()
	policy.CreatedAt = time.Now()
	if err := s.repo.Save(ctx, policy); err != nil {
		return err
	}
	return nil
}
