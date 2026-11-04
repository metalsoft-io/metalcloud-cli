package firmware_catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func nvidiaFixtureCatalog(t *testing.T) *VendorCatalog {
	t.Helper()
	vc, err := NewVendorCatalogFromCreateOptions(FirmwareCatalogCreateOptions{
		Name:                    "NVIDIA HGX H100 25.06.1",
		Vendor:                  VendorNvidia,
		UpdateType:              UpdateTypeOffline,
		VendorLocalCatalogPath:  filepath.Join(nvidiaFixtureDir, "manifest.json"),
		VendorLocalBinariesPath: nvidiaFixtureDir,
		RepoBaseUrl:             "https://repo.example.com/nvidia/hgx-h100/25.06.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	return vc
}

func TestProcessNvidiaCatalogBuildsTheBundleAndItsPackages(t *testing.T) {
	vc := nvidiaFixtureCatalog(t)

	if err := vc.ProcessVendorCatalog(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := *vc.CatalogInfo.VendorId; got != "HGX H100 8-GPU@25.06.1" {
		t.Errorf("vendorId = %q", got)
	}
	if got := strings.Join(vc.CatalogInfo.VendorServerTypesSupported, ","); got != "NVIDIA DGXH100,Supermicro SYS-821GE-TNHR" {
		t.Errorf("vendorServerTypesSupported = %q", got)
	}
	if got := vc.CatalogInfo.VendorReleaseTimestamp.Format("2006-01-02"); got != "2025-06-12" {
		t.Errorf("vendorReleaseTimestamp = %q", got)
	}
	if len(vc.Binaries) != 2 {
		t.Fatalf("expected the bundle and the retimer package, got %d binaries", len(vc.Binaries))
	}

	bundle := vc.Binaries[0]
	if *bundle.ExternalId != nvidiaFixtureBundleFile || !*bundle.IsBundle || *bundle.Sha256 != nvidiaFixtureBundleSha {
		t.Errorf("bundle binary: externalId %q, isBundle %v, sha256 %q", *bundle.ExternalId, *bundle.IsBundle, *bundle.Sha256)
	}
	if bundle.VendorDownloadUrl != "https://repo.example.com/nvidia/hgx-h100/25.06.1/"+nvidiaFixtureBundleFile {
		t.Errorf("bundle download URL = %q", bundle.VendorDownloadUrl)
	}
	if bundle.Name != "NVIDIA HGX H100 8-GPU 25.06.1" || *bundle.PackageVersion != "25.06.1" || bundle.UpdateSeverity != "recommended" {
		t.Errorf("bundle name %q, version %q, severity %q", bundle.Name, *bundle.PackageVersion, bundle.UpdateSeverity)
	}
	if !bundle.RebootRequired {
		t.Error("a bundle with a chassis power cycle requires a reboot")
	}
	if len(bundle.BundleSteps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(bundle.BundleSteps))
	}
	gpu := bundle.BundleSteps[2]
	if gpu.Id != "gpu" || gpu.ComponentClass != sdk.SERVERCOMPONENTCLASS_GPU || gpu.Activation != sdk.FIRMWAREBUNDLESTEPACTIVATION_CHASSIS_POWER_CYCLE ||
		gpu.VerifyVersion == nil || *gpu.VerifyVersion != "96.00.9F.00.01" || gpu.DependsOnStep.Get() == nil || *gpu.DependsOnStep.Get() != "hmc" {
		t.Errorf("gpu step = %+v", gpu)
	}
	// The first step depends on nothing: sent as null
	steps, _ := json.Marshal(bundle.BundleSteps[0])
	if !strings.Contains(string(steps), `"dependsOnStep":null`) {
		t.Errorf("erot step = %s", steps)
	}
	systems, _ := json.Marshal(bundle.VendorSupportedSystems)
	if string(systems) != `[{"id":"NVIDIA DGXH100"},{"id":"Supermicro SYS-821GE-TNHR"}]` {
		t.Errorf("vendorSupportedSystems = %s", systems)
	}

	retimer := vc.Binaries[1]
	if *retimer.ExternalId != nvidiaFixtureRetimer || *retimer.IsBundle || *retimer.Sha256 != nvidiaFixtureRetimerSha || len(retimer.BundleSteps) != 0 {
		t.Errorf("retimer binary: externalId %q, isBundle %v, sha256 %q", *retimer.ExternalId, *retimer.IsBundle, *retimer.Sha256)
	}
	if retimer.Name != "NVIDIA HGX H100 8-GPU 25.06.1 (retimer)" || !retimer.RebootRequired {
		t.Errorf("retimer name %q, reboot %v", retimer.Name, retimer.RebootRequired)
	}
}

func TestProcessNvidiaCatalogRejections(t *testing.T) {
	tests := []struct {
		name     string
		change   func(vc *VendorCatalog)
		expected string
	}{
		{"download binaries", func(vc *VendorCatalog) { vc.DownloadBinaries = true }, "--download-binaries"},
		{"vendor URL", func(vc *VendorCatalog) { vc.CatalogInfo.VendorUrl = sdk.PtrString("https://example.com") }, "--vendor-url"},
		{"no repository", func(vc *VendorCatalog) { vc.RepoBaseUrl = "" }, "has no downloadUrls entry"},
		{"no manifest", func(vc *VendorCatalog) { vc.VendorLocalCatalogPath = "" }, "--vendor-local-catalog-path"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vc := nvidiaFixtureCatalog(t)
			test.change(vc)
			err := vc.ProcessVendorCatalog(context.Background())
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Fatalf("expected an error mentioning %q, got %v", test.expected, err)
			}
		})
	}
}

