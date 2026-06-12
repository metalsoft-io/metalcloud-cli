package cmd

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

// runCLIFormat is like runCLI but sets the output format via viper before
// executing. viper.Set takes precedence over persistent flag defaults, so this
// reliably exercises the formatter path for each format.
func runCLIFormat(t *testing.T, srv *httptest.Server, format string, args ...string) (string, error) {
	t.Helper()
	viper.Set(system.ConfigEndpoint, srv.URL)
	viper.Set(system.ConfigApiKey, "test-key")
	viper.Set(formatter.ConfigFormat, format)
	system.AllowDevelop = true

	var execErr error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(args)
		execErr = rootCmd.Execute()
	})
	return out, execErr
}

func TestSubnetList_FormatCSV(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "csv", "subnet", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// CSV output contains commas separating fields and does not start with JSON brackets.
	if strings.HasPrefix(strings.TrimSpace(out), "{") || strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected non-JSON (CSV) output, got: %s", out)
	}
	if !strings.Contains(out, ",") {
		t.Errorf("expected CSV output to contain commas, got: %s", out)
	}
}

func TestSubnetList_FormatYAML(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "yaml", "subnet", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// YAML output uses "key: value" notation; should not be valid JSON.
	if json.Valid([]byte(strings.TrimSpace(out))) {
		t.Errorf("expected YAML output, but got valid JSON: %s", out)
	}
	if !strings.Contains(out, ":") {
		t.Errorf("expected YAML output to contain ':', got: %s", out)
	}
}

func TestSubnetList_FormatText(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "text", "subnet", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Table renderer uses box-drawing characters or + borders.
	hasTableBorder := strings.Contains(out, "│") || strings.Contains(out, "+") || strings.Contains(out, "|")
	if !hasTableBorder {
		t.Errorf("expected text table output with borders, got: %s", out)
	}
}

func TestSubnetList_FormatMarkdown(t *testing.T) {
	srv := newSubnetTestServer()
	defer srv.Close()

	out, err := runCLIFormat(t, srv, "md", "subnet", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Markdown tables use | as column separators.
	if !strings.Contains(out, "|") {
		t.Errorf("expected markdown table output containing '|', got: %s", out)
	}
}

