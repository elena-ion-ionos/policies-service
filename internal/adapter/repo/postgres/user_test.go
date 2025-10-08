package postgres

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ionos-cloud/go-paaskit/observability/paaslog"
	"github.com/ionos-cloud/go-sample-service/internal/config"
	"github.com/ionos-cloud/go-sample-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserRepoImpl_SaveAndFindByID(t *testing.T) {
	ctx := context.Background()

	dbName := "test-db"
	user := "postgres"
	password := "postgres"
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	u, err := url.Parse(connStr)
	host, port := u.Hostname(), u.Port()

	paaslog.InfoCf(ctx, "Test db connection string: %s", connStr)

	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Database{
		Database: dbName,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		SslMode:  "disable",
	}
	db := config.MustNewDB(cfg)
	repo := NewUserRepo(db)
	userId := uuid.New()

	err = repo.Save(ctx, &model.User{
		UserID:         userId,
		ContractNumber: 12345,
	})
	if err != nil {
		t.Fatal(err)
	}

	userDb, err := repo.FindByID(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, userDb)
}
