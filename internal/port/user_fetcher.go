package port

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"
)

type UserFetcher interface {
	Fetch(ctx context.Context) (*model.User, error)
}
