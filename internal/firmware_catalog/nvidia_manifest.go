package firmware_catalog

// NVIDIA publishes no machine-readable firmware catalog. HGX and DGX bundles are downloaded manually from
// the NVIDIA Enterprise Support Portal or NGC and described to MetalSoft by a JSON manifest that sits next
// to the package files. This file reads that manifest, validates it and verifies the package files.
//
// Only the Go standard library is used here, so the manifest handling does not depend on the SDK.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

// NvidiaManifestSchemaVersion is the only manifest schema version currently supported.
const NvidiaManifestSchemaVersion = 1

// NvidiaManifestReleaseDateLayout is the layout of the manifest releaseDate field.
const NvidiaManifestReleaseDateLayout = "2006-01-02"

// Activation values of a bundle step: what has to happen for the new firmware to become active.
const (
	NvidiaActivationNone              = "none"
	NvidiaActivationHostReset         = "host_reset"
	NvidiaActivationBmcReset          = "bmc_reset"
	NvidiaActivationChassisPowerCycle = "chassis_power_cycle"
	// NvidiaActivationHgxReset resets the GPU baseboard management controller (the HGX HMC), for HGX
	// firmware that needs no chassis power cycle.
	NvidiaActivationHgxReset = "hgx_reset"
)

// ValidNvidiaActivations lists the supported bundle step activation values.
var ValidNvidiaActivations = []string{
	NvidiaActivationNone,
	NvidiaActivationHostReset,
	NvidiaActivationBmcReset,
	NvidiaActivationChassisPowerCycle,
	NvidiaActivationHgxReset,
}

// ValidServerComponentClasses mirrors ServerComponentClass in the MetalSoft inventory service
// (libs/dto/server.component.dto.ts). Keep the two lists in sync.
var ValidServerComponentClasses = []string{
	"bios",
	"bmc",
	"nic",
	"raid",
	"gpu",
	"gpu_erot",
	"hgx_hmc",
	"hgx_fpga",
	"nvswitch",
	"pcie_switch",
	"pcie_retimer",
	"other",
}

var nvidiaSha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

// NvidiaBundleManifest describes one NVIDIA firmware bundle.
type NvidiaBundleManifest struct {
	SchemaVersion   int                `json:"schemaVersion"`
	BundleName      string             `json:"bundleName"`
	BundleVersion   string             `json:"bundleVersion"`
	BundleFile      string             `json:"bundleFile"`
	TargetPlatforms []string           `json:"targetPlatforms"`
	ReleaseDate     string             `json:"releaseDate"`
	ReleaseNotesUrl string             `json:"releaseNotesUrl,omitempty"`
	UpdateSeverity  string             `json:"updateSeverity,omitempty"`
	DownloadUrls    map[string]string  `json:"downloadUrls,omitempty"`
	Steps           []NvidiaBundleStep `json:"steps"`
	// Prerequisites are versions other components must already have, e.g. the minimum host BIOS and BMC per
	// platform from the NVIDIA release notes. They are copied to the bundle binary.
	Prerequisites []NvidiaPrerequisite `json:"prerequisites,omitempty"`
}

// NvidiaPrerequisite is a version another component of the server must have before the bundle is scheduled.
type NvidiaPrerequisite struct {
	// Platform limits the prerequisite to one of the bundle's targetPlatforms; empty means every platform
	Platform       string `json:"platform,omitempty"`
	ComponentClass string `json:"componentClass"`
	ComponentModel string `json:"componentModel,omitempty"`
	ComponentName  string `json:"componentName,omitempty"`
	MinimumVersion string `json:"minimumVersion"`
}

// NvidiaBundleStep is one ordered step of a bundle.
type NvidiaBundleStep struct {
	Id             string  `json:"id"`
	ComponentClass string  `json:"componentClass"`
	TargetPattern  string  `json:"targetPattern"`
	PackageFile    string  `json:"packageFile"`
	PackageVersion string  `json:"packageVersion"`
	Sha256         string  `json:"sha256"`
	Activation     string  `json:"activation"`
	DependsOnStep  *string `json:"dependsOnStep"`
	// VerifyVersion is the version the component reports in FirmwareInventory once the step is active
	// (e.g. the GPU VBIOS version of a bundle). Optional; when set, the update fails if the BMC reports
	// another version after activation.
	VerifyVersion string `json:"verifyVersion,omitempty"`
}

