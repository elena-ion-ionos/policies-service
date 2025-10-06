package service

import (
	"net/http"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"
	userapi "github.com/ionos-cloud/go-sample-service/internal/api"
	"github.com/ionos-cloud/go-sample-service/internal/config"
	"github.com/ionos-cloud/go-sample-service/internal/model"
)

var _ userapi.ServerInterface = (*UserApi)(nil)

type UserApi struct {
	Helper crud.ReaderWriter[
		model.User,
		*userapi.Metadata,
		userapi.User,
	]
}

func (u UserApi) PostUsers(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

func (u UserApi) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me
	panic("implement me")
}

func (u UserApi) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me
	panic("implement me")
}

func (u UserApi) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me
	panic("implement me")
}

func MustNewWebServerUser(cfg *config.Service) *UserApi {
	s := &UserApi{}

	return s
}
