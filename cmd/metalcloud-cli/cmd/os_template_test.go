package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func osTemplateFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "name": "Ubuntu 22.04", "visibility": "public",
		"status": "active", "revision": 1, "createdBy": 1,
		"createdAt": "2024-01-01T00:00:00Z",
		"device":    map[string]interface{}{"type": "server", "bootMode": "uefi", "architecture": "x86_64"},
		"install":   map[string]interface{}{"method": "oob", "driveType": "local_drive", "readyMethod": "wait_for_power_off"},
		"os": map[string]interface{}{
			"name": "Ubuntu", "version": "22.04",
			"credential": map[string]interface{}{"username": "root", "passwordType": "plain"},
		},
		"links": []interface{}{},
	}
}

func TestOsTemplateList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(osTemplateFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "os-template", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestOsTemplateList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(osTemplateFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "os-template", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestOsTemplateList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "os-template", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

func TestOsTemplateList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(osTemplateFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "os-template", "list")
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

func TestOsTemplateGet_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(osTemplateFixture(1))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "os-template", "get", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestOsTemplateGet_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "os-template", "get", "1"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- os-template create ---

func TestOSTemplateCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(osTemplateFixture(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ost-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"template":{"name":"Ubuntu 22.04","device":{"type":"server","bootMode":"uefi","architecture":"x86_64"},"install":{"method":"oob","driveType":"local_drive","readyMethod":"wait_for_power_off"},"imageBuild":{"required":true},"os":{"name":"Ubuntu","version":"22.04","credential":{"username":"root","passwordType":"plain"}}}}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "os-template", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- os-template update ---

func TestOSTemplateUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(osTemplateFixture(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ost-update-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"template":{"name":"Ubuntu 22.04 Updated","device":{"type":"server","bootMode":"uefi","architecture":"x86_64"},"install":{"method":"oob","driveType":"local_drive","readyMethod":"wait_for_power_off"},"imageBuild":{"required":true},"os":{"name":"Ubuntu","version":"22.04","credential":{"username":"root","passwordType":"plain"}}}}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "os-template", "update", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- os-template delete ---

func TestOSTemplateDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/os-templates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(osTemplateFixture(1))
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "os-template", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
