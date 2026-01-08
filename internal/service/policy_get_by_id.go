package service

import (
	"context"
	"github.com/ionos-cloud/go-paaskit/api/paashttp"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (u PoliciesApi) GetPoliciesId(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	paashttp.Handle("GetPolicies All Policies", w, r, func(ctx context.Context) error {
		policy, err := u.GetPolicyByIdController.GetPolicyById(ctx, id)
		paaslog.InfoCf(ctx, "get policy by id")
		err = checkError(err)
		if err != nil {
			return err
		}

		return u.Helper.WriteOne(ctx, w, policy)
	})
}
