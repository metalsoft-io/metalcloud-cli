//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var returnedRe = regexp.MustCompile(`Returned (\d+) out of (\d+) records`)

// assertListCommand is a shared helper for list commands without flags that
// promise ALL records.  It verifies exit 0, valid JSON on stdout, the
// pagination summary on stderr, and — critically — that returned == total
// (X == Y in "Returned X out of Y records").  X < Y means the command
// silently truncated, which is the regression class this guards against.
func assertListCommand(t *testing.T, args ...string) {
	t.Helper()
	stdout, stderr, err := runCLI(t, args...)
	if err != nil {
		t.Errorf("%v failed: %v\nstdout: %s\nstderr: %s", args, err, stdout, stderr)
		return
	}
	if !json.Valid([]byte(stdout)) {
		t.Errorf("%v: stdout is not valid JSON: %s", args, stdout)
	}
	m := returnedRe.FindStringSubmatch(stderr)
	if m == nil {
		t.Errorf("%v: expected pagination summary in stderr, got: %s", args, stderr)
		return
	}
	if m[1] != m[2] {
		t.Errorf("%v: claims to return ALL records but returned %s out of %s", args, m[1], m[2])
	}
}

// assertListCommandNoPage is for list commands that don't use FetchAllPages
// and therefore produce no pagination summary on stderr.
func assertListCommandNoPage(t *testing.T, args ...string) {
	t.Helper()
	stdout, stderr, err := runCLI(t, args...)
	if err != nil {
		t.Errorf("%v failed: %v\nstdout: %s\nstderr: %s", args, err, stdout, stderr)
		return
	}
	if !json.Valid([]byte(stdout)) {
		t.Errorf("%v: stdout is not valid JSON: %s", args, stdout)
	}
	_ = stderr
}

// ---------------------------------------------------------------------------
// Global list commands (no required positional argument)
// ---------------------------------------------------------------------------

func TestSubnetList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "subnet", "list", "-f", "json")
}

func TestAccountList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "account", "list", "-f", "json")
}

func TestFabricList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "fabric", "list", "-f", "json")
}

func TestDnsZoneList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "dns-zone", "list", "-f", "json")
}

func TestSecretList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "secret", "list", "-f", "json")
}

func TestVariableList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "variable", "list", "-f", "json")
}

func TestRoleList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommandNoPage(t, "-x", "role", "list", "-f", "json")
}

func TestSiteList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "site", "list", "-f", "json")
}

func TestInfrastructureList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "infrastructure", "list", "-f", "json")
}

func TestStorageList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "storage", "list", "-f", "json")
}

func TestResourcePoolList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "resource-pool", "list", "-f", "json")
}

func TestServerTypeList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "server-type", "list", "-f", "json")
}

func TestNetworkDeviceList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "network-device", "list", "-f", "json")
}

func TestEndpointList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "endpoint", "list", "-f", "json")
}

func TestEventList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// Use --limit to avoid fetching millions of event records.
	assertListCommandNoPage(t, "-x", "event", "list", "--limit", "10", "-f", "json")
}

func TestExtensionList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "extension", "list", "-f", "json")
}

func TestLogicalNetworkList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "logical-network", "list", "-f", "json")
}

func TestLogicalNetworkProfileList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "logical-network-profile", "list", "-f", "json")
}

func TestOsTemplateList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "os-template", "list", "-f", "json")
}

func TestTemplateAssetList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "template-asset", "list", "-f", "json")
}

func TestFirmwareCatalogList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "firmware-catalog", "list", "-f", "json")
}

func TestFirmwareBaselineList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "firmware-baseline", "list", "-f", "json")
}

func TestFirmwareBinaryList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "firmware-binary", "list", "-f", "json")
}

func TestFirmwarePolicyList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	stdout, stderr, err := runCLI(t, "-x", "firmware-policy", "list", "-f", "json")
	if err != nil {
		if strings.Contains(stderr, "404") {
			t.Skip("firmware-policy endpoint not available on this deployment")
		}
		t.Errorf("firmware-policy list failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
		return
	}
	if !json.Valid([]byte(stdout)) {
		t.Errorf("stdout is not valid JSON: %s", stdout)
	}
}

