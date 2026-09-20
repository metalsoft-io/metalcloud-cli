package ai

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// aiResponse satisfies both required properties of sdk.AIGenerateResponse.
var aiResponse = map[string]any{
	"result": "There are 3 available servers.",
	"steps":  "1. list servers\n2. filter by status",
}

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Body   string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{Method: req.Method, Path: req.URL.Path, Body: string(body)})
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
		"/api/v2/ai/generate": rec.wrap(testutils.JSONHandler(http.StatusOK, aiResponse)),
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

// TestAIGenerateSendsBody pins the gotcha that the request body must always be
// set: skipping the optional SDK setter sends a literal "null" the API rejects.
func TestAIGenerateSendsBody(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return AIGenerate(ctx, "dc1", "which servers are available?") })

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/ai/generate" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(got.Body), &body); err != nil {
		t.Fatalf("body should be a JSON object, got %q", got.Body)
	}
	if body["datacenter"] != "dc1" || body["prompt"] != "which servers are available?" {
		t.Errorf("body should carry datacenter and prompt, got %s", got.Body)
	}
	if !strings.Contains(out, "There are 3 available servers.") {
		t.Errorf("json output should carry the result: %s", out)
	}
}

// TestAIGenerateTextOutputIsProse pins that table formats print the answer as
// plain text rather than squeezing prose into a table cell.
func TestAIGenerateTextOutputIsProse(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	viper.Set(formatter.ConfigFormat, "text")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := run(t, func() error { return AIGenerate(ctx, "dc1", "hello") })
	if strings.Contains(out, "|") {
		t.Errorf("text output should not be a table: %s", out)
	}
	if !strings.HasPrefix(out, "There are 3 available servers.") {
		t.Errorf("text output should start with the answer: %s", out)
	}
	if !strings.Contains(out, "1. list servers") {
		t.Errorf("text output should include the steps: %s", out)
	}
}

func TestAIGenerateRejectsEmptyPrompt(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := AIGenerate(ctx, "dc1", "   "); err == nil || !strings.Contains(err.Error(), "prompt cannot be empty") {
		t.Fatalf("expected an empty prompt error, got %v", err)
	}
}

func TestAIGenerateRequiresDatacenter(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := AIGenerate(ctx, "", "hello"); err == nil || !strings.Contains(err.Error(), "datacenter is required") {
		t.Fatalf("expected a missing datacenter error, got %v", err)
	}
}

// TestAIGenerateFromConfig verifies a config document supplies both fields, and
// that a document carrying only the prompt still works with --datacenter.
func TestAIGenerateFromConfig(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return AIGenerateFromConfig(ctx, "", []byte(`{"datacenter":"dc2","prompt":"list infrastructures"}`))
	})
	if !strings.Contains(rec.last().Body, `"datacenter":"dc2"`) {
		t.Errorf("config document should supply the datacenter, got %s", rec.last().Body)
	}

	run(t, func() error {
		return AIGenerateFromConfig(ctx, "dc3", []byte(`{"prompt":"list servers"}`))
	})
	body := rec.last().Body
	if !strings.Contains(body, `"datacenter":"dc3"`) || !strings.Contains(body, `"prompt":"list servers"`) {
		t.Errorf("the --datacenter flag should fill in what the document omits, got %s", body)
	}
}

func TestAIGenerateServerError(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/ai/generate": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := AIGenerate(ctx, "dc1", "hello"); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}
