package infrastructure

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// infraExtrasRecorder captures the requests received by the mock server.
type infraExtrasRecorder struct {
	mu   sync.Mutex
	reqs []infraExtrasRequest
}

type infraExtrasRequest struct {
	Method string
	Path   string
	Body   string
}

func (r *infraExtrasRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, infraExtrasRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *infraExtrasRecorder) last() infraExtrasRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func infraExtrasItem() map[string]interface{} {
	return map[string]interface{}{
		"id":               12,
		"revision":         3,
		"label":            "my-infra",
		"siteId":           1,
		"datacenterName":   "dc1",
		"userIdOwner":      1,
		"serviceStatus":    "active",
		"designIsLocked":   0,
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-02T00:00:00Z",
		"config": map[string]interface{}{
			"revision":         3,
			"label":            "my-infra",
			"deployStatus":     "not_started",
			"deployType":       "create",
			"datacenterName":   "dc1",
			"siteId":           1,
			"updatedTimestamp": "2024-01-02T00:00:00Z",
		},
	}
}

func TestInfrastructureUpdateMetadata(t *testing.T) {
	rec := &infraExtrasRecorder{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/infrastructures", rec.wrap(
		testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{infraExtrasItem()}, 1, 1))))
	mux.HandleFunc("/api/v2/infrastructures/12/metadata", rec.wrap(
		testutils.JSONHandler(http.StatusOK, infraExtrasItem())))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureUpdateMetadata(ctx, "12", []byte(`{"name":"Production","tags":["prod"]}`)); err != nil {
		t.Fatalf("InfrastructureUpdateMetadata: expected nil error, got: %v", err)
	}

	last := rec.last()
	if last.Method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", last.Method)
	}
	if !strings.HasSuffix(last.Path, "/metadata") {
		t.Errorf("expected the metadata endpoint, got %s", last.Path)
	}
	for _, want := range []string{`"name":"Production"`, `"tags":["prod"]`} {
		if !strings.Contains(last.Body, want) {
			t.Errorf("expected request body to contain %s, got %s", want, last.Body)
		}
	}
}

func TestInfrastructureUpdateMetadata_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK,
			testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureUpdateMetadata(ctx, "99", []byte(`{"name":"x"}`)); err == nil {
		t.Error("expected error for a missing infrastructure, got nil")
	}
}

func infraUtilizationSummary() map[string]interface{} {
	return map[string]interface{}{
		"startTimestamp": "2024-01-01T00:00:00Z",
		"endTimestamp":   "2024-01-08T00:00:00Z",
		"infrastructures": map[string]interface{}{
			"12": map[string]interface{}{
				"infrastructureId":            12,
				"infrastructureLabel":         "my-infra",
				"infrastructureServiceStatus": "active",
				// The API returns a list here, while the typed SDK model
				// declares a string - hence the raw parsing.
				"tags": []interface{}{"prod"},
			},
		},
		"resourceUtilization": map[string]interface{}{
			"instance": map[string]interface{}{
				"quantity":        10,
				"measurementUnit": "hours",
			},
			"drive": map[string]interface{}{
				"quantity":        5,
				"measurementUnit": "GB-hours",
			},
		},
		"internet": map[string]interface{}{
			"upload": map[string]interface{}{
				"quantity":        1,
				"measurementUnit": "GB",
			},
			"download": map[string]interface{}{
				"quantity":        2,
				"measurementUnit": "GB",
			},
		},
	}
}

func TestInfrastructureGetUtilizationSummary_SendsBody(t *testing.T) {
	rec := &infraExtrasRecorder{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/infrastructures/actions/get/resource-utilization-summarized", rec.wrap(
		testutils.JSONHandler(http.StatusOK, infraUtilizationSummary())))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetUtilizationSummary(ctx, 3, start, end, []int{12}); err != nil {
		t.Fatalf("InfrastructureGetUtilizationSummary: expected nil error, got: %v", err)
	}

	last := rec.last()
	if last.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", last.Method)
	}
	// The body is always sent; an unset body would be serialised as `null`
	// and rejected by the API with 400.
	for _, want := range []string{`"userIdOwner":3`, `"startTimestamp":"2024-01-01T00:00:00Z"`, `"infrastructureIds":[12]`} {
		if !strings.Contains(last.Body, want) {
			t.Errorf("expected request body to contain %s, got %s", want, last.Body)
		}
	}
}

func TestInfrastructureGetUtilizationSummary_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/actions/get/resource-utilization-summarized": testutils.ErrorHandler(http.StatusBadRequest, "bad request"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	err := InfrastructureGetUtilizationSummary(ctx, 3, time.Now(), time.Now(), nil)
	if err == nil {
		t.Error("expected error for HTTP 400, got nil")
	}
}

func TestBuildUtilizationSummaryRecords(t *testing.T) {
	// The tags list is the field that breaks the typed SDK model (it declares
	// a string), which is why the response is parsed raw.
	summary := map[string]any{
		"infrastructures": map[string]any{
			"12": map[string]any{
				"infrastructureId":            float64(12),
				"infrastructureLabel":         "my-infra",
				"infrastructureServiceStatus": "active",
				"tags":                        []any{"prod"},
			},
		},
		"resourceUtilization": map[string]any{
			"instance": map[string]any{"quantity": float64(10), "measurementUnit": "hours"},
			"drive":    map[string]any{"quantity": float64(5), "measurementUnit": "GB-hours"},
		},
		"internet": map[string]any{
			"upload":   map[string]any{"quantity": float64(1), "measurementUnit": "GB"},
			"download": map[string]any{"quantity": float64(2), "measurementUnit": "GB"},
		},
	}

	records := buildUtilizationSummaryRecords(summary)
	if len(records) != 5 {
		t.Fatalf("expected 5 rows (1 infrastructure + 2 resources + 2 internet), got %d", len(records))
	}
	if records[0].Kind != "Infrastructure" || records[0].InfrastructureLabel != "my-infra" {
		t.Errorf("unexpected first row: %+v", records[0])
	}
	// Resource rows are sorted by key: drive before instance.
	if records[1].Resource != "drive" || records[2].Resource != "instance" {
		t.Errorf("expected resource rows sorted by name, got %q then %q", records[1].Resource, records[2].Resource)
	}
	if records[3].Kind != "Internet" || records[3].Resource != "upload" {
		t.Errorf("unexpected internet upload row: %+v", records[3])
	}
}
