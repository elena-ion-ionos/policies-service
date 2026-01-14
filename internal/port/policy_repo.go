package port

import (
	"context"
	"github.com/ionos-cloud/policies-service/internal/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

//go:generate mockgen -source=policy_repo.go -destination=../mocks/port/mock_policy_repo.go -package=mocks

// portul e o interfata ceva abstract care va definii operatii de write and read ce va comunica cu surse externe: ex:
// Ex: Poate sa scrie sau sa citeasca din fisiere, sau din baze de date
// interfata ce defineste operatiile cu baza de date
type PolicyRepo interface {
	Save(ctx context.Context, policy *model.Policy) error
	GetPolicies(ctx context.Context) ([]*model.Policy, error)
	GetPolicyById(ctx context.Context, id openapi_types.UUID) (*model.Policy, error)
	DeletePolicyById(ctx context.Context, id string) error
}
