package controller

import (
	"context"
	"github.com/ionos-cloud/policies-service/internal/port"
)

//type RegisterRule interface {
//	// RegisterRule registers a new user and notifies them.
//	RegisterPolicy(ctx context.Context, user *model.Policy) error
//}

type DeletePolicyById struct {
	repo      port.PolicyRepo
	notifiers []port.Notifier
}

func NewDeletePolicyByIdCtrl(repo port.PolicyRepo, notifiers ...port.Notifier) (*DeletePolicyById, error) {
	return &DeletePolicyById{repo: repo, notifiers: notifiers}, nil
}

func (s *DeletePolicyById) DeletePolicyById(ctx context.Context, id string) error {
	err := s.repo.DeletePolicyById(ctx, id)
	return err
}