// NvidiaManifestError lists every problem found in a manifest or its package files.
// Each issue about a step starts with `Step "<id>":`.
type NvidiaManifestError struct {
	Issues []string
}

func (e *NvidiaManifestError) Error() string {
	return "invalid NVIDIA bundle manifest:\n  - " + strings.Join(e.Issues, "\n  - ")
}

// LoadNvidiaManifest reads and validates a manifest file. It does not look at the package files,
// see VerifyNvidiaPackageFiles.
func LoadNvidiaManifest(manifestPath string) (*NvidiaBundleManifest, error) {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read NVIDIA bundle manifest: %w", err)
	}

	return ParseNvidiaManifest(content)
}

// ParseNvidiaManifest decodes and validates manifest JSON. Unknown fields are rejected so that typos
// (for example "dependsOn" instead of "dependsOnStep") do not go unnoticed.
func ParseNvidiaManifest(content []byte) (*NvidiaBundleManifest, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()

	var manifest NvidiaBundleManifest
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse NVIDIA bundle manifest: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// Validate checks the manifest fields and the relations between its steps. All problems are reported at once.
func (m *NvidiaBundleManifest) Validate() error {
	issues := []string{}
	addIssue := func(format string, args ...any) {
		issues = append(issues, fmt.Sprintf(format, args...))
	}

	if m.SchemaVersion != NvidiaManifestSchemaVersion {
		addIssue("schemaVersion must be %d, got %d", NvidiaManifestSchemaVersion, m.SchemaVersion)
	}
	if strings.TrimSpace(m.BundleName) == "" {
		addIssue("bundleName is required")
	}
	if strings.TrimSpace(m.BundleVersion) == "" {
		addIssue("bundleVersion is required")
	}
	if len(m.TargetPlatforms) == 0 {
		addIssue("targetPlatforms must list at least one platform (server vendor and model, e.g. \"NVIDIA DGXH100\")")
	}
	for i, platform := range m.TargetPlatforms {
		if strings.TrimSpace(platform) == "" {
			addIssue("targetPlatforms[%d] is empty", i)
		}
	}
	if m.ReleaseDate == "" {
		addIssue("releaseDate is required (YYYY-MM-DD)")
	} else if _, err := time.Parse(NvidiaManifestReleaseDateLayout, m.ReleaseDate); err != nil {
		addIssue("releaseDate %q is not a date in YYYY-MM-DD format", m.ReleaseDate)
	}
	if m.ReleaseNotesUrl != "" && !isHttpUrl(m.ReleaseNotesUrl) {
		addIssue("releaseNotesUrl %q is not an http(s) URL", m.ReleaseNotesUrl)
	}
	if m.UpdateSeverity != "" && !slices.Contains(
		[]string{UpdateSeverityCritical, UpdateSeverityRecommended, UpdateSeverityOptional, UpdateSeverityUnknown},
		m.UpdateSeverity,
	) {
		addIssue("updateSeverity %q must be one of critical, recommended, optional, unknown", m.UpdateSeverity)
	}
	if m.BundleFile == "" {
		addIssue("bundleFile is required")
	} else if err := validatePackageFilePath(m.BundleFile); err != nil {
		addIssue("bundleFile: %v", err)
	}

	if len(m.Steps) == 0 {
		addIssue("steps must contain at least one step")
	}

	stepIds := map[string]bool{}
	fileChecksums := map[string]string{}
	fileFirstStep := map[string]string{}
	bundleFileUsed := false

	for i, step := range m.Steps {
		label := fmt.Sprintf("Step %q", step.Id)
		if step.Id == "" {
			label = fmt.Sprintf("Step #%d", i+1)
			addIssue("%s: id is required", label)
		} else if stepIds[step.Id] {
			addIssue("%s: duplicate step id", label)
		}
		stepIds[step.Id] = true

		if !slices.Contains(ValidServerComponentClasses, step.ComponentClass) {
			addIssue("%s: unknown component class %q, must be one of %s",
				label, step.ComponentClass, strings.Join(ValidServerComponentClasses, ", "))
		}

		if step.TargetPattern == "" {
			addIssue("%s: targetPattern is required", label)
		} else if _, err := regexp.Compile(step.TargetPattern); err != nil {
			addIssue("%s: targetPattern %q is not a valid regular expression: %v", label, step.TargetPattern, err)
		}

		if step.PackageVersion == "" {
			addIssue("%s: packageVersion is required", label)
		}

		if len(step.VerifyVersion) > 255 {
			addIssue("%s: verifyVersion must be at most 255 characters", label)
		}

		if !slices.Contains(ValidNvidiaActivations, step.Activation) {
			addIssue("%s: unknown activation %q, must be one of %s",
				label, step.Activation, strings.Join(ValidNvidiaActivations, ", "))
		}

		checksumValid := nvidiaSha256Regex.MatchString(step.Sha256)
		if !checksumValid {
			addIssue("%s: sha256 must be 64 hexadecimal characters", label)
		}

		if step.PackageFile == "" {
			addIssue("%s: packageFile is required", label)
		} else if err := validatePackageFilePath(step.PackageFile); err != nil {
			addIssue("%s: packageFile: %v", label, err)
		} else {
			if step.PackageFile == m.BundleFile {
				bundleFileUsed = true
			}
			// Steps that share a package file must agree on its checksum
			if checksumValid {
				checksum := strings.ToLower(step.Sha256)
				if previous, seen := fileChecksums[step.PackageFile]; seen && previous != checksum {
					addIssue("%s: sha256 for %s differs from the one declared by step %q",
						label, step.PackageFile, fileFirstStep[step.PackageFile])
				} else if !seen {
					fileChecksums[step.PackageFile] = checksum
					fileFirstStep[step.PackageFile] = step.Id
				}
			}
		}
	}

	if m.BundleFile != "" && len(m.Steps) > 0 && !bundleFileUsed {
		addIssue("bundleFile %s is not the packageFile of any step", m.BundleFile)
	}

	for file, downloadUrl := range sortedMap(m.DownloadUrls) {
		if _, used := fileFirstStep[file]; !used && file != m.BundleFile {
			addIssue("downloadUrls: %s is not the packageFile of any step", file)
		}
		if !isHttpUrl(downloadUrl) {
			addIssue("downloadUrls: %q for %s is not an http(s) URL", downloadUrl, file)
		}
	}

	for i, prerequisite := range m.Prerequisites {
		label := fmt.Sprintf("prerequisites[%d]", i)
		if !slices.Contains(ValidServerComponentClasses, prerequisite.ComponentClass) {
			addIssue("%s: unknown component class %q, must be one of %s",
				label, prerequisite.ComponentClass, strings.Join(ValidServerComponentClasses, ", "))
		}
		if strings.TrimSpace(prerequisite.MinimumVersion) == "" {
			addIssue("%s: minimumVersion is required", label)
		}
		if prerequisite.Platform != "" && !slices.ContainsFunc(m.TargetPlatforms, func(platform string) bool {
			return SameFirmwarePlatform(platform, prerequisite.Platform)
		}) {
			addIssue("%s: platform %q is not one of targetPlatforms", label, prerequisite.Platform)
		}
	}

	issues = append(issues, findNvidiaStepDependencyIssues(m.Steps)...)

	if len(issues) > 0 {
		return &NvidiaManifestError{Issues: issues}
	}

	return nil
}

// PackageFiles returns the distinct package files of the bundle, the bundle file first and the others in step order.
func (m *NvidiaBundleManifest) PackageFiles() []string {
	files := []string{}
	if m.BundleFile != "" {
		files = append(files, m.BundleFile)
	}
	for _, step := range m.Steps {
		if step.PackageFile != "" && !slices.Contains(files, step.PackageFile) {
			files = append(files, step.PackageFile)
		}
	}
	return files
}

// PackageFileSha256 returns the checksum the steps declare for a package file.
func (m *NvidiaBundleManifest) PackageFileSha256(packageFile string) string {
	for _, step := range m.Steps {
		if step.PackageFile == packageFile {
			return strings.ToLower(step.Sha256)
		}
	}
	return ""
}

// VerifyNvidiaPackageFiles checks that every package file of the manifest is a regular file inside
// binariesDir and that its SHA-256 matches the manifest. Every problem names the steps using the file.
func VerifyNvidiaPackageFiles(manifest *NvidiaBundleManifest, binariesDir string) error {
	root, err := filepath.Abs(binariesDir)
	if err != nil {
		return fmt.Errorf("invalid binaries directory %s: %w", binariesDir, err)
	}
	if resolvedRoot, err := filepath.EvalSymlinks(root); err == nil {
		root = resolvedRoot
	} else {
		return fmt.Errorf("binaries directory %s is not accessible: %w", binariesDir, err)
	}

	issues := []string{}
	for _, packageFile := range manifest.PackageFiles() {
		stepsLabel := nvidiaStepsLabel(manifest, packageFile)

		localPath, err := resolvePackageFile(root, packageFile)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", stepsLabel, err))
			continue
		}

		actual, err := fileSha256(localPath)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", stepsLabel, err))
			continue
		}

		expected := manifest.PackageFileSha256(packageFile)
		if actual != expected {
			issues = append(issues, fmt.Sprintf("%s: checksum mismatch for %s: manifest has %s, file has %s",
				stepsLabel, packageFile, expected, actual))
		}
	}

	if len(issues) > 0 {
		return &NvidiaManifestError{Issues: issues}
	}

	return nil
}

