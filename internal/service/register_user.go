package service

import (
	"context"
	"time"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/config"
	"github.com/ionos-cloud/go-sample-service/internal/controller"
	"github.com/ionos-cloud/go-sample-service/internal/port"
)

type RegisterUser interface {
	Serve(ctx context.Context) error
	ServeOne(ctx context.Context) (bool, error)
}
type registerUserImpl struct {
	fetcher port.UserFetcher
	ctrl    controller.RegisterUser
	cfg     *config.RegisterUser
}

func MustNewRegisterUser(cfg *config.RegisterUser, fetcher port.UserFetcher, ctrl controller.RegisterUser) *registerUserImpl {
	s := &registerUserImpl{
		cfg:     cfg,
		fetcher: fetcher,
		ctrl:    ctrl,
	}

	return s
}

// Serve starts the background worker for key management, it
// checks the backend key queue for non-active backend keys and
// manages them.
func (s *registerUserImpl) Serve(ctx context.Context) error {
	paaslog.InfoCf(ctx, "starting to handle key operations")

	for {
		stop, err := s.ServeOne(ctx)
		if err != nil {
			paaslog.ErrorCf(ctx, "ServeOne failed, err: %v", err)
		}

		if stop {
			return err
		}

		readOpsTimeoutMs := time.Duration(10) * time.Millisecond // default read ops timeout
		// continue to a new batch only if one of the conditions below is satisfied
		select {
		case <-ctx.Done(): // in case of context cancellation break the loop
			return ctx.Err()
		case <-time.After(readOpsTimeoutMs * time.Millisecond): // get more work every configured interval
		}
	}
}

// Return true in case it should stop, false in case it should continue
func (s *registerUserImpl) ServeOne(ctx context.Context) (bool, error) {
	user, err := s.fetcher.Fetch(ctx)
	if err != nil {
		paaslog.ErrorCf(ctx, "failed to read a new batch of keys, err: %v", err)
		return false, err
	}

	err = s.ctrl.RegisterUser(ctx, user)
	if err != nil {
		paaslog.ErrorCf(ctx, "cuid updater service err: %v, user: %v", err, user)
	}

	return false, nil
}
