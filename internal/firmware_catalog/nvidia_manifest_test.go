package firmware_catalog

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	nvidiaFixtureDir        = "testdata/nvidia/hgx-h100"
	nvidiaFixtureBundleFile = "nvfw_HGX_H100x8_250612.fwpkg"
	nvidiaFixtureBundleSha  = "f013d7e18fbde4f8f4cd8fce030ec358d5b50c02d2f711ce3dca0179dfe979c2"
	nvidiaFixtureRetimer    = "nvfw_HGX_H100x8_retimer_250612.fwpkg"
	nvidiaFixtureRetimerSha = "77b6cc60811f5b21ffe0e559364a59f2b6158f5ad10a3bbc2279e712bea75b6d"
)

// loadNvidiaFixture decodes the fixture manifest without validating it, so tests can break it on purpose.
func loadNvidiaFixture(t *testing.T) *NvidiaBundleManifest {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(nvidiaFixtureDir, "manifest.json"))
	if err != nil {
		t.Fatalf("failed to read fixture manifest: %v", err)
	}
	var manifest NvidiaBundleManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		t.Fatalf("failed to decode fixture manifest: %v", err)
	}
	return &manifest
}

// copyNvidiaFixture copies the fixture directory to a temporary directory the test can modify.
func copyNvidiaFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	entries, err := os.ReadDir(nvidiaFixtureDir)
	if err != nil {
		t.Fatalf("failed to list fixture directory: %v", err)
	}
	for _, entry := range entries {
		content, err := os.ReadFile(filepath.Join(nvidiaFixtureDir, entry.Name()))
		if err != nil {
			t.Fatalf("failed to read fixture %s: %v", entry.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(dir, entry.Name()), content, 0644); err != nil {
			t.Fatalf("failed to copy fixture %s: %v", entry.Name(), err)
		}
	}
	return dir
}

func manifestIssues(t *testing.T, err error) []string {
	t.Helper()
	if err == nil {
		t.Fatal("expected the manifest to be rejected, got no error")
	}
	var manifestErr *NvidiaManifestError
	if !errors.As(err, &manifestErr) {
		t.Fatalf("expected *NvidiaManifestError, got %T: %v", err, err)
	}
	return manifestErr.Issues
}

func assertSingleIssue(t *testing.T, err error, expected string) {
	t.Helper()
	issues := manifestIssues(t, err)
	if len(issues) != 1 || !strings.Contains(issues[0], expected) {
		t.Fatalf("expected exactly one issue containing %q, got %q", expected, issues)
	}
}

func nvidiaStringPtr(value string) *string { return &value }

func TestLoadNvidiaManifestFixture(t *testing.T) {
	manifest, err := LoadNvidiaManifest(filepath.Join(nvidiaFixtureDir, "manifest.json"))
	if err != nil {
		t.Fatalf("expected the fixture manifest to be valid, got: %v", err)
	}

	if manifest.BundleName != "HGX H100 8-GPU" || manifest.BundleVersion != "25.06.1" {
		t.Errorf("unexpected bundle %q %q", manifest.BundleName, manifest.BundleVersion)
	}
	if len(manifest.Steps) != 4 {
		t.Errorf("expected 4 steps, got %d", len(manifest.Steps))
	}

	if manifest.Steps[1].Activation != NvidiaActivationHgxReset {
		t.Errorf("expected the HMC step to activate with %s, got %q", NvidiaActivationHgxReset, manifest.Steps[1].Activation)
	}
	if manifest.Steps[2].VerifyVersion != "96.00.9F.00.01" || manifest.Steps[0].VerifyVersion != "" {
		t.Errorf("unexpected verifyVersion values %q %q", manifest.Steps[2].VerifyVersion, manifest.Steps[0].VerifyVersion)
	}

	files := manifest.PackageFiles()
	if strings.Join(files, ",") != nvidiaFixtureBundleFile+","+nvidiaFixtureRetimer {
		t.Errorf("expected the bundle file first and the retimer package second, got %v", files)
	}
	if manifest.PackageFileSha256(nvidiaFixtureRetimer) != nvidiaFixtureRetimerSha {
		t.Errorf("unexpected retimer checksum %s", manifest.PackageFileSha256(nvidiaFixtureRetimer))
	}

	if err := VerifyNvidiaPackageFiles(manifest, nvidiaFixtureDir); err != nil {
		t.Fatalf("expected the fixture package files to match the manifest, got: %v", err)
	}
}

