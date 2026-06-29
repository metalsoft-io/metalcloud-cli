package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

// ---------------------------------------------------------------------------
// User fixture helpers
// ---------------------------------------------------------------------------

// fullUserJSON is the minimal /api/v2/user JSON that passes strict SDK
// validation (User + UserConfiguration both have required-property checks).
// Individual permission keys are embedded as AdditionalProperties and are
// read by GetUserPermissions to decide which commands remain visible.
const fullUserJSON = `{
  "id":1,"revision":1,"email":"test@test.com","displayName":"Test User",
  "emailStatus":"active","language":"en","brand":"default",
  "isBrandManager":false,"lastLoginTimestamp":"2024-01-01T00:00:00Z",
  "lastLoginType":"password","isBlocked":false,"passwordChangeRequired":false,
  "accessLevel":"admin","isBillable":false,"isTestingMode":false,
  "authenticatorMustChange":false,"authenticatorCreatedTimestamp":"2024-01-01T00:00:00Z",
  "excludeFromReports":false,"isTestAccount":true,"isArchived":false,
  "isDatastorePublisher":false,"provider":"local",
  "passwordLastChangedTimestamp":"2024-01-01T00:00:00Z",
  "franchise":"default","createdTimestamp":"2024-01-01T00:00:00Z",
  "planType":"default","isSuspended":false,"authenticatorEnabled":false,
  "meta":{},
  "config":{
    "revision":1,"displayName":"Test User","emailStatus":"active",
    "language":"en","brand":"default","isBrandManager":false,
    "lastLoginTimestamp":"2024-01-01T00:00:00Z","lastLoginType":"password",
    "isBlocked":false,"passwordChangeRequired":false,"accessLevel":"admin",
    "isBillable":false,"isTestingMode":false,"authenticatorMustChange":false,
    "authenticatorCreatedTimestamp":"2024-01-01T00:00:00Z",
    "excludeFromReports":false,"isTestAccount":true,"isArchived":false,
    "isDatastorePublisher":false,"provider":"local",
    "passwordLastChangedTimestamp":"2024-01-01T00:00:00Z"
  },
  "permissions":{
    "rolePermissions":[],
    "servers_read":true,"server_types_read":true,
    "sites_read":true,"sites_write":true,
    "infrastructures_read":true,"infrastructures_write":true,
    "storage_read":true,"switches_read":true,"events_read":true,
    "firmware_baselines_read":true,"firmware_baselines_write":true,
    "templates_read":true,"templates_write":true,
    "extensions_read":true,"extension_instances_read":true,
    "job_queue_read":true,"network_profiles_read":true,
    "logical_networks_read":true,"vm_types_read":true,
    "users_read":true,"users_write":true,
    "subnets_read":true,"subnets_write":true,
    "accounts_read":true,"accounts_write":true,
    "roles_read":true,"roles_write":true,
    "variables_and_secrets_read":true,"variables_and_secrets_write":true,
    "firmware_upgrade_read":true,"firmware_upgrade_write":true,
    "custom_iso_read":true,"custom_iso_write":true,
    "buckets_read":true,"buckets_write":true,
    "drives_read":true,"drives_write":true,
    "file_shares_read":true,"file_shares_write":true,
    "resource_pools_read":true,"resource_pools_write":true,
    "network_endpoint_groups_read":true,"network_endpoint_groups_write":true,
    "network_device_configuration_templates_read":true,
    "device_auth_providers_read":true,"device_auth_providers_write":true,
    "server_default_credentials_read":true,"server_default_credentials_write":true,
    "cleanup_policies_read":true,"cleanup_policies_write":true,
    "server_instances_read":true,"server_instances_write":true,
    "server_instance_groups_read":true,"server_instance_groups_write":true,
    "vm_pools_read":true,"vm_pools_write":true,
    "vm_instance_groups_read":true,"vm_instance_groups_write":true
  }
}`

// allPerms lists the common permissions used across tests; newMux grants them all.
var allPerms = []string{
	"servers_read", "server_types_read", "sites_read",
	"infrastructures_read", "storage_read", "switches_read", "events_read",
	"firmware_baselines_read", "firmware_upgrade_read",
	"templates_read", "extensions_read", "extension_instances_read",
	"job_queue_read", "network_profiles_read", "logical_networks_read",
	"users_read", "subnets_read", "accounts_read", "roles_read",
	"variables_and_secrets_read", "custom_iso_read",
	"buckets_read", "drives_read", "file_shares_read", "resource_pools_read",
	"network_endpoint_groups_read",
	"network_device_configuration_templates_read",
	"device_auth_providers_read",
	"server_default_credentials_read",
}

