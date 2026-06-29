package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// networkDeviceDefaultSecretsFixture satisfies sdk.NetworkDeviceDefaultSecrets required fields.
func networkDeviceDefaultSecretsFixture(id float32) map[string]interface{} {
	return map[string]interface{}{
		"id":                       id,
		"siteId":                   1.0,
		"macAddressOrSerialNumber": "AA:BB:CC:DD:EE:FF",
		"secretName":               "admin_password",
		"createdTimestamp":         "2024-01-01T00:00:00Z",
		"updatedTimestamp":         "2024-01-01T00:00:00Z",
		"links":                    []interface{}{},
	}
}

func newNetworkDeviceDefaultSecretsServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-devices/default-secrets", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(networkDeviceDefaultSecretsFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(networkDeviceDefaultSecretsFixture(1)))
		})
		mux.HandleFunc("/api/v2/network-devices/default-secrets/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(networkDeviceDefaultSecretsFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(networkDeviceDefaultSecretsFixture(1))
		})
	}))
}

// --- network-device default-secrets list ---

func TestNetworkDeviceDefaultSecretsList_HappyPath(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-device", "default-secrets", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "AA:BB:CC:DD:EE:FF") {
		t.Fatalf("expected MAC address in output, got: %s", out)
	}
}

func TestNetworkDeviceDefaultSecretsList_Alias(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "network-device", "default-secrets", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- network-device default-secrets get ---

func TestNetworkDeviceDefaultSecretsGet_HappyPath(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-device", "default-secrets", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "AA:BB:CC:DD:EE:FF") {
		t.Fatalf("expected MAC address in output, got: %s", out)
	}
}

func TestNetworkDeviceDefaultSecretsGet_Alias(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsGet_NoArgs(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- network-device default-secrets create ---

func TestNetworkDeviceDefaultSecretsCreate_HappyPath(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "create",
		"--site-id", "1",
		"--mac-or-serial", "AA:BB:CC:DD:EE:FF",
		"--secret-name", "admin_password",
		"--secret-value", "s3cur3",
	); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsCreate_NoEndpoint(t *testing.T) {
	// With no endpoint the request should fail even with all flags provided.
	if _, err := runCLI(t, nil, "network-device", "default-secrets", "create",
		"--site-id", "1",
		"--mac-or-serial", "AA:BB:CC:DD:EE:FF",
		"--secret-name", "admin_password",
		"--secret-value", "s3cur3",
	); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- network-device default-secrets delete ---

func TestNetworkDeviceDefaultSecretsDelete_HappyPath(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsDelete_NoArgs(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "default-secrets", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

func TestNetworkDeviceDefaultSecretsList_Formats(t *testing.T) {
	srv := newNetworkDeviceDefaultSecretsServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "network-device", "default-secrets", "list")
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
