package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// makeLogicalNetworkFixture mirrors the fixture from internal/logical_network tests.
// config is a LogicalNetworkConfig with its own required fields.
func makeLogicalNetworkFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "label": "ln-label", "name": "ln-name",
		"annotations": map[string]interface{}{},
		"createdAt":   "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
		"revision": 1, "kind": "vlan", "fabricId": 2,
		"infrastructureId":                   nil,
		"serviceStatus":                      "active",
		"lastAppliedLogicalNetworkProfileId": nil,
		"lastLogicalNetworkProfileAppliedAt": "2024-01-01T00:00:00Z",
		"config": map[string]interface{}{
			"id": 1, "deployType": "none", "deployStatus": "idle",
			"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
			"revision": 1, "kind": "vlan",
		},
	}
}

// Required: id, label, name, annotations, createdAt, updatedAt, revision, kind, fabricId
func makeLogicalNetworkProfileFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "label": "cisco-profile", "name": "Cisco Profile",
		"annotations": map[string]interface{}{},
		"createdAt":   "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
		"revision": 1, "kind": "vlan", "fabricId": 1,
		"links": []interface{}{},
	}
}

// --- logical-network list ---

func TestLogicalNetworkList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestLogicalNetworkList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "logical-network", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- logical-network-profile list ---

func TestLogicalNetworkProfileList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "logical-network-profile", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- format tests ---

func TestLogicalNetworkList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "logical-network", "list")
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

func TestLogicalNetworkProfileList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "logical-network-profile", "list")
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

// --- logical-network create ---

func TestLogicalNetworkCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeLogicalNetworkFixture(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ln-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"kind":"vlan","fabricId":1,"label":"test-net"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "logical-network", "create", "vlan", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- logical-network delete ---

func TestLogicalNetworkDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			// GET — return fixture so delete can read the revision
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeLogicalNetworkFixture(1))
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "logical-network", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