// ---------------------------------------------------------------------------
// JSON helpers
// ---------------------------------------------------------------------------

// jsonResponse writes v as JSON with status to w.
func jsonResponse(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// paginatedMeta returns a meta block for single-page list responses.
func paginatedMeta() map[string]interface{} {
	return map[string]interface{}{
		"currentPage":  1,
		"totalPages":   1,
		"itemsPerPage": 100,
	}
}

// paginatedList wraps items in the standard paginated envelope (map form).
func paginatedList(items ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": items,
		"meta": map[string]interface{}{
			"currentPage":  1,
			"totalPages":   1,
			"itemsPerPage": 100,
			"totalItems":   len(items),
		},
	}
}

// paginatedListJSON wraps raw JSON item strings in the standard paginated
// envelope and returns the result as a JSON string. Used by test files that
// build mockSrv route maps (map[string]string) with pre-encoded JSON items.
func paginatedListJSON(items ...string) string {
	dataItems := make([]string, 0, len(items))
	for _, item := range items {
		dataItems = append(dataItems, item)
	}
	data := "["
	for i, item := range dataItems {
		if i > 0 {
			data += ","
		}
		data += item
	}
	data += "]"
	return `{"data":` + data + `,"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100,"totalItems":` +
		func() string {
			n := len(items)
			if n == 0 {
				return "0"
			}
			digits := ""
			for n > 0 {
				digits = string(rune('0'+n%10)) + digits
				n /= 10
			}
			return digits
		}() + `}}`
}

// ---------------------------------------------------------------------------
// Test server constructors
// ---------------------------------------------------------------------------

// newTestServer builds an httptest.Server whose mux handles /api/v2/user
// with the full fixture (all permissions) plus any extra handler provided.
// AllowDevelop is set true so version is never validated.
func newTestServer(t *testing.T, extra http.Handler) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fullUserJSON))
	})
	if extra != nil {
		mux.Handle("/", extra)
	}
	return httptest.NewServer(mux)
}

// newMux builds a ServeMux that handles /api/v2/user with the given permission
// keys plus any routes registered by extra. Use httptest.NewServer(newMux(...))
// to create a server.
func newMux(userPerms []string, extra func(mux *http.ServeMux)) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fullUserJSON))
	})
	if extra != nil {
		extra(mux)
	}
	return mux
}

// mockSrv creates a test server with /api/v2/user plus routes from the routes
// map (path -> raw JSON body). Used by subnet_test.go and similar.
func mockSrv(routes map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fullUserJSON))
	})
	for path, body := range routes {
		body := body
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		})
	}
	return httptest.NewServer(mux)
}

// ---------------------------------------------------------------------------
// CLI execution helpers
// ---------------------------------------------------------------------------

// captureStdout redirects os.Stdout around fn and returns what was written.
// The formatter uses fmt.Printf so rootCmd.SetOut is insufficient.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

// runCLI sets up viper for srv (nil means empty endpoint) and executes rootCmd
// with args. Returns captured stdout and the execution error.
func runCLI(t *testing.T, srv *httptest.Server, args ...string) (string, error) {
	t.Helper()

	if srv != nil {
		viper.Set(system.ConfigEndpoint, srv.URL)
	} else {
		viper.Set(system.ConfigEndpoint, "")
	}
	viper.Set(system.ConfigApiKey, "test-key")
	viper.Set(formatter.ConfigFormat, "json")
	system.AllowDevelop = true

	var errBuf bytes.Buffer
	rootCmd.SetErr(&errBuf)

	var execErr error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(args)
		execErr = rootCmd.Execute()
	})
	return out, execErr
}

// execCLI is an alias used by subnet_test.go and similar files that were
// written before the captureStdout pattern was established. It sets up viper
// and runs rootCmd, returning buf output (cobra's SetOut/SetErr) and error.
// Note: JSON/YAML output from the formatter goes to os.Stdout (not captured
// here); execCLI is only useful for commands that write to cobra's output or
// for error-path tests.
func execCLI(t *testing.T, srv *httptest.Server, args ...string) (string, error) {
	t.Helper()
	if srv != nil {
		viper.Set(system.ConfigEndpoint, srv.URL)
	}
	viper.Set(system.ConfigApiKey, "test-key")
	viper.Set(formatter.ConfigFormat, "json")
	system.AllowDevelop = true

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

