package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// accountItem deliberately omits "limits": the real API does not return it even
// though the SDK Account model marks it required (regression: account archive
// failed with "no value given for required property limits").
var accountItem = map[string]interface{}{
	"id": 1.0, "name": "acme", "revision": 1.0,
	"config": map[string]interface{}{"revision": 1.0, "name": "acme"},
}

func newAccountTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestAccountArchive(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
		mux.HandleFunc("/api/v2/accounts/1/actions/archive", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err := runCLI(t, srv, "account", "archive", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccountList(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "account", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountListAlias(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "accounts", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountGet(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "account", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountGetRequiresArg(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "account", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestAccountHelp(t *testing.T) {
	out, err := runCLI(t, nil, "account", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "account") {
		t.Errorf("expected help output to contain 'account', got: %s", out)
	}
}

func TestAccountCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, accountItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "account-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"acme"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "account", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestAccountUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
		mux.HandleFunc("/api/v2/accounts/1/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, accountItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "account-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"acme-updated"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "account", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestAccountList_Formats(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "account", "list")
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

// TestAccountList_Archived verifies --archived sends filter.archived to the API
// (the API excludes archived accounts by default).
func TestAccountList_Archived(t *testing.T) {
	gotFilter := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			gotFilter = r.URL.Query().Get("filter.archived")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
	}))
	defer srv.Close()

	// No-flag case first: the in-process test harness shares cobra flag
	// globals across invocations, so --archived would leak into later runs.
	if _, err := runCLI(t, srv, "account", "list"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFilter != "" {
		t.Errorf("expected no filter.archived without --archived, got %q", gotFilter)
	}

	if _, err := runCLI(t, srv, "account", "list", "--archived"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFilter != "$eq:1" {
		t.Errorf("expected filter.archived=$eq:1, got %q", gotFilter)
	}
}

// acctQuotaLimits returns a QuotaProfileLimits payload with every required
// property populated.
func acctQuotaLimits(serverGroups int) map[string]interface{} {
	return map[string]interface{}{
		"infrastructureServerGroupMaxCount":                serverGroups,
		"infrastructureDriveMaxCount":                      10,
		"infrastructureFileShareMaxCount":                  10,
		"infrastructureBucketMaxCount":                     10,
		"infrastructureVmInstanceGroupMaxCount":            10,
		"infrastructureContainerInstanceGroupMaxCount":     10,
		"serverGroupInstancesMaxCount":                     10,
		"serverGroupInstancesMinCount":                     1,
		"vmInstanceGroupVmInstancesMaxCount":               10,
		"containerInstanceGroupContainerInstancesMaxCount": 10,
		"vmInstanceMaxDiskSizeMbytes":                      1024,
		"containerInstanceMaxDiskSizeMbytes":               1024,
		"driveMaxSizeMbytes":                               1024,
		"driveMinSizeMbytes":                               1,
		"fileShareMinSizeGb":                               1,
		"fileShareMaxSizeGb":                               100,
		"bucketMinSizeGb":                                  1,
		"bucketMaxSizeGb":                                  100,
		"showOperatingSystemImagesTab":                     true,
		"showTemplateAssetsView":                           true,
		"userResourceServerTypeNameToMaxCount":             map[string]interface{}{},
		"userSshKeysCountMax":                              5,
		"showLegacyPages":                                  false,
		"showEliChatBot":                                   false,
		"enableCustomRaidConfiguration":                    false,
		"enableInfrastructureVmInstance":                   true,
		"enableInfrastructureContainerInstance":            true,
		"enableInfrastructureExtensions":                   true,
		"allowedInfrastructureExtensions":                  []interface{}{},
		"allowedServerTypes":                               []interface{}{},
		"allowedSites":                                     []interface{}{},
		"allowedLogicalNetworkProfiles":                    []interface{}{},
		"allowedPreCreatedLogicalNetworks":                 []interface{}{},
	}
}

func TestAccountUnarchive(t *testing.T) {
	gotIfMatch := ""
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
		mux.HandleFunc("/api/v2/accounts/1/actions/unarchive", func(w http.ResponseWriter, r *http.Request) {
			gotIfMatch = r.Header.Get("If-Match")
			jsonResponse(w, http.StatusOK, accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	if _, err := runCLI(t, srv, "account", "unarchive", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotIfMatch != "1" {
		t.Errorf("expected If-Match 1, got %q", gotIfMatch)
	}
}

func TestAccountConfigGet(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts/1/config", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"revision": 1, "name": "acme", "code": "ACME",
			})
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "account", "config", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountQuotaBreakdown(t *testing.T) {
	gotIncludeUsage := ""
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts/1/quota-limits-breakdown", func(w http.ResponseWriter, r *http.Request) {
			gotIncludeUsage = r.URL.Query().Get("includeUsage")
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"effective": acctQuotaLimits(5),
				"account":   acctQuotaLimits(7),
			})
		})
	}))
	defer srv.Close()

	// No-flag case first: runCLI resets the shared cobra flags, but asserting
	// the default here keeps the two cases independent.
	if _, err := runCLI(t, srv, "account", "quota-breakdown", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotIncludeUsage != "" {
		t.Errorf("expected no includeUsage without the flag, got %q", gotIncludeUsage)
	}

	if _, err := runCLI(t, srv, "account", "quota", "1", "--include-usage"); err != nil {
		t.Fatalf("alias quota: unexpected error: %v", err)
	}
	if gotIncludeUsage != "true" {
		t.Errorf("expected includeUsage=true, got %q", gotIncludeUsage)
	}
}

// TestAccountQuotaBreakdown_Formats checks the flattened table renders in the
// text formats while JSON and YAML keep the nested object.
func TestAccountQuotaBreakdown_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts/1/quota-limits-breakdown", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"effective": acctQuotaLimits(5),
				"account":   acctQuotaLimits(7),
			})
		})
	}))
	defer srv.Close()

	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "account", "quota-breakdown", "1")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, "infrastructureServerGroupMaxCount") {
				t.Errorf("format csv: expected one row per limit, got: %s", out)
			}
		})
	}
}

func TestAccountQuotaBreakdownRequiresArg(t *testing.T) {
	if _, err := runCLI(t, nil, "account", "quota-breakdown"); err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}
