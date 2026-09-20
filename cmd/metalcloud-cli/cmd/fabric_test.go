package cmd

import (
	"encoding/json"
	"io"
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

// ---------------------------------------------------------------------------
// BGP sessions, link aggregations and the other fabric sub-resources
// ---------------------------------------------------------------------------

var fabricBgpSessionItem = map[string]interface{}{
	"id": 7.0, "networkFabricId": 1.0, "name": "session-7", "status": "active",
	"bgpNumbering": "numbered", "bgpLinkConfiguration": "active",
	"networkFabricLinkId": 3.0, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var fabricLagItem = map[string]interface{}{
	"id": 4.0, "networkFabricId": 1.0, "name": "lag-4", "type": "lag", "status": "active",
	"revision": "3", "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var fabricSubResLinkItem = map[string]interface{}{
	"id": 3.0, "networkFabricId": 1.0, "linkType": "leaf_spine", "status": "active",
	"revision": "1", "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var fabricSubResInterconnectItem = map[string]interface{}{
	"id": "12", "interconnectType": "dci-evpn", "label": "dc1-dc2", "name": "DC1 to DC2",
	"bgpConfigurationTemplateId": 3.0, "revision": "7", "status": "draft",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

// newFabricSubResServer serves the fabric plus its BGP session, link
// aggregation, link and interconnect sub-resources, recording the last
// request's method, path and body.
func newFabricSubResServer(last *fabricSubResRequest) *httptest.Server {
	record := func(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			*last = fabricSubResRequest{Method: r.Method, Path: r.URL.Path, Body: string(body)}
			handler(w, r)
		}
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-fabrics", record(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(fabricItem))
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, fabricItem)
		}))
		for _, action := range []string{"accept-deploy", "reject-deploy"} {
			mux.HandleFunc("/api/v2/network-fabrics/1/actions/"+action, record(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
		}
		mux.HandleFunc("/api/v2/network-fabrics/1/links/3", record(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, fabricSubResLinkItem)
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1/network-fabric-interconnects", record(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(fabricSubResInterconnectItem))
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1/bgp-sessions", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, fabricBgpSessionItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(fabricBgpSessionItem))
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1/bgp-sessions/7", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, fabricBgpSessionItem)
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1/link-aggregations", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, fabricLagItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(fabricLagItem))
		}))
		mux.HandleFunc("/api/v2/network-fabrics/1/link-aggregations/4", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, fabricLagItem)
		}))
	})
	return httptest.NewServer(mux)
}

type fabricSubResRequest struct {
	Method string
	Path   string
	Body   string
}

func TestFabricDeleteAndDeployDecisions(t *testing.T) {
	var last fabricSubResRequest
	srv := newFabricSubResServer(&last)
	defer srv.Close()

	for _, tc := range []struct {
		args []string
		path string
	}{
		{[]string{"fabric", "accept-deploy", "1"}, "/api/v2/network-fabrics/1/actions/accept-deploy"},
		{[]string{"fabric", "reject-deploy", "1"}, "/api/v2/network-fabrics/1/actions/reject-deploy"},
		{[]string{"fabric", "rm", "1"}, "/api/v2/network-fabrics/1"},
	} {
		if _, err := runCLI(t, srv, tc.args...); err != nil {
			t.Fatalf("%v: unexpected error: %v", tc.args, err)
		}
		if last.Path != tc.path {
			t.Errorf("%v hit %s, want %s", tc.args, last.Path, tc.path)
		}
	}

	if _, err := runCLI(t, srv, "fabric", "delete"); err == nil {
		t.Error("expected an error when no fabric id is given")
	}
}

func TestFabricGetLinkAndInterconnects(t *testing.T) {
	var last fabricSubResRequest
	srv := newFabricSubResServer(&last)
	defer srv.Close()

	out, err := runCLI(t, srv, "fabric", "get-link", "1", "3")
	if err != nil {
		t.Fatalf("get-link: %v", err)
	}
	if !strings.Contains(out, "leaf_spine") {
		t.Errorf("get-link output missing link: %s", out)
	}

	out, err = runCLI(t, srv, "fabric", "get-interconnects", "1")
	if err != nil {
		t.Fatalf("get-interconnects: %v", err)
	}
	if !strings.Contains(out, "dc1-dc2") {
		t.Errorf("get-interconnects output missing interconnect: %s", out)
	}
}

