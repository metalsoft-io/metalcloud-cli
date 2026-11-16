package user

import (
	"context"
	"testing"
	"time"
)

func TestResolveApiTokenExpiry(t *testing.T) {
	ctx := context.Background()

	got, err := resolveApiTokenExpiry(ctx, "2027-06-30")
	if err != nil {
		t.Fatalf("date: unexpected error: %v", err)
	}
	if want := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("date: got %v, want %v", got, want)
	}

	got, err = resolveApiTokenExpiry(ctx, "2027-01-01T12:30:00+02:00")
	if err != nil {
		t.Fatalf("rfc3339: unexpected error: %v", err)
	}
	if want := time.Date(2027, 1, 1, 10, 30, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("rfc3339: got %v, want %v", got, want)
	}

	if _, err := resolveApiTokenExpiry(ctx, "next tuesday"); err == nil {
		t.Error("expected error for invalid expiry")
	}
}
