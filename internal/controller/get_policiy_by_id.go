package controller

import (
	"context"
	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/ionos-cloud/policies-service/internal/port"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

//type RegisterRule interface {
//	// RegisterRule registers a new user and notifies them.
//	RegisterPolicy(ctx context.Context, user *model.Policy) error
//}

type GetPolicyById struct {
	repo      port.PolicyRepo
	notifiers []port.Notifier
}

func NewGetPolicyByIdCtrl(repo port.PolicyRepo, notifiers ...port.Notifier) (*GetPolicyById, error) {
	return &GetPolicyById{repo: repo, notifiers: notifiers}, nil
}

func (s *GetPolicyById) GetPolicyById(ctx context.Context, id openapi_types.UUID) (*model.Policy, error) {
	policies, err := s.repo.GetPolicyById(ctx, id)
	return policies, err
}
