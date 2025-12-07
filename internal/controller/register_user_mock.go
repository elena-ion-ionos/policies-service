package controller

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/go-paaskit/service/contract"
	"github.com/ionos-cloud/policies-service/internal/model"
)

type KeysCleanerMock struct {
	OnCleanUserKeys     func(ctx context.Context, user *model.User) error
	OnCleanContractKeys func(ctx context.Context, contractNumber contract.Number) error
}

// CleanUserKeys delete all user keys
func (c *KeysCleanerMock) CleanUserKeys(ctx context.Context, user *model.User) error {
	if c.OnCleanUserKeys == nil {
		return fmt.Errorf("OnCleanUserKeys not set")
	}

	return c.OnCleanUserKeys(ctx, user)
}

// CleanContractKeys delete all contract keys
func (c *KeysCleanerMock) CleanContractKeys(ctx context.Context, contractNumber contract.Number) error {
	if c.OnCleanContractKeys == nil {
		return fmt.Errorf("OnCleanContractKeys not set")
	}

	return c.OnCleanContractKeys(ctx, contractNumber)
}
