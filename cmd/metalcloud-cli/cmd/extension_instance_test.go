package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// extensionInstanceFixture satisfies sdk.ExtensionInstance required fields.
func extensionInstanceFixture(id float32) map[string]interface{} {
	return map[string]interface{}{
		"id":                  id,
		"revision":            1.0,
		"label":               "my-ext-instance",
		"automaticManagement": 0.0,
		"updatedTimestamp":    "2024-01-01T00:00:00Z",
		"infrastructureId":    1.0,
		"infrastructure": map[string]interface{}{
			"id":    1.0,
			"label": "my-infra",
		},
		"extensionId":    1.0,
		"serviceStatus":  "active",
		"links":          []interface{}{},
		"inputVariables": []interface{}{},
		"outputVariables": []interface{}{},
		"config": map[string]interface{}{
			"revision":            1.0,
			"label":               "my-ext-instance",
			"automaticManagement": 0.0,
			"deployType":          "deploy",
			"deployStatus":        "active",
			"updatedTimestamp":    "2024-01-01T00:00:00Z",
		},
	}
}

func newExtensionInstanceServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// Infrastructure resolution (GetInfrastructureByIdOrLabel)
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(infrastructureItem))
		})
		// extension-instance list
		mux.HandleFunc("/api/v2/infrastructures/1/extension-instances", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(extensionInstanceFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(extensionInstanceFixture(1)))
		})
		// extension-instance get / delete / update
		mux.HandleFunc("/api/v2/extension-instances/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(extensionInstanceFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(extensionInstanceFixture(1))
		})
	}))
}

// List tests are in extension_test.go (TestExtensionInstanceList_*).

func TestExtensionInstanceList_NoArgs(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "list"); err == nil {
		t.Fatal("expected error when no infrastructure arg given")
	}
}

// --- extension-instance get ---

func TestExtensionInstanceGet_HappyPath(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "extension-instance", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "my-ext-instance") {
		t.Fatalf("expected instance label in output, got: %s", out)
	}
}

func TestExtensionInstanceGet_Alias(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestExtensionInstanceGet_NoArgs(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- extension-instance create ---

func TestExtensionInstanceCreate_HappyPath(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ext-inst-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"extensionId":1,"label":"my-ext-instance"}`)
	f.Close()

	if _, err := runCLI(t, srv, "extension-instance", "create", "1", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestExtensionInstanceCreate_MissingFlag(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "create", "1"); err == nil {
		t.Fatal("expected error when neither --config-source nor --extension-id given")
	}
}

// --- extension-instance delete ---

func TestExtensionInstanceDelete_HappyPath(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestExtensionInstanceDelete_NoArgs(t *testing.T) {
	srv := newExtensionInstanceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "extension-instance", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}
