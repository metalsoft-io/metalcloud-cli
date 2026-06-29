package network_device_default_secrets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func decodeItems(jsonStrings []string) []any {
	out := make([]any, len(jsonStrings))
	for i, s := range jsonStrings {
		var v any
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			panic(fmt.Sprintf("decodeItems: invalid JSON at index %d: %v", i, err))
		}
		out[i] = v
	}
	return out
}

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

const secretItem = `{
	"id":1,"siteId":1,"macAddressOrSerialNumber":"AA:BB:CC:DD:EE:01",
	"secretName":"admin-password","createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

func secretListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// --- List (fetch-all, no flags) ---

func TestNetworkDeviceDefaultSecretsList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": secretListHandler(http.StatusOK, []string{secretItem, secretItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": secretListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"siteId":1,"macAddressOrSerialNumber":"AA:BB:CC:DD:EE:%02d",
				"secretName":"secret-%d","createdTimestamp":"2024-01-01T00:00:00Z",
				"updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-devices/default-secrets", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- List (single-page path: page=1, limit=5) ---

func TestNetworkDeviceDefaultSecretsList_SinglePage(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": secretListHandler(http.StatusOK, []string{secretItem}, 1, 3),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 1, 5); err != nil {
		t.Errorf("expected nil error for page/limit path, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_SinglePage_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": testutils.ErrorHandler(http.StatusInternalServerError, "error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 1, 5); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

// --- Get ---

func TestNetworkDeviceDefaultSecretsGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": testutils.RawHandler(http.StatusOK, secretItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Create ---

func TestNetworkDeviceDefaultSecretsCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, secretItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsCreate(ctx, 1, "AA:BB:CC:DD:EE:01", "admin-password", "s3cr3t"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsCreate(ctx, 1, "AA:BB:CC:DD:EE:01", "name", "val"); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- Update ---

func TestNetworkDeviceDefaultSecretsUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, secretItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsUpdate(ctx, "1", "new-secret-value"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsUpdate(ctx, "99", "value"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- Delete ---

func TestNetworkDeviceDefaultSecretsDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}
