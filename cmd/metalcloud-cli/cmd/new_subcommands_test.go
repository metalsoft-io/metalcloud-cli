package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Shared recorder
// ---------------------------------------------------------------------------

// newSubRecorder records the requests the CLI sends to the mock API so the
// tests can assert on method, path and body.
type newSubRecorder struct {
	mu   sync.Mutex
	reqs []newSubRequest
}

type newSubRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *newSubRecorder) handler(status int, body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		raw, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, newSubRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Header: req.Header.Clone(),
			Body:   string(raw),
		})
		r.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}
}

func (r *newSubRecorder) find(method, path string) *newSubRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.reqs {
		if r.reqs[i].Method == method && r.reqs[i].Path == path {
			return &r.reqs[i]
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// storage
// ---------------------------------------------------------------------------

var storIfaceItem = map[string]interface{}{
	"id":            3.0,
	"revision":      7.0,
	"storageId":     1.0,
	"name":          "if-a",
	"nodeIds":       []string{"node-1"},
	"protocols":     []string{"iscsi"},
	"isUplink":      true,
	"useForDeploys": false,
}

var storScopedUserItem = map[string]interface{}{
	"id":               5.0,
	"revision":         2.0,
	"storageId":        1.0,
	"infrastructureId": 10.0,
	"username":         "scoped-user",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-02T00:00:00Z",
}

func TestStorageInterfaceCommands(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages/1/interfaces", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(storIfaceItem))
		})
		mux.HandleFunc("/api/v2/storages/1/interfaces/3", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, storIfaceItem)
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "storage", "interfaces", "1")
	if err != nil {
		t.Fatalf("storage interfaces: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "if-a") {
		t.Fatalf("expected if-a in output, got: %s", out)
	}

	out, err = runCLI(t, srv, "storage", "interface", "1", "3")
	if err != nil {
		t.Fatalf("storage interface: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "if-a") {
		t.Fatalf("expected if-a in output, got: %s", out)
	}
}

func TestStorageScopedAccessUserCommands(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages/1/scoped-access-users", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(storScopedUserItem))
		})
		mux.HandleFunc("/api/v2/storages/1/scoped-access-users/5", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, storScopedUserItem)
		})
		mux.HandleFunc("/api/v2/storages/1/scoped-access-users/5/credentials", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"username": "scoped-user", "password": "secret"})
		})
	}))
	defer srv.Close()

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"list", []string{"storage", "scoped-access-users", "1"}, "scoped-user"},
		{"get", []string{"storage", "scoped-access-user", "1", "5"}, "scoped-user"},
		{"credentials", []string{"storage", "scoped-access-user-credentials", "1", "5"}, "secret"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := runCLI(t, srv, tc.args...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("expected %q in output, got: %s", tc.want, out)
			}
		})
	}
}

func TestStorageStatisticsCommand(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"maintenanceCount": 1, "experimentalCount": 0, "lowSpaceCount": 2,
				"usedSpace": 400, "freeSpace": 600,
				"types":        map[string]interface{}{"netapp": 2},
				"pendingCount": 0, "readyCount": 1, "activeCount": 3,
			})
		})
		mux.HandleFunc("/api/v2/storages/1/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"totalSpaceGB": 100, "usedSpaceGB": 40, "freeSpaceGB": 60,
			})
		})
	}))
	defer srv.Close()

	t.Run("global", func(t *testing.T) {
		out, err := runCLI(t, srv, "storage", "statistics")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if !strings.Contains(out, "lowSpaceCount") {
			t.Fatalf("expected global statistics in output, got: %s", out)
		}
	})

	t.Run("perStorage", func(t *testing.T) {
		out, err := runCLI(t, srv, "storage", "statistics", "1")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if !strings.Contains(out, "totalSpaceGB") {
			t.Fatalf("expected per-storage statistics in output, got: %s", out)
		}
	})
}

