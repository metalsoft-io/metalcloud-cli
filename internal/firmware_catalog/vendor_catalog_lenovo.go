package firmware_catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const (
	lenovoContentServiceUrl = "https://support.lenovo.com/services/ContentService/"

	lenovoSoftwareUpdateComponentXcc  = "XCC"
	lenovoSoftwareUpdateComponentUefi = "UEFI"
	lenovoSoftwareUpdateComponentLxpm = "LXPM"

	lenovoSoftwareUpdateTypeFix        = "Fix"
	lenovoSoftwareUpdateTypeInstallXML = "InstallXML"
	lenovoSoftwareUpdateTypeReadMe     = "ReadMe"
)

type lenovoSoftwareUpdateFile struct {
	Type        string `json:"Type"`
	Description string `json:"Description"`
	URL         string `json:"URL"`
	FileHash    string `json:"FileHash"`
}

type lenovoSoftwareUpdate struct {
	FixID            string                     `json:"FixID"`
	ComponentID      string                     `json:"ComponentID"`
	Files            []lenovoSoftwareUpdateFile `json:"Files"`
	RequisitesFixIDs []string                   `json:"RequisitesFixIDs"`
	Version          string
	UpdateKey        string
}

type lenovoCatalog struct {
	ResultCode     *int                   `json:"ResultCode"`
	Message        *string                `json:"Message"`
	Data           []*lenovoSoftwareUpdate `json:"Data"`
	FixIdsNotFound []string               `json:"FixIdsNotFound"`
}

func (vc *VendorCatalog) processLenovoCatalog(ctx context.Context) error {
	// If VendorSystemsFilterEx is empty but VendorSystemsFilter has entries,
	// populate from the basic filter with empty serial numbers as fallback.
	// This handles the case where --vendor-systems is provided without --server-types,
	// or when no matching servers are found in the MetalSoft inventory.
	if len(vc.VendorSystemsFilterEx) == 0 && len(vc.VendorSystemsFilter) > 0 {
		vc.VendorSystemsFilterEx = make(map[string]string, len(vc.VendorSystemsFilter))
		for _, system := range vc.VendorSystemsFilter {
			vc.VendorSystemsFilterEx[system] = ""
		}
	}

	if len(vc.VendorSystemsFilterEx) == 0 {
		return fmt.Errorf("no vendor systems filter provided - use --vendor-systems or --server-types flag")
	}

	for serverType, serverSerialNumber := range vc.VendorSystemsFilterEx {
		catalog, err := vc.readLenovoCatalog(serverType, serverSerialNumber)
		if err != nil {
			return err
		}

		if catalog.ResultCode != nil && *catalog.ResultCode != 0 {
			message := "unknown error"
			if catalog.Message != nil {
				message = *catalog.Message
			}
			return fmt.Errorf("lenovo catalog API returned error for machine type %s: code=%d, message=%s", serverType, *catalog.ResultCode, message)
		}

		logger.Get().Debug().Msgf("Lenovo catalog for machine type %s: %d software updates found", serverType, len(catalog.Data))

		skippedComponent := 0
		skippedNoFixUrl := 0
		for _, softwareUpdate := range catalog.Data {
			if softwareUpdate.ComponentID != lenovoSoftwareUpdateComponentXcc &&
				softwareUpdate.ComponentID != lenovoSoftwareUpdateComponentUefi &&
				softwareUpdate.ComponentID != lenovoSoftwareUpdateComponentLxpm {
				logger.Get().Debug().Msgf("Skipping software update %s with unsupported component ID: %s", softwareUpdate.FixID, softwareUpdate.ComponentID)
				skippedComponent++
				continue
			}

			downloadUrl := ""
			description := ""
			infoUrl := ""
			for _, file := range softwareUpdate.Files {
				if file.Type == lenovoSoftwareUpdateTypeFix {
					downloadUrl = file.URL
					continue
				}
				if file.Type == lenovoSoftwareUpdateTypeInstallXML {
					description = file.Description
					continue
				}
				if file.Type == lenovoSoftwareUpdateTypeReadMe {
					infoUrl = file.URL
					continue
				}
			}

			if downloadUrl == "" {
				logger.Get().Warn().Msgf("no firmware fix was found for software update %s", softwareUpdate.FixID)
				skippedNoFixUrl++
				continue
			}

			if description == "" {
				description = softwareUpdate.FixID
			}

			componentVendorConfiguration := map[string]any{
				"requires": softwareUpdate.RequisitesFixIDs,
			}

			supportedDevices := []map[string]interface{}{
				{
					"type": softwareUpdate.UpdateKey,
				},
			}

			supportedSystems := []map[string]interface{}{
				{
					"machineType":  serverType,
					"serialNumber": serverSerialNumber,
				},
			}

			firmwareBinary := sdk.FirmwareBinary{
				ExternalId:             sdk.PtrString(softwareUpdate.FixID),
				Name:                   description,
				VendorInfoUrl:          &infoUrl,
				VendorDownloadUrl:      downloadUrl,
				CacheDownloadUrl:       nil, //	Will be set after the binary is downloaded
				PackageId:              sdk.PtrString(softwareUpdate.FixID),
				PackageVersion:         sdk.PtrString(softwareUpdate.Version),
				RebootRequired:         true,
				UpdateSeverity:         sdk.FIRMWAREBINARYUPDATESEVERITY_UNKNOWN,
				VendorSupportedDevices: supportedDevices,
				VendorSupportedSystems: supportedSystems,
				VendorReleaseTimestamp: nil,
				Vendor:                 componentVendorConfiguration,
			}

			vc.Binaries = append(vc.Binaries, &firmwareBinary)
		}

		logger.Get().Info().Msgf("Lenovo catalog for machine type %s: %d binaries added, %d skipped (unsupported component), %d skipped (no download URL)",
			serverType, len(catalog.Data)-skippedComponent-skippedNoFixUrl, skippedComponent, skippedNoFixUrl)
	}

	return nil
}

