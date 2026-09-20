package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
)

// extSysFixture satisfies every required property of sdk.ExternalSystem:
// id, label, name, annotations, revision, createdAt and updatedAt.
func extSysFixture() map[string]interface{} {
	return map[string]interface{}{
		"id": "12", "label": "my-external-system", "name": "My External System",
		"annotations": map[string]interface{}{"owner": "platform-team"},
		"revision":    "3",
		"createdAt":   "2024-01-01T00:00:00Z",
		"updatedAt":   "2024-01-02T00:00:00Z",
	}
}

type extSysRecordedRequest struct {
	Method  string
	Path    string
	IfMatch string
	Body    string
}

func extSysServer(t *testing.T, last *extSysRecordedRequest) *httptest.Server {
	t.Helper()
	record := func(r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*last = extSysRecordedRequest{
			Method: r.Method, Path: r.URL.Path,
			IfMatch: r.Header.Get("If-Match"), Body: string(body),
		}
	}
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/external-systems", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, extSysFixture())
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(extSysFixture()))
		})
		mux.HandleFunc("/api/v2/external-systems/12", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, extSysFixture())
		})
	}))
}

func TestExternalSystemList_HappyPath(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "external-system", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "my-external-system") {
		t.Errorf("list output missing the external system: %s", out)
	}
}

func TestExternalSystemGet_HappyPath(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "ext-system", "get", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "My External System") {
		t.Errorf("get output missing the name: %s", out)
	}
	if last.Path != "/api/v2/external-systems/12" {
		t.Errorf("unexpected get path %s", last.Path)
	}
}

func TestExternalSystemGet_InvalidId(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "external-system", "get", "abc"); err == nil {
		t.Fatal("a non-numeric ID should fail")
	}
}

func TestExternalSystemConfigExample_HappyPath(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "external-system", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label"`) {
		t.Errorf("config example missing the label field: %s", out)
	}
}

func TestExternalSystemCreate_FromFlags(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	_, err := runCLI(t, srv, "external-system", "create",
		"--label", "new-system", "--name", "New System",
		"--annotations", `{"owner":"platform-team"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodPost {
		t.Fatalf("expected POST, got %s", last.Method)
	}
	if !strings.Contains(last.Body, `"label":"new-system"`) ||
		!strings.Contains(last.Body, `"name":"New System"`) ||
		!strings.Contains(last.Body, `"owner":"platform-team"`) {
		t.Errorf("create body should carry the flag values, got %s", last.Body)
	}
}

func TestExternalSystemCreate_FromFile(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "external-system.json")
	if err := os.WriteFile(path, []byte(`{"label":"from-file","name":"From File"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "external-system", "create", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(last.Body, `"label":"from-file"`) {
		t.Errorf("create body should come from the file, got %s", last.Body)
	}
}

func TestExternalSystemCreate_RequiresLabelOrConfigSource(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "external-system", "create"); err == nil {
		t.Fatal("create without --label or --config-source should fail")
	}
}

func TestExternalSystemCreate_InvalidAnnotations(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	_, err := runCLI(t, srv, "external-system", "create", "--label", "x", "--annotations", "not-json")
	if err == nil || !strings.Contains(err.Error(), "invalid --annotations") {
		t.Fatalf("expected an invalid annotations error, got %v", err)
	}
}

func TestExternalSystemUpdate_SendsIfMatch(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "update.json")
	if err := os.WriteFile(path, []byte(`{"name":"Renamed"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "external-system", "update", "12", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodPatch {
		t.Errorf("update should use PATCH, got %s", last.Method)
	}
	if last.IfMatch != "3" {
		t.Errorf("update should send the fetched revision as If-Match, got %q", last.IfMatch)
	}
}

func TestExternalSystemDelete_HappyPath(t *testing.T) {
	var last extSysRecordedRequest
	srv := extSysServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "external-system", "rm", "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", last.Method)
	}
}

func TestExternalSystemCommand_Wiring(t *testing.T) {
	for _, sub := range []string{"list", "get", "config-example", "create", "update", "delete"} {
		found, _, err := externalSystemCmd.Find([]string{sub})
		if err != nil || found == externalSystemCmd {
			t.Fatalf("external-system %s is not registered", sub)
		}
		if found.Annotations[system.REQUIRED_PERMISSION] == "" {
			t.Errorf("external-system %s must carry a permission annotation", sub)
		}
		if !found.SilenceUsage {
			t.Errorf("external-system %s must set SilenceUsage", sub)
		}
	}
}