func TestServerList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// ServerList uses raw-body parsing to work around SDK bug — no FetchAllPages, no pagination summary.
	assertListCommandNoPage(t, "-x", "server", "list", "-f", "json")
}

func TestServerDefaultCredentialsList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "server-default-credentials", "list", "-f", "json")
}

func TestServerCleanupPolicyList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "server-cleanup-policy", "list", "-f", "json")
}

func TestServerRegistrationProfileList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "server-registration-profile", "list", "-f", "json")
}

func TestCronJobList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// CronJobList uses raw-body parsing — no FetchAllPages, no pagination summary.
	assertListCommandNoPage(t, "-x", "cron-job", "list", "-f", "json")
}

func TestJobList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// Use --limit to avoid fetching millions of job records.
	assertListCommandNoPage(t, "-x", "job", "list", "--limit", "10", "-f", "json")
}

func TestJobGroupList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// Use --limit to avoid fetching millions of job-group records.
	assertListCommandNoPage(t, "-x", "job-group", "list", "--limit", "10", "-f", "json")
}

func TestVmPoolList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "vm-pool", "list", "-f", "json")
}

func TestVmTypeList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "vm-type", "list", "-f", "json")
}

func TestNetworkDeviceConfigurationTemplateList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "network-configuration", "device-template", "list", "-f", "json")
}

func TestNetworkDeviceLinkAggregationTemplateList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "network-configuration", "link-aggregation-template", "list", "-f", "json")
}

func TestNetworkDeviceDefaultSecretsList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "network-device", "default-secrets", "list", "-f", "json")
}

func TestCustomIsoList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "custom-iso", "list", "-f", "json")
}

func TestSiteDeviceAuthProviderList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "site", "device-auth-provider", "list", "-f", "json")
}

func TestPermissionList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// permission list calls a non-paginated endpoint.
	assertListCommandNoPage(t, "-x", "permission", "list", "-f", "json")
}

func TestUserList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertListCommand(t, "-x", "user", "list", "-f", "json")
}

func TestAuthLdapMappingList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// LDAP mapping list calls a non-paginated endpoint.
	assertListCommandNoPage(t, "-x", "auth", "ldap", "mapping-list", "-f", "json")
}

// ---------------------------------------------------------------------------
// Infrastructure-scoped list commands
// These require an infrastructure ID/label.  They are skipped when no
// infrastructures exist.
// ---------------------------------------------------------------------------

func TestBucketList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping bucket list")
	}
	assertListCommand(t, "-x", "bucket", "list", infraID, "-f", "json")
}

func TestDriveList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping drive list")
	}
	assertListCommand(t, "-x", "drive", "list", infraID, "-f", "json")
}

func TestFileShareList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping file-share list")
	}
	assertListCommand(t, "-x", "file-share", "list", infraID, "-f", "json")
}

func TestServerInstanceGroupList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping server-instance-group list")
	}
	assertListCommand(t, "-x", "server-instance-group", "list", infraID, "-f", "json")
}

func TestServerInstanceList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping server-instance list")
	}
	assertListCommand(t, "-x", "server-instance", "list", infraID, "-f", "json")
}

func TestVmInstanceGroupList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping vm-instance-group list")
	}
	assertListCommand(t, "-x", "vm-instance-group", "list", infraID, "-f", "json")
}

func TestExtensionInstanceList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping extension-instance list")
	}
	assertListCommand(t, "-x", "extension-instance", "list", infraID, "-f", "json")
}

// ---------------------------------------------------------------------------
// Pagination flag variations: --page, --limit, --page+--limit
// Covers every command that supports these flags to ensure single-page and
// multi-page-limited paths return valid JSON without panicking.
// ---------------------------------------------------------------------------

