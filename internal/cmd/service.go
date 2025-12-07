package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/router"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/server"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/policies-service/internal/config"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	Name string
	Cfg  config.Service

	wg sync.WaitGroup

	// main context of the service
	Ctx       context.Context
	ctxCancel context.CancelFunc

	// context used to start additional routines
	routineCtx       context.Context
	routineCtxCancel context.CancelFunc

	// make sure some operations are done one time only
	startHealthCheckServerOnce   sync.Once
	startObservabilityServerOnce sync.Once
	shutdownOnce                 sync.Once
}

func MustNewService(ctx context.Context, name string, cfg config.Service) *Service {
	svc, err := NewService(ctx, name, cfg)
	if err != nil {
		paaslog.Fatalf("create service %s error: %v", name, err)
	}

	return svc
}

func NewService(ctx context.Context, name string, cfg config.Service) (*Service, error) {
	serviceCtx, serviceCxtCancel := context.WithCancel(ctx)
	routineCtx, routineCtxCancel := context.WithCancel(ctx)

	svc := Service{
		Name:             name,
		Cfg:              cfg,
		wg:               sync.WaitGroup{},
		routineCtx:       routineCtx,
		routineCtxCancel: routineCtxCancel,
		Ctx:              serviceCtx,
		ctxCancel:        serviceCxtCancel,
	}

	// configure listening for termination signals which will trigger shutdown of the service
	ListenSignals(func(sig os.Signal) {
		svc.TriggerShutdown(fmt.Sprintf("got signal %s", sig))
	})

	return &svc, nil
}

// StartHealthCheckServer will start a health check server on the Cfg.HealthCheckAddr using the provided handler.
// The server will be started one time no matter how many times this function is called.
// The server will be automatically shutdown when the service is shutdown.
func (s *Service) StartHealthCheckServer(health healthcheck.Handler) {
	s.startHealthCheckServerOnce.Do(func() {
		srv := server.CreateHTTPServer(s.Cfg.HealthCheckAddr, health)

		s.StartRoutine("health check server", func(ctx context.Context, triggerShutdown func(reason string)) {
			if err := ListenAndServe(s.routineCtx, srv /*shutdownTimeout*/, 30*time.Second); err != nil {
				paaslog.Errorf("health check server got error : %v", err)
				triggerShutdown("health check server error")
			}
		})

		paaslog.Infof("started health check server on addr %s", s.Cfg.HealthCheckAddr)
	})
}

// StartObservabilityServer will start a metrics server on the Cfg.MetricsAddr.
// NOTE: The server will be started one time no matter how many times this function is called.
// The server will be automatically shutdown when the service is shutdown.
func (s *Service) StartObservabilityServer() {
	s.startObservabilityServerOnce.Do(func() {
		mux := router.NewObservabilityRouter()
		srv := server.CreateHTTPServer(s.Cfg.MetricsAddr, mux)
		// Set a very high write timeout to support long-running pprof profiling runs.
		srv.WriteTimeout = 60 * time.Second

		s.StartRoutine("observability server", func(ctx context.Context, triggerShutdown func(reason string)) {
			if err := ListenAndServe(s.routineCtx, srv /*shutdownTimeout*/, 30*time.Second); err != nil {
				paaslog.Errorf("observability server got error : %v", err)
				triggerShutdown("observability server error")
			}
		})

		paaslog.Infof("started observability server on port %d", s.Cfg.MetricsPort)
	})
}

// TriggerShutdown starts the shutdown procedure of the service. This will take care of gracefully shutting down the
// service and all goroutines started including health check and observability servers.
// The shutdown sequence will be done one time only.
func (s *Service) TriggerShutdown(reason string) {
	s.shutdownOnce.Do(func() {
		ctx := s.Ctx
		paaslog.InfoCf(ctx, "service %s shutdown triggered due to reason: `%s`", s.Name, reason)

		// first shutdown routines
		s.routineCtxCancel()

		// secondly shutdown service
		s.ctxCancel()
	})
}

