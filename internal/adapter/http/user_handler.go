package http

import (
	"encoding/json"
	"net/http"

	"github.com/ionos-cloud/go-sample-service/internal/model"
	"github.com/ionos-cloud/go-sample-service/internal/port"
)

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.service.RegisterUser(r.Context(), &user); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
