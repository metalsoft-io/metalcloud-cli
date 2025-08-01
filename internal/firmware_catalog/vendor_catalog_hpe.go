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
		tempFile, err := os.CreateTemp("", "dell_catalog_*.xml")
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
	json.Unmarshal(byteValue, &packages)

	for packageKey, packageInfo := range packages {
		// We only check for components that are of type firmware
		if !strings.HasSuffix(packageKey, "fwpkg") {
			continue
		}

		if packageInfo.DeviceClass == "" || packageInfo.DeviceClass == "null" || packageInfo.Target == nil || len(packageInfo.Target) == 0 {
			continue
		}

		includedBinary := false
		supportedDevices := []map[string]interface{}{}
		for _, target := range packageInfo.Target {
			if len(vc.VendorSystemsFilter) == 0 || slices.Contains(vc.VendorSystemsFilter, target) {
				includedBinary = true

				supportedDevices = append(supportedDevices, map[string]interface{}{
					"DeviceClass":            packageInfo.DeviceClass,
					"Target":                 target,
					"MinimumVersionRequired": packageInfo.MinimumActiveVersion,
				})
			}
		}

		if !includedBinary {
			continue
		}

		downloadUrl.Path = path.Join(downloadPath, packageKey)
		packageDownloadUrl := downloadUrl.String()

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