func TestNvidiaManifestPrerequisites(t *testing.T) {
	manifest := loadNvidiaFixture(t)
	manifest.Prerequisites = []NvidiaPrerequisite{
		{Platform: "NVIDIA DGX H100", ComponentClass: "bios", MinimumVersion: "01.02.05"},
		{ComponentClass: "bmc", MinimumVersion: "23.06.00"},
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("a platform written with spaces is the same platform: %v", err)
	}

	repoUrl, _ := url.Parse("https://repo.example.com/nvidia")
	binaries := nvidiaBinaries(manifest, repoUrl, mustParseDate(t, manifest.ReleaseDate))
	prerequisites, _ := json.Marshal(binaries[0].Prerequisites)
	if string(prerequisites) != `[{"componentClass":"bios","minimumVersion":"01.02.05","platform":"NVIDIA DGX H100"},{"componentClass":"bmc","minimumVersion":"23.06.00"}]` {
		t.Errorf("bundle prerequisites = %s", prerequisites)
	}

	manifest.Prerequisites = []NvidiaPrerequisite{
		{Platform: "NVIDIA DGXB200", ComponentClass: "bios", MinimumVersion: "1"},
		{ComponentClass: "psu", MinimumVersion: ""},
	}
	issues := manifestIssues(t, manifest.Validate())
	for _, expected := range []string{
		`prerequisites[0]: platform "NVIDIA DGXB200" is not one of targetPlatforms`,
		`prerequisites[1]: unknown component class "psu"`,
		`prerequisites[1]: minimumVersion is required`,
	} {
		if !containsIssue(issues, expected) {
			t.Errorf("expected an issue starting with %q, got %v", expected, issues)
		}
	}
}

// fakeFirmwareApi keeps catalogs and binaries in memory, for the create-or-update import.
type fakeFirmwareApi struct {
	mu       sync.Mutex
	catalogs []map[string]any
	binaries []map[string]any
	calls    []string
}

func (f *fakeFirmwareApi) server() *httptest.Server {
	write := func(w http.ResponseWriter, status int, body any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}
	page := func(data []map[string]any) map[string]any {
		return map[string]any{"data": data, "meta": testutils.PaginatedMeta(1, 1, 100), "links": map[string]any{}}
	}
	filterValue := func(r *http.Request, key string) string {
		return strings.TrimPrefix(r.URL.Query().Get("filter."+key), "$eq:")
	}
	saved := func(body map[string]any, id int) map[string]any {
		body["id"] = id
		body["links"] = []any{}
		body["createdTimestamp"] = "2026-11-04T10:00:00Z"
		return body
	}

	mux := http.NewServeMux()
	handle := func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.calls = append(f.calls, r.Method+" "+r.URL.Path)
		var body map[string]any
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		collection, id, _ := strings.Cut(strings.TrimPrefix(r.URL.Path, "/api/v2/firmware/"), "/")
		items := &f.catalogs
		if collection == "binary" {
			items = &f.binaries
		}
		switch {
		case r.Method == http.MethodGet:
			result := []map[string]any{}
			for _, item := range *items {
				if collection == "catalog" && (item["vendor"] != filterValue(r, "vendor") || item["vendorId"] != filterValue(r, "vendorId")) {
					continue
				}
				if collection == "binary" && (strconv.Itoa(item["catalogId"].(int)) != filterValue(r, "catalogId") || item["externalId"] != filterValue(r, "externalId")) {
					continue
				}
				result = append(result, item)
			}
			write(w, http.StatusOK, page(result))
		case r.Method == http.MethodPost:
			if collection == "binary" {
				body["catalogId"] = int(body["catalogId"].(float64))
			}
			*items = append(*items, saved(body, len(*items)+1))
			write(w, http.StatusCreated, (*items)[len(*items)-1])
		case r.Method == http.MethodPut:
			index, _ := strconv.Atoi(id)
			if collection == "binary" {
				body["catalogId"] = int(body["catalogId"].(float64))
			}
			(*items)[index-1] = saved(body, index)
			write(w, http.StatusOK, (*items)[index-1])
		}
	}
	mux.HandleFunc("/api/v2/firmware/", handle)
	return httptest.NewServer(mux)
}