// findNvidiaStepDependencyIssues reports dependencies on unknown steps, steps that depend on themselves
// and dependency cycles. Each step has at most one dependency, so following dependsOnStep either ends
// or loops; each cycle is reported once, starting from the step that comes first in the manifest.
func findNvidiaStepDependencyIssues(steps []NvidiaBundleStep) []string {
	issues := []string{}
	dependsOn := map[string]string{}
	for _, step := range steps {
		if step.Id == "" {
			continue
		}
		if step.DependsOnStep != nil {
			dependsOn[step.Id] = *step.DependsOnStep
		} else if _, exists := dependsOn[step.Id]; !exists {
			dependsOn[step.Id] = ""
		}
	}

	for _, step := range steps {
		if step.Id == "" || step.DependsOnStep == nil {
			continue
		}
		dependency := *step.DependsOnStep
		switch {
		case dependency == "":
			issues = append(issues, fmt.Sprintf("Step %q: dependsOnStep is empty, use null for no dependency", step.Id))
		case dependency == step.Id:
			issues = append(issues, fmt.Sprintf("Step %q: depends on itself", step.Id))
		default:
			if _, exists := dependsOn[dependency]; !exists {
				issues = append(issues, fmt.Sprintf("Step %q: depends on unknown step %q", step.Id, dependency))
			}
		}
	}

	reported := map[string]bool{}
	for _, step := range steps {
		if step.Id == "" {
			continue
		}
		visited := []string{}
		current := step.Id
		for current != "" && !slices.Contains(visited, current) {
			if _, exists := dependsOn[current]; !exists {
				break
			}
			visited = append(visited, current)
			current = dependsOn[current]
		}

		start := slices.Index(visited, current)
		if current == "" || start < 0 || len(visited)-start < 2 {
			continue // no cycle, or a self dependency (reported above)
		}

		cycle := visited[start:]
		key := slices.Clone(cycle)
		sort.Strings(key)
		if reported[strings.Join(key, "|")] {
			continue
		}
		reported[strings.Join(key, "|")] = true
		issues = append(issues, fmt.Sprintf("Step %q: dependency cycle %s -> %s",
			cycle[0], strings.Join(cycle, " -> "), cycle[0]))
	}

	return issues
}

