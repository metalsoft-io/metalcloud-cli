package system

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func newTestContext(t *testing.T, handler http.HandlerFunc) context.Context {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: server.URL}}
	return context.WithValue(context.Background(), api.ApiClientContextKey, sdk.NewAPIClient(cfg))
}

// The SDK's CurrentUser model requires many properties; a response missing any
// of them must not break the startup user fetch.
func TestGetUserPermissions_MinimalUserBody(t *testing.T) {
	ctx := newTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/user" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":42,"email":"a@b.c","accessLevel":"admin","unknownField":{"x":1}}`))
	})

	userId, accessLevel, permissions, err := GetUserPermissions(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userId != "42" || accessLevel != "admin" || permissions != nil {
		t.Fatalf("got userId=%q accessLevel=%q permissions=%v", userId, accessLevel, permissions)
	}
}

func TestGetUserPermissions_HTTPError(t *testing.T) {
	ctx := newTestContext(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	})

	if _, _, _, err := GetUserPermissions(ctx); err == nil {
		t.Fatal("expected error for 401 response")
	}
}
