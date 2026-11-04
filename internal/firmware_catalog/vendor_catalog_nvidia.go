package firmware_catalog

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// processNvidiaCatalog builds the catalog of one NVIDIA HGX/DGX bundle version from its manifest: one bundle binary
// carrying every step, plus one binary for each other package file the steps use, so that every file is uploaded
// and cached next to the bundle. See NVIDIA_BUNDLE_IMPORT.md.
func (vc *VendorCatalog) processNvidiaCatalog(ctx context.Context) error {
	if vc.VendorLocalCatalogPath == "" {
		return fmt.Errorf("an NVIDIA catalog is built from a bundle manifest: set --vendor-local-catalog-path to the manifest file")
	}
	if vc.VendorLocalBinariesPath == "" {
		return fmt.Errorf("an NVIDIA catalog needs the folder with the package files: set --vendor-local-binaries-path")
	}
	if vc.CatalogInfo.VendorUrl != nil && *vc.CatalogInfo.VendorUrl != "" {
		return fmt.Errorf("NVIDIA publishes no online catalog: use --vendor-local-catalog-path with a bundle manifest instead of --vendor-url")
	}
	if vc.DownloadBinaries {
		return fmt.Errorf("NVIDIA bundles are not downloaded by the CLI: put the package files in --vendor-local-binaries-path instead of using --download-binaries")
	}

	manifest, err := LoadNvidiaManifest(vc.VendorLocalCatalogPath)
	if err != nil {
		return err
	}
	if err := VerifyNvidiaPackageFiles(manifest, vc.VendorLocalBinariesPath); err != nil {
		return err
	}
	for _, warning := range manifest.Warnings() {
		logger.Get().Warn().Msg(warning)
	}

	var repoUrl *url.URL
	if vc.RepoBaseUrl != "" {
		repoUrl, err = url.Parse(vc.RepoBaseUrl)
		if err != nil {
			return fmt.Errorf("unable to parse repo base URL: %v", err)
		}
	}
	// The API requires a download URL for every binary
	for _, packageFile := range manifest.PackageFiles() {
		if manifest.DownloadUrls[packageFile] == "" && repoUrl == nil {
			return fmt.Errorf("%s has no downloadUrls entry in the manifest: set --repo-base-url, the repository the package files are served from", packageFile)
		}
	}

	releaseDate, err := time.Parse(NvidiaManifestReleaseDateLayout, manifest.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid releaseDate %q: %w", manifest.ReleaseDate, err)
	}

	vc.CatalogInfo.VendorId = sdk.PtrString(nvidiaCatalogVendorId(manifest))
	vc.CatalogInfo.VendorReleaseTimestamp = &releaseDate
	vc.CatalogInfo.VendorServerTypesSupported = slices.Clone(manifest.TargetPlatforms)
	vc.CatalogInfo.VendorConfiguration = map[string]interface{}{
		"schemaVersion": manifest.SchemaVersion,
		"bundleName":    manifest.BundleName,
		"bundleVersion": manifest.BundleVersion,
	}
	if manifest.ReleaseNotesUrl != "" {
		vc.CatalogInfo.VendorConfiguration["releaseNotesUrl"] = manifest.ReleaseNotesUrl
	}

	vc.Binaries = nvidiaBinaries(manifest, repoUrl, releaseDate)
	return nil
}

// nvidiaCatalogVendorId identifies the catalog of a bundle version, so that importing it again updates it.
func nvidiaCatalogVendorId(manifest *NvidiaBundleManifest) string {
	return manifest.BundleName + "@" + manifest.BundleVersion
}