// validatePackageFilePath accepts relative, slash separated paths that stay inside the bundle directory.
func validatePackageFilePath(packageFile string) error {
	if strings.Contains(packageFile, "\\") {
		return fmt.Errorf("%q must use forward slashes", packageFile)
	}
	if path.IsAbs(packageFile) || filepath.IsAbs(packageFile) || filepath.VolumeName(packageFile) != "" {
		return fmt.Errorf("%q must be a path relative to the bundle directory", packageFile)
	}
	cleaned := path.Clean(packageFile)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("%q points outside the bundle directory", packageFile)
	}
	if cleaned != packageFile {
		return fmt.Errorf("%q is not a clean path, use %q", packageFile, cleaned)
	}
	return nil
}

// resolvePackageFile returns the local path of a package file, refusing files outside root (including through symlinks).
func resolvePackageFile(root string, packageFile string) (string, error) {
	localPath := filepath.Join(root, filepath.FromSlash(packageFile))

	resolved, err := filepath.EvalSymlinks(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("package file %s not found in %s", packageFile, root)
		}
		return "", fmt.Errorf("package file %s is not accessible: %w", packageFile, err)
	}

	relative, err := filepath.Rel(root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("package file %s resolves outside %s", packageFile, root)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("package file %s is not accessible: %w", packageFile, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("package file %s is not a regular file", packageFile)
	}

	return resolved, nil
}

