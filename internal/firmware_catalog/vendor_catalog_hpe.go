package firmware_catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type StringOrSlice []string

func (sos *StringOrSlice) UnmarshalJSON(data []byte) error {
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*sos = slice
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*sos = []string{single}
		return nil
	}

	return fmt.Errorf("target can be either a string or an array of strings, but was neither")
}

type hpCatalogTemplate struct {
	Date                 string        `json:"date"`
	Description          string        `json:"description"`
	DeviceClass          string        `json:"deviceclass"`
	MinimumActiveVersion string        `json:"minimum_active_version"`
	RebootRequired       string        `json:"reboot_required"`
	Target               StringOrSlice `json:"target"`
	Version              string        `json:"version"`
}

func (vc *VendorCatalog) processHpeCatalog(ctx context.Context) error {
	catalogUrl := ""
	downloadPath := ""
	var downloadUrl *url.URL
	var err error
	if vc.CatalogInfo.VendorUrl != nil {
		downloadUrl, err = url.Parse(*vc.CatalogInfo.VendorUrl)
		if err != nil {
			return fmt.Errorf("invalid catalog URL: %v", err)
		}

		downloadPath = path.Dir(path.Dir(downloadUrl.Path))
		catalogUrl = *vc.CatalogInfo.VendorUrl
	}

	if catalogUrl == "" && vc.VendorLocalCatalogPath == "" {
		return fmt.Errorf("no catalog source provided")
	}

	localPath := vc.VendorLocalCatalogPath

	if localPath == "" {
		// Create a temporary file to download the catalog
		tempFile, err := os.CreateTemp("", "hp_catalog_*.json")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		localPath = tempFile.Name()

		// Download catalog from URL
		err = downloadCatalog(catalogUrl, localPath, vc.VendorToken)
		if err != nil {
			return fmt.Errorf("failed to download catalog: %v", err)
		}
	}

	jsonFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var packages map[string]hpCatalogTemplate
	err = json.Unmarshal(byteValue, &packages)
	if err != nil {
		return fmt.Errorf("failed to parse catalog: %v", err)
	}

	logger.Get().Debug().Msgf("Parsed %d packages from the catalog", len(packages))
	logger.Get().Debug().Msgf("Download path: %s", downloadPath)
	logger.Get().Debug().Msgf("Local path: %s", localPath)
	logger.Get().Debug().Msgf("Vendor systems filter: %v", vc.VendorSystemsFilter)

	for packageKey, packageInfo := range packages {
		// We only check for components that are of type firmware
		if !strings.HasSuffix(packageKey, "fwpkg") {
			logger.Get().Debug().Msgf("Skipping package %s - Type is not Firmware", packageKey)
			continue
		}

		// Skip if device class is empty or null, or if there are no devices
		if packageInfo.DeviceClass == "" || packageInfo.DeviceClass == "null" || packageInfo.Target == nil || len(packageInfo.Target) == 0 {
			logger.Get().Debug().Msgf("Skipping package %s - no DeviceClass, Device or Target", packageKey)
			continue
		}

		includedBinary := false
		supportedDevices := []map[string]interface{}{}
		for _, target := range packageInfo.Target {
			if len(vc.VendorSystemsFilter) == 0 || slices.Contains(vc.VendorSystemsFilter, target) {
				includedBinary = true

				supportedDevices = append(supportedDevices, map[string]interface{}{
					"id":    target, //.Target,
					"model": target, //.DeviceName,
					// "DeviceClass":            packageInfo.DeviceClass,
					// "Target":                 target,
					// "MinimumVersionRequired": packageInfo.MinimumActiveVersion,
				})
			}
		}

		if !includedBinary {
			logger.Get().Debug().Msgf("Skipping package %s - targets %v not included in the vendor systems filter", packageKey, packageInfo.Target)
			continue
		}

		packageDownloadUrl := ""
		if downloadPath != "" {
			downloadUrl.Path = path.Join(downloadPath, packageKey)
			packageDownloadUrl = downloadUrl.String()
		} else {
			packageDownloadUrl = "https://not-supported.local/" + packageKey
		}

		firmwareBinary := sdk.FirmwareBinary{
			ExternalId:             sdk.PtrString(packageKey),
			Name:                   packageKey,
			VendorInfoUrl:          nil,
			VendorDownloadUrl:      packageDownloadUrl,
			CacheDownloadUrl:       nil, //	Will be set after the binary is downloaded
			PackageId:              sdk.PtrString(packageKey),
			PackageVersion:         sdk.PtrString(packageInfo.Version),
			RebootRequired:         packageInfo.RebootRequired == "yes",
			UpdateSeverity:         sdk.FIRMWAREBINARYUPDATESEVERITY_UNKNOWN,
			VendorSupportedDevices: supportedDevices,
			VendorSupportedSystems: []map[string]interface{}{},
			VendorReleaseTimestamp: nil,
			Vendor:                 map[string]any{},
		}

		vc.Binaries = append(vc.Binaries, &firmwareBinary)
	}

	return nil
}

