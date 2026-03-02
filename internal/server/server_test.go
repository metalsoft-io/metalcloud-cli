package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

// TestMain configures the output format used by formatter.PrintResult so that
// tests don't fail with "format not supported yet" when viper has no value set.
func TestMain(m *testing.M) {
	viper.Set(formatter.ConfigFormat, "json")
	m.Run()
}

// setupTestContext creates a context with an SDK client pointed at the given test server URL.
func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// serverListLinksArrayJSON is a server list response where each server has
// "links" as an empty array — the shape that triggered the original SDK
// decode bug.
const serverListLinksArrayJSON = `{
	"data": [
		{
			"serverId": 1,
			"siteId": 2,
			"serverTypeId": 3,
			"serverUUID": "uuid-1",
			"serialNumber": "SN-001",
			"managementAddress": "10.0.0.1",
			"vendor": "Dell",
			"model": "R640",
			"serverStatus": "active",
			"revision": 1,
			"links": [],
			"serverMetricsMetadata": []
		}
	]
}`

// serverListLinksMapJSON is a server list response where "links" is a map —
// the normal, spec-compliant shape.
const serverListLinksMapJSON = `{
	"data": [
		{
			"serverId": 1,
			"siteId": 2,
			"serverTypeId": 3,
			"serverUUID": "uuid-1",
			"serialNumber": "SN-001",
			"managementAddress": "10.0.0.1",
			"vendor": "Dell",
			"model": "R640",
			"serverStatus": "active",
			"revision": 1,
			"links": {"self": "/api/v2/servers/1"},
			"serverMetricsMetadata": []
		}
	]
}`

// newServerListHandler returns an http.HandlerFunc that serves server list
// responses for paths containing "/servers", and 404s for everything else.
func newServerListHandler(body string, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/servers") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			fmt.Fprint(w, body)
			return
		}
		http.NotFound(w, r)
	}
}

func TestServerList_LinksAsArray(t *testing.T) {
	ts := httptest.NewServer(newServerListHandler(serverListLinksArrayJSON, http.StatusOK))
	defer ts.Close()

	ctx := setupTestContext(ts.URL)

	err := ServerList(ctx, false, nil, nil)
	if err != nil {
		t.Errorf("ServerList with links-as-array: expected nil error, got: %v", err)
	}
}

func TestServerList_LinksAsMap(t *testing.T) {
	ts := httptest.NewServer(newServerListHandler(serverListLinksMapJSON, http.StatusOK))
	defer ts.Close()

	ctx := setupTestContext(ts.URL)

	err := ServerList(ctx, false, nil, nil)
	if err != nil {
		t.Errorf("ServerList with links-as-map: expected nil error, got: %v", err)
	}
}

func TestServerList_HttpError(t *testing.T) {
	const errorBody = `{"message":"not found"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, errorBody)
	}))
	defer ts.Close()

	ctx := setupTestContext(ts.URL)

	err := ServerList(ctx, false, nil, nil)
	if err == nil {
		t.Error("ServerList with 404 response: expected error, got nil")
	}
}

func TestServerIdentify_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/identify-server") {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	ctx := setupTestContext(ts.URL)

	t.Run("ValidServerId", func(t *testing.T) {
		err := ServerIdentify(ctx, "1")
		if err != nil {
			t.Errorf("ServerIdentify with valid id: expected nil error, got: %v", err)
		}
	})

	t.Run("InvalidServerId", func(t *testing.T) {
		err := ServerIdentify(ctx, "not-a-number")
		if err == nil {
			t.Error("ServerIdentify with invalid id: expected error, got nil")
		}
	})
}
