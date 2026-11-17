package os_template

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"gopkg.in/yaml.v3"
)

// escapePath returns a relative path that climbs out of any temp directory and lands on target.
func escapePath(target string) string {
	return strings.Repeat("../", 32) + strings.TrimPrefix(filepath.ToSlash(target), "/")
}

func minimalTemplateCreate() sdk.OSTemplateCreate {
	return sdk.OSTemplateCreate{
		Name:    "Archived",
		Device:  sdk.OSTemplateDevice{Type: "server", BootMode: "uefi", Architecture: "x86_64"},
		Install: sdk.OSTemplateInstall{Method: "oob", DriveType: "local_drive", ReadyMethod: "wait_for_power_off"},
		Os: sdk.OSTemplateOs{
			Name: "Ubuntu", Version: "22.04",
			Credential: sdk.OSTemplateOsCredential{Username: "root", PasswordType: "plain"},
		},
	}
}

// writeImportArchive builds an import archive the same way export does: template.yaml marshalled
// from OsTemplateCreateOptions, plus the given files under assets/.
func writeImportArchive(t *testing.T, assets []sdk.TemplateAssetCreate, files map[string]string) string {
	t.Helper()

	yamlData, err := yaml.Marshal(OsTemplateCreateOptions{Template: minimalTemplateCreate(), TemplateAssets: assets})
	if err != nil {
		t.Fatalf("marshal template.yaml: %v", err)
	}

	archivePath := filepath.Join(t.TempDir(), "template.zip")
	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	zw := zip.NewWriter(f)
	entries := map[string]string{templateFileName: string(yamlData)}
	for name, content := range files {
		entries["assets/"+name] = content
	}
	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("create zip entry %s: %v", name, err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("write zip entry %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close archive: %v", err)
	}
	return archivePath
}

// importServer records how many API calls were made and the asset bodies that were posted.
type importServer struct {
	calls  atomic.Int32
	assets []sdk.TemplateAssetCreate
}

func (s *importServer) start(t *testing.T) string {
	t.Helper()
	tmpl := osTemplateFixture(10)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates": func(w http.ResponseWriter, r *http.Request) {
			s.calls.Add(1)
			testutils.JSONHandler(201, tmpl)(w, r)
		},
		"/api/v2/template-assets": func(w http.ResponseWriter, r *http.Request) {
			s.calls.Add(1)
			body, _ := io.ReadAll(r.Body)
			var asset sdk.TemplateAssetCreate
			if err := json.Unmarshal(body, &asset); err != nil {
				t.Errorf("decode asset body: %v", err)
			}
			s.assets = append(s.assets, asset)
			testutils.JSONHandler(201, map[string]any{
				"id": len(s.assets), "templateId": 10, "usage": asset.Usage, "revision": 1,
				"createdBy": 1, "createdAt": "2024-01-01T00:00:00Z",
				"file": map[string]any{"name": asset.File.Name, "mimeType": asset.File.MimeType, "templatingEngine": false, "path": asset.File.Path},
			})(w, r)
		},
	})
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestValidateAssetFileName(t *testing.T) {
	valid := []string{"ks.cfg", "sub/dir/file.txt", "..hidden", "file..name"}
	for _, name := range valid {
		if err := validateAssetFileName(name); err != nil {
			t.Errorf("validateAssetFileName(%q) unexpected error: %v", name, err)
		}
	}

	invalid := []string{"", ".", "..", "../secret", "a/../../secret", "/etc/passwd", escapePath("/etc/passwd")}
	for _, name := range invalid {
		if err := validateAssetFileName(name); err == nil {
			t.Errorf("validateAssetFileName(%q) expected error, got nil", name)
		}
	}
}

func TestOsTemplateImport_HappyPath(t *testing.T) {
	archive := writeImportArchive(t,
		[]sdk.TemplateAssetCreate{
			{Usage: "build_component", File: sdk.TemplateAssetFile{Name: "ks.cfg", MimeType: "text/plain", Path: "/ks.cfg"}},
			{Usage: "install_component", File: sdk.TemplateAssetFile{Name: "os.iso", MimeType: "application/octet-stream", Path: "/os.iso", Url: sdk.PtrString("https://example.com/os.iso")}},
		},
		map[string]string{"ks.cfg": "kickstart-content"},
	)

	s := &importServer{}
	ctx := testutils.SetupTestContext(s.start(t))
	if err := OsTemplateImport(ctx, archive, "Imported", ""); err != nil {
		t.Fatalf("OsTemplateImport() unexpected error: %v", err)
	}

	if len(s.assets) != 2 {
		t.Fatalf("expected 2 assets posted, got %d", len(s.assets))
	}
	got := s.assets[0].File.ContentBase64
	want := base64.StdEncoding.EncodeToString([]byte("kickstart-content"))
	if got == nil || *got != want {
		t.Errorf("asset content = %v, want %q", got, want)
	}
	if s.assets[1].File.ContentBase64 != nil {
		t.Errorf("URL asset should not carry content, got %q", *s.assets[1].File.ContentBase64)
	}
}

