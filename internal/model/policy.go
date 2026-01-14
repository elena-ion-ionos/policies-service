package model

import (
	"github.com/google/uuid"
	"time"
)

type Policy struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	Prefix    string    `json:"prefix"`
	Action    string    `json:"action"`
	Time      string    `json:"time"`
	Metadata  *Metadata `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
