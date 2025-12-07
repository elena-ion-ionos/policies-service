package model

import (
	"time"

	"github.com/ionos-cloud/go-paaskit/api/paastype"
)

type Metadata struct {
	// CreatedBy Unique name of the identity that created the resource.
	CreatedBy *string

	// CreatedByUserId Unique id of the identity that created the resource.
	CreatedByUserId *string

	// CreatedDate The ISO 8601 creation timestamp.
	CreatedDate *time.Time

	// LastModifiedBy Unique name of the identity that last modified the resource.
	LastModifiedBy *string

	// LastModifiedByUserId Unique id of the identity that last modified the resource.
	LastModifiedByUserId *string

	// LastModifiedDate The ISO 8601 modified timestamp.
	LastModifiedDate *time.Time
}

var _ paastype.Metadata = (*Metadata)(nil)

func (m *Metadata) SetCreated(createdBy *string, createdByUserId *string, createdDate *time.Time) {
	m.CreatedBy = createdBy
	m.CreatedByUserId = createdByUserId
	m.CreatedDate = createdDate
}

func (m *Metadata) SetLastModified(lastModifiedBy *string, lastModifiedByUserId *string, lastModifiedDate *time.Time) {
	m.LastModifiedBy = lastModifiedBy
	m.LastModifiedByUserId = lastModifiedByUserId
	m.LastModifiedDate = lastModifiedDate
}

func (m *Metadata) SetResourceURN(urn *string) {
	// not implemented
}
