package model

import (
	"github.com/ionos-cloud/go-paaskit/service/contract"

	"github.com/google/uuid"
)

type User struct {
	ContractNumber contract.Number
	UserID         uuid.UUID
	Email          string
	Phone          string
}