func TestStorageUpdateCommand_SendsIfMatch(t *testing.T) {
	rec := &newSubRecorder{}
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages/1", rec.handler(http.StatusOK, storageItem))
	}))
	defer srv.Close()

	configPath := filepath.Join(t.TempDir(), "storage-update.json")
	if err := os.WriteFile(configPath, []byte(`{"inMaintenance":1}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if _, err := runCLI(t, srv, "storage", "update", "1", "--config-source", configPath); err != nil {
		t.Fatalf("storage update: expected no error, got: %v", err)
	}

	patch := rec.find(http.MethodPatch, "/api/v2/storages/1")
	if patch == nil {
		t.Fatal("expected a PATCH on /api/v2/storages/1")
	}
	if got := patch.Header.Get("If-Match"); got != "1" {
		t.Errorf("expected If-Match 1, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// jobs
// ---------------------------------------------------------------------------

func TestJobCmdIssueCommand(t *testing.T) {
	rec := &newSubRecorder{}
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/jobs/12/actions/issue-command", rec.handler(http.StatusNoContent, nil))
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "job", "issue-command", "12", "--command", "kill", "--execute-immediately"); err != nil {
		t.Fatalf("job issue-command: expected no error, got: %v", err)
	}

	req := rec.find(http.MethodPost, "/api/v2/jobs/12/actions/issue-command")
	if req == nil {
		t.Fatal("expected a POST on the issue-command action")
	}
	if got := strings.TrimSpace(req.Body); got != `{"command":"kill","executeImmediately":true}` {
		t.Errorf("unexpected request body: %s", got)
	}
}

func TestJobCmdGetArchived(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/jobs/archive/12", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"jobId": 12, "type": "deploy", "status": "finished", "functionName": "fn1",
				"infrastructureId": 3, "jobGroupId": 4,
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "job", "get-archived", "12")
	if err != nil {
		t.Fatalf("job get-archived: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "fn1") {
		t.Fatalf("expected fn1 in output, got: %s", out)
	}
}

func TestJobCmdScheduledJobFunctions(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/scheduled-jobs/supported-functions", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []interface{}{
				map[string]interface{}{"name": "cleanup", "description": "Cleans up"},
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "job", "scheduled-job-functions")
	if err != nil {
		t.Fatalf("job scheduled-job-functions: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "cleanup") {
		t.Fatalf("expected cleanup in output, got: %s", out)
	}
}

func TestJobGroupCmdStatistics(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/job-groups/7/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"groupId": 7, "groupType": "infrastructure_deploy",
				"groupCreatedTimestamp": "2024-01-01T00:00:00Z", "groupCompletedTimestamp": "2024-01-01T01:00:00Z",
				"jobsThrownError": 1, "jobsCompleted": 4, "jobsTotal": 5,
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "job-group", "statistics", "7")
	if err != nil {
		t.Fatalf("job-group statistics: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "infrastructure_deploy") {
		t.Fatalf("expected the group type in output, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// resource pools
// ---------------------------------------------------------------------------

func TestResourcePoolCmdListForUser(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/resource-pools/user/42", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []interface{}{
				map[string]interface{}{
					"resourcePoolId": 1, "resourcePoolLabel": "pool-a", "resourcePoolDescription": "Pool A",
				},
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "resource-pool", "list-for-user", "42")
	if err != nil {
		t.Fatalf("resource-pool list-for-user: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "pool-a") {
		t.Fatalf("expected pool-a in output, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// sites
// ---------------------------------------------------------------------------

func TestSiteCmdStatisticsAndRegistryUrls(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/sites/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"siteControllerSeenAliveStatus": map[string]interface{}{
					"may_be_offline": 1, "offline": 2, "was_seen_connected_very_recently": 5,
				},
				"sitesTotalCount": 8, "sitesActiveCount": 6,
				// Per-site resource counters, the shape the live API returns.
				"sitesResourceCount": []interface{}{
					map[string]interface{}{"id": 1, "serversCount": 4, "infrastructuresCount": 36},
				},
			})
		})
		mux.HandleFunc("/api/v2/sites/controllers/actions/get/registry-urls", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []string{"registry.metalsoft.dev"})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "site", "statistics")
	if err != nil {
		t.Fatalf("site statistics: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "sitesTotalCount") {
		t.Fatalf("expected the site statistics in output, got: %s", out)
	}

	out, err = runCLI(t, srv, "site", "registry-urls")
	if err != nil {
		t.Fatalf("site registry-urls: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "registry.metalsoft.dev") {
		t.Fatalf("expected the registry URL in output, got: %s", out)
	}
}

// ---------------------------------------------------------------------------
// dns
// ---------------------------------------------------------------------------

var dnsExtrasRecordSet = map[string]interface{}{
	"id": 456, "name": "www", "type": "A", "records": []string{"10.0.0.1"}, "ttl": 300,
	"status": "active", "zoneName": "example.com", "siteId": 1, "infrastructureId": 1,
	"zoneId": 123, "revision": 1, "createdBy": 1, "createdAt": "2024-01-01T00:00:00Z",
	"links": []interface{}{},
}

func TestDnsZoneCmdNameserversAndRecords(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/dns-zones/123/nameservers", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []string{"ns1.example.com"})
		})
		mux.HandleFunc("/api/v2/dns-zones/123/recordsets/456", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, dnsExtrasRecordSet)
		})
		mux.HandleFunc("/api/v2/dns-recordsets", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(dnsExtrasRecordSet))
		})
	}))
	defer srv.Close()

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"nameservers", []string{"dns-zone", "nameservers", "123"}, "ns1.example.com"},
		{"record", []string{"dns-zone", "record", "123", "456"}, "www"},
		{"globalRecords", []string{"dns-zone", "records"}, "www"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := runCLI(t, srv, tc.args...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("expected %q in output, got: %s", tc.want, out)
			}
		})
	}
}
