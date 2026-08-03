package os_template

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const directoryTemplateYaml = `template:
  name: Ubuntu 24.04 Test
  label: ubuntu-2404-test
  device:
    type: server
    bootmode: uefi
    architecture: x86_64
  install:
    method: oob
    drivetype: local_drive
    readymethod: wait_for_power_off
  os:
    name: Ubuntu
    version: "24.04"
    credential:
      username: root
      passwordtype: plain
  visibility: public
templateassets:
  - usage: build_source_image
    file:
      name: ubuntu-24.04.iso
      mimetype: application/octet-stream
      url: http://repo.local/ubuntu-24.04.iso
      path: /ubuntu-24.04.iso
  - usage: build_component
    file:
      name: autoinstall.yaml
      mimetype: text/plain
      templatingengine: true
      path: /autoinstall.yaml
`

const directoryAssetContent = "autoinstall: content\n"

// newTemplateDirectory creates a local directory holding one template in the layout used by the
// template repositories, along with files that must be ignored while reading it.
func newTemplateDirectory(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	templateDir := filepath.Join(dir, "Ubuntu", "24.04", "test")
	if err := os.MkdirAll(templateDir, 0o755); err != nil {
		t.Fatalf("failed to create template directory: %v", err)
	}

	files := map[string]string{
		filepath.Join(templateDir, templateFileName):   directoryTemplateYaml,
		filepath.Join(templateDir, "autoinstall.yaml"): directoryAssetContent,
		filepath.Join(templateDir, readMeFileName):     "# ignored",
		// Files outside the <vendor>/<os>/<version> layout must be ignored
		filepath.Join(dir, readMeFileName):                      "# ignored",
		filepath.Join(dir, "Ubuntu", "24.04", templateFileName): "ignored: true",
		filepath.Join(templateDir, "extra", templateFileName):   "ignored: true",
	}
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("failed to write %s: %v", path, err)
		}
	}

	return dir
}

