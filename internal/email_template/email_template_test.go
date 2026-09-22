package email_template

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// emailTemplateItem satisfies every required property of sdk.EmailTemplate:
// id, revision, name, subject, text and html.
var emailTemplateItem = map[string]any{
	"id": 4, "revision": 9, "name": "infrastructure-deployed",
	"subject": "Deployment finished", "description": "Sent after a deploy",
	"text": "plain body", "html": "<p>html body</p>",
}

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{Method: req.Method, Path: req.URL.Path, Header: req.Header.Clone(), Body: string(body)})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *recorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newServer(rec *recorder) *httptest.Server {
	return testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/email-templates": rec.wrap(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(http.StatusCreated, emailTemplateItem)(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{emailTemplateItem}, 1, 1))(w, r)
		}),
		"/api/v2/email-templates/infrastructure-deployed": rec.wrap(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, emailTemplateItem)(w, r)
		}),
	})
}

func run(t *testing.T, fn func() error) string {
	t.Helper()
	var err error
	out := testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestEmailTemplateListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EmailTemplateList(ctx, nil) })
	if !strings.Contains(out, "infrastructure-deployed") {
		t.Errorf("list output missing the template: %s", out)
	}

	out = run(t, func() error { return EmailTemplateGet(ctx, "infrastructure-deployed") })
	if !strings.Contains(out, `"subject":"Deployment finished"`) {
		t.Errorf("get output missing the subject: %s", out)
	}
	if got := rec.last().Path; got != "/api/v2/email-templates/infrastructure-deployed" {
		t.Errorf("templates are addressed by name, got path %s", got)
	}
}

func TestEmailTemplateGetRejectsEmptyName(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EmailTemplateGet(ctx, "   "); err == nil || !strings.Contains(err.Error(), "cannot be empty") {
		t.Fatalf("expected an empty-name error, got %v", err)
	}
}

func TestEmailTemplateConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("")

	out := run(t, func() error { return EmailTemplateConfigExample(ctx) })
	if !strings.Contains(out, `"name"`) || !strings.Contains(out, `"html"`) {
		t.Errorf("config example should carry the create fields: %s", out)
	}
}

func TestEmailTemplateCreate(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	config := []byte(`{"name":"welcome","subject":"Hi","text":"t","html":"<p>h</p>"}`)
	run(t, func() error { return EmailTemplateCreate(ctx, config) })

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/email-templates" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(got.Body), &body)
	if body["name"] != "welcome" {
		t.Errorf("create body should carry the template name, got %s", got.Body)
	}
}

func TestEmailTemplateUpdateSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return EmailTemplateUpdate(ctx, "infrastructure-deployed", []byte(`{"subject":"New subject"}`))
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/email-templates/infrastructure-deployed" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if ifMatch := got.Header.Get("If-Match"); ifMatch != "9" {
		t.Errorf("expected If-Match 9 (the fetched revision), got %q", ifMatch)
	}
	if !strings.Contains(got.Body, `"subject":"New subject"`) {
		t.Errorf("update body should carry the changed fields, got %s", got.Body)
	}
}

func TestEmailTemplateDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EmailTemplateDelete(ctx, "infrastructure-deployed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", got.Method)
	}
}

func TestEmailTemplateGetServerError(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/email-templates/missing": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EmailTemplateGet(ctx, "missing"); err == nil {
		t.Fatal("expected an error for a missing template, got nil")
	}
}
