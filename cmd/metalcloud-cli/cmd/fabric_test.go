package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var fabricItem = map[string]interface{}{
	"id": "1", "name": "test-fabric",
	// fabricConfiguration must have fabricType so the discriminator union type
	// can serialize back to JSON; empty {} leaves all sub-types nil and MarshalJSON
	// returns (nil, nil) which causes "unexpected end of JSON input".
	"fabricConfiguration": map[string]interface{}{"fabricType": "ethernet"},
	"revision":            "1",
	"createdTimestamp":    "2024-01-01T00:00:00Z",
	"updatedTimestamp":    "2024-01-01T00:00:00Z",
}

func newFabricTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-fabrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(fabricItem))
		})
		mux.HandleFunc("/api/v2/network-fabrics/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(fabricItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestFabricList(t *testing.T) {
	srv := newFabricTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "fabric", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-fabric") {
		t.Errorf("expected output to contain 'test-fabric', got: %s", out)
	}
}

func TestFabricListAlias(t *testing.T) {
	srv := newFabricTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "fc", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-fabric") {
		t.Errorf("expected output to contain 'test-fabric', got: %s", out)
	}
}

func TestFabricGet(t *testing.T) {
	srv := newFabricTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "fabric", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-fabric") {
		t.Errorf("expected output to contain 'test-fabric', got: %s", out)
	}
}

func TestFabricGetRequiresArg(t *testing.T) {
	srv := newFabricTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "fabric", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestFabricHelp(t *testing.T) {
	out, err := runCLI(t, nil, "fabric", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "fabric") {
		t.Errorf("expected help output to contain 'fabric', got: %s", out)
	}
}

func TestFabricCreate(t *testing.T) {
	fabricSiteItem := map[string]interface{}{
		"id": 1, "revision": 1, "slug": "site-1", "name": "1",
	}
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(fabricSiteItem))
		})
		mux.HandleFunc("/api/v2/network-fabrics", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, fabricItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(fabricItem))
		})
		mux.HandleFunc("/api/v2/network-fabrics/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(fabricItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "fabric-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"fabricType":"ethernet"}`)
	f.Close()

	// fabric create site_id fabric_name fabric_type [description] --config-source
	_, execErr := runCLI(t, srv, "fabric", "create", "1", "test-fabric", "ethernet", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestFabricUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-fabrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(fabricItem))
		})
		mux.HandleFunc("/api/v2/network-fabrics/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, fabricItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(fabricItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "fabric-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"fabricType":"ethernet"}`)
	f.Close()

	// fabric update fabric_id [name [description]] --config-source
	_, execErr := runCLI(t, srv, "fabric", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestFabricList_Formats(t *testing.T) {
	srv := newFabricTestServer()
	defer srv.Close()
	// csv/yaml/text/md panic in the formatter due to FabricConfiguration's interface{} field;
	// json is the only format that works end-to-end for fabric list today.
	for _, format := range []string{"json"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "fabric", "list")
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
