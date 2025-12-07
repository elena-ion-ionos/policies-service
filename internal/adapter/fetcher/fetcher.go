package fetcher

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/policies-service/internal/model"
)

type fetcherImpl struct{}

// NewFetcher creates a new instance of FetcherImpl.
func NewFetcher() *fetcherImpl {
	return &fetcherImpl{}
}

func (r *fetcherImpl) Fetch(ctx context.Context) (*model.Policy, error) {
	// This is a placeholder implementation.
	return nil, fmt.Errorf("not implemented")
}
