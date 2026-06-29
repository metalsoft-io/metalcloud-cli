package os_template

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func osTemplateFixture(id int) map[string]any {
	return map[string]any{
		"id":         id,
		"name":       "Ubuntu 22.04",
		"visibility": "public",
		"status":     "active",
		"revision":   1,
		"createdBy":  1,
		"createdAt":  "2024-01-01T00:00:00Z",
		"device":     map[string]any{"type": "server", "bootMode": "uefi", "architecture": "x86_64"},
		"install":    map[string]any{"method": "oob", "driveType": "local_drive", "readyMethod": "wait_for_power_off"},
		"os":         map[string]any{"name": "Ubuntu", "version": "22.04", "credential": map[string]any{"username": "root", "passwordType": "plain"}},
		"imageBuild": map[string]any{"required": false},
	}
}

func TestOsTemplateList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		osTemplateFixture(1),
		osTemplateFixture(2),
	}
	srv := testutils.MultiPageServer("/api/v2/os-templates", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateList(ctx); err != nil {
		t.Fatalf("OsTemplateList() unexpected error: %v", err)
	}
}

func TestOsTemplateList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateList(ctx); err == nil {
		t.Fatal("OsTemplateList() expected error, got nil")
	}
}

func TestOsTemplateList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/os-templates", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateList(ctx); err != nil {
		t.Fatalf("OsTemplateList() unexpected error on empty: %v", err)
	}
}

func TestOsTemplateList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = osTemplateFixture(start + i)
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/os-templates", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateList(ctx); err != nil {
		t.Fatalf("OsTemplateList() pagination error: %v", err)
	}
}

func TestOsTemplateGet_HappyPath(t *testing.T) {
	tmpl := osTemplateFixture(5)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/5": testutils.JSONHandler(200, tmpl),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGet(ctx, "5"); err != nil {
		t.Fatalf("OsTemplateGet() unexpected error: %v", err)
	}
}

func TestOsTemplateGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGet(ctx, "99"); err == nil {
		t.Fatal("OsTemplateGet() expected error for not found, got nil")
	}
}

func TestOsTemplateGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateGet(ctx, "not-a-number"); err == nil {
		t.Fatal("OsTemplateGet() expected error for invalid ID, got nil")
	}
}

func TestOsTemplateDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/7": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateDelete(ctx, "7"); err != nil {
		t.Fatalf("OsTemplateDelete() unexpected error: %v", err)
	}
}

func TestOsTemplateDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/7": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateDelete(ctx, "7"); err == nil {
		t.Fatal("OsTemplateDelete() expected error, got nil")
	}
}

func TestGetOsTemplateByIdOrLabel_HappyPath(t *testing.T) {
	tmpl := osTemplateFixture(12)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/12": testutils.JSONHandler(200, tmpl),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	result, err := GetOsTemplateByIdOrLabel(ctx, "12")
	if err != nil {
		t.Fatalf("GetOsTemplateByIdOrLabel() unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("GetOsTemplateByIdOrLabel() returned nil result")
	}
}

func TestGetOsTemplateByIdOrLabel_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if _, err := GetOsTemplateByIdOrLabel(ctx, "99"); err == nil {
		t.Fatal("GetOsTemplateByIdOrLabel() expected error for not found, got nil")
	}
}

func TestOsTemplateCreate_HappyPath(t *testing.T) {
	tmpl := osTemplateFixture(10)
	assetResp := map[string]any{
		"id":         1,
		"templateId": 10,
		"usage":      "build_component",
		"revision":   1,
		"createdBy":  1,
		"createdAt":  "2024-01-01T00:00:00Z",
		"file":       map[string]any{"name": "test.xml", "mimeType": "text/plain", "templatingEngine": false, "path": "/test.xml"},
	}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates":    testutils.JSONHandler(201, tmpl),
		"/api/v2/template-assets": testutils.JSONHandler(201, assetResp),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	opts := OsTemplateCreateOptions{
		Template: sdk.OSTemplateCreate{
			Name: "Test Template",
			Device: sdk.OSTemplateDevice{
				Type:         "server",
				BootMode:     "uefi",
				Architecture: "x86_64",
			},
			Install: sdk.OSTemplateInstall{
				Method:      "oob",
				DriveType:   "local_drive",
				ReadyMethod: "wait_for_power_off",
			},
			Os: sdk.OSTemplateOs{
				Name:    "Ubuntu",
				Version: "22.04",
				Credential: sdk.OSTemplateOsCredential{
					Username:     "root",
					PasswordType: "plain",
				},
			},
		},
		TemplateAssets: []sdk.TemplateAssetCreate{
			{
				Usage: "build_component",
				File: sdk.TemplateAssetFile{
					Name:             "test.xml",
					MimeType:         "text/plain",
					TemplatingEngine: true,
					Path:             "/test.xml",
				},
			},
		},
	}
	if err := OsTemplateCreate(ctx, opts); err != nil {
		t.Fatalf("OsTemplateCreate() unexpected error: %v", err)
	}
}

