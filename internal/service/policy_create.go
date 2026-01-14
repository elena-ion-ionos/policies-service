package service

import (
	"context"
	"github.com/ionos-cloud/go-paaskit/api/paashttp"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"
	"net/http"

	"github.com/ionos-cloud/policies-service/internal/model"
)

func (l *PoliciesApi) PostPolicies(w http.ResponseWriter, r *http.Request) {
	paashttp.Handle("GetPolicies All Policies", w, r, func(ctx context.Context) error {
		policy, err := l.loadRequestBody(ctx, r)
		if err != nil {
			return err
		}
		err = l.CreatePolicyController.RegisterPolicy(ctx, policy)
		if err != nil {
			return err
		}
		return l.Helper.WriteOne(ctx, w, policy, crud.WithStatusCode(http.StatusCreated))
	})
}

func (l *PoliciesApi) loadRequestBody(ctx context.Context, r *http.Request) (*model.Policy, error) {
	if r.Body == http.NoBody {
		return &model.Policy{}, nil
	}
	requestedPolicy, err := l.Helper.ReadOne(ctx, r)
	err = checkError(err)
	if err != nil {
		return nil, err
	}

	return requestedPolicy, nil
}
