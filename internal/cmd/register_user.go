package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/heptiolabs/healthcheck"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/adapter/fetcher"
	"github.com/ionos-cloud/go-sample-service/internal/adapter/notification"
	"github.com/ionos-cloud/go-sample-service/internal/adapter/repo/postgres"
	"github.com/ionos-cloud/go-sample-service/internal/config"
	"github.com/ionos-cloud/go-sample-service/internal/controller"
	"github.com/ionos-cloud/go-sample-service/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	RegisterUserName = "register-user"
)

func newRegisterUserCtrl(dbConn *sqlx.DB) (controller.RegisterUser, error) {
	paaslog.Infof("creating worker controller")

	userRepo := postgres.NewUserRepo(dbConn)
	emailNotifier := &notification.EmailNotifier{}
	smsNotifier := &notification.SMSNotifier{}

	return controller.NewRegisterUser(userRepo, emailNotifier, smsNotifier)
}

func RunRegisterUser(ctx context.Context, cfg *config.RegisterUser) error {
	svc := MustNewService(ctx, RegisterUserName, cfg.Service)

	dbConn := config.MustNewDB(cfg.Database)
	ctrl, _ := newRegisterUserCtrl(dbConn)
	fetcher := fetcher.NewFetcher()
	registerUserService := service.MustNewRegisterUser(cfg, fetcher, ctrl)

	svc.StartObservabilityServer()

	health := healthcheck.NewHandler()
	if err := ConfigureHealthCheckHandler(health, dbConn, cfg.Service); err != nil {
		paaslog.Fatalf("failed to configure health check handler: %v", err)
	}

	svc.StartHealthCheckServer(health)

	svc.StartRoutine(RegisterUserName, func(ctx context.Context, triggerShutdown func(reason string)) {
		if err := registerUserService.Serve(ctx); err != nil {
			paaslog.ErrorCf(ctx, "%s Serve had error: %v, terminating", RegisterUserName, err)
			triggerShutdown(fmt.Sprintf("%s Serve had error", RegisterUserName))
		}
	})

	// this will block until a term signal is caught or service has error
	svc.Wait()

	if err := dbConn.Close(); err != nil {
		paaslog.Warnf("db close error: %v", err)
	}

	return nil
}

func RegisterUserFunc() *cobra.Command {
	var cfg config.RegisterUser

	cmd := &cobra.Command{
		Use:   RegisterUserName,
		Short: "starts the worker service for handling key operations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.InitViperFlags(cmd, os.Args); err != nil {
				paaslog.Fatalf("got error initializing config flags, err: %v", err)
			}

			if err := RunRegisterUser(context.Background(), &cfg); err != nil {
				paaslog.Fatalf("got error while running worker, err: %v", err)
			}
		},
	}

	cfg.AddFlags(cmd)

	return cmd
}
