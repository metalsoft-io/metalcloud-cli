package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
)

// cfgDocument is a trimmed platform configuration document.
var cfgDocument = map[string]interface{}{
	"platform": map[string]interface{}{"brandName": "MetalSoft"},
}

type cfgRecordedRequest struct {
	Method string
	Query  string
	Body   string
}

// cfgServer serves /api/v2/config and /api/v2/config/platform, recording the
// last request so the tests can assert on method, query and body.
func cfgServer(t *testing.T, last *cfgRecordedRequest) *httptest.Server {
	t.Helper()
	record := func(r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*last = cfgRecordedRequest{Method: r.Method, Query: r.URL.RawQuery, Body: string(body)}
	}
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/config", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, cfgDocument)
		})
		mux.HandleFunc("/api/v2/config/platform", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, cfgDocument["platform"])
		})
	}))
}

func TestConfigurationGet_HappyPath(t *testing.T) {
	var last cfgRecordedRequest
	srv := cfgServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "configuration", "get")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "MetalSoft") {
		t.Errorf("get output missing the configuration: %s", out)
	}
	if last.Method != http.MethodGet || last.Query != "" {
		t.Errorf("unexpected request %s ?%s", last.Method, last.Query)
	}
}

func TestConfigurationGet_WithService(t *testing.T) {
	var last cfgRecordedRequest
	srv := cfgServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "configuration", "get", "platform"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Query != "filter=platform" {
		t.Errorf("expected filter query parameter, got %q", last.Query)
	}
}

func TestConfigurationReplaceAndUpdate(t *testing.T) {
	var last cfgRecordedRequest
	srv := cfgServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "platform.json")
	if err := os.WriteFile(path, []byte(`{"brandName":"Acme"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "configuration", "replace", "platform", "--config-source", path); err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	if last.Method != http.MethodPut {
		t.Errorf("replace should use PUT, got %s", last.Method)
	}
	var body map[string]interface{}
	_ = json.Unmarshal([]byte(last.Body), &body)
	if body["brandName"] != "Acme" {
		t.Errorf("replace body should carry the document, got %s", last.Body)
	}

	if _, err := runCLI(t, srv, "configuration", "update", "platform", "--config-source", path); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if last.Method != http.MethodPatch {
		t.Errorf("update should use PATCH, got %s", last.Method)
	}
}

func TestConfigurationWrite_UnknownService(t *testing.T) {
	var last cfgRecordedRequest
	srv := cfgServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "x.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := runCLI(t, srv, "configuration", "update", "nope", "--config-source", path)
	if err == nil || !strings.Contains(err.Error(), "unknown configuration service") {
		t.Fatalf("expected an unknown service error, got %v", err)
	}
}

func TestConfigurationWrite_RequiresConfigSource(t *testing.T) {
	var last cfgRecordedRequest
	srv := cfgServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "configuration", "replace", "platform"); err == nil {
		t.Fatal("replace without --config-source should fail")
	}
}

// TestConfigurationCommand_Wiring checks the sub-commands, their permissions and
// that no alias collides with the root --config flag's name.
func TestConfigurationCommand_Wiring(t *testing.T) {
	for _, sub := range []string{"get", "replace", "update"} {
		found, _, err := configurationCmd.Find([]string{sub})
		if err != nil || found == configurationCmd {
			t.Fatalf("configuration %s is not registered", sub)
		}
		if found.Annotations[system.REQUIRED_PERMISSION] == "" {
			t.Errorf("configuration %s must carry a permission annotation", sub)
		}
		if !found.SilenceUsage {
			t.Errorf("configuration %s must set SilenceUsage", sub)
		}
	}

	for _, alias := range configurationCmd.Aliases {
		if alias == "config" {
			t.Error("'config' must not be used as an alias: the root command owns the --config flag")
		}
	}
}
