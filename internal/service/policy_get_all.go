package service

import (
	"context"
	"net/http"

	"github.com/ionos-cloud/go-paaskit/api/paashttp"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/policies-service/internal/model"
)

func (u PoliciesApi) GetPolicies(w http.ResponseWriter, r *http.Request) {
	paashttp.Handle("GetPolicies All Policies", w, r, func(ctx context.Context) error {
		policies, err := u.GetPolicyController.GetPolicies(ctx)
		paaslog.InfoCf(ctx, "get all policies")
		rules := make([]*model.Policy, 0)
		if err != nil {
			return err
		}
		for _, policy := range policies {
			rules = append(rules, policy)
		}

		return u.Helper.WriteMany(ctx, w, rules, crud.NewPaginationFromRequest(r))
	})
}