func TestGetDirectoryTemplateAssets(t *testing.T) {
	dir := newTemplateDirectory(t)

	assets, err := getDirectoryTemplateAssets(dir)
	if err != nil {
		t.Fatalf("getDirectoryTemplateAssets() unexpected error: %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("getDirectoryTemplateAssets() returned %d templates, expected 1: %v", len(assets), assets)
	}

	template, ok := assets["Ubuntu/24.04/test"]
	if !ok {
		t.Fatalf("getDirectoryTemplateAssets() did not return template Ubuntu/24.04/test, got %v", assets)
	}

	if template.SourcePath != "Ubuntu/24.04/test" {
		t.Errorf("SourcePath = %q, expected Ubuntu/24.04/test", template.SourcePath)
	}
	if template.SourceContent != directoryTemplateYaml {
		t.Errorf("SourceContent = %q, expected the template.yaml content", template.SourceContent)
	}

	if len(template.Assets) != 1 {
		t.Fatalf("template has %d assets, expected 1 (README.md must be skipped): %v", len(template.Assets), template.Assets)
	}

	expectedContent := base64.StdEncoding.EncodeToString([]byte(directoryAssetContent))
	if template.Assets["autoinstall.yaml"].ContentBase64 != expectedContent {
		t.Errorf("asset content = %q, expected %q", template.Assets["autoinstall.yaml"].ContentBase64, expectedContent)
	}

	// The template definition and its local assets must be linked together
	if err := processTemplateContent(&template); err != nil {
		t.Fatalf("processTemplateContent() unexpected error: %v", err)
	}
	if template.OsTemplate.Template.Name != "Ubuntu 24.04 Test" {
		t.Errorf("template name = %q, expected Ubuntu 24.04 Test", template.OsTemplate.Template.Name)
	}
	if len(template.OsTemplate.TemplateAssets) != 2 {
		t.Fatalf("template has %d asset definitions, expected 2", len(template.OsTemplate.TemplateAssets))
	}
	if template.OsTemplate.TemplateAssets[0].File.ContentBase64 != nil {
		t.Error("URL based asset should not have its content populated")
	}
	if template.OsTemplate.TemplateAssets[1].File.ContentBase64 == nil ||
		*template.OsTemplate.TemplateAssets[1].File.ContentBase64 != expectedContent {
		t.Error("local asset should have its content populated from the directory")
	}
	expectedChecksum := fmt.Sprintf("%x", sha256.Sum256([]byte(expectedContent)))
	if template.OsTemplate.TemplateAssets[1].File.Checksum == nil ||
		*template.OsTemplate.TemplateAssets[1].File.Checksum != expectedChecksum {
		t.Error("local asset should have its checksum populated")
	}
}

func TestGetDirectoryTemplateAssets_NotADirectory(t *testing.T) {
	if _, err := getDirectoryTemplateAssets(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("getDirectoryTemplateAssets() expected error for missing directory, got nil")
	}

	file := filepath.Join(t.TempDir(), templateFileName)
	if err := os.WriteFile(file, []byte("template: {}"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if _, err := getDirectoryTemplateAssets(file); err == nil {
		t.Fatal("getDirectoryTemplateAssets() expected error for a file path, got nil")
	}
}

func TestOsTemplateListDirectory_HappyPath(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateListDirectory(ctx, newTemplateDirectory(t)); err != nil {
		t.Fatalf("OsTemplateListDirectory() unexpected error: %v", err)
	}
}

func TestOsTemplateListDirectory_MissingDirectory(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := OsTemplateListDirectory(ctx, filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("OsTemplateListDirectory() expected error for missing directory, got nil")
	}
}

func TestOsTemplateCreateFromDirectory_HappyPath(t *testing.T) {
	var createdTemplate sdk.OSTemplateCreate
	createdAssets := make([]sdk.TemplateAssetCreate, 0, 2)

	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates": func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(&createdTemplate); err != nil {
				t.Errorf("failed to decode template request: %v", err)
			}
			testutils.JSONHandler(201, osTemplateFixture(7))(w, r)
		},
		"/api/v2/template-assets": func(w http.ResponseWriter, r *http.Request) {
			var asset sdk.TemplateAssetCreate
			if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
				t.Errorf("failed to decode asset request: %v", err)
			}
			createdAssets = append(createdAssets, asset)
			testutils.JSONHandler(201, map[string]any{
				"id":         1,
				"templateId": 7,
				"usage":      asset.Usage,
				"revision":   1,
				"createdBy":  1,
				"createdAt":  "2024-01-01T00:00:00Z",
				"file":       map[string]any{"name": asset.File.Name, "mimeType": asset.File.MimeType, "templatingEngine": false, "path": asset.File.Path},
			})(w, r)
		},
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	dir := newTemplateDirectory(t)

	err := OsTemplateCreateFromDirectory(ctx, "Ubuntu/24.04/test", dir, "My Ubuntu", "my-ubuntu", "http://repo.local/custom.iso")
	if err != nil {
		t.Fatalf("OsTemplateCreateFromDirectory() unexpected error: %v", err)
	}

	if createdTemplate.Name != "My Ubuntu" {
		t.Errorf("created template name = %q, expected the overridden name My Ubuntu", createdTemplate.Name)
	}
	if createdTemplate.Label == nil || *createdTemplate.Label != "my-ubuntu" {
		t.Errorf("created template label = %v, expected the overridden label my-ubuntu", createdTemplate.Label)
	}
	if createdTemplate.Visibility == nil || *createdTemplate.Visibility != "private" {
		t.Errorf("created template visibility = %v, expected private", createdTemplate.Visibility)
	}

	if len(createdAssets) != 2 {
		t.Fatalf("created %d assets, expected 2", len(createdAssets))
	}
	for _, asset := range createdAssets {
		switch asset.Usage {
		case "build_source_image":
			if asset.File.Url == nil || *asset.File.Url != "http://repo.local/custom.iso" {
				t.Errorf("source image URL = %v, expected the overridden source ISO", asset.File.Url)
			}
		case "build_component":
			expectedContent := base64.StdEncoding.EncodeToString([]byte(directoryAssetContent))
			if asset.File.ContentBase64 == nil || *asset.File.ContentBase64 != expectedContent {
				t.Errorf("asset content = %v, expected the content read from the directory", asset.File.ContentBase64)
			}
		default:
			t.Errorf("unexpected asset usage %q", asset.Usage)
		}
	}
}

func TestOsTemplateCreateFromDirectory_ApiErrorNamesSource(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/os-templates": testutils.ErrorHandler(400, "Password is required for PLAIN or HASHED password type"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	dir := newTemplateDirectory(t)

	err := OsTemplateCreateFromDirectory(ctx, "Ubuntu/24.04/test", dir, "", "", "")
	if err == nil {
		t.Fatal("OsTemplateCreateFromDirectory() expected error, got nil")
	}

	// The API only reports what is wrong with the definition, so the error must say which template it was
	if !strings.Contains(err.Error(), "Ubuntu/24.04/test") {
		t.Errorf("error %q should name the source template", err)
	}
	if !strings.Contains(err.Error(), dir) {
		t.Errorf("error %q should name the source directory", err)
	}
	if !strings.Contains(err.Error(), "Password is required") {
		t.Errorf("error %q should keep the message reported by the API", err)
	}
}

func TestOsTemplateCreateFromDirectory_TemplateNotFound(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")

	err := OsTemplateCreateFromDirectory(ctx, "Ubuntu/22.04/missing", newTemplateDirectory(t), "", "", "")
	if err == nil {
		t.Fatal("OsTemplateCreateFromDirectory() expected error for missing template, got nil")
	}
}

func TestOsTemplateCreateFromDirectory_MissingDirectory(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")

	err := OsTemplateCreateFromDirectory(ctx, "Ubuntu/24.04/test", filepath.Join(t.TempDir(), "missing"), "", "", "")
	if err == nil {
		t.Fatal("OsTemplateCreateFromDirectory() expected error for missing directory, got nil")
	}
}
