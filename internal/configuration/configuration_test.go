package configuration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// configurationDocument is a trimmed but representative platform configuration:
// nested objects under per-service keys.
var configurationDocument = map[string]any{
	"platform": map[string]any{
		"brandName":  "MetalSoft",
		"apiBaseUrl": "https://qa01.metalcloud.io",
	},
	"notification": map[string]any{
		"smtp": map[string]any{"host": "smtp.example.com", "port": 25},
	},
}

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method   string
	Path     string
	RawQuery string
	Body     string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{
			Method:   req.Method,
			Path:     req.URL.Path,
			RawQuery: req.URL.RawQuery,
			Body:     string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *recorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newServer(rec *recorder) *httptest.Server {
	return testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/config":          rec.wrap(testutils.JSONHandler(http.StatusOK, configurationDocument)),
		"/api/v2/config/platform": rec.wrap(testutils.JSONHandler(http.StatusOK, configurationDocument["platform"])),
		"/api/v2/config/tunnel":   rec.wrap(testutils.NoContentHandler()),
	})
}

func run(t *testing.T, fn func() error) string {
	t.Helper()
	var err error
	out := testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestConfigurationGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return ConfigurationGet(ctx, "") })
	if !strings.Contains(out, "smtp.example.com") {
		t.Errorf("get output missing a nested value: %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/config" || got.RawQuery != "" {
		t.Errorf("unexpected request %s?%s", got.Path, got.RawQuery)
	}
}

func TestConfigurationGetSendsFilterQuery(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error { return ConfigurationGet(ctx, "platform") })

	if got := rec.last(); got.RawQuery != "filter=platform" {
		t.Errorf("a service argument should become the filter query parameter, got %q", got.RawQuery)
	}
}

// TestConfigurationGetTextFormatIsYaml pins the decision to render the nested
// configuration document as YAML for table formats.
func TestConfigurationGetTextFormatIsYaml(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	viper.Set(formatter.ConfigFormat, "text")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := run(t, func() error { return ConfigurationGet(ctx, "") })
	if !strings.Contains(out, "host: smtp.example.com") {
		t.Errorf("text output should be YAML, got: %s", out)
	}
}

func TestConfigurationReplaceUsesPut(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return ConfigurationReplace(ctx, "platform", []byte(`{"brandName":"Acme"}`))
	})

	got := rec.last()
	if got.Method != http.MethodPut || got.Path != "/api/v2/config/platform" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"brandName":"Acme"`) {
		t.Errorf("the configuration document should be sent verbatim, got %s", got.Body)
	}
}

func TestConfigurationUpdateUsesPatch(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return ConfigurationUpdate(ctx, "platform", []byte(`{"brandName":"Acme"}`))
	})

	if got := rec.last(); got.Method != http.MethodPatch {
		t.Errorf("update should use PATCH, got %s", got.Method)
	}
}

// TestConfigurationUpdateAcceptsYaml verifies YAML input is converted to JSON.
// utils.UnmarshalContent only falls back to YAML when the output format is not
// pinned to json, so the format is switched for this test.
func TestConfigurationUpdateAcceptsYaml(t *testing.T) {
	rec := &recorder{}
	srv := newServer(rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	viper.Set(formatter.ConfigFormat, "text")
	defer viper.Set(formatter.ConfigFormat, "json")

	run(t, func() error {
		return ConfigurationUpdate(ctx, "platform", []byte("brandName: Acme\n"))
	})

	if got := rec.last(); !strings.Contains(got.Body, `"brandName":"Acme"`) {
		t.Errorf("YAML input should be sent as JSON, got %s", got.Body)
	}
}

// TestConfigurationWriteNoContent covers a 204 answer, which carries no body.
func TestConfigurationWriteNoContent(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ConfigurationReplace(ctx, "tunnel", []byte(`{"shared_secret":"x"}`)); err != nil {
		t.Fatalf("a 204 response should not be an error, got: %v", err)
	}
}

func TestConfigurationUnknownService(t *testing.T) {
	srv := newServer(&recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	err := ConfigurationUpdate(ctx, "not-a-service", []byte(`{}`))
	if err == nil || !strings.Contains(err.Error(), "unknown configuration service") {
		t.Fatalf("expected an unknown service error, got %v", err)
	}
}

func TestValidateServiceNormalizes(t *testing.T) {
	service, err := ValidateService("  Platform ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if service != "platform" {
		t.Errorf("expected 'platform', got %q", service)
	}

	if _, err := ValidateService(""); err == nil {
		t.Error("an empty service should be rejected")
	}
}
