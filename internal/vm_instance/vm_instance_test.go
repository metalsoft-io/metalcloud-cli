package vm_instance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

// setupTestContext creates a context wired to the given test server URL.
func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// TestVMInstanceList tests the VMInstanceList function under multiple response shapes.
func TestVMInstanceList(t *testing.T) {
	// PrintResult requires a known format; use "json" so it never touches the
	// table renderer and stays side-effect-free in tests.
	viper.Set(formatter.ConfigFormat, "json")

	t.Run("LinksAsArray", func(t *testing.T) {
		body := `{
			"data": [
				{
					"id": 1,
					"label": "vm-1",
					"infrastructureId": 123,
					"groupId": 10,
					"serviceStatus": "active",
					"typeId": 5,
					"diskSizeGB": 100,
					"ramGB": 16,
					"cpuCores": 8,
					"createdTimestamp": "2024-01-01T00:00:00Z",
					"updatedTimestamp": "2024-01-01T00:00:00Z",
					"links": []
				}
			]
		}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/infrastructures/") && strings.Contains(r.URL.Path, "/vm-instances") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := VMInstanceList(ctx, "123")
		if err != nil {
			t.Errorf("expected nil error with links as array, got: %v", err)
		}
	})

	t.Run("LinksAsMap", func(t *testing.T) {
		body := `{
			"data": [
				{
					"id": 1,
					"label": "vm-1",
					"infrastructureId": 123,
					"groupId": 10,
					"serviceStatus": "active",
					"typeId": 5,
					"diskSizeGB": 100,
					"ramGB": 16,
					"cpuCores": 8,
					"createdTimestamp": "2024-01-01T00:00:00Z",
					"updatedTimestamp": "2024-01-01T00:00:00Z",
					"links": {"self": "/api/v2/vm-instances/1"}
				}
			]
		}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/infrastructures/") && strings.Contains(r.URL.Path, "/vm-instances") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := VMInstanceList(ctx, "123")
		if err != nil {
			t.Errorf("expected nil error with links as map, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not found"}`))
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := VMInstanceList(ctx, "123")
		if err == nil {
			t.Error("expected an error for HTTP 404, got nil")
		}
	})
}

// TestVMInstanceGetCredentials tests the VMInstanceGetCredentials function.
func TestVMInstanceGetCredentials(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		body := `{
			"username": "admin",
			"initialPassword": "secret123",
			"publicSshKey": "ssh-rsa AAAA..."
		}`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/credentials") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := VMInstanceGetCredentials(ctx, "123", "1")
		if err != nil {
			t.Errorf("expected nil error for credentials success, got: %v", err)
		}
	})
}
