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

// emailTplFixture satisfies every required property of sdk.EmailTemplate:
// id, revision, name, subject, text and html.
func emailTplFixture() map[string]interface{} {
	return map[string]interface{}{
		"id": 4, "revision": 9, "name": "infrastructure-deployed",
		"subject": "Deployment finished", "description": "Sent after a deploy",
		"text": "plain body", "html": "<p>html body</p>",
	}
}

type emailTplRecordedRequest struct {
	Method  string
	Path    string
	IfMatch string
	Body    string
}

func emailTplServer(t *testing.T, last *emailTplRecordedRequest) *httptest.Server {
	t.Helper()
	record := func(r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*last = emailTplRecordedRequest{
			Method: r.Method, Path: r.URL.Path,
			IfMatch: r.Header.Get("If-Match"), Body: string(body),
		}
	}
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/email-templates", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, emailTplFixture())
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(emailTplFixture()))
		})
		mux.HandleFunc("/api/v2/email-templates/infrastructure-deployed", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, emailTplFixture())
		})
	}))
}

func TestEmailTemplateList_HappyPath(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "email-template", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "infrastructure-deployed") {
		t.Errorf("list output missing the template: %s", out)
	}
}

func TestEmailTemplateGet_HappyPath(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "email-templates", "get", "infrastructure-deployed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Deployment finished") {
		t.Errorf("get output missing the subject: %s", out)
	}
	if last.Path != "/api/v2/email-templates/infrastructure-deployed" {
		t.Errorf("templates are addressed by name, got %s", last.Path)
	}
}

func TestEmailTemplateConfigExample_HappyPath(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	out, err := runCLI(t, srv, "email-template", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"html"`) {
		t.Errorf("config example missing the html field: %s", out)
	}
}

func TestEmailTemplateCreate_FromFile(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "template.json")
	if err := os.WriteFile(path, []byte(`{"name":"welcome","subject":"Hi","text":"t","html":"<p>h</p>"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "email-template", "create", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodPost || !strings.Contains(last.Body, `"name":"welcome"`) {
		t.Errorf("unexpected create request %s %s", last.Method, last.Body)
	}
}

func TestEmailTemplateUpdate_SendsIfMatch(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "update.json")
	if err := os.WriteFile(path, []byte(`{"subject":"New subject"}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if _, err := runCLI(t, srv, "email-template", "update", "infrastructure-deployed", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodPatch {
		t.Errorf("update should use PATCH, got %s", last.Method)
	}
	if last.IfMatch != "9" {
		t.Errorf("update should send the fetched revision as If-Match, got %q", last.IfMatch)
	}
}

func TestEmailTemplateDelete_HappyPath(t *testing.T) {
	var last emailTplRecordedRequest
	srv := emailTplServer(t, &last)
	defer srv.Close()

	if _, err := runCLI(t, srv, "email-template", "rm", "infrastructure-deployed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", last.Method)
	}
}

func TestEmailTemplateCommand_Wiring(t *testing.T) {
	for _, sub := range []string{"list", "get", "config-example", "create", "update", "delete"} {
		found, _, err := emailTemplateCmd.Find([]string{sub})
		if err != nil || found == emailTemplateCmd {
			t.Fatalf("email-template %s is not registered", sub)
		}
		if found.Annotations[system.REQUIRED_PERMISSION] == "" {
			t.Errorf("email-template %s must carry a permission annotation", sub)
		}
		if !found.SilenceUsage {
			t.Errorf("email-template %s must set SilenceUsage", sub)
		}
	}
}
