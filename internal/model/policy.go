package model

import (
	"github.com/google/uuid"
	"time"
)

type Policy struct {
	ID        uuid.UUID
	Name      string
	Prefix    string
	Action    string
	Time      string
	CreatedAt time.Time
}