// getHpeFirmwareTargets retrieves the component target UUIDs from the firmware
// inventory of every registered HPE server whose server type appears in
// ServerTypesFilter. The HPE catalog identifies firmware packages by these
// UUIDs (the "target" field on each catalog entry), so they are needed to
// drive catalog filtering when the user has not supplied them via
// --vendor-systems.
func (vc *VendorCatalog) getHpeFirmwareTargets(ctx context.Context) ([]string, error) {
	client := api.GetApiClient(ctx)

	serverTypes, httpRes, err := client.ServerTypeAPI.GetServerTypes(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	seen := map[string]struct{}{}
	targets := []string{}

	for _, serverType := range serverTypes.Data {
		if !slices.Contains(vc.ServerTypesFilter, serverType.Label) {
			continue
		}

		servers, httpRes, err := client.ServerAPI.GetServers(ctx).
			FilterServerTypeId([]string{fmt.Sprintf("%d", int(serverType.Id))}).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, err
		}

		for _, server := range servers.Data {
			if server.Vendor == nil || !isHpeVendor(*server.Vendor) {
				continue
			}

			logger.Get().Debug().Msgf("Retrieving firmware inventory for HPE server %d (type %s)", int(server.ServerId), serverType.Label)

			inventory, httpRes, err := client.ServerFirmwareAPI.GetServerFirmwareInventory(ctx, server.ServerId).Execute()
			if err := response_inspector.InspectResponse(httpRes, err); err != nil {
				logger.Get().Warn().Err(err).Msgf("Failed to retrieve firmware inventory for server %d - skipping", int(server.ServerId))
				continue
			}

			for _, entry := range inventory {
				for _, uuid := range extractHpeFirmwareTargets(entry) {
					if _, exists := seen[uuid]; !exists {
						seen[uuid] = struct{}{}
						targets = append(targets, uuid)
					}
				}
			}
		}
	}

	return targets, nil
}

// extractHpeFirmwareTargets pulls component target UUIDs out of a single
// Redfish SoftwareInventory entry returned by iLO. HPE places these under
// Oem.Hpe.Targets; older iLO firmware uses Oem.Hp.Targets.
func extractHpeFirmwareTargets(entry map[string]interface{}) []string {
	oem, ok := entry["Oem"].(map[string]interface{})
	if !ok {
		return nil
	}
	for _, key := range []string{"Hpe", "Hp"} {
		sub, ok := oem[key].(map[string]interface{})
		if !ok {
			continue
		}
		raw, ok := sub["Targets"].([]interface{})
		if !ok {
			continue
		}
		result := []string{}
		for _, t := range raw {
			if s, ok := t.(string); ok && s != "" {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// isHpeVendor reports whether a server's reported Vendor field corresponds to
// the HPE firmware catalog vendor. The MetalSoft API may report either "HP" or
// "HPE" depending on the source, while the catalog vendor enum value is "hp".
func isHpeVendor(serverVendor string) bool {
	v := strings.ToLower(serverVendor)
	return v == "hp" || v == "hpe"
}
