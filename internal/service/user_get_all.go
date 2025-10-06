package service

import (
	"context"
	"net/http"

	"github.com/ionos-cloud/go-paaskit/api/paashttp"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/model"
)

func (u UserApi) GetUsers(w http.ResponseWriter, r *http.Request) {
	paashttp.Handle("Get All Keys", w, r, func(ctx context.Context) error {
		paaslog.InfoCf(ctx, "get all keys as admin")
		users := make([]*model.User, 0)

		return u.Helper.WriteMany(ctx, w, users, crud.NewPaginationFromRequest(r))
	})
}
