package job

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

// jobActionRecorder captures the requests received by the mock server so the
// tests can assert on method, path, headers and body.
type jobActionRecorder struct {
	mu   sync.Mutex
	reqs []jobActionRequest
}

type jobActionRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *jobActionRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, jobActionRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Header: req.Header.Clone(),
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *jobActionRecorder) last() jobActionRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newJobActionServer(t *testing.T, rec *jobActionRecorder, path string, status int) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(path, rec.wrap(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	return httptest.NewServer(mux)
}

func TestJobSkip(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/skip", http.StatusNoContent)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobSkip(ctx, "12"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		last := rec.last()
		if last.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", last.Method)
		}
		// The skip action carries no body: the SDK request has no body setter.
		if strings.TrimSpace(last.Body) != "" {
			t.Errorf("expected empty body, got %q", last.Body)
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ctx := testutils.SetupTestContext("http://127.0.0.1:1")
		if err := JobSkip(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/skip", http.StatusConflict)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobSkip(ctx, "12"); err == nil {
			t.Error("expected error for HTTP 409, got nil")
		}
	})
}

func TestJobRetry(t *testing.T) {
	t.Run("DefaultBody", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/retry", http.StatusNoContent)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobRetry(ctx, "12", false); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		last := rec.last()
		if last.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", last.Method)
		}
		// sdk.NewJobRetryInfo() defaults retryEvenIfSuccessful to false, so the
		// body is always a real object and never a literal `null`.
		if got := strings.TrimSpace(last.Body); got != `{"retryEvenIfSuccessful":false}` {
			t.Errorf("unexpected request body: %s", got)
		}
	})

	t.Run("RetryEvenIfSuccessful", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/retry", http.StatusNoContent)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobRetry(ctx, "12", true); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if got := strings.TrimSpace(rec.last().Body); got != `{"retryEvenIfSuccessful":true}` {
			t.Errorf("unexpected request body: %s", got)
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ctx := testutils.SetupTestContext("http://127.0.0.1:1")
		if err := JobRetry(ctx, "not-a-number", false); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})
}

func TestJobIssueCommand(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/issue-command", http.StatusNoContent)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		commandInfo := sdk.JobCommandInfo{
			Command:            sdk.PtrString("kill"),
			ExecuteImmediately: sdk.PtrBool(true),
		}
		if err := JobIssueCommand(ctx, "12", commandInfo); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		last := rec.last()
		if last.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", last.Method)
		}
		if got := strings.TrimSpace(last.Body); got != `{"command":"kill","executeImmediately":true}` {
			t.Errorf("unexpected request body: %s", got)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		rec := &jobActionRecorder{}
		ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/issue-command", http.StatusBadRequest)
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobIssueCommand(ctx, "12", sdk.JobCommandInfo{}); err == nil {
			t.Error("expected error for HTTP 400, got nil")
		}
	})
}

func TestJobKill(t *testing.T) {
	rec := &jobActionRecorder{}
	ts := newJobActionServer(t, rec, "/api/v2/jobs/12/actions/issue-command", http.StatusNoContent)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := JobKill(ctx, "12"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if got := strings.TrimSpace(rec.last().Body); got != `{"command":"kill","executeImmediately":true}` {
		t.Errorf("unexpected request body: %s", got)
	}
}

func TestJobGetArchived(t *testing.T) {
	// The archived job fixture deliberately omits `links`, which the typed SDK
	// model requires: JobGetArchived parses the body raw.
	const archivedJob = `{
		"jobId": 12, "type": "deploy", "status": "finished", "functionName": "fn1",
		"infrastructureId": 3, "jobGroupId": 4,
		"callCount": 1, "retryMax": 3, "retryCount": 0, "retryMinSeconds": 5,
		"requiresConfirmation": false, "options": {},
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z"
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs/archive/12": testutils.RawHandler(http.StatusOK, archivedJob),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGetArchived(ctx, "12"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs/archive/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGetArchived(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ctx := testutils.SetupTestContext("http://127.0.0.1:1")
		if err := JobGetArchived(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})
}

func TestJobGroupStatistics(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/job-groups/7/statistics": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
				"groupId":                 7,
				"groupType":               "infrastructure_deploy",
				"groupCreatedTimestamp":   "2024-01-01T00:00:00Z",
				"groupCompletedTimestamp": "2024-01-01T01:00:00Z",
				"jobsThrownError":         1,
				"jobsCompleted":           4,
				"jobsTotal":               5,
			}),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGroupStatistics(ctx, "7"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ctx := testutils.SetupTestContext("http://127.0.0.1:1")
		if err := JobGroupStatistics(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})
}

func TestScheduledJobSupportedFunctions(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/scheduled-jobs/supported-functions": testutils.JSONHandler(http.StatusOK, []interface{}{
			map[string]interface{}{
				"name":        "cleanup",
				"description": "Cleans up stale resources",
				"paramsSchema": map[string]interface{}{
					"properties": map[string]interface{}{
						"olderThanDays": map[string]interface{}{"type": "number"},
						"dryRun":        map[string]interface{}{"type": "boolean"},
					},
				},
			},
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ScheduledJobSupportedFunctions(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFormatScheduledJobParams(t *testing.T) {
	if got := formatScheduledJobParams(nil); got != "" {
		t.Errorf("expected empty string for nil schema, got %q", got)
	}

	if got := formatScheduledJobParams(map[string]interface{}{"type": "object"}); got != "" {
		t.Errorf("expected empty string for schema without properties, got %q", got)
	}

	got := formatScheduledJobParams(map[string]interface{}{
		"properties": map[string]interface{}{"b": 1, "a": 2},
	})
	if got != "a, b" {
		t.Errorf("expected sorted property names, got %q", got)
	}
}
