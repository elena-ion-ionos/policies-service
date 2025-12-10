package api

import (
	"github.com/ionos-cloud/policies-service/internal/model"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"
	"github.com/ionos-cloud/go-paaskit/api/paastype"
)

type (
	PolicyApi    = paastype.Item[*Metadata, Policy]
	ReaderWriter = crud.ReaderWriter[model.Policy, *Metadata, Policy]
)

type PolicyConverter struct {
	Host string
}

func (c *PolicyConverter) withHost(Host string) {
	c.Host = Host
}

// Implement the required methods for the Converter interface

func (c *PolicyConverter) ConvertToModel(a *PolicyApi, m *model.Policy) error {
	m.Name = a.Properties.Name
	m.Action = a.Properties.Action
	m.Prefix = a.Properties.Prefix
	m.Time = a.Properties.Time
	return nil
}

func (c *PolicyConverter) ConvertToApi(m *model.Policy, a *PolicyApi) error {
	a.Properties.Name = m.Name
	a.Properties.Action = m.Action
	a.Properties.Prefix = m.Prefix
	a.Properties.Time = m.Time
	return nil
}

func NewReaderWriter(host string) ReaderWriter {
	converter := new(PolicyConverter)
	converter.withHost(host)
	return crud.NewReaderWriter(converter)
}
