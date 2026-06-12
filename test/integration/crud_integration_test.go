//go:build integration_crud

package integration

// CRUD integration tests that CREATE and DELETE real resources.
//
// Requirements:
//   - METALCLOUD_ENDPOINT  — API endpoint URL
//   - METALCLOUD_API_KEY   — API key
//   - METALCLOUD_TEST_SITE_ID — site ID to create resources under (where required)
//
// Run: go test -tags integration_crud ./test/integration/
//
// WARNING: these tests create and immediately delete resources.
// Use a dedicated test environment — never run against production.

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

// uniqueName returns a name with a timestamp suffix to avoid 409 conflicts on retries.
func uniqueName(base string) string {
	return fmt.Sprintf("%s-%d", base, time.Now().UnixMilli()%100000)
}

func envOrSkip(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("%s required for this test", key)
	}
	return v
}

func testSiteID(t *testing.T) string {
	t.Helper()
	id := os.Getenv("METALCLOUD_TEST_SITE_ID")
	if id == "" {
		t.Skip("METALCLOUD_TEST_SITE_ID required for CRUD integration tests")
	}
	return id
}

// extractID parses a JSON object and returns an ID field as a string.
// Tries "id" first, then "resourcePoolId", then any *Id field.
func extractID(t *testing.T, jsonStr string) string {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		t.Fatalf("extractID: unmarshal failed: %v\nraw: %s", err, jsonStr)
	}
	for _, key := range []string{"id", "resourcePoolId", "policyId", "profileId"} {
		if v, ok := obj[key]; ok {
			switch vt := v.(type) {
			case float64:
				return fmt.Sprintf("%.0f", vt)
			case string:
				return vt
			}
		}
	}
	t.Fatalf("extractID: no id field found in: %s", jsonStr)
	return ""
}

// ---------------------------------------------------------------------------
// Subnet CRUD
// ---------------------------------------------------------------------------

func TestSubnetCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	testSiteID(t) // skip if site ID not set

	payload := `{"name":"test-ci-subnet","networkAddress":"192.168.254.0","prefixLength":24,"isPool":false}`
	stdout, stderr, err := runCLI(t, "-x", "subnet", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("subnet create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "subnet", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "subnet", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("subnet get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Secret CRUD
// ---------------------------------------------------------------------------

func TestSecretCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	payload := `{"name":"test-ci-secret","value":"integration-test-value"}`
	stdout, stderr, err := runCLI(t, "-x", "secret", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("secret create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "secret", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "secret", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("secret get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Variable CRUD
// ---------------------------------------------------------------------------

func TestVariableCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	payload := `{"name":"test-ci-variable","value":{"env":"ci","purpose":"integration-test"}}`
	stdout, stderr, err := runCLI(t, "-x", "variable", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("variable create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "variable", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "variable", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("variable get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Role CRUD
// ---------------------------------------------------------------------------

func TestRoleCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	payload := `{"label":"test-ci-role","permissions":[]}`
	stdout, stderr, err := runCLI(t, "-x", "role", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("role create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "role", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "role", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("role get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// DNS Zone CRUD
// ---------------------------------------------------------------------------

func TestDnsZoneCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	payload := `{"zoneName":"test-ci.internal","isDefault":false,"nameServers":["ns1.test.com"]}`
	stdout, stderr, err := runCLI(t, "-x", "dns-zone", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("dns-zone create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "dns-zone", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "dns-zone", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("dns-zone get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Firmware Baseline CRUD
// ---------------------------------------------------------------------------

func TestFirmwareBaselineCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	payload := `{"name":"test-ci-baseline"}`
	stdout, stderr, err := runCLI(t, "-x", "firmware-baseline", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("firmware-baseline create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "firmware-baseline", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "firmware-baseline", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("firmware-baseline get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Resource Pool CRUD
// ---------------------------------------------------------------------------

func TestResourcePoolCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	label := uniqueName("test-ci-pool")
	stdout, stderr, err := runCLI(t, "-x", "resource-pool", "create",
		"--label", label,
		"--description", "CI test pool",
		"-f", "json",
	)
	if err != nil {
		t.Fatalf("resource-pool create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "resource-pool", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "resource-pool", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("resource-pool get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Server Cleanup Policy CRUD
// ---------------------------------------------------------------------------

func TestServerCleanupPolicyCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	stdout, stderr, err := runCLI(t, "-x", "server-cleanup-policy", "create",
		"--label", uniqueName("test-ci-policy"),
		"--raid-one-drive", "RAID0",
		"--raid-two-drives", "RAID1",
		"--raid-even-drives", "RAID10",
		"--raid-odd-drives", "RAID5",
		"--skip-raid-actions", "provisioning",
		"-f", "json",
	)
	if err != nil {
		t.Fatalf("server-cleanup-policy create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "server-cleanup-policy", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "server-cleanup-policy", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("server-cleanup-policy get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Cron Job CRUD
// ---------------------------------------------------------------------------

func TestCronJobCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	// Valid functionName values are deployment-specific. Pass via METALCLOUD_TEST_CRON_FUNCTION.
	fn := os.Getenv("METALCLOUD_TEST_CRON_FUNCTION")
	if fn == "" {
		t.Skip("METALCLOUD_TEST_CRON_FUNCTION required — set to a valid cron function name for this deployment (e.g. 'syncUsers')")
	}
	payload := `{"label":"` + uniqueName("test-ci-cron") + `","functionName":"` + fn + `","params":[],"schedule":"0 * * * *","waitForCompletion":0,"lifetimeSeconds":60,"disabled":1}`
	stdout, stderr, err := runCLI(t, "-x", "cron-job", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if err != nil {
		t.Fatalf("cron-job create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "cron-job", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "cron-job", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("cron-job get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// VM Type CRUD
// ---------------------------------------------------------------------------

func TestVmTypeCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)

	vmTypeName := uniqueName("test-ci-vmtype")
	payload := `{"name":"` + vmTypeName + `","cpuCores":2,"ramGB":4}`
	_, stderr, err := runCLI(t, "-x", "vm-type", "create", "--config-source", writeTempJSON(t, payload))
	if err != nil {
		t.Fatalf("vm-type create failed: %v\nstderr: %s", err, stderr)
	}

	// vm-type create logs but doesn't print JSON — find the created type via list
	listOut, _, err := runCLI(t, "-x", "vm-type", "list", "-f", "json")
	if err != nil {
		t.Fatalf("vm-type list failed: %v", err)
	}
	id := extractIDByField(t, listOut, "name", vmTypeName)

	t.Cleanup(func() {
		runCLI(t, "-x", "vm-type", "delete", id) //nolint:errcheck
	})

	stdout, _, err := runCLI(t, "-x", "vm-type", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("vm-type get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Infrastructure CRUD
// ---------------------------------------------------------------------------

func TestInfrastructureCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	siteID := envOrSkip(t, "METALCLOUD_TEST_SITE_ID")

	stdout, stderr, err := runCLI(t, "-x", "infrastructure", "create", siteID, "test-ci-infra", "-f", "json")
	if err != nil {
		t.Fatalf("infrastructure create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "infrastructure", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "infrastructure", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("infrastructure get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Logical Network Profile CRUD
// ---------------------------------------------------------------------------

func TestLogicalNetworkProfileCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	fabricID := envOrSkip(t, "METALCLOUD_TEST_FABRIC_ID")

	payload := fmt.Sprintf(`{"fabricId":"%s"}`, fabricID)
	stdout, stderr, err := runCLI(t, "-x", "logical-network-profile", "create", "vlan",
		"--config-source", writeTempJSON(t, payload),
		"-f", "json",
	)
	if err != nil {
		t.Fatalf("logical-network-profile create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "logical-network-profile", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "logical-network-profile", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("logical-network-profile get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Endpoint CRUD
// ---------------------------------------------------------------------------

func TestEndpointCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	siteIDStr := envOrSkip(t, "METALCLOUD_TEST_SITE_ID")

	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		t.Fatalf("METALCLOUD_TEST_SITE_ID must be numeric, got %q: %v", siteIDStr, err)
	}

	payload := fmt.Sprintf(`{"siteId":%d,"name":"test-ci-endpoint","label":"test-ci-endpoint"}`, siteID)
	stdout, stderr, errRun := runCLI(t, "-x", "endpoint", "create", "--config-source", writeTempJSON(t, payload), "-f", "json")
	if errRun != nil {
		t.Fatalf("endpoint create failed: %v\nstdout: %s\nstderr: %s", errRun, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "endpoint", "delete", id) //nolint:errcheck
	})

	stdout, _, errRun = runCLI(t, "-x", "endpoint", "get", id, "-f", "json")
	if errRun != nil {
		t.Errorf("endpoint get after create failed: %v\nstdout: %s", errRun, stdout)
	}
}

// ---------------------------------------------------------------------------
// Fabric CRUD
// ---------------------------------------------------------------------------

func TestFabricCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	siteID := envOrSkip(t, "METALCLOUD_TEST_SITE_ID")

	// fabric create takes positional args: site_id fabric_name fabric_type
	// The fabric-type-specific configuration body is supplied via --config-source.
	fabricConfig := `{"fabricType":"ethernet"}`
	stdout, stderr, err := runCLI(t, "-x", "fabric", "create",
		siteID, uniqueName("test-ci-fabric"), "ethernet",
		"--config-source", writeTempJSON(t, fabricConfig),
		"-f", "json",
	)
	if err != nil {
		t.Fatalf("fabric create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "fabric", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "fabric", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("fabric get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Server Instance Group CRUD
// ---------------------------------------------------------------------------

func TestServerInstanceGroupCRUD_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	infraID := envOrSkip(t, "METALCLOUD_TEST_INFRA_ID")
	serverTypeID := envOrSkip(t, "METALCLOUD_TEST_SERVER_TYPE_ID")

	// create takes positional args: infrastructure_id_or_label label server_type_id instance_count
	stdout, stderr, err := runCLI(t, "-x", "server-instance-group", "create",
		infraID, "test-ci-sig", serverTypeID, "1",
		"-f", "json",
	)
	if err != nil {
		t.Fatalf("server-instance-group create failed: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	id := extractID(t, stdout)
	t.Cleanup(func() {
		runCLI(t, "-x", "server-instance-group", "delete", id) //nolint:errcheck
	})

	stdout, _, err = runCLI(t, "-x", "server-instance-group", "get", id, "-f", "json")
	if err != nil {
		t.Errorf("server-instance-group get after create failed: %v\nstdout: %s", err, stdout)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// writeTempJSON writes payload to a temporary file and returns its path.
// The file is cleaned up automatically when the test ends.
func writeTempJSON(t *testing.T, payload string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "crud-*.json")
	if err != nil {
		t.Fatalf("writeTempJSON: %v", err)
	}
	if _, err := f.WriteString(payload); err != nil {
		t.Fatalf("writeTempJSON write: %v", err)
	}
	f.Close()
	return f.Name()
}
