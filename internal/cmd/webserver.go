package cmd

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/policies-service/internal/adapter/dbrepo"
	policiesApi "github.com/ionos-cloud/policies-service/internal/api"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/heptiolabs/healthcheck"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/account"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/activity"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/cors"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/health"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/identity"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/metric"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/router"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/server"
	"github.com/ionos-cloud/go-paaskit/api/paashttp/tracing"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	//"github.com/ionos-cloud/policies-service/internal/api"
	"github.com/ionos-cloud/policies-service/internal/config"
	"github.com/ionos-cloud/policies-service/internal/service"
	"github.com/spf13/cobra"
)

var WebserverUsername = "webserver"

func WebserverUser() *cobra.Command {
	var cfg config.Webserver

	cmd := &cobra.Command{
		Use:   WebserverUsername,
		Short: "starts the worker service for handling key operations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.InitViperFlags(cmd, os.Args); err != nil {
				paaslog.Fatalf("got error initializing config flags, err: %v", err)
			}

			if err := RunWebserverUser(context.Background(), &cfg); err != nil {
				paaslog.Fatalf("got error while running worker, err: %v", err)
			}
		},
	}

	cfg.AddFlags(cmd)

	return cmd
}

func RunWebserverUser(ctx context.Context, cfg *config.Webserver) error {
	svc := MustNewService(ctx, WebserverUsername, cfg.Service)

	dbConn := config.MustNewDB(cfg.Service.Database)
	_, policyRepo := dbrepo.MustCreateFromConfig(cfg.Service.Database)
	//trebuie sa mi creez repoul de database
	// ctrl, _ := newRegisterLifecycleCtrl(cfg, dbConn)
	registerUserService := service.MustNewWebServerUser(&cfg.Service, cfg.ServerHost, policyRepo, nil)

	svc.StartObservabilityServer()

	health := healthcheck.NewHandler()
	if err := ConfigureHealthCheckHandler(health, dbConn, cfg.Service); err != nil {
		paaslog.Fatalf("failed to configure health check handler: %v", err)
	}

	svc.StartHealthCheckServer(health)

	svc.StartRoutine(WebserverUsername, func(ctx context.Context, triggerShutdown func(reason string)) {
		httpServer, err := CreateHttpServer(cfg, registerUserService)
		if err != nil {
			paaslog.Fatalf("failed creating http server, error: %v", err)
		}

		if err := ListenAndServe(svc.Ctx, httpServer, 30*time.Second); err != nil {
			paaslog.ErrorCf(ctx, "%s Serve had error: %v, terminating", WebserverUsername, err)
			triggerShutdown(fmt.Sprintf("%s Serve had error", WebserverUsername))
		}
	})
	// this will block until a term signal is caught or service has error
	svc.Wait()

	if err := dbConn.Close(); err != nil {
		paaslog.Warnf("db close error: %v", err)
	}

	return nil
}

func CreateHttpServer(cfg *config.Webserver, userApi *service.PoliciesApi) (*http.Server, error) {
	tracing.Configure()
	newRouter := router.New()
	newRouter.Use(cors.Middleware)
	newRouter.Use(health.Middleware)
	newRouter.Use(metric.Middleware)
	newRouter.Use(account.AuthCompatibilityMiddleware)
	//Enable if you want to use IONOS token negotiator
	//client := createHttpClient(cfg)
	//middleware, err := setupIonosTokenNegotiator(client, cfg)
	//if err != nil {
	//	return nil, err
	//}
	//newRouter.Use(middleware)
	//Enable if you want to use IONOS auto refresh handler
	//customerMiddleware, err := setupIonosAutoRefreshHandler(client, cfg)
	//if err != nil {
	//	return nil, err
	//}
	//newRouter.Use(customerMiddleware)
	newRouter.Use(identity.AuthCompatibilityMiddleware)
	newRouter.Use(activity.Middleware)
	return setupTheHttpServer(newRouter, userApi, cfg.Port), nil
}

func setupTheHttpServer(router *chi.Mux, mgmt policiesApi.ServerInterface, port int) *http.Server {
	policiesApi.HandlerWithOptions(mgmt, policiesApi.ChiServerOptions{
		BaseRouter: router,
	})

	return server.CreateHTTPServer(fmt.Sprintf(":%d", port), router)
}
