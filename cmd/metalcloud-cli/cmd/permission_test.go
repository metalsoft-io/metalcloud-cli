package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// permission_test.go covers:
//   permission list — exits 0

var permissionItem = map[string]interface{}{
	"id":    "servers_read",
	"name":  "servers_read",
	"label": "Servers Read",
}

func newPermissionTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/permissions", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(permissionItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestPermissionList(t *testing.T) {
	srv := newPermissionTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "permission", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionListAlias(t *testing.T) {
	srv := newPermissionTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "permissions", "ls")
	if err != nil {
		t.Fatalf("unexpected error running 'permissions ls' alias: %v", err)
	}
}

func TestPermissionList_Formats(t *testing.T) {
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			srv := newPermissionTestServer()
			defer srv.Close()
			out, err := runCLIFormat(t, srv, format, "permission", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
		})
	}
}
