package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// infrastructureItem satisfies sdk.Infrastructure required fields.
// InfrastructureList uses FetchAllPages (paginated GET /api/v2/infrastructures).
// InfrastructureGet calls GetInfrastructureByIdOrLabel (search then match by id/label).
var infrastructureItem = map[string]interface{}{
	"id":               1.0,
	"revision":         1.0,
	"label":            "my-infra",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"serviceStatus":    "active",
	"datacenterName":   "dc-01",
	"siteId":           1.0,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"designIsLocked":   0.0,
	"config":           map[string]interface{}{"revision": 1.0},
}

func newInfrastructureTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(infrastructureItem))
		})
	}))
}

// --- infrastructure list ---

func TestInfrastructureList_HappyPath(t *testing.T) {
	srv := newInfrastructureTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "infrastructure", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "my-infra") {
		t.Fatalf("expected my-infra in output, got: %s", out)
	}
}

func TestInfrastructureList_Alias(t *testing.T) {
	srv := newInfrastructureTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "infrastructure", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestInfrastructureList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "infrastructure", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- infrastructure get ---

func TestInfrastructureGet_HappyPath(t *testing.T) {
	srv := newInfrastructureTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "infrastructure", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "my-infra") {
		t.Fatalf("expected my-infra in output, got: %s", out)
	}
}

func TestInfrastructureGet_NoArgs(t *testing.T) {
	srv := newInfrastructureTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "infrastructure", "get"); err == nil {
		t.Fatal("expected error when no args given to infrastructure get")
	}
}

func TestInfrastructureCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, infrastructureItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(infrastructureItem))
		})
	}))
	defer srv.Close()

	// infrastructure create site_id label
	_, err := runCLI(t, srv, "infrastructure", "create", "1", "my-infra")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfrastructureUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(infrastructureItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, infrastructureItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(infrastructureItem)
		})
	}))
	defer srv.Close()

	// infrastructure update infrastructure_id [new_label]
	_, err := runCLI(t, srv, "infrastructure", "update", "1", "my-infra-updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfrastructureDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(infrastructureItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(infrastructureItem)
		})
	}))
	defer srv.Close()

	_, err := runCLI(t, srv, "infrastructure", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfrastructureList_Formats(t *testing.T) {
	srv := newInfrastructureTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "infrastructure", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}