// assertPaginated verifies that a list command with explicit pagination flags
// exits 0, returns a valid JSON array, and never returns more records than
// the requested --limit (when one is present in args).
func assertPaginated(t *testing.T, args ...string) {
	t.Helper()
	stdout, stderr, err := runCLI(t, args...)
	if err != nil {
		t.Errorf("%v failed: %v\nstdout: %s\nstderr: %s", args, err, stdout, stderr)
		return
	}
	var items []map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(stdout), &items); jsonErr != nil {
		t.Errorf("%v: stdout is not a JSON array: %v\nstdout: %s", args, jsonErr, stdout)
		return
	}
	// Find --limit value in args, if any.
	for i, a := range args {
		if a == "--limit" && i+1 < len(args) {
			var limit int
			fmt.Sscanf(args[i+1], "%d", &limit)
			if limit > 0 && len(items) > limit {
				t.Errorf("%v: returned %d records, more than --limit %d", args, len(items), limit)
			}
		}
	}
	_ = stderr
}

func TestEventList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "event", "list", "--limit", "5", "-f", "json")
}
func TestEventList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "event", "list", "--page", "1", "--limit", "5", "-f", "json")
}
func TestEventList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "event", "list", "--limit", "150", "-f", "json")
}

// TestEventList_PageWithLimitOver100_Integration: --page 3 --limit 120 must
// return up to 120 records (window 241..360), not the API's 100/page cap.
func TestEventList_PageWithLimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	stdout, stderr, err := runCLI(t, "-x", "event", "list", "--page", "3", "--limit", "120", "-f", "json")
	if err != nil {
		t.Fatalf("failed: %v\nstderr: %s", err, stderr)
	}
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &items); err != nil {
		t.Fatalf("stdout is not a JSON array: %v", err)
	}
	// On a deployment with >360 events the window must be exactly 120 records.
	if len(items) != 120 {
		t.Errorf("expected 120 records for --page 3 --limit 120, got %d", len(items))
	}
}

// Same check for limit-only: exactly N records returned when more exist.
func TestEventList_LimitExactCount_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	stdout, stderr, err := runCLI(t, "-x", "event", "list", "--limit", "150", "-f", "json")
	if err != nil {
		t.Fatalf("failed: %v\nstderr: %s", err, stderr)
	}
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &items); err != nil {
		t.Fatalf("stdout is not a JSON array: %v", err)
	}
	if len(items) != 150 {
		t.Errorf("expected exactly 150 records for --limit 150, got %d", len(items))
	}
}

func TestLogicalNetworkList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "logical-network", "list", "--limit", "5", "-f", "json")
}
func TestLogicalNetworkList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "logical-network", "list", "--page", "1", "--limit", "5", "-f", "json")
}
func TestLogicalNetworkList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "logical-network", "list", "--limit", "150", "-f", "json")
}
func TestLogicalNetworkList_PageWithLimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "logical-network", "list", "--page", "2", "--limit", "120", "-f", "json")
}

func TestServerDefaultCredentialsList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "server-default-credentials", "list", "--limit", "5", "-f", "json")
}
func TestServerDefaultCredentialsList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "server-default-credentials", "list", "--page", "1", "-f", "json")
}
func TestServerDefaultCredentialsList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "server-default-credentials", "list", "--limit", "150", "-f", "json")
}
func TestServerDefaultCredentialsList_PageLimit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "server-default-credentials", "list", "--page", "1", "--limit", "5", "-f", "json")
}

func TestVmTypeList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "vm-type", "list", "--limit", "5", "-f", "json")
}
func TestVmTypeList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "vm-type", "list", "--page", "1", "-f", "json")
}
func TestVmTypeList_PageLimit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "vm-type", "list", "--page", "1", "--limit", "5", "-f", "json")
}
func TestVmTypeList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "vm-type", "list", "--limit", "150", "-f", "json")
}

func TestResourcePoolList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "resource-pool", "list", "--limit", "5", "-f", "json")
}
func TestResourcePoolList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "resource-pool", "list", "--page", "1", "-f", "json")
}
func TestResourcePoolList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "resource-pool", "list", "--limit", "150", "-f", "json")
}

func TestNetworkDeviceDefaultSecretsList_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "network-device", "default-secrets", "list", "--limit", "5", "-f", "json")
}
func TestNetworkDeviceDefaultSecretsList_Page_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "network-device", "default-secrets", "list", "--page", "1", "-f", "json")
}
func TestNetworkDeviceDefaultSecretsList_LimitOver100_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	assertPaginated(t, "-x", "network-device", "default-secrets", "list", "--limit", "150", "-f", "json")
}