func TestOsTemplateImport_RejectsAssetPathTraversal(t *testing.T) {
	secret := filepath.Join(t.TempDir(), "id_ed25519")
	if err := os.WriteFile(secret, []byte("PRIVATE KEY"), 0600); err != nil {
		t.Fatalf("write secret: %v", err)
	}

	names := map[string]string{
		"relative traversal": escapePath(secret),
		"absolute path":      secret,
		"nested traversal":   "sub/" + escapePath(secret),
	}
	for label, name := range names {
		t.Run(label, func(t *testing.T) {
			archive := writeImportArchive(t,
				[]sdk.TemplateAssetCreate{{Usage: "build_component", File: sdk.TemplateAssetFile{Name: name, MimeType: "text/plain", Path: "/asset"}}},
				nil,
			)

			s := &importServer{}
			ctx := testutils.SetupTestContext(s.start(t))
			err := OsTemplateImport(ctx, archive, "Malicious", "")
			if err == nil {
				t.Fatal("OsTemplateImport() expected error for traversing asset name, got nil")
			}
			if n := s.calls.Load(); n != 0 {
				t.Errorf("expected no API calls before rejecting the archive, got %d", n)
			}
			for _, a := range s.assets {
				if a.File.ContentBase64 != nil {
					t.Errorf("asset content was uploaded: %q", *a.File.ContentBase64)
				}
			}
		})
	}
}

func exportServer(t *testing.T, assetName string, content string) string {
	t.Helper()
	assetFixture := map[string]any{
		"id": 1, "templateId": 7, "usage": "build_component", "revision": 1,
		"createdBy": 1, "createdAt": "2024-01-01T00:00:00Z",
		"file": map[string]any{"name": assetName, "mimeType": "text/plain", "templatingEngine": false, "path": "/asset"},
	}
	fullAsset := map[string]any{}
	for k, v := range assetFixture {
		fullAsset[k] = v
	}
	fullAsset["file"] = map[string]any{
		"name": assetName, "mimeType": "text/plain", "templatingEngine": false, "path": "/asset",
		"contentBase64": base64.StdEncoding.EncodeToString([]byte(content)),
	}

	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates/7":    testutils.JSONHandler(200, osTemplateFixture(7)),
		"/api/v2/template-assets":   testutils.JSONHandler(200, testutils.PaginatedResponse([]any{assetFixture}, 1, 1)),
		"/api/v2/template-assets/1": testutils.JSONHandler(200, fullAsset),
	})
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestOsTemplateExport_RejectsAssetPathTraversal(t *testing.T) {
	target := filepath.Join(t.TempDir(), "planted.sh")

	ctx := testutils.SetupTestContext(exportServer(t, escapePath(target), "malicious"))
	outputPath := filepath.Join(t.TempDir(), "out.zip")
	if err := OsTemplateExport(ctx, "7", outputPath); err == nil {
		t.Fatal("OsTemplateExport() expected error for traversing asset name, got nil")
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("file outside the export directory was written: %s (stat err: %v)", target, err)
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Errorf("archive should not be produced on failure, stat err: %v", err)
	}
}

func TestOsTemplateExportImport_RoundTrip(t *testing.T) {
	ctx := testutils.SetupTestContext(exportServer(t, "ks.cfg", "kickstart-content"))
	archive := filepath.Join(t.TempDir(), "out.zip")
	testutils.CaptureStdout(t, func() {
		if err := OsTemplateExport(ctx, "7", archive); err != nil {
			t.Fatalf("OsTemplateExport() unexpected error: %v", err)
		}
	})

	s := &importServer{}
	ctx = testutils.SetupTestContext(s.start(t))
	if err := OsTemplateImport(ctx, archive, "Round Trip", ""); err != nil {
		t.Fatalf("OsTemplateImport() unexpected error: %v", err)
	}
	if len(s.assets) != 1 || s.assets[0].File.ContentBase64 == nil ||
		*s.assets[0].File.ContentBase64 != base64.StdEncoding.EncodeToString([]byte("kickstart-content")) {
		t.Fatalf("round-tripped asset content mismatch: %+v", s.assets)
	}
}
