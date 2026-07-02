package network_device_configuration_template

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func decodeItems(jsonStrings []string) []any {
	out := make([]any, len(jsonStrings))
	for i, s := range jsonStrings {
		var v any
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			panic(fmt.Sprintf("decodeItems: invalid JSON at index %d: %v", i, err))
		}
		out[i] = v
	}
	return out
}

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

const bgpTmplItem = `{
	"id":1,
	"action":"add-neighbor",
	"networkType":"underlay",
	"networkDeviceDriver":"sonic_enterprise",
	"networkDevicePosition":"leaf",
	"remoteNetworkDevicePosition":"spine",
	"bgpNumbering":"numbered",
	"bgpLinkConfiguration":"active",
	"executionType":"cli",
	"libraryLabel":"my-template",
	"configuration":"Y21k",
	"createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

func bgpTmplListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// --- List ---

func TestNetworkDeviceConfigurationTemplateList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, []string{bgpTmplItem, bgpTmplItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"action":"add-neighbor","networkType":"underlay",
				"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
				"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
				"bgpLinkConfiguration":"active","executionType":"cli",
				"libraryLabel":"label-%d","configuration":"Y21k",
				"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-device-bgp-configuration-templates", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- Get ---

func TestNetworkDeviceConfigurationTemplateGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": testutils.RawHandler(http.StatusOK, bgpTmplItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Create ---

func TestNetworkDeviceConfigurationTemplateCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{
		"action":"add-neighbor","networkType":"underlay",
		"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
		"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
		"bgpLinkConfiguration":"active","executionType":"cli",
		"libraryLabel":"my-template","configuration":"Y21k"
	}`)
	if err := NetworkDeviceConfigurationTemplateCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateCreate(ctx, []byte(`{}`)); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- Import library ---

// validTmplDescriptor is a single template descriptor as it would live in a
// library directory. It deliberately carries its own libraryLabel to prove the
// import overrides it with the one passed on the command line.
const validTmplDescriptor = `{
	"action":"add-neighbor","networkType":"underlay",
	"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
	"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
	"bgpLinkConfiguration":"active","executionType":"cli",
	"libraryLabel":"from-file","configuration":"Y21k"
}`

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestImportLibrary_HappyPath(t *testing.T) {
	var posts int64
	var labels []string
	var mu sync.Mutex
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			atomic.AddInt64(&posts, 1)
			body, _ := io.ReadAll(r.Body)
			var payload struct {
				LibraryLabel string `json:"libraryLabel"`
			}
			_ = json.Unmarshal(body, &payload)
			mu.Lock()
			labels = append(labels, payload.LibraryLabel)
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	dir := t.TempDir()
	writeFile(t, dir, "01-underlay.json", validTmplDescriptor)
	writeFile(t, dir, "02-overlay.json", validTmplDescriptor)
	writeFile(t, dir, "notes.txt", "ignored - wrong extension")

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateImportLibrary(ctx, "spectrumx", dir, false); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got := atomic.LoadInt64(&posts); got != 2 {
		t.Fatalf("expected 2 POSTs (txt ignored), got %d", got)
	}
	for _, l := range labels {
		if l != "spectrumx" {
			t.Fatalf("expected libraryLabel overridden to 'spectrumx', got %q", l)
		}
	}
}

func TestImportLibrary_DryRunMakesNoCalls(t *testing.T) {
	var posts int64
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt64(&posts, 1)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	dir := t.TempDir()
	writeFile(t, dir, "t.yaml", validTmplDescriptor)

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateImportLibrary(ctx, "spectrumx", dir, true); err != nil {
		t.Fatalf("expected nil error on dry-run, got: %v", err)
	}
	if got := atomic.LoadInt64(&posts); got != 0 {
		t.Fatalf("expected 0 POSTs on dry-run, got %d", got)
	}
}

func TestImportLibrary_NoMatchingFiles(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	dir := t.TempDir()
	writeFile(t, dir, "readme.md", "not a template")

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateImportLibrary(ctx, "spectrumx", dir, false); err == nil {
		t.Error("expected error when no template files present, got nil")
	}
}

func TestImportLibrary_MissingDir(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateImportLibrary(ctx, "spectrumx", filepath.Join(t.TempDir(), "nope"), false); err == nil {
		t.Error("expected error for missing directory, got nil")
	}
}

func TestImportLibrary_BadFileCountedAsFailure(t *testing.T) {
	var posts int64
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt64(&posts, 1)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	dir := t.TempDir()
	writeFile(t, dir, "good.json", validTmplDescriptor)
	writeFile(t, dir, "bad.json", "{not valid json")

	ctx := testutils.SetupTestContext(ts.URL)
	err := NetworkDeviceConfigurationTemplateImportLibrary(ctx, "spectrumx", dir, false)
	if err == nil {
		t.Fatal("expected error because one file failed, got nil")
	}
	// The good file must still have been imported despite the bad one.
	if got := atomic.LoadInt64(&posts); got != 1 {
		t.Fatalf("expected the 1 good file to import, got %d POSTs", got)
	}
}

// --- Export library / list libraries ---

// tmplInLib renders a list item carrying the given id and libraryLabel.
func tmplInLib(id int, lib string) string {
	return fmt.Sprintf(`{
		"id":%d,"action":"add-neighbor","networkType":"underlay",
		"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
		"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
		"bgpLinkConfiguration":"active","executionType":"cli",
		"libraryLabel":"%s","preparation":"cHJlcA==","configuration":"Y21k",
		"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"
	}`, id, lib)
}

func TestExportLibrary_HappyPathAndRoundTrip(t *testing.T) {
	items := []string{tmplInLib(1, "alpha"), tmplInLib(2, "alpha"), tmplInLib(3, "beta")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, items, 1, 1),
	})
	defer ts.Close()

	dir := filepath.Join(t.TempDir(), "out") // must be created by the export
	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateExportLibrary(ctx, "alpha", dir); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read out dir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 exported files (only 'alpha'), got %d", len(entries))
	}

	// Every exported descriptor must carry the exported library and re-parse as
	// a create body (no id/timestamps leak in that would break re-import).
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("%s not valid JSON: %v", e.Name(), err)
		}
		if m["libraryLabel"] != "alpha" {
			t.Errorf("%s: expected libraryLabel 'alpha', got %v", e.Name(), m["libraryLabel"])
		}
		if _, leaked := m["id"]; leaked {
			t.Errorf("%s: exported descriptor leaked 'id'", e.Name())
		}
		if _, leaked := m["createdTimestamp"]; leaked {
			t.Errorf("%s: exported descriptor leaked 'createdTimestamp'", e.Name())
		}
	}
}

func TestExportLibrary_NoTemplatesInLibrary(t *testing.T) {
	items := []string{tmplInLib(1, "alpha")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, items, 1, 1),
	})
	defer ts.Close()

	dir := t.TempDir()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateExportLibrary(ctx, "does-not-exist", dir); err == nil {
		t.Error("expected error when library has no templates, got nil")
	}
}

func TestListLibraries_HappyPath(t *testing.T) {
	items := []string{tmplInLib(1, "alpha"), tmplInLib(2, "alpha"), tmplInLib(3, "beta")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, items, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateListLibraries(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestListLibraries_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": testutils.ErrorHandler(http.StatusInternalServerError, "boom"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateListLibraries(ctx); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

// --- Update ---

func TestNetworkDeviceConfigurationTemplateUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateUpdate(ctx, "1", []byte(`{"libraryLabel":"updated"}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateUpdate(ctx, "99", []byte(`{}`)); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- Delete ---

func TestNetworkDeviceConfigurationTemplateDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}
