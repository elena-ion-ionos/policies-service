package fetcher

import (
	"context"
	"errors"
	"testing"
)

func TestFetcherImpl_Fetch_NotImplemented(t *testing.T) {
	f := NewFetcher()
	_, err := f.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, err) || err.Error() != "not implemented" {
		t.Fatalf("expected 'not implemented' error, got: %v", err)
	}
}
