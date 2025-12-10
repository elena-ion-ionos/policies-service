package controller

import (
	"context"
	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/ionos-cloud/policies-service/internal/port"
)

//type RegisterRule interface {
//	// RegisterRule registers a new user and notifies them.
//	RegisterPolicy(ctx context.Context, user *model.Policy) error
//}

type GetPolicies struct {
	repo      port.PolicyRepo
	notifiers []port.Notifier
}

func NewGetPoliciesCtrl(repo port.PolicyRepo, notifiers ...port.Notifier) (*GetPolicies, error) {
	return &GetPolicies{repo: repo, notifiers: notifiers}, nil
}

func (s *GetPolicies) GetPolicies(ctx context.Context) ([]*model.Policy, error) {
	policies, err := s.repo.Get(ctx)
	return policies, err
}
