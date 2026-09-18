package job

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

func init() {
	viper.Set(formatter.ConfigFormat, "text")
}

const scheduledJobListResponseLinksArray = `{
	"data": [
		{
			"id": 1,
			"label": "test-cron",
			"description": "test description",
			"functionName": "testFunc",
			"params": [],
			"schedule": "*/5 * * * *",
			"waitForCompletion": 0,
			"lifetimeSeconds": 3600,
			"disabled": 0,
			"links": []
		}
	],
	"meta": {
		"currentPage": 1,
		"totalPages": 1,
		"itemsPerPage": 100,
		"totalItems": 1
	}
}`

const scheduledJobListResponseLinksMap = `{
	"data": [
		{
			"id": 1,
			"label": "test-cron",
			"description": "test description",
			"functionName": "testFunc",
			"params": [],
			"schedule": "*/5 * * * *",
			"waitForCompletion": 0,
			"lifetimeSeconds": 3600,
			"disabled": 0,
			"links": [{"rel": "self", "href": "http://example.com/scheduled-jobs/1"}]
		}
	],
	"meta": {
		"currentPage": 1,
		"totalPages": 1,
		"itemsPerPage": 100,
		"totalItems": 1
	}
}`

const scheduledJobSingleResponseLinksArray = `{
	"id": 1,
	"label": "test-cron",
	"description": "test description",
	"functionName": "testFunc",
	"params": [],
	"schedule": "*/5 * * * *",
	"waitForCompletion": 0,
	"lifetimeSeconds": 3600,
	"disabled": 0,
	"links": []
}`

// isScheduledJobListPath returns true when the path ends at /scheduled-jobs with no
// further segments (i.e. it is the collection endpoint, not a single-item
// endpoint like /scheduled-jobs/1).
func isScheduledJobListPath(path string) bool {
	idx := strings.LastIndex(path, "/scheduled-jobs")
	if idx == -1 {
		return false
	}
	// Everything after "/scheduled-jobs" must be empty (possibly a trailing slash).
	suffix := strings.TrimSuffix(path[idx+len("/scheduled-jobs"):], "/")
	return suffix == ""
}

func TestScheduledJobList(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/scheduled-jobs") && isScheduledJobListPath(r.URL.Path) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(scheduledJobListResponseLinksArray))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ScheduledJobList(ctx)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("LinksAsMap", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/scheduled-jobs") && isScheduledJobListPath(r.URL.Path) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(scheduledJobListResponseLinksMap))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ScheduledJobList(ctx)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/scheduled-jobs") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ScheduledJobList(ctx)
		if err == nil {
			t.Error("expected an error for HTTP 500, got nil")
		}
	})
}

func TestScheduledJobGet(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/scheduled-jobs/1") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(scheduledJobSingleResponseLinksArray))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := ScheduledJobGet(ctx, "1")
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}
