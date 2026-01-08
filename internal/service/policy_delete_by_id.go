package service

import (
	"context"
	"github.com/ionos-cloud/go-paaskit/api/paashttp"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"net/http"
)

func (u PoliciesApi) DeletePoliciesId(w http.ResponseWriter, r *http.Request, id string) {
	paashttp.Handle("Delete policy", w, r, func(ctx context.Context) error {
		err := u.DeletePolicyByIdController.DeletePolicyById(ctx, id)
		paaslog.InfoCf(ctx, "delete policy by id")
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}