// fileSha256 streams a file through SHA-256 and returns the lowercase hex digest.
func fileSha256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// nvidiaStepsLabel names the steps that use a package file, e.g. `Step "hmc"` or `Steps "erot", "hmc"`.
func nvidiaStepsLabel(manifest *NvidiaBundleManifest, packageFile string) string {
	ids := []string{}
	for _, step := range manifest.Steps {
		if step.PackageFile == packageFile {
			ids = append(ids, fmt.Sprintf("%q", step.Id))
		}
	}
	switch len(ids) {
	case 0:
		return fmt.Sprintf("Bundle file %s", packageFile)
	case 1:
		return "Step " + ids[0]
	default:
		return "Steps " + strings.Join(ids, ", ")
	}
}

func isHttpUrl(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

// sortedMap iterates a string map in key order, so issues come out in a stable order.
func sortedMap(values map[string]string) func(yield func(string, string) bool) {
	return func(yield func(string, string) bool) {
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if !yield(key, values[key]) {
				return
			}
		}
	}
}

// Warnings lists what is valid but likely wrong in a validated manifest. A step without a verifyVersion whose
// packageVersion is the bundleVersion gives the inventory no version in the component's own numbering: it then
// compares the bundle release (e.g. 25.06.1) with what the component reports (e.g. a GPU VBIOS 96.00.89.00.01), and
// offers upgrades as downgrades. Set verifyVersion to the version the component reports once the step is active.
func (m *NvidiaBundleManifest) Warnings() []string {
	var warnings []string
	bundleVersion := strings.TrimSpace(m.BundleVersion)
	for _, step := range m.Steps {
		if strings.TrimSpace(step.VerifyVersion) != "" || bundleVersion == "" {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(step.PackageVersion), bundleVersion) {
			warnings = append(warnings, fmt.Sprintf(
				"Step %q: packageVersion is the bundle version %s and there is no verifyVersion; components are compared with %s instead of their own firmware version. Set verifyVersion to the version the %s components report once the step is active",
				step.Id, bundleVersion, bundleVersion, step.ComponentClass))
		}
	}
	return warnings
}

// SameFirmwarePlatform compares platform names as the inventory does: ignoring case and whitespace, so
// "NVIDIA DGX H100" (NVIDIA's name) is the server "NVIDIA" / "DGXH100" (what its BMC reports). Hyphens and
// suffixes count.
func SameFirmwarePlatform(a string, b string) bool {
	key := func(value string) string {
		return strings.ToLower(strings.Join(strings.Fields(value), ""))
	}
	return key(a) != "" && key(a) == key(b)
}