func (f *fakeFirmwareApi) callsWith(prefix string) []string {
	result := []string{}
	for _, call := range f.calls {
		if strings.HasPrefix(call, prefix) {
			result = append(result, call)
		}
	}
	return result
}

func importFixture(t *testing.T, ctx context.Context, change func(vc *VendorCatalog)) error {
	t.Helper()
	vc := nvidiaFixtureCatalog(t)
	if change != nil {
		change(vc)
	}
	if err := vc.ProcessVendorCatalog(ctx); err != nil {
		t.Fatal(err)
	}
	return vc.CreateMetalsoftCatalog(ctx)
}

func TestNvidiaImportCreatesThenUpdates(t *testing.T) {
	api := &fakeFirmwareApi{}
	ts := api.server()
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := importFixture(t, ctx, nil); err != nil {
		t.Fatalf("first import: %v", err)
	}
	if got := strings.Join(api.callsWith("POST"), ", "); got != "POST /api/v2/firmware/catalog, POST /api/v2/firmware/binary, POST /api/v2/firmware/binary" {
		t.Fatalf("the first import creates the catalog and both binaries, got %s", got)
	}
	bundle := api.binaries[0]
	if bundle["isBundle"] != true || bundle["sha256"] != nvidiaFixtureBundleSha || len(bundle["bundleSteps"].([]any)) != 4 {
		t.Errorf("bundle sent = %v", bundle)
	}
	if bundle["cacheDownloadUrl"] != "https://repo.example.com/nvidia/hgx-h100/25.06.1/"+nvidiaFixtureBundleFile {
		t.Errorf("cacheDownloadUrl = %v", bundle["cacheDownloadUrl"])
	}

	api.calls = nil
	if err := importFixture(t, ctx, nil); err != nil {
		t.Fatalf("second import: %v", err)
	}
	if posts := api.callsWith("POST"); len(posts) != 0 {
		t.Errorf("the second import only updates, got %v", posts)
	}
	if got := strings.Join(api.callsWith("PUT"), ", "); got != "PUT /api/v2/firmware/catalog/1, PUT /api/v2/firmware/binary/1, PUT /api/v2/firmware/binary/2" {
		t.Errorf("updates = %s", got)
	}
	if len(api.catalogs) != 1 || len(api.binaries) != 2 {
		t.Errorf("expected the same catalog and binaries, got %d catalogs, %d binaries", len(api.catalogs), len(api.binaries))
	}
}

func TestNvidiaImportRefusesAChangedPackageWithoutForce(t *testing.T) {
	api := &fakeFirmwareApi{}
	ts := api.server()
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := importFixture(t, ctx, nil); err != nil {
		t.Fatal(err)
	}
	// The catalog has another file under the same name
	api.binaries[1]["sha256"] = strings.Repeat("0", 64)

	api.calls = nil
	err := importFixture(t, ctx, nil)
	if err == nil || !strings.Contains(err.Error(), nvidiaFixtureRetimer) || !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected the changed file to be named, got %v", err)
	}
	if changes := append(api.callsWith("PUT"), api.callsWith("POST")...); len(changes) != 0 {
		t.Errorf("nothing is changed before the checksums are checked, got %v", changes)
	}

	if err := importFixture(t, ctx, func(vc *VendorCatalog) { vc.Force = true }); err != nil {
		t.Fatalf("--force replaces the file: %v", err)
	}
	if api.binaries[1]["sha256"] != nvidiaFixtureRetimerSha {
		t.Errorf("sha256 after --force = %v", api.binaries[1]["sha256"])
	}
}

func TestNvidiaServerTypesFilterIsIgnored(t *testing.T) {
	vc := nvidiaFixtureCatalog(t)
	vc.ServerTypesFilter = []string{"M.8.8.2"}
	// No API in the context: a lookup would fail
	if err := vc.ProcessVendorCatalog(context.Background()); err != nil {
		t.Fatalf("the server types filter is not looked up for NVIDIA: %v", err)
	}
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()
	date, err := time.Parse(NvidiaManifestReleaseDateLayout, value)
	if err != nil {
		t.Fatal(err)
	}
	return date
}

func containsIssue(issues []string, prefix string) bool {
	for _, issue := range issues {
		if strings.HasPrefix(issue, prefix) {
			return true
		}
	}
	return false
}
