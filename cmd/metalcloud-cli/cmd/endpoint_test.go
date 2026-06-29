package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// endpointItem satisfies sdk.Endpoint required fields.
var endpointItem = map[string]interface{}{
	"id":               "1",
	"revision":         "1",
	"siteId":           1,
	"name":             "ep-01",
	"label":            "ep-01",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newEndpointTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// endpoint list — FetchAllPages / GetEndpoints
		mux.HandleFunc("/api/v2/endpoints", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(endpointItem))
		})
		// endpoint get — GetEndpointById
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
}

// --- endpoint list ---

func TestEndpointList_HappyPath(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ep-01") {
		t.Fatalf("expected ep-01 in output, got: %s", out)
	}
}

func TestEndpointList_Alias(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestEndpointList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "endpoint", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- endpoint get ---

func TestEndpointGet_HappyPath(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ep-01") {
		t.Fatalf("expected ep-01 in output, got: %s", out)
	}
}

func TestEndpointGet_NoArgs(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "get"); err == nil {
		t.Fatal("expected error when no args given to endpoint get")
	}
}

// --- endpoint create ---

func TestEndpointCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ep-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"siteId":1,"name":"ep-new","label":"ep-new"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- endpoint update ---

func TestEndpointUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ep-update-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"ep-updated"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "update", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- endpoint delete ---

func TestEndpointDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
