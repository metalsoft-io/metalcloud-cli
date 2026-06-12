package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// deviceAuthProviderFixture satisfies sdk.DeviceAuthProvider required fields.
func deviceAuthProviderFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":              id,
		"siteId":          1,
		"label":           "tacacs-provider-1",
		"name":            "TACACS Provider",
		"kind":            "tacacs",
		"ipAddress":       "10.0.0.1",
		"port":            49,
		"username":        "admin",
		"hasSharedSecret": true,
		"hasPassword":     true,
		"status":          "active",
		"revision":        1,
		"createdBy":       1,
		"createdAt":        "2024-01-01T00:00:00Z",
		"updatedAt":        "2024-01-01T00:00:00Z",
		"links":           []interface{}{},
	}
}

func newDeviceAuthProviderServer(extra func(mux *http.ServeMux)) *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/device-auth-providers", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(deviceAuthProviderFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(deviceAuthProviderFixture(1)))
		})
		mux.HandleFunc("/api/v2/device-auth-providers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(deviceAuthProviderFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(deviceAuthProviderFixture(1))
		})
		if extra != nil {
			extra(mux)
		}
	}))
}

// --- site device-auth-provider list ---

func TestDeviceAuthProviderList_HappyPath(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "device-auth-provider", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "tacacs-provider-1") {
		t.Fatalf("expected provider label in output, got: %s", out)
	}
}

func TestDeviceAuthProviderList_Alias(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestDeviceAuthProviderList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "site", "device-auth-provider", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- site device-auth-provider get ---

func TestDeviceAuthProviderGet_HappyPath(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "device-auth-provider", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "tacacs-provider-1") {
		t.Fatalf("expected provider label in output, got: %s", out)
	}
}

func TestDeviceAuthProviderGet_Alias(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestDeviceAuthProviderGet_NoArgs(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- site device-auth-provider create ---

func TestDeviceAuthProviderCreate_HappyPath(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "dap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"label":"tacacs-new","name":"New","siteId":1,"kind":"tacacs","ipAddress":"10.0.0.2","port":49,"sharedSecret":"s","username":"admin"}`)
	f.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "create", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeviceAuthProviderCreate_MissingFlag(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "create"); err == nil {
		t.Fatal("expected error when --config-source is missing")
	}
}

// --- site device-auth-provider delete ---

func TestDeviceAuthProviderDelete_HappyPath(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeviceAuthProviderDelete_NoArgs(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()

	if _, err := runCLI(t, srv, "site", "device-auth-provider", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

func TestDeviceAuthProviderList_Formats(t *testing.T) {
	srv := newDeviceAuthProviderServer(nil)
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "site", "device-auth-provider", "list")
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