// nvidiaBinaries returns the bundle binary, from bundleFile, followed by one binary for each other package file.
func nvidiaBinaries(manifest *NvidiaBundleManifest, repoUrl *url.URL, releaseDate time.Time) []*sdk.FirmwareBinary {
	severity := manifest.UpdateSeverity
	if severity == "" {
		severity = UpdateSeverityUnknown
	}
	downloadUrl := func(packageFile string) string {
		if manifestUrl := manifest.DownloadUrls[packageFile]; manifestUrl != "" {
			return manifestUrl
		}
		return repoUrl.JoinPath(packageFile).String()
	}
	platforms := []sdk.ServerFirmwareBinaryVendorSupportedSystems{}
	for _, platform := range manifest.TargetPlatforms {
		platforms = append(platforms, sdk.ServerFirmwareBinaryVendorSupportedSystems{
			AdditionalProperties: map[string]interface{}{"id": platform},
		})
	}

	newBinary := func(packageFile string, name string, steps []NvidiaBundleStep) *sdk.FirmwareBinary {
		binary := &sdk.FirmwareBinary{
			ExternalId:             sdk.PtrString(packageFile),
			Name:                   name,
			PackageId:              sdk.PtrString(manifest.BundleName),
			PackageVersion:         sdk.PtrString(manifest.BundleVersion),
			VendorDownloadUrl:      downloadUrl(packageFile),
			RebootRequired:         nvidiaRebootRequired(steps),
			UpdateSeverity:         severity,
			VendorSupportedDevices: nvidiaSupportedDevices(steps),
			VendorSupportedSystems: platforms,
			VendorReleaseTimestamp: &releaseDate,
			Sha256:                 sdk.PtrString(manifest.PackageFileSha256(packageFile)),
		}
		if manifest.ReleaseNotesUrl != "" {
			binary.VendorInfoUrl = sdk.PtrString(manifest.ReleaseNotesUrl)
		}
		return binary
	}

	bundle := newBinary(
		manifest.BundleFile,
		fmt.Sprintf("NVIDIA %s %s", manifest.BundleName, manifest.BundleVersion),
		manifest.Steps,
	)
	bundle.IsBundle = sdk.PtrBool(true)
	bundle.BundleSteps = nvidiaBundleSteps(manifest.Steps)
	bundle.Prerequisites = nvidiaPrerequisites(manifest.Prerequisites)
	binaries := []*sdk.FirmwareBinary{bundle}

	for _, packageFile := range manifest.PackageFiles() {
		if packageFile == manifest.BundleFile {
			continue
		}
		steps := slices.DeleteFunc(slices.Clone(manifest.Steps), func(step NvidiaBundleStep) bool {
			return step.PackageFile != packageFile
		})
		stepIds := []string{}
		for _, step := range steps {
			stepIds = append(stepIds, step.Id)
		}
		binary := newBinary(
			packageFile,
			fmt.Sprintf("NVIDIA %s %s (%s)", manifest.BundleName, manifest.BundleVersion, strings.Join(stepIds, ", ")),
			steps,
		)
		binary.IsBundle = sdk.PtrBool(false)
		binaries = append(binaries, binary)
	}

	return binaries
}

func nvidiaBundleSteps(steps []NvidiaBundleStep) []sdk.FirmwareBundleStep {
	result := []sdk.FirmwareBundleStep{}
	for _, step := range steps {
		bundleStep := sdk.FirmwareBundleStep{
			Id:             step.Id,
			ComponentClass: sdk.ServerComponentClass(step.ComponentClass),
			TargetPattern:  step.TargetPattern,
			PackageFile:    step.PackageFile,
			PackageVersion: step.PackageVersion,
			Sha256:         strings.ToLower(step.Sha256),
			Activation:     sdk.FirmwareBundleStepActivation(step.Activation),
			// null when the step depends on no other step
			DependsOnStep: *sdk.NewNullableString(step.DependsOnStep),
		}
		if step.VerifyVersion != "" {
			bundleStep.VerifyVersion = sdk.PtrString(step.VerifyVersion)
		}
		result = append(result, bundleStep)
	}
	return result
}

func nvidiaPrerequisites(prerequisites []NvidiaPrerequisite) []sdk.FirmwarePrerequisite {
	if len(prerequisites) == 0 {
		return nil
	}
	result := []sdk.FirmwarePrerequisite{}
	for _, prerequisite := range prerequisites {
		entry := sdk.FirmwarePrerequisite{
			ComponentClass: sdk.ServerComponentClass(prerequisite.ComponentClass),
			MinimumVersion: prerequisite.MinimumVersion,
		}
		if prerequisite.Platform != "" {
			entry.Platform = sdk.PtrString(prerequisite.Platform)
		}
		if prerequisite.ComponentModel != "" {
			entry.ComponentModel = sdk.PtrString(prerequisite.ComponentModel)
		}
		if prerequisite.ComponentName != "" {
			entry.ComponentName = sdk.PtrString(prerequisite.ComponentName)
		}
		result = append(result, entry)
	}
	return result
}

