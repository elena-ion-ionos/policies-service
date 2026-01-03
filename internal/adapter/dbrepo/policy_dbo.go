package dbrepo

import (
	"github.com/google/uuid"
	"github.com/ionos-cloud/policies-service/internal/model"
	"time"
)

type PolicyDBO struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Action    string    `db:"action"`
	Time      string    `db:"time"`
	Prefix    string    `db:"prefix"`
	CreatedAt time.Time `db:"created_at"`
}

func NewPolicyFromPolicyDBO(keyDBO PolicyDBO) model.Policy {
	return model.Policy{
		ID:        keyDBO.ID,
		Name:      keyDBO.Name,
		Prefix:    keyDBO.Prefix,
		Action:    keyDBO.Action,
		Time:      keyDBO.Time,
		CreatedAt: keyDBO.CreatedAt,
	}
}