// ---------------------------------------------------------------------------
// VM pool sub-commands: VMs, cluster hosts, cluster host VMs
// ---------------------------------------------------------------------------

func TestVmPoolGetVMs_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	poolID := firstVMPoolID(t)
	if poolID == "" {
		t.Skip("no VM pools available; skipping vm-pool vms")
	}
	assertListCommand(t, "-x", "vm-pool", "vms", poolID, "-f", "json")
}

func TestVmPoolGetVMs_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	poolID := firstVMPoolID(t)
	if poolID == "" {
		t.Skip("no VM pools available; skipping vm-pool vms --limit")
	}
	assertPaginated(t, "-x", "vm-pool", "vms", poolID, "--limit", "5", "-f", "json")
}

func TestVmPoolGetClusterHosts_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	poolID := firstVMPoolID(t)
	if poolID == "" {
		t.Skip("no VM pools available; skipping vm-pool cluster-hosts")
	}
	assertListCommand(t, "-x", "vm-pool", "cluster-hosts", poolID, "-f", "json")
}

func TestVmPoolGetClusterHosts_Limit_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	poolID := firstVMPoolID(t)
	if poolID == "" {
		t.Skip("no VM pools available; skipping vm-pool cluster-hosts --limit")
	}
	assertPaginated(t, "-x", "vm-pool", "cluster-hosts", poolID, "--limit", "5", "-f", "json")
}

// ---------------------------------------------------------------------------
// VMInstanceList
// ---------------------------------------------------------------------------

func TestVmInstanceList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping vm-instance list")
	}
	assertListCommand(t, "-x", "vm-instance", "list", infraID, "-f", "json")
}

// ---------------------------------------------------------------------------
// Drive/file-share snapshots: flat API — no pagination; verify no crash on empty.
// ---------------------------------------------------------------------------

func TestDriveSnapshotList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping drive snapshot list")
	}
	// Drive snapshot list requires a valid drive ID within the infra.
	// We skip rather than fail when no drives exist — the important check is
	// that the command doesn't truncate or crash when it returns results.
	stdout, stderr, err := runCLI(t, "-x", "drive", "snapshot", "list", infraID, "1", "-f", "json")
	if err != nil {
		if strings.Contains(stderr, "404") || strings.Contains(stderr, "not found") {
			t.Skip("no drive with ID 1 in this infra; skipping drive snapshot list")
		}
		t.Fatalf("drive snapshot list failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !json.Valid([]byte(stdout)) {
		t.Errorf("stdout is not valid JSON: %s", stdout)
	}
}

func TestFileShareSnapshotList_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available; skipping file-share snapshot list")
	}
	stdout, stderr, err := runCLI(t, "-x", "file-share", "snapshot", "list", infraID, "1", "-f", "json")
	if err != nil {
		if strings.Contains(stderr, "404") || strings.Contains(stderr, "not found") {
			t.Skip("no file share with ID 1 in this infra; skipping file-share snapshot list")
		}
		t.Fatalf("file-share snapshot list failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !json.Valid([]byte(stdout)) {
		t.Errorf("stdout is not valid JSON: %s", stdout)
	}
}

// assertListFormat runs a list command with a specific output format and verifies
// the output is non-empty. For json format it additionally validates JSON syntax.
func assertListFormat(t *testing.T, format string, args ...string) {
	t.Helper()
	fullArgs := append(args, "-f", format)
	stdout, stderr, err := runCLI(t, fullArgs...)
	if err != nil {
		t.Errorf("%v -f %s: failed: %v\nstdout:%s\nstderr:%s", args, format, err, stdout, stderr)
		return
	}
	// Table-based formats (csv/text/md) legitimately render nothing for an
	// empty result set; json/yaml always produce output.
	if stdout == "" && (format == "json" || format == "yaml") {
		t.Errorf("%v -f %s: empty output", args, format)
	}
	if format == "json" && stdout != "" && !json.Valid([]byte(stdout)) {
		t.Errorf("%v -f %s: invalid JSON: %s", args, format, stdout)
	}
	_ = stderr
}

func TestSubnetList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "subnet", "list")
		})
	}
}

func TestAccountList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "account", "list")
		})
	}
}