// Search the lenovo support site for the server firmware update information. A JSON response is returned and is saved in the local catalog path folder from the raw config file.
func (vc *VendorCatalog) readLenovoCatalog(machineType string, serialNumber string) (*lenovoCatalog, error) {
	if machineType == "" {
		return nil, fmt.Errorf("machine type must be specified when searching for a lenovo catalog")
	}

	catalog := lenovoCatalog{}

	var catalogFileName string
	if serialNumber != "" {
		catalogFileName = fmt.Sprintf("lenovo_%s_%s.json", machineType, serialNumber)
	} else {
		catalogFileName = fmt.Sprintf("lenovo_%s.json", machineType)
	}
	localCatalogPath := filepath.Join(vc.VendorLocalCatalogPath, catalogFileName)

	fileExists := false
	info, err := os.Stat(localCatalogPath)
	if !os.IsNotExist(err) {
		fileExists = !info.IsDir()
	}

	if fileExists {
		logger.Get().Info().Msgf("Reading local Lenovo catalog %s", localCatalogPath)

		content, err := os.ReadFile(localCatalogPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(content, &catalog)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Get().Info().Msgf("Download Lenovo catalog for %s", machineType)

		response, err := downloadLenovoFirmwareUpdates(machineType, serialNumber)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(response), &catalog)
		if err != nil {
			return nil, err
		}
	}

	return &catalog, nil
}

func downloadLenovoFirmwareUpdates(machineType string, serialNumber string) (string, error) {
	targetInfos := map[string]string{
		"MachineType": machineType,
	}
	if serialNumber != "" {
		targetInfos["SerialNumber"] = serialNumber
	}

	searchParams := map[string]interface{}{
		"Category":            "",
		"FixIds":              "",
		"IsIncludeData":       "true",
		"IsIncludeMetaData":   "true",
		"IsIncludeRequisites": "true",
		"IsLatest":            "true",
		"QueryType":           "SUP",
		"SelectSupersedes":    "3",
		"SubmitterName":       "",
		"SubmitterVersion":    "",
		"TargetInfos":         []map[string]string{targetInfos},
		"XmlUpdateType":       "",
	}

	jsonParams, err := json.Marshal(searchParams)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, lenovoContentServiceUrl+"SearchDrivers", bytes.NewBuffer(jsonParams))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lenovo catalog API returned HTTP %d for machine type %s: %s", resp.StatusCode, machineType, string(responseBody))
	}

	return string(responseBody), nil
}
