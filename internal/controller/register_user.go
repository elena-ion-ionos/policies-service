package controller

import (
	"context"

	"github.com/ionos-cloud/go-sample-service/internal/model"
	"github.com/ionos-cloud/go-sample-service/internal/port"
)

type RegisterUser interface {
	// RegisterUser registers a new user and notifies them.
	RegisterUser(ctx context.Context, user *model.User) error
}

type registerUserImpl struct {
	repo      port.UserRepository
	notifiers []port.Notifier
}

func NewRegisterUser(repo port.UserRepository, notifiers ...port.Notifier) (*registerUserImpl, error) {
	return &registerUserImpl{repo: repo, notifiers: notifiers}, nil
}

func (s *registerUserImpl) RegisterUser(ctx context.Context, user *model.User) error {
	if err := s.repo.Save(ctx, user); err != nil {
		return err
	}
	for _, n := range s.notifiers {
		_ = n.Notify(ctx, user, "Welcome!")
	}
	return nil
}
