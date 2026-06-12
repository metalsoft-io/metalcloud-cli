package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// bgpConfigTemplateFixture satisfies sdk.NetworkDeviceBGPConfigurationTemplate required fields.
func bgpConfigTemplateFixture(id float32) map[string]interface{} {
	return map[string]interface{}{
		"id":                          id,
		"action":                      "add-global-config",
		"networkType":                 "underlay",
		"networkDeviceDriver":         "junos",
		"networkDevicePosition":       "all",
		"remoteNetworkDevicePosition": "all",
		"bgpNumbering":                "numbered",
		"bgpLinkConfiguration":        "disabled",
		"executionType":               "cli",
		"libraryLabel":                "my-lib",
		"configuration":               "Y29uZmln",
		"createdTimestamp":            time.Now().UTC().Format(time.RFC3339),
		"updatedTimestamp":            time.Now().UTC().Format(time.RFC3339),
		"links":                       []interface{}{},
	}
}

func newNetworkDeviceConfigTemplateServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-device-bgp-configuration-templates", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(bgpConfigTemplateFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(bgpConfigTemplateFixture(1)))
		})
		mux.HandleFunc("/api/v2/network-device-bgp-configuration-templates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(bgpConfigTemplateFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(bgpConfigTemplateFixture(1))
		})
	}))
}

// --- network-configuration device-template list ---

func TestNetworkDeviceConfigTemplateList_HappyPath(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-configuration", "device-template", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "my-lib") {
		t.Fatalf("expected libraryLabel in output, got: %s", out)
	}
}

func TestNetworkDeviceConfigTemplateList_Alias(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestNetworkDeviceConfigTemplateList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "network-configuration", "device-template", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- network-configuration device-template get ---

func TestNetworkDeviceConfigTemplateGet_HappyPath(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-configuration", "device-template", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "my-lib") {
		t.Fatalf("expected libraryLabel in output, got: %s", out)
	}
}

func TestNetworkDeviceConfigTemplateGet_Alias(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestNetworkDeviceConfigTemplateGet_NoArgs(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- network-configuration device-template create ---

func TestNetworkDeviceConfigTemplateCreate_HappyPath(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ndct-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"action":"add-global-config","networkType":"underlay","networkDeviceDriver":"junos","networkDevicePosition":"all","remoteNetworkDevicePosition":"all","bgpNumbering":"numbered","bgpLinkConfiguration":"disabled","executionType":"cli","libraryLabel":"my-lib","configuration":"Y29uZmln"}`)
	f.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "create", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestNetworkDeviceConfigTemplateCreate_MissingFlag(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "create"); err == nil {
		t.Fatal("expected error when --config-source is missing")
	}
}

// --- network-configuration device-template delete ---

func TestNetworkDeviceConfigTemplateDelete_HappyPath(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestNetworkDeviceConfigTemplateDelete_NoArgs(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "device-template", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

func TestNetworkDeviceConfigTemplateList_Formats(t *testing.T) {
	srv := newNetworkDeviceConfigTemplateServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "network-configuration", "device-template", "list")
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
