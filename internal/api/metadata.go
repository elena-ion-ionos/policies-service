package api

import (
	"time"

	"github.com/ionos-cloud/go-paaskit/api/paastype"
)

var _ paastype.Metadata = (*Metadata)(nil)

func (m *Metadata) SetCreated(createdBy *string, createdByUserId *string, createdDate *time.Time) {
	//m.CreatedBy = createdBy
	//m.CreatedByUserId = createdByUserId
	//m.CreatedDate = createdDate
}

func (m *Metadata) SetLastModified(lastModifiedBy *string, lastModifiedByUserId *string, lastModifiedDate *time.Time) {
	//m.LastModifiedBy = lastModifiedBy
	//m.LastModifiedByUserId = lastModifiedByUserId
	//m.LastModifiedDate = lastModifiedDate
}

func (m *Metadata) SetResourceURN(urn *string) {
	//m.ResourceURN = urn
}
