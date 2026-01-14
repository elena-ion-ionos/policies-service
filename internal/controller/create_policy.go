package controller

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/ionos-cloud/policies-service/internal/port"
)

type CreatePolicy struct {
	repo port.PolicyRepo
}

func NewCreatePolicyCtrl(repo port.PolicyRepo) (*CreatePolicy, error) {
	return &CreatePolicy{repo: repo}, nil
}

func (s *CreatePolicy) RegisterPolicy(ctx context.Context, policy *model.Policy) error {
	policy.ID = uuid.New()
	policy.CreatedAt = time.Now()
	if err := s.repo.Save(ctx, policy); err != nil {
		return err
	}
	return nil
}