// nvidiaSupportedDevices lists the component class and target pattern of each step, for operators reading the
// binary; the inventory matches bundles on their steps.
func nvidiaSupportedDevices(steps []NvidiaBundleStep) []sdk.VendorSupportedDevice {
	devices := []sdk.VendorSupportedDevice{}
	for _, step := range steps {
		devices = append(devices, sdk.VendorSupportedDevice{
			AdditionalProperties: map[string]interface{}{"id": step.ComponentClass, "model": step.TargetPattern},
		})
	}
	return devices
}

// nvidiaRebootRequired is true when a step needs the host restarted or the chassis power cycled.
func nvidiaRebootRequired(steps []NvidiaBundleStep) bool {
	return slices.ContainsFunc(steps, func(step NvidiaBundleStep) bool {
		return step.Activation == NvidiaActivationHostReset || step.Activation == NvidiaActivationChassisPowerCycle
	})
}

// findNvidiaCatalog returns the catalog an earlier import of the same bundle version created, or nil.
func (vc *VendorCatalog) findNvidiaCatalog(ctx context.Context) (*sdk.FirmwareCatalog, error) {
	if vc.CatalogInfo.VendorId == nil {
		return nil, nil
	}
	client := api.GetApiClient(ctx)
	catalogs, httpRes, err := client.FirmwareCatalogAPI.GetFirmwareCatalogs(ctx).
		FilterVendor([]string{"$eq:" + VendorNvidia}).
		FilterVendorId([]string{"$eq:" + *vc.CatalogInfo.VendorId}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	for i := range catalogs.Data {
		if catalogs.Data[i].VendorId != nil && *catalogs.Data[i].VendorId == *vc.CatalogInfo.VendorId {
			return &catalogs.Data[i], nil
		}
	}
	return nil, nil
}

// checkExistingNvidiaBinaries returns the binaries of the catalog this import updates, by external id (the package
// file). A package file whose checksum changed is refused, before anything is uploaded or changed, unless Force is
// set: the same bundle version should not ship another file.
func (vc *VendorCatalog) checkExistingNvidiaBinaries(ctx context.Context, catalogId int64) (map[string]*sdk.FirmwareBinary, error) {
	client := api.GetApiClient(ctx)
	existing := map[string]*sdk.FirmwareBinary{}
	changed := []string{}
	for _, binary := range vc.Binaries {
		binaries, httpRes, err := client.FirmwareBinaryAPI.GetFirmwareBinaries(ctx).
			FilterCatalogId([]string{fmt.Sprintf("$eq:%d", catalogId)}).
			FilterExternalId([]string{"$eq:" + *binary.ExternalId}).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, err
		}
		for i := range binaries.Data {
			found := &binaries.Data[i]
			if found.ExternalId == nil || *found.ExternalId != *binary.ExternalId {
				continue
			}
			existing[*binary.ExternalId] = found
			if found.Sha256 != nil && binary.Sha256 != nil && !strings.EqualFold(*found.Sha256, *binary.Sha256) {
				changed = append(changed, fmt.Sprintf("%s (binary %d has sha256 %s, the file has %s)",
					*binary.ExternalId, int(found.Id), *found.Sha256, *binary.Sha256))
			}
			break
		}
	}
	if len(changed) > 0 {
		if !vc.Force {
			return nil, fmt.Errorf("catalog %d already has these package files with another checksum: %s. "+
				"A bundle version should not change; use --force to replace them", int(catalogId), strings.Join(changed, "; "))
		}
		logger.Get().Warn().Msgf("Replacing package files with another checksum (--force): %s", strings.Join(changed, "; "))
	}
	return existing, nil
}