func TestFabricList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "fabric", "list")
		})
	}
}

func TestDnsZoneList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "dns-zone", "list")
		})
	}
}

func TestSecretList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "secret", "list")
		})
	}
}

func TestVariableList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "variable", "list")
		})
	}
}

func TestRoleList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "role", "list")
		})
	}
}

func TestSiteList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "site", "list")
		})
	}
}

func TestInfrastructureList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "infrastructure", "list")
		})
	}
}

func TestStorageList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "storage", "list")
		})
	}
}

func TestResourcePoolList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "resource-pool", "list")
		})
	}
}

func TestServerTypeList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-type", "list")
		})
	}
}

func TestNetworkDeviceList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "network-device", "list")
		})
	}
}

func TestEndpointList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "endpoint", "list")
		})
	}
}

func TestEventList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "event", "list", "--limit", "5")
		})
	}
}

func TestExtensionList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "extension", "list")
		})
	}
}

func TestLogicalNetworkList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "logical-network", "list")
		})
	}
}

func TestLogicalNetworkProfileList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "logical-network-profile", "list")
		})
	}
}

func TestOsTemplateList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "os-template", "list")
		})
	}
}

func TestTemplateAssetList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "template-asset", "list")
		})
	}
}

func TestFirmwareCatalogList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "firmware-catalog", "list")
		})
	}
}

func TestFirmwareBaselineList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "firmware-baseline", "list")
		})
	}
}

func TestFirmwareBinaryList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "firmware-binary", "list")
		})
	}
}

func TestFirmwarePolicyList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	// Probe once: firmware-policy endpoint is missing on some deployments.
	if _, stderr, err := runCLI(t, "-x", "firmware-policy", "list", "-f", "json"); err != nil && strings.Contains(stderr, "404") {
		t.Skip("firmware-policy endpoint not available on this deployment")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "firmware-policy", "list")
		})
	}
}

func TestServerList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server", "list")
		})
	}
}

func TestServerDefaultCredentialsList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-default-credentials", "list")
		})
	}
}

func TestServerCleanupPolicyList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-cleanup-policy", "list")
		})
	}
}

func TestServerRegistrationProfileList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-registration-profile", "list")
		})
	}
}

func TestCronJobList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "cron-job", "list")
		})
	}
}

func TestJobList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "job", "list", "--limit", "5")
		})
	}
}

func TestJobGroupList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "job-group", "list", "--limit", "5")
		})
	}
}

func TestVmPoolList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "vm-pool", "list")
		})
	}
}

func TestVmTypeList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "vm-type", "list")
		})
	}
}

func TestNetworkDeviceConfigurationTemplateList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "network-configuration", "device-template", "list")
		})
	}
}

func TestNetworkDeviceLinkAggregationTemplateList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "network-configuration", "link-aggregation-template", "list")
		})
	}
}

func TestNetworkDeviceDefaultSecretsList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "network-device", "default-secrets", "list")
		})
	}
}

func TestCustomIsoList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "custom-iso", "list")
		})
	}
}

func TestSiteDeviceAuthProviderList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "site", "device-auth-provider", "list")
		})
	}
}

func TestPermissionList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "permission", "list")
		})
	}
}

func TestUserList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "user", "list")
		})
	}
}

func TestAuthLdapMappingList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "auth", "ldap", "mapping-list")
		})
	}
}

func TestBucketList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "bucket", "list", infraID)
		})
	}
}

func TestDriveList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "drive", "list", infraID)
		})
	}
}

func TestFileShareList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "file-share", "list", infraID)
		})
	}
}

func TestServerInstanceGroupList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-instance-group", "list", infraID)
		})
	}
}

func TestServerInstanceList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "server-instance", "list", infraID)
		})
	}
}

func TestVmInstanceGroupList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "vm-instance-group", "list", infraID)
		})
	}
}

func TestExtensionInstanceList_AllFormats_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := firstInfraID(t)
	if infraID == "" {
		t.Skip("no infrastructures available")
	}
	for _, fmt := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(fmt, func(t *testing.T) {
			assertListFormat(t, fmt, "-x", "extension-instance", "list", infraID)
		})
	}
}
