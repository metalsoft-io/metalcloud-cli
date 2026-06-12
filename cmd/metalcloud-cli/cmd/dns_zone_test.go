package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os" //nolint:goimports
	"strings"
	"testing"
)

var dnsZoneItem = map[string]interface{}{
	"id": 1.0, "label": "test", "zoneName": "test.example.com",
	"zoneType": "primary", "soaEmail": "admin@example.com",
	"soaSerial": 2024010101.0, "ttl": 300.0,
	"nameServers": []interface{}{"ns1.example.com"},
	"isDefault":   false, "status": "active", "revision": 1.0,
	"createdBy": 1.0, "createdAt": "2024-01-01T00:00:00Z",
}

func newDnsZoneTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/dns-zones", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(dnsZoneItem))
		})
		mux.HandleFunc("/api/v2/dns-zones/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(dnsZoneItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestDnsZoneList(t *testing.T) {
	srv := newDnsZoneTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "dns-zone", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected output to contain 'test.example.com', got: %s", out)
	}
}

func TestDnsZoneListAlias(t *testing.T) {
	srv := newDnsZoneTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "dns", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected output to contain 'test.example.com', got: %s", out)
	}
}

func TestDnsZoneGet(t *testing.T) {
	srv := newDnsZoneTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "dns-zone", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test.example.com") {
		t.Errorf("expected output to contain 'test.example.com', got: %s", out)
	}
}

func TestDnsZoneGetRequiresArg(t *testing.T) {
	srv := newDnsZoneTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "dns-zone", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestDnsZoneHelp(t *testing.T) {
	out, err := runCLI(t, nil, "dns-zone", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dns") {
		t.Errorf("expected help output to contain 'dns', got: %s", out)
	}
}

func TestDNSZoneCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/dns-zones", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, dnsZoneItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(dnsZoneItem))
		})
		mux.HandleFunc("/api/v2/dns-zones/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(dnsZoneItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// Use --zone-name flags path (no --config-source needed)
	_, execErr := runCLI(t, srv, "dns-zone", "create",
		"--zone-name", "test.example.com",
		"--name-servers", "ns1.example.com")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestDNSZoneUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/dns-zones", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(dnsZoneItem))
		})
		mux.HandleFunc("/api/v2/dns-zones/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, dnsZoneItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(dnsZoneItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "dns-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"description":"updated"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "dns-zone", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestDNSZoneDelete(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/dns-zones", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(dnsZoneItem))
		})
		mux.HandleFunc("/api/v2/dns-zones/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(dnsZoneItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, execErr := runCLI(t, srv, "dns-zone", "delete", "1")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
