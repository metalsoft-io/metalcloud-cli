package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Required: jobId, type, status, functionName, callCount, retryMax, retryCount,
//           retryMinSeconds, requiresConfirmation, options, createdTimestamp,
//           updatedTimestamp, links
func jobFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"jobId": id, "type": "deploy", "status": "finished",
		"functionName": "TestFunc", "callCount": 1, "retryMax": 3,
		"retryCount": 0, "retryMinSeconds": 5, "requiresConfirmation": false,
		"options":          map[string]interface{}{},
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"links":            []interface{}{},
	}
}

// Required: id, type, description, createdTimestamp
func jobGroupFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "type": "deploy", "description": "Deploy job group",
		"createdTimestamp": "2024-01-01T00:00:00Z", "links": []interface{}{},
	}
}

// Required: id, label, functionName, params, schedule, waitForCompletion (float32),
//           lifetimeSeconds, disabled (float32)
func cronJobFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "label": "nightly-task", "functionName": "CleanupFunc",
		"params": []interface{}{}, "schedule": "0 2 * * *",
		"waitForCompletion": 0, "lifetimeSeconds": 3600, "disabled": 0,
		"links": []interface{}{},
	}
}

// --- job list ---

func TestJobList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "job", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestJobList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "job", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestJobList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "job", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- job-group list ---

func TestJobGroupList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/job-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobGroupFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "job-group", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestJobGroupList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/job-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobGroupFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "job-group", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestJobGroupList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "job-group", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- cron-job list ---

func TestCronJobList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/cron-jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(cronJobFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "cron-job", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCronJobList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/cron-jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(cronJobFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "cron-job", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestCronJobList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "cron-job", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- format tests ---

func TestJobList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "job", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestJobGroupList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/job-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(jobGroupFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "job-group", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestCronJobList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/cron-jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(cronJobFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "cron-job", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

// --- cron-job create ---

func TestCronJobCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/cron-jobs", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(cronJobFixture(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "cronjob-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"label":"nightly","functionName":"CleanupFunc","params":[],"schedule":"0 2 * * *","waitForCompletion":0,"lifetimeSeconds":3600,"disabled":0}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "cron-job", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- cron-job delete ---

func TestCronJobDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/cron-jobs/1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "cron-job", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}
