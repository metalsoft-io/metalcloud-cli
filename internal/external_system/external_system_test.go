package external_system

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// externalSystemItem satisfies every required property of sdk.ExternalSystem:
// id, label, name, annotations, revision, createdAt and updatedAt. Note that id
// and revision are strings in the model while the endpoints take an int64 id.
var externalSystemItem = map[string]any{
	"id": "12", "label": "my-external-system", "name": "My External System",
	"annotations": map[string]any{"owner": "platform-team"},
	"revision":    "3",
	"createdAt":   "2024-01-01T00:00:00Z",
	"updatedAt":   "2024-01-02T00:00:00Z",
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
		"/api/v2/external-systems": rec.wrap(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(http.StatusCreated, externalSystemItem)(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{externalSystemItem}, 1, 1))(w, r)
		}),
		"/api/v2/external-systems/12": rec.wrap(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(http.StatusOK, externalSystemItem)(w, r)
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

func TestExternalSystemListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return ExternalSystemList(ctx, []string{"my-external-system"}) })
	if !strings.Contains(out, "my-external-system") {
		t.Errorf("list output missing the external system: %s", out)
	}

	out = run(t, func() error { return ExternalSystemGet(ctx, "12") })
	if !strings.Contains(out, `"name":"My External System"`) {
		t.Errorf("get output missing the name: %s", out)
	}
	if got := rec.last().Path; got != "/api/v2/external-systems/12" {
		t.Errorf("unexpected get path %s", got)
	}
}

func TestExternalSystemGetInvalidId(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	err := ExternalSystemGet(ctx, "not-a-number")
	if err == nil || !strings.Contains(err.Error(), "invalid external system ID") {
		t.Fatalf("expected an invalid ID error, got %v", err)
	}
}

func TestExternalSystemConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("")

	out := run(t, func() error { return ExternalSystemConfigExample(ctx) })
	if !strings.Contains(out, `"label"`) || !strings.Contains(out, `"annotations"`) {
		t.Errorf("config example should carry the create fields: %s", out)
	}
}

func TestExternalSystemCreate(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return ExternalSystemCreate(ctx, sdk.CreateExternalSystem{Label: "new-system", Name: "New System"})
	})

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/external-systems" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"label":"new-system"`) {
		t.Errorf("create body should carry the label, got %s", got.Body)
	}
}

func TestExternalSystemUpdateSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error { return ExternalSystemUpdate(ctx, "12", []byte(`{"name":"Renamed"}`)) })

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/external-systems/12" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if ifMatch := got.Header.Get("If-Match"); ifMatch != "3" {
		t.Errorf("expected If-Match 3 (the fetched revision), got %q", ifMatch)
	}
	if !strings.Contains(got.Body, `"name":"Renamed"`) {
		t.Errorf("update body should carry the changed fields, got %s", got.Body)
	}
}

func TestExternalSystemDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ExternalSystemDelete(ctx, "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := rec.last()
	if got.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", got.Method)
	}
	// The SDK exposes no IfMatch setter on the delete request.
	if ifMatch := got.Header.Get("If-Match"); ifMatch != "" {
		t.Errorf("delete should not send If-Match, got %q", ifMatch)
	}
}

func TestExternalSystemGetServerError(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/external-systems/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ExternalSystemGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for a missing external system, got nil")
	}
}
