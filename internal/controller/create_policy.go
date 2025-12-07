package controller

import (
	"context"

	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/ionos-cloud/policies-service/internal/port"
)

type RegisterRule interface {
	// RegisterRule registers a new user and notifies them.
	RegisterRule(ctx context.Context, user *model.Policy) error
}

type CreatePolicy struct {
	repo      port.PolicyRepo
	notifiers []port.Notifier
}

func NewPolicyCtrl(repo port.PolicyRepo, notifiers ...port.Notifier) (*CreatePolicy, error) {
	return &CreatePolicy{repo: repo, notifiers: notifiers}, nil
}

func (s *CreatePolicy) RegisterRule(ctx context.Context, rule *model.Policy) error {
	if err := s.repo.Save(ctx, rule); err != nil {
		return err
	}
	return nil
}
