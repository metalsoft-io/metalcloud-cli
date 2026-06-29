package job

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestJobList(t *testing.T) {
	const jobListPage1 = `{
		"data": [
			{"jobId": 1, "type": "deploy", "status": "finished", "functionName": "fn1",
			 "callCount": 1, "retryMax": 3, "retryCount": 0, "retryMinSeconds": 5,
			 "requiresConfirmation": false, "options": {},
			 "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
			 "links": []}
		],
		"meta": {"currentPage": 1, "totalPages": 1, "itemsPerPage": 100}
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs": testutils.RawHandler(http.StatusOK, jobListPage1),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobList(ctx, ListFlags{}); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError500", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobList(ctx, ListFlags{}); err == nil {
			t.Error("expected error for HTTP 500, got nil")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobList(ctx, ListFlags{}); err != nil {
			t.Errorf("expected nil error for empty list, got: %v", err)
		}
	})
}

func TestJobGet(t *testing.T) {
	const jobSingle = `{
		"jobId": 1, "type": "deploy", "status": "finished", "functionName": "fn1",
		"callCount": 1, "retryMax": 3, "retryCount": 0, "retryMinSeconds": 5,
		"requiresConfirmation": false, "options": {},
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
		"links": []
	}`

	t.Run("HappyPath", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs/1": testutils.RawHandler(http.StatusOK, jobSingle),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGet(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{
			"/api/v2/jobs/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGet(ctx, "99"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
		defer ts.Close()

		ctx := testutils.SetupTestContext(ts.URL)
		if err := JobGet(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})
}
