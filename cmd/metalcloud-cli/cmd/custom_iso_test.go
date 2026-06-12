package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var customIsoItem = map[string]interface{}{
	"id":               1.0,
	"label":            "test-iso",
	"name":             "Test ISO",
	"type":             "standard",
	"isPublic":         0.0,
	"accessUrl":        "http://example.com/test.iso",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newCustomIsoTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/custom-isos/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(customIsoItem)
		})
		mux.HandleFunc("/api/v2/custom-isos", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(customIsoItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestCustomIsoList(t *testing.T) {
	srv := newCustomIsoTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "custom-iso", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-iso") {
		t.Errorf("expected output to contain 'test-iso', got: %s", out)
	}
}

func TestCustomIsoGet(t *testing.T) {
	srv := newCustomIsoTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "custom-iso", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-iso") {
		t.Errorf("expected output to contain 'test-iso', got: %s", out)
	}
}

func TestCustomIsoGetRequiresArg(t *testing.T) {
	srv := newCustomIsoTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "custom-iso", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestCustomIsoList_Formats(t *testing.T) {
	// customIsoRaw.IsPublic is *bool, so the list path requires a JSON bool.
	// The shared customIsoItem uses 0.0 (float) to satisfy the SDK GET path
	// (sdk.CustomIso.IsPublic is float32), so we use a separate inline fixture here.
	isoListItem := map[string]interface{}{
		"id": 1.0, "label": "test-iso", "name": "Test ISO",
		"type": "standard", "isPublic": false,
		"accessUrl":        "http://example.com/test.iso",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/custom-isos", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(isoListItem))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "custom-iso", "list")
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
