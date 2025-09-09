package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ionos-cloud/go-sample-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type fakeDB struct{}

func (f *fakeDB) Exec(query string, args ...interface{}) (interface{}, error)      { return nil, nil }
func (f *fakeDB) Get(dest interface{}, query string, args ...interface{}) error    { return nil }
func (f *fakeDB) Select(dest interface{}, query string, args ...interface{}) error { return nil }

func TestUserRepoImpl_SaveAndFindByID(t *testing.T) {
	// Use a real *sqlx.DB or a mock; here we use nil for brevity since methods are stubbed
	repo := NewUserRepo((*sqlx.DB)(nil))
	ctx := context.Background()
	userID := uuid.New()
	user := &model.User{UserID: userID}

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save() error = %v, want nil", err)
	}

	got, err := repo.FindByID(ctx, userID)
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}
	if got.UserID != userID {
		t.Errorf("FindByID() got UserID = %v, want %v", got.UserID, userID)
	}
}
