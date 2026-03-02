package server_instance

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

func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

func init() {
	viper.Set(formatter.ConfigFormat, "text")
}

func TestServerInstanceList(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/server-instances") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"data": [
						{
							"id": 1,
							"label": "si-1",
							"infrastructureId": 123,
							"groupId": 10,
							"serviceStatus": "active",
							"createdTimestamp": "2024-01-01T00:00:00Z",
							"updatedTimestamp": "2024-01-01T00:00:00Z",
							"links": []
						}
					]
				}`))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ServerInstanceList(ctx, "123")
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("LinksAsMap", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/server-instances") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"data": [
						{
							"id": 1,
							"label": "si-1",
							"infrastructureId": 123,
							"groupId": 10,
							"serviceStatus": "active",
							"createdTimestamp": "2024-01-01T00:00:00Z",
							"updatedTimestamp": "2024-01-01T00:00:00Z",
							"links": {"self": "http://example.com/server-instances/1"}
						}
					]
				}`))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ServerInstanceList(ctx, "123")
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/server-instances") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error": "not found"}`))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ServerInstanceList(ctx, "123")
		if err == nil {
			t.Error("expected an error for HTTP 404, got nil")
		}
	})
}
