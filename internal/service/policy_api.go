package service

import (
	"errors"
	"github.com/ionos-cloud/go-paaskit/api/paastype"
	policiesApi "github.com/ionos-cloud/policies-service/internal/api"
	"github.com/ionos-cloud/policies-service/internal/controller"
	"github.com/ionos-cloud/policies-service/internal/port"
	"net/http"

	"github.com/ionos-cloud/go-paaskit/api/paashttp/crud"

	"github.com/ionos-cloud/policies-service/internal/config"
	"github.com/ionos-cloud/policies-service/internal/model"
)

// ELENA
// The selected line is a compile-time assertion in Go.
// It checks that the PoliciesApi type implements the ServerInterface interface from the policiesApi package.
// If PoliciesApi does not implement all required methods of ServerInterface, the code will fail to compile,
// helping catch interface implementation errors early.
// Compile-time check for interface implementation

var _ policiesApi.ServerInterface = (*PoliciesApi)(nil)

// PoliciesApi aggregates controllers and helpers for policy operations
type PoliciesApi struct {
	CreatePolicyController     *controller.CreatePolicy
	GetPolicyController        *controller.GetPolicies
	GetPolicyByIdController    *controller.GetPolicyById
	DeletePolicyByIdController *controller.DeletePolicyById
	Helper                     crud.ReaderWriter[
		model.Policy,
		*policiesApi.Metadata,
		policiesApi.Policy,
	]
	ServerHost string
}

func (l PoliciesApi) PutPoliciesId(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

// aici trebuie sa pasez parametri pentru controller
func MustNewWebServerUser(cfg *config.Service, serverHost string, repo port.PolicyRepo, notifier port.Notifier) *PoliciesApi {
	createPolicyCtrl, err := controller.NewCreatePolicyCtrl(repo)
	getPolicyCtrl, err := controller.NewGetPoliciesCtrl(repo, notifier)
	getPolicyByIDCtrl, err := controller.NewGetPolicyByIdCtrl(repo, notifier)
	deletePolicyByIdCtrl, err := controller.NewDeletePolicyByIdCtrl(repo, notifier)
	if err != nil {

	}
	s := &PoliciesApi{
		CreatePolicyController:     createPolicyCtrl,
		GetPolicyController:        getPolicyCtrl,
		GetPolicyByIdController:    getPolicyByIDCtrl,
		DeletePolicyByIdController: deletePolicyByIdCtrl,
		Helper:                     policiesApi.NewReaderWriter(serverHost),
		ServerHost:                 "",
	}

	return s
}

func checkError(err error) error {
	//todo elena sa vedem daca aici se face handler la get by id
	switch err := err.(type) {
	case nil:
		return nil
	case *paastype.Error:
		return paastype.NewError(err.HTTPStatus, err.PrimaryErrorCode(), err.Error())
	default:
		switch {
		case errors.Is(err, model.ErrNotFound), errors.Is(err, model.ErrIdNotFound): // policy not exists
			return paastype.NewError(404, "policy-not-found", err.Error())
		case errors.Is(err, model.ErrLimitExceeded): // limit exceeded
			return paastype.NewError(429, "too-many-requests", err.Error())
		case errors.Is(err, model.ErrUnprocessableEntity): // unprocessable entity / access denied
			return paastype.NewError(403, "access-denied", err.Error())
		case errors.Is(err, model.ErrFilterInvalid): // filter invalid
			return paastype.NewError(400, "filter-invalid", err.Error())
		default:
			return paastype.NewError(500, "internal-server-error", err.Error())
		}
	}
}
