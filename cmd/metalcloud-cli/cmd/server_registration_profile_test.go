package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// server_registration_profile_test.go covers get, create, delete — list is
// already covered in server_creds_test.go.

func newRegistrationProfileTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/registration-profiles/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				_ = json.NewEncoder(w).Encode(serverRegistrationProfileItem)
			}
		})
		mux.HandleFunc("/api/v2/servers/registration-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(serverRegistrationProfileItem)
			default:
				_ = json.NewEncoder(w).Encode(paginatedList(serverRegistrationProfileItem))
			}
		})
	})
	return httptest.NewServer(mux)
}

func TestServerRegistrationProfileGet(t *testing.T) {
	srv := newRegistrationProfileTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-registration-profile", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "default-profile") {
		t.Errorf("expected output to contain 'default-profile', got: %s", out)
	}
}

func TestServerRegistrationProfileGetRequiresArg(t *testing.T) {
	srv := newRegistrationProfileTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-registration-profile", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestServerRegistrationProfileCreate(t *testing.T) {
	srv := newRegistrationProfileTestServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "srp-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"name":"default-profile","settings":{}}`)
	f.Close()

	_, err = runCLI(t, srv, "server-registration-profile", "create",
		"--name", "default-profile",
		"--config-source", f.Name(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerRegistrationProfileDelete(t *testing.T) {
	srv := newRegistrationProfileTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-registration-profile", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerRegistrationProfileDeleteRequiresArg(t *testing.T) {
	srv := newRegistrationProfileTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-registration-profile", "delete")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

// srpLookupServer serves the search/for-server/system-defaults routes.
func srpLookupServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/registration-profiles/search", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverRegistrationProfileItem)
		})
		mux.HandleFunc("/api/v2/servers/registration-profiles/search/for-server/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverRegistrationProfileItem)
		})
		mux.HandleFunc("/api/v2/servers/registration-profiles/system-defaults", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"registerCredentials":         "user",
				"defaultVirtualMediaProtocol": "HTTPS",
				"enableTpm":                   true,
			})
		})
	})
	return httptest.NewServer(mux)
}

func TestServerRegistrationProfileSearch(t *testing.T) {
	srv := srpLookupServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-registration-profile", "search", "--site-id", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"id"`) {
		t.Fatalf("expected the profile in the output, got: %s", out)
	}
}

func TestServerRegistrationProfileSearchRequiresSiteId(t *testing.T) {
	srv := srpLookupServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-registration-profile", "search"); err == nil {
		t.Fatal("expected an error when --site-id is missing")
	}
}

func TestServerRegistrationProfileForServer(t *testing.T) {
	srv := srpLookupServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-registration-profile", "for-server", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"id"`) {
		t.Fatalf("expected the profile in the output, got: %s", out)
	}
}

func TestServerRegistrationProfileForServerRequiresArg(t *testing.T) {
	srv := srpLookupServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-registration-profile", "for-server"); err == nil {
		t.Fatal("expected an error when no server ID is given")
	}
}

func TestServerRegistrationProfileSystemDefaults(t *testing.T) {
	srv := srpLookupServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-registration-profile", "system-defaults")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HTTPS") {
		t.Fatalf("expected the default settings in the output, got: %s", out)
	}
}
