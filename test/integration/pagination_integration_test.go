//go:build integration

package integration

import (
	"regexp"
	"testing"
)

var paginationSummaryRe = regexp.MustCompile(`Returned \d+ out of \d+ records`)

// TestPaginationSummary_Integration verifies that list commands emit a
// "Returned X out of Y records" summary line on stderr.
func TestPaginationSummary_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	_, stderr, err := runCLI(t, "-x", "subnet", "list", "-f", "json")
	if err != nil {
		t.Fatalf("subnet list failed: %v\nstderr: %s", err, stderr)
	}
	if !paginationSummaryRe.MatchString(stderr) {
		t.Errorf("expected pagination summary matching %q, got: %s", paginationSummaryRe, stderr)
	}
}

// TestPaginationSummary_Server_Integration: server list uses raw-body parsing
// (SDK Links bug workaround) with no pagination loop — it prints no summary.
// Verify it at least exits 0 with non-empty valid output.
func TestPaginationSummary_Server_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	stdout, stderr, err := runCLI(t, "-x", "server", "list", "-f", "json")
	if err != nil {
		t.Fatalf("server list failed: %v\nstderr: %s", err, stderr)
	}
	if stdout == "" {
		t.Error("server list: empty output")
	}
}

// TestPaginationSummary_Infrastructure_Integration checks infrastructure list
// pagination.
func TestPaginationSummary_Infrastructure_Integration(t *testing.T) {
	skipIfNoEndpoint(t)
	_, stderr, err := runCLI(t, "-x", "infrastructure", "list", "-f", "json")
	if err != nil {
		t.Fatalf("infrastructure list failed: %v\nstderr: %s", err, stderr)
	}
	if !paginationSummaryRe.MatchString(stderr) {
		t.Errorf("expected pagination summary matching %q, got: %s", paginationSummaryRe, stderr)
	}
}
