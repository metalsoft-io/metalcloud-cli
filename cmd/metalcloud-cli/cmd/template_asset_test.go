package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// template_asset_test.go covers:
//   template-asset list
//   template-asset get <id>
//   template-asset create --config-source <json>
//   template-asset delete <id>

var templateAssetItem = map[string]interface{}{
	"id":         1,
	"templateId": 10,
	"usage":      "bootstrap",
	"revision":   1,
	"createdBy":  1,
	"createdAt":  "2024-01-01T00:00:00Z",
	"file": map[string]interface{}{
		"name":             "bootstrap.sh",
		"mimeType":         "text/x-shellscript",
		"templatingEngine": false,
		"path":             "/assets/bootstrap.sh",
	},
	"tags": []interface{}{},
}

func newTemplateAssetTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/template-assets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				_ = json.NewEncoder(w).Encode(templateAssetItem)
			}
		})
		mux.HandleFunc("/api/v2/template-assets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(templateAssetItem)
			default:
				_ = json.NewEncoder(w).Encode(paginatedList(templateAssetItem))
			}
		})
	})
	return httptest.NewServer(mux)
}

func TestTemplateAssetList(t *testing.T) {
	srv := newTemplateAssetTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "template-asset", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "bootstrap") {
		t.Errorf("expected output to contain 'bootstrap', got: %s", out)
	}
}

func TestTemplateAssetGet(t *testing.T) {
	srv := newTemplateAssetTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "template-asset", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "bootstrap") {
		t.Errorf("expected output to contain 'bootstrap', got: %s", out)
	}
}

func TestTemplateAssetGetRequiresArg(t *testing.T) {
	srv := newTemplateAssetTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "template-asset", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestTemplateAssetDelete(t *testing.T) {
	srv := newTemplateAssetTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "template-asset", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplateAssetDeleteRequiresArg(t *testing.T) {
	srv := newTemplateAssetTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "template-asset", "delete")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}
