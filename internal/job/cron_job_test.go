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

const cronJobListResponseLinksArray = `{
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
	]
}`

const cronJobListResponseLinksMap = `{
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
			"links": {"self": "http://example.com/cron-jobs/1"}
		}
	]
}`

const cronJobSingleResponseLinksArray = `{
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

// isCronJobListPath returns true when the path ends at /cron-jobs with no
// further segments (i.e. it is the collection endpoint, not a single-item
// endpoint like /cron-jobs/1).
func isCronJobListPath(path string) bool {
	idx := strings.LastIndex(path, "/cron-jobs")
	if idx == -1 {
		return false
	}
	// Everything after "/cron-jobs" must be empty (possibly a trailing slash).
	suffix := strings.TrimSuffix(path[idx+len("/cron-jobs"):], "/")
	return suffix == ""
}

func TestCronJobList(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/cron-jobs") && isCronJobListPath(r.URL.Path) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(cronJobListResponseLinksArray))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := CronJobList(ctx)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("LinksAsMap", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/cron-jobs") && isCronJobListPath(r.URL.Path) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(cronJobListResponseLinksMap))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := CronJobList(ctx)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/cron-jobs") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "internal server error"}`))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := CronJobList(ctx)
		if err == nil {
			t.Error("expected an error for HTTP 500, got nil")
		}
	})
}

func TestCronJobGet(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/cron-jobs/1") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(cronJobSingleResponseLinksArray))
				return
			}
			http.NotFound(w, r)
		}))
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		err := CronJobGet(ctx, "1")
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}
