//go:build integration || integration_crud

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func skipIfNoEndpoint(t *testing.T) {
	t.Helper()
	if os.Getenv("METALCLOUD_ENDPOINT") == "" || os.Getenv("METALCLOUD_API_KEY") == "" {
		t.Skip("METALCLOUD_ENDPOINT and METALCLOUD_API_KEY required for integration tests")
	}
}

func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	binary := os.Getenv("METALCLOUD_CLI_BINARY")
	if binary == "" {
		binary = "../../metalcloud-cli"
	}
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(),
		"METALCLOUD_ENDPOINT="+os.Getenv("METALCLOUD_ENDPOINT"),
		"METALCLOUD_API_KEY="+os.Getenv("METALCLOUD_API_KEY"),
		"METALCLOUD_FORMAT=json",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// extractIDByField scans a JSON array for an item where fieldName == fieldValue
// and returns its "id" (or first *Id field).
func extractIDByField(t *testing.T, jsonStr, fieldName, fieldValue string) string {
	t.Helper()
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &items); err != nil {
		t.Fatalf("extractIDByField: unmarshal failed: %v", err)
	}
	for _, item := range items {
		if v, ok := item[fieldName]; ok && fmt.Sprintf("%v", v) == fieldValue {
			for _, key := range []string{"id", "resourcePoolId", "policyId"} {
				if id, ok := item[key]; ok {
					switch vt := id.(type) {
					case float64:
						return fmt.Sprintf("%.0f", vt)
					case string:
						return vt
					}
				}
			}
		}
	}
	t.Fatalf("extractIDByField: no item with %s=%q found", fieldName, fieldValue)
	return ""
}

// firstInfraID runs "infrastructure list" and returns the ID of the first
// result as a string, or "" if there are no infrastructures.
func firstInfraID(t *testing.T) string {
	t.Helper()
	stdout, _, err := runCLI(t, "-x", "infrastructure", "list", "-f", "json")
	if err != nil || !json.Valid([]byte(stdout)) {
		return ""
	}
	var items []map[string]interface{}
	if json.Unmarshal([]byte(stdout), &items) != nil || len(items) == 0 {
		return ""
	}
	switch v := items[0]["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	}
	return ""
}

// firstVMPoolID runs "vm-pool list" and returns the ID of the first result as
// a string, or "" if there are no VM pools.
func firstVMPoolID(t *testing.T) string {
	t.Helper()
	stdout, _, err := runCLI(t, "-x", "vm-pool", "list", "-f", "json")
	if err != nil || !json.Valid([]byte(stdout)) {
		return ""
	}
	var items []map[string]interface{}
	if json.Unmarshal([]byte(stdout), &items) != nil || len(items) == 0 {
		return ""
	}
	switch v := items[0]["id"].(type) {
	case float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	}
	return ""
}