// Wait will block until the service and all its goroutines have shutdown completely.
func (s *Service) Wait() {
	paaslog.InfoCf(s.Ctx, "service %s started, entering wait", s.Name)

	s.wg.Wait()

	paaslog.InfoCf(s.Ctx, "service %s wait done", s.Name)
}

// StartRoutine can be used to start additional go routines which will be managed by the service.
// The goroutines flow should be defined by the provided function parameter.
func (s *Service) StartRoutine(name string, f func(ctx context.Context, triggerShutdown func(reason string))) {
	s.wg.Add(1)
	go func() {
		defer func() {
			s.TriggerShutdown(fmt.Sprintf("routine: %s, reason: clean exit", name))
			s.wg.Done()
			paaslog.DebugCf(s.routineCtx, "routine shutdown done: %s", name)
		}()

		triggerShutdown := func(reason string) {
			s.TriggerShutdown(fmt.Sprintf("routine: %s, reason: %s", name, reason))
		}

		f(s.routineCtx, triggerShutdown)
	}()
}

// ListenSignals will listen to OS signals (currently SIGINT, SIGKILL, SIGTERM)
// and will trigger the callback when signal are received from OS.
func ListenSignals(callback func(signal os.Signal)) {
	go func() {
		// we use buffered to mitigate losing the signal
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

		sig := <-sigchan
		close(sigchan)
		if callback != nil {
			callback(sig)
		}
	}()
}

// ListenAndServe is a wrapper over http.ListenAndServe but bounds the lifetime of the http server to the given context
// and also handles the graceful shutdown of the server. It should be called in a goroutine because it is a blocking call.
// ListenAndServe will return when the server will be stopped either because the given context was cancelled or it had
// an error. The error which lead to the shutdown of the server will be returned.
func ListenAndServe(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	paaslog.InfoCf(ctx, "starting server on addr: %s", server.Addr)
	serverErr := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// When Shutdown is called, [Serve], [ListenAndServe], and
			// [ListenAndServeTLS] immediately return [ErrServerClosed]. Make sure the
			// program doesn't exit and waits instead for Shutdown to return.

			paaslog.WarnCf(ctx, "ListenAndServe on addr: %s got error: %v", server.Addr, err)
			serverErr <- err
		}

		// ListenAndServe finished, write nil to unblock
		serverErr <- nil
	}()

	var err error
	select {
	case <-ctx.Done():
		// TODO: is it ok to use ctx which is cancelled as parent ???
		ctxShutdown, cancelCtxShutdown := context.WithTimeout(ctx, shutdownTimeout)
		defer cancelCtxShutdown()
		err := server.Shutdown(ctxShutdown)
		// Shutdown returns the context's error, otherwise it returns any error returned from closing the Server's underlying Listener(s).
		if err == nil {
			paaslog.InfoCf(ctxShutdown, "shutdown server on addr: %s was successfull", server.Addr)
		} else if errors.Is(err, ctx.Err()) {
			//  context's error
			paaslog.WarnCf(ctx, "ListenAndServe on addr: %s stopped because of context error: %v", server.Addr, ctx.Err())
		} else {
			// error returned from closing the [Server]'s underlying Listener(s).
			paaslog.ErrorCf(ctxShutdown, "shutdown server on addr: %s got error: %v", server.Addr, err)
		}
	case err = <-serverErr:
		return err
	}

	return err
}

func ConfigureHealthCheckHandler(health healthcheck.Handler, db *sqlx.DB, cfg config.Service) error {
	if cfg.HealthCheckMaxAllowedGoroutines <= 0 {
		return fmt.Errorf("health-check-max-allowed-goroutines must be greater than 0")
	}
	// It's fine to use a higher goroutine number here because v2 worker runs on multiple concurrent goroutines
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(cfg.HealthCheckMaxAllowedGoroutines))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	// This checks we can get a connection to the db in a reasonable time. 3 sec might look too much but the system is
	// still responsive even with 3 sec db timeout, TODO: consider lowering it once performance is stable.
	if cfg.HealthCheckDbPingTimeoutSec == 0 {
		return fmt.Errorf("health-check-db-ping-timeout-sec must be greater than 0")
	}
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db.DB, cfg.HealthCheckDbPingTimeoutSec))

	return nil
}
