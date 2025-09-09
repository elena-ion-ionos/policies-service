package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/adapter/fetcher"
	"github.com/ionos-cloud/go-sample-service/internal/adapter/notification"
	dbRepo "github.com/ionos-cloud/go-sample-service/internal/adapter/repo/postgres"
	"github.com/ionos-cloud/go-sample-service/internal/config"
	"github.com/ionos-cloud/go-sample-service/internal/controller"
	"github.com/ionos-cloud/go-sample-service/internal/service"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	registerUserName = "register-user"
)

func RegisterUser() *cobra.Command {
	var cfg config.RegisterUser

	cmd := &cobra.Command{
		Use:   registerUserName,
		Short: "starts the worker service for handling key operations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.InitViperFlags(cmd, os.Args); err != nil {
				paaslog.Fatalf("got error initializing config flags, err: %v", err)
			}

			if err := runRegisterUser(context.Background(), &cfg); err != nil {
				paaslog.Fatalf("got error while running worker, err: %v", err)
			}
		},
	}

	cfg.AddFlags(cmd)

	return cmd
}

func newRegisterUserCtrl(cfg *config.RegisterUser, dbConn *sqlx.DB) (controller.RegisterUser, error) {
	paaslog.Infof("creating worker controller")

	userRepo := dbRepo.NewUserRepo(dbConn)
	emailNotifier := &notification.EmailNotifier{}
	smsNotifier := &notification.SMSNotifier{}

	return controller.NewRegisterUser(userRepo, emailNotifier, smsNotifier)
}

func runRegisterUser(ctx context.Context, cfg *config.RegisterUser) error {
	svc := MustNewService(ctx, registerUserName, cfg.Service)

	dbConn := config.MustNewDB(cfg.Database)
	ctrl, _ := newRegisterUserCtrl(cfg, dbConn)
	fetcher := fetcher.NewFetcher()
	registerUserService := service.MustNewRegisterUser(cfg, fetcher, ctrl)

	svc.StartObservabilityServer()

	health := healthcheck.NewHandler()
	if err := cfg.ConfigureHealthCheckHandler(health, dbConn); err != nil {
		paaslog.Fatalf("failed to configure health check handler: %v", err)
	}

	svc.StartHealthCheckServer(health)

	svc.StartRoutine(registerUserName, func(ctx context.Context, triggerShutdown func(reason string)) {
		if err := registerUserService.Serve(ctx); err != nil {
			paaslog.ErrorCf(ctx, "%s Serve had error: %v, terminating", registerUserName, err)
			triggerShutdown(fmt.Sprintf("%s Serve had error", registerUserName))
		}
	})

	// this will block until a term signal is caught or service has error
	svc.Wait()

	if err := dbConn.Close(); err != nil {
		paaslog.Warnf("db close error: %v", err)
	}

	return nil
}
