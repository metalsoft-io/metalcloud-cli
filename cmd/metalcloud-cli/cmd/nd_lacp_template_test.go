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

// lacpTemplateFixture satisfies sdk.NetworkDeviceLinkAggregationConfigurationTemplate required fields.
func lacpTemplateFixture(id float32) map[string]interface{} {
	return map[string]interface{}{
		"id":                  id,
		"action":              "create",
		"aggregationType":     "lag",
		"networkDeviceDriver": "junos",
		"executionType":       "cli",
		"libraryLabel":        "lacp-lib",
		"configuration":       "Y29uZmln",
		"createdTimestamp":    time.Now().UTC().Format(time.RFC3339),
		"updatedTimestamp":    time.Now().UTC().Format(time.RFC3339),
		"links":               []interface{}{},
	}
}

func newLacpTemplateServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-device-link-aggregation-configuration-templates", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(lacpTemplateFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(lacpTemplateFixture(1)))
		})
		mux.HandleFunc("/api/v2/network-device-link-aggregation-configuration-templates/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(lacpTemplateFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(lacpTemplateFixture(1))
		})
	}))
}

// --- network-configuration link-aggregation-template list ---

func TestLinkAggTemplateList_HappyPath(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "lacp-lib") {
		t.Fatalf("expected libraryLabel in output, got: %s", out)
	}
}

func TestLinkAggTemplateList_Alias(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestLinkAggTemplateList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "network-configuration", "link-aggregation-template", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- network-configuration link-aggregation-template get ---

func TestLinkAggTemplateGet_HappyPath(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "lacp-lib") {
		t.Fatalf("expected libraryLabel in output, got: %s", out)
	}
}

func TestLinkAggTemplateGet_Alias(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestLinkAggTemplateGet_NoArgs(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- network-configuration link-aggregation-template create ---

func TestLinkAggTemplateCreate_HappyPath(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "lacp-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"action":"create","aggregationType":"lag","networkDeviceDriver":"junos","executionType":"cli","libraryLabel":"lacp-lib","configuration":"Y29uZmln"}`)
	f.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "create", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLinkAggTemplateCreate_MissingFlag(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "create"); err == nil {
		t.Fatal("expected error when --config-source is missing")
	}
}

// --- network-configuration link-aggregation-template delete ---

func TestLinkAggTemplateDelete_HappyPath(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLinkAggTemplateDelete_NoArgs(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-configuration", "link-aggregation-template", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

func TestLacpTemplateList_Formats(t *testing.T) {
	srv := newLacpTemplateServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "network-configuration", "link-aggregation-template", "list")
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