func TestNvidiaManifestValidationRejections(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(m *NvidiaBundleManifest)
		expected string
	}{
		{
			name:     "unknown component class",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[1].ComponentClass = "warp_drive" },
			expected: `Step "hmc": unknown component class "warp_drive"`,
		},
		{
			name:     "invalid target pattern",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[2].TargetPattern = "^HGX_FW_GPU_(" },
			expected: `Step "gpu": targetPattern "^HGX_FW_GPU_(" is not a valid regular expression`,
		},
		{
			name:     "unknown activation",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[0].Activation = "reboot_twice" },
			expected: `Step "erot": unknown activation "reboot_twice"`,
		},
		{
			name:     "verifyVersion too long",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[2].VerifyVersion = strings.Repeat("9", 256) },
			expected: `Step "gpu": verifyVersion must be at most 255 characters`,
		},
		{
			name:     "malformed sha256",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].Sha256 = "abc" },
			expected: `Step "retimer": sha256 must be 64 hexadecimal characters`,
		},
		{
			name:     "missing package version",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].PackageVersion = "" },
			expected: `Step "retimer": packageVersion is required`,
		},
		{
			name:     "duplicate step id",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].Id = "gpu" },
			expected: `Step "gpu": duplicate step id`,
		},
		{
			name:     "dependency on an unknown step",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].DependsOnStep = nvidiaStringPtr("fpga") },
			expected: `Step "retimer": depends on unknown step "fpga"`,
		},
		{
			name:     "step depending on itself",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[1].DependsOnStep = nvidiaStringPtr("hmc") },
			expected: `Step "hmc": depends on itself`,
		},
		{
			name:     "dependency cycle",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[0].DependsOnStep = nvidiaStringPtr("gpu") },
			expected: `Step "erot": dependency cycle erot -> gpu -> hmc -> erot`,
		},
		{
			name: "steps disagreeing on the checksum of a shared file",
			mutate: func(m *NvidiaBundleManifest) {
				m.Steps[1].Sha256 = strings.Repeat("0", 64)
			},
			expected: `Step "hmc": sha256 for ` + nvidiaFixtureBundleFile + ` differs from the one declared by step "erot"`,
		},
		{
			name:     "package file outside the bundle directory",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].PackageFile = "../" + nvidiaFixtureRetimer },
			expected: `Step "retimer": packageFile: "../` + nvidiaFixtureRetimer + `" points outside the bundle directory`,
		},
		{
			name:     "absolute package file",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].PackageFile = "/tmp/" + nvidiaFixtureRetimer },
			expected: `Step "retimer": packageFile: "/tmp/` + nvidiaFixtureRetimer + `" must be a path relative to the bundle directory`,
		},
		{
			name:     "backslash in package file",
			mutate:   func(m *NvidiaBundleManifest) { m.Steps[3].PackageFile = `fw\` + nvidiaFixtureRetimer },
			expected: `Step "retimer": packageFile: "fw\\` + nvidiaFixtureRetimer + `" must use forward slashes`,
		},
		{
			name:     "bundle file not used by any step",
			mutate:   func(m *NvidiaBundleManifest) { m.BundleFile = "nvfw_other.fwpkg" },
			expected: "bundleFile nvfw_other.fwpkg is not the packageFile of any step",
		},
		{
			name:     "no target platforms",
			mutate:   func(m *NvidiaBundleManifest) { m.TargetPlatforms = nil },
			expected: "targetPlatforms must list at least one platform",
		},
		{
			name:     "release date in another format",
			mutate:   func(m *NvidiaBundleManifest) { m.ReleaseDate = "12/06/2025" },
			expected: `releaseDate "12/06/2025" is not a date in YYYY-MM-DD format`,
		},
		{
			name:     "unsupported schema version",
			mutate:   func(m *NvidiaBundleManifest) { m.SchemaVersion = 2 },
			expected: "schemaVersion must be 1, got 2",
		},
		{
			name:     "unknown update severity",
			mutate:   func(m *NvidiaBundleManifest) { m.UpdateSeverity = "urgent" },
			expected: `updateSeverity "urgent" must be one of`,
		},
		{
			name: "download URL that is not http(s)",
			mutate: func(m *NvidiaBundleManifest) {
				m.DownloadUrls = map[string]string{nvidiaFixtureRetimer: "ftp://example.com/retimer.fwpkg"}
			},
			expected: `downloadUrls: "ftp://example.com/retimer.fwpkg" for ` + nvidiaFixtureRetimer + ` is not an http(s) URL`,
		},
		{
			name: "download URL for a file no step uses",
			mutate: func(m *NvidiaBundleManifest) {
				m.DownloadUrls = map[string]string{"nvfw_other.fwpkg": "https://example.com/other.fwpkg"}
			},
			expected: "downloadUrls: nvfw_other.fwpkg is not the packageFile of any step",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manifest := loadNvidiaFixture(t)
			test.mutate(manifest)
			assertSingleIssue(t, manifest.Validate(), test.expected)
		})
	}
}

func TestNvidiaManifestReportsAllIssuesAtOnce(t *testing.T) {
	manifest := loadNvidiaFixture(t)
	manifest.Steps[0].ComponentClass = "warp_drive"
	manifest.Steps[2].Activation = "reboot_twice"

	issues := manifestIssues(t, manifest.Validate())
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %q", issues)
	}
	if !strings.HasPrefix(issues[0], `Step "erot":`) || !strings.HasPrefix(issues[1], `Step "gpu":`) {
		t.Errorf("expected issues for steps erot and gpu in manifest order, got %q", issues)
	}
}

func TestParseNvidiaManifestRejectsUnknownFields(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(nvidiaFixtureDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(content), `"dependsOnStep": "erot"`, `"dependsOn": "erot"`, 1)

	_, err = ParseNvidiaManifest([]byte(broken))
	if err == nil || !strings.Contains(err.Error(), `unknown field "dependsOn"`) {
		t.Fatalf("expected an unknown field error, got: %v", err)
	}
}

func TestVerifyNvidiaPackageFiles(t *testing.T) {
	t.Run("checksum mismatch names the step", func(t *testing.T) {
		dir := copyNvidiaFixture(t)
		if err := os.WriteFile(filepath.Join(dir, nvidiaFixtureRetimer), []byte("tampered"), 0644); err != nil {
			t.Fatal(err)
		}

		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), dir)
		assertSingleIssue(t, err, `Step "retimer": checksum mismatch for `+nvidiaFixtureRetimer+`: manifest has `+nvidiaFixtureRetimerSha)
	})

	t.Run("checksum mismatch of a shared file names every step using it", func(t *testing.T) {
		dir := copyNvidiaFixture(t)
		if err := os.WriteFile(filepath.Join(dir, nvidiaFixtureBundleFile), []byte("tampered"), 0644); err != nil {
			t.Fatal(err)
		}

		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), dir)
		assertSingleIssue(t, err, `Steps "erot", "hmc", "gpu": checksum mismatch for `+nvidiaFixtureBundleFile)
	})

	t.Run("missing package file", func(t *testing.T) {
		dir := copyNvidiaFixture(t)
		if err := os.Remove(filepath.Join(dir, nvidiaFixtureRetimer)); err != nil {
			t.Fatal(err)
		}

		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), dir)
		assertSingleIssue(t, err, `Step "retimer": package file `+nvidiaFixtureRetimer+` not found`)
	})

	t.Run("package file that is a directory", func(t *testing.T) {
		dir := copyNvidiaFixture(t)
		if err := os.Remove(filepath.Join(dir, nvidiaFixtureRetimer)); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(dir, nvidiaFixtureRetimer), 0755); err != nil {
			t.Fatal(err)
		}

		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), dir)
		assertSingleIssue(t, err, `Step "retimer": package file `+nvidiaFixtureRetimer+` is not a regular file`)
	})

	t.Run("symlink pointing outside the bundle directory", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("creating symlinks needs extra privileges on Windows")
		}
		dir := copyNvidiaFixture(t)
		outside := filepath.Join(t.TempDir(), nvidiaFixtureRetimer)
		if err := os.Rename(filepath.Join(dir, nvidiaFixtureRetimer), outside); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(dir, nvidiaFixtureRetimer)); err != nil {
			t.Fatal(err)
		}

		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), dir)
		assertSingleIssue(t, err, `Step "retimer": package file `+nvidiaFixtureRetimer+` resolves outside`)
	})

	t.Run("missing binaries directory", func(t *testing.T) {
		err := VerifyNvidiaPackageFiles(loadNvidiaFixture(t), filepath.Join(t.TempDir(), "missing"))
		if err == nil || !strings.Contains(err.Error(), "is not accessible") {
			t.Fatalf("expected an inaccessible directory error, got: %v", err)
		}
	})
}

func TestNvidiaFixtureChecksumsAreCurrent(t *testing.T) {
	// Guards the constants above: if a fixture file is edited, its checksum has to be updated everywhere.
	for file, expected := range map[string]string{
		nvidiaFixtureBundleFile: nvidiaFixtureBundleSha,
		nvidiaFixtureRetimer:    nvidiaFixtureRetimerSha,
	} {
		actual, err := fileSha256(filepath.Join(nvidiaFixtureDir, file))
		if err != nil {
			t.Fatal(err)
		}
		if actual != expected {
			t.Errorf("fixture %s has checksum %s, expected %s", file, actual, expected)
		}
	}
}
