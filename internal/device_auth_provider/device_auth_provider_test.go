package device_auth_provider

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeProvider(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "prov-1", "name": "Provider 1",
		"kind": "radius", "ipAddress": "10.0.0.1", "port": 1812,
		"username": "admin", "hasSharedSecret": true, "hasPassword": false,
		"status": "active", "siteId": 1, "revision": 1,
		"createdBy": 1, "createdAt": "2024-01-01T00:00:00Z",
	}
}

func TestDeviceAuthProviderList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeProvider(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/device-auth-providers": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDeviceAuthProviderList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/device-auth-providers": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderList(ctx, nil, nil, nil); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestDeviceAuthProviderList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/device-auth-providers": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("expected nil error on empty list, got: %v", err)
	}
}

func TestDeviceAuthProviderList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeProvider(i + 1)
		page2[i] = makeProvider(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeProvider(i + 201)
	}

	ts := testutils.MultiPageServer("/api/v2/device-auth-providers", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestDeviceAuthProviderGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		// GetDeviceAuthProviderByIdOrLabel with numeric id calls GetDeviceAuthProviderById
		"/api/v2/device-auth-providers/1": testutils.JSONHandler(200, makeProvider(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderGet(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDeviceAuthProviderGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/device-auth-providers/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderGet(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestDeviceAuthProviderCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/device-auth-providers": testutils.JSONHandler(201, makeProvider(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// sharedSecret is required by the SDK's CreateDeviceAuthProvider struct.
	config := []byte(`{"label":"prov-2","name":"Provider 2","kind":"radius","ipAddress":"10.0.0.2","port":1812,"username":"admin","siteId":1,"sharedSecret":"s3cr3t"}`)
	if err := DeviceAuthProviderCreate(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDeviceAuthProviderDelete_HappyPath(t *testing.T) {
	getHandler := testutils.JSONHandler(200, makeProvider(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		// GetDeviceAuthProviderById (GET) then DeleteDeviceAuthProvider (DELETE) — same path.
		"/api/v2/device-auth-providers/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceAuthProviderDelete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