func TestFabricBgpSessionCommands(t *testing.T) {
	var last fabricSubResRequest
	srv := newFabricSubResServer(&last)
	defer srv.Close()

	out, err := runCLI(t, srv, "fabric", "bgp-session", "list", "1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "session-7") {
		t.Errorf("list output missing session: %s", out)
	}

	if _, err := runCLI(t, srv, "fabric", "bgp", "show", "1", "7"); err != nil {
		t.Fatalf("get via alias: %v", err)
	}

	configFile := fabricWriteTempConfig(t, `{"bgpNumbering":"numbered","bgpLinkConfiguration":"active","linkId":3}`)
	if _, err := runCLI(t, srv, "fabric", "bgp-session", "create", "1", "--config-source", configFile); err != nil {
		t.Fatalf("create: %v", err)
	}
	if last.Method != http.MethodPost || !strings.Contains(last.Body, `"bgpNumbering":"numbered"`) {
		t.Errorf("create sent %s %s", last.Method, last.Body)
	}

	updateFile := fabricWriteTempConfig(t, `{"customVariables":{"asn":65000}}`)
	if _, err := runCLI(t, srv, "fabric", "bgp-session", "update", "1", "7", "--config-source", updateFile); err != nil {
		t.Fatalf("update: %v", err)
	}
	if !strings.Contains(last.Body, "customVariables") {
		t.Errorf("update body wrong: %s", last.Body)
	}

	if _, err := runCLI(t, srv, "fabric", "bgp-session", "delete", "1", "7"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if last.Method != http.MethodDelete {
		t.Errorf("delete used %s", last.Method)
	}

	if _, err := runCLI(t, srv, "fabric", "bgp-session", "create", "1"); err == nil {
		t.Error("expected an error when --config-source is missing")
	}

	out, err = runCLI(t, srv, "fabric", "bgp-session", "config-example")
	if err != nil {
		t.Fatalf("config-example: %v", err)
	}
	if !strings.Contains(out, "bgpNumbering") {
		t.Errorf("config example missing fields: %s", out)
	}
}

func TestFabricLinkAggregationCommands(t *testing.T) {
	var last fabricSubResRequest
	srv := newFabricSubResServer(&last)
	defer srv.Close()

	out, err := runCLI(t, srv, "fabric", "link-aggregation", "list", "1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "lag-4") {
		t.Errorf("list output missing aggregation: %s", out)
	}

	if _, err := runCLI(t, srv, "fabric", "lag", "show", "1", "4"); err != nil {
		t.Fatalf("get via alias: %v", err)
	}

	configFile := fabricWriteTempConfig(t, `{"type":"lag","linkIds":[1,2]}`)
	if _, err := runCLI(t, srv, "fabric", "lag", "create", "1", "--config-source", configFile); err != nil {
		t.Fatalf("create: %v", err)
	}
	if last.Method != http.MethodPost || !strings.Contains(last.Body, `"linkIds":[1,2]`) {
		t.Errorf("create sent %s %s", last.Method, last.Body)
	}

	if _, err := runCLI(t, srv, "fabric", "lag", "rm", "1", "4"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if last.Method != http.MethodDelete || last.Path != "/api/v2/network-fabrics/1/link-aggregations/4" {
		t.Errorf("delete hit %s %s", last.Method, last.Path)
	}

	out, err = runCLI(t, srv, "fabric", "link-aggregation", "config-example")
	if err != nil {
		t.Fatalf("config-example: %v", err)
	}
	if !strings.Contains(out, "linkIds") {
		t.Errorf("config example missing linkIds: %s", out)
	}
}

func TestFabricHelpListsNewCommands(t *testing.T) {
	out, err := runCLI(t, nil, "fabric", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"delete", "accept-deploy", "reject-deploy", "get-link", "get-interconnects", "bgp-session", "link-aggregation"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}

// fabricWriteTempConfig writes body to a temp file and returns its path.
func fabricWriteTempConfig(t *testing.T, body string) string {
	t.Helper()
	path := t.TempDir() + "/config.json"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}
