package resource_pool

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// rpExtrasRecorder captures the requests received by the mock server.
type rpExtrasRecorder struct {
	mu   sync.Mutex
	reqs []rpExtrasRequest
}

type rpExtrasRequest struct {
	Method string
	Path   string
	Body   string
}

func (r *rpExtrasRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, rpExtrasRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *rpExtrasRecorder) last() rpExtrasRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func rpExtrasItem(id int) map[string]interface{} {
	return map[string]interface{}{
		"resourcePoolId":          id,
		"resourcePoolLabel":       "pool-a",
		"resourcePoolDescription": "Pool A",
	}
}

func TestResourcePoolUpdate(t *testing.T) {
	rec := &rpExtrasRecorder{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/resource-pools/123", rec.wrap(
		testutils.JSONHandler(http.StatusOK, rpExtrasItem(123))))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolUpdate(ctx, "123", []byte(`{"resourcePoolLabel":"Production Pool"}`)); err != nil {
		t.Fatalf("ResourcePoolUpdate: expected nil error, got: %v", err)
	}

	last := rec.last()
	if last.Method != http.MethodPut {
		t.Errorf("expected PUT, got %s", last.Method)
	}
	if got := strings.TrimSpace(last.Body); got != `{"resourcePoolLabel":"Production Pool"}` {
		t.Errorf("unexpected request body: %s", got)
	}
}

func TestResourcePoolUpdate_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ResourcePoolUpdate(ctx, "not-a-number", []byte(`{}`)); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestResourcePoolUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/123": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolUpdate(ctx, "123", []byte(`{}`)); err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}

func TestResourcePoolListForUser(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/42": testutils.JSONHandler(http.StatusOK,
			[]interface{}{rpExtrasItem(1), rpExtrasItem(2)}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolListForUser(ctx, "42"); err != nil {
		t.Errorf("ResourcePoolListForUser: expected nil error, got: %v", err)
	}
}

func TestResourcePoolListForUser_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ResourcePoolListForUser(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid user ID, got nil")
	}
}

func TestResourcePoolListForUser_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/42": testutils.ErrorHandler(http.StatusForbidden, "forbidden"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolListForUser(ctx, "42"); err == nil {
		t.Error("expected error for HTTP 403, got nil")
	}
}