func TestOsTemplateCreate_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	opts := OsTemplateCreateOptions{
		Template: sdk.OSTemplateCreate{Name: "bad"},
	}
	if err := OsTemplateCreate(ctx, opts); err == nil {
		t.Fatal("OsTemplateCreate() expected error, got nil")
	}
}

func TestOsTemplateUpdate_HappyPath(t *testing.T) {
	tmpl := osTemplateFixture(5)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/5": testutils.JSONHandler(200, tmpl),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	opts := OsTemplateUpdateOptions{}
	if err := OsTemplateUpdate(ctx, "5", opts); err != nil {
		t.Fatalf("OsTemplateUpdate() unexpected error: %v", err)
	}
}

func TestOsTemplateUpdate_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateUpdate(ctx, "not-a-number", OsTemplateUpdateOptions{}); err == nil {
		t.Fatal("OsTemplateUpdate() expected error for invalid id, got nil")
	}
}

func TestOsTemplateSetStatus_HappyPath(t *testing.T) {
	tmpl := osTemplateFixture(3)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/3": testutils.JSONHandler(200, tmpl),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateSetStatus(ctx, "3", "inactive"); err != nil {
		t.Fatalf("OsTemplateSetStatus() unexpected error: %v", err)
	}
}

func TestOsTemplateSetStatus_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateSetStatus(ctx, "99", "inactive"); err == nil {
		t.Fatal("OsTemplateSetStatus() expected error for not found, got nil")
	}
}

func TestOsTemplateSetStatus_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateSetStatus(ctx, "bad-id", "active"); err == nil {
		t.Fatal("OsTemplateSetStatus() expected error for invalid id, got nil")
	}
}

func TestOsTemplateGetAssets_HappyPath(t *testing.T) {
	assetsResp := map[string]any{
		"data": []map[string]any{
			{
				"id":         1,
				"templateId": 5,
				"usage":      "build_component",
				"revision":   1,
				"createdBy":  1,
				"createdAt":  "2024-01-01T00:00:00Z",
				"file":       map[string]any{"name": "test.xml", "mimeType": "text/plain", "templatingEngine": false, "path": "/test.xml"},
			},
		},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets": testutils.JSONHandler(200, assetsResp),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGetAssets(ctx, "5"); err != nil {
		t.Fatalf("OsTemplateGetAssets() unexpected error: %v", err)
	}
}

func TestOsTemplateGetAssets_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGetAssets(ctx, "5"); err == nil {
		t.Fatal("OsTemplateGetAssets() expected error, got nil")
	}
}

func TestOsTemplateGetAssets_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateGetAssets(ctx, "bad-id"); err == nil {
		t.Fatal("OsTemplateGetAssets() expected error for invalid id, got nil")
	}
}

func TestOsTemplateGetCredentials_HappyPath(t *testing.T) {
	creds := map[string]any{"username": "root", "passwordType": "plain", "password": "secret"}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/5/credentials": testutils.JSONHandler(200, creds),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGetCredentials(ctx, "5"); err != nil {
		t.Fatalf("OsTemplateGetCredentials() unexpected error: %v", err)
	}
}

func TestOsTemplateGetCredentials_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/5/credentials": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateGetCredentials(ctx, "5"); err == nil {
		t.Fatal("OsTemplateGetCredentials() expected error, got nil")
	}
}

func TestOsTemplateGetCredentials_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateGetCredentials(ctx, "bad-id"); err == nil {
		t.Fatal("OsTemplateGetCredentials() expected error for invalid id, got nil")
	}
}

func TestOsTemplateListRepo_InvalidURL(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	// Passing a non-local, non-valid git URL should fail gracefully
	err := OsTemplateListRepo(ctx, "http://localhost:1/nonexistent.git", "", "")
	if err == nil {
		t.Fatal("OsTemplateListRepo() expected error for invalid repo URL, got nil")
	}
}

func TestOsTemplateExport_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateExport(ctx, "bad-id", ""); err == nil {
		t.Fatal("OsTemplateExport() expected error for invalid id, got nil")
	}
}

func TestOsTemplateExport_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := OsTemplateExport(ctx, "99", ""); err == nil {
		t.Fatal("OsTemplateExport() expected error for not found, got nil")
	}
}

func TestOsTemplateImport_InvalidArchive(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateImport(ctx, "/tmp/nonexistent-archive-xyz.zip", "TestTemplate", ""); err == nil {
		t.Fatal("OsTemplateImport() expected error for non-existent archive, got nil")
	}
}
