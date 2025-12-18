package firmware_catalog

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// supermicroFirmwareComponent represents a firmware component in the Supermicro catalog
type supermicroFirmwareComponent struct {
	ComponentType string `json:"componentType"`
	Model         string `json:"model"`
	Version       string `json:"version"`
	ReleaseDate   string `json:"releaseDate"`
	FileName      string `json:"fileName"`
	FileSize      int64  `json:"fileSize"`
	Description   string `json:"description"`
	DownloadURL   string `json:"downloadURL"`
	MD5           string `json:"md5"`
	SHA256        string `json:"sha256"`
}

// supermicroCatalog represents the Supermicro firmware catalog structure
type supermicroCatalog struct {
	CatalogVersion string                        `json:"catalogVersion"`
	ReleaseDate    string                        `json:"releaseDate"`
	Components     []supermicroFirmwareComponent `json:"components"`
}

func (vc *VendorCatalog) processSupermicroCatalog(ctx context.Context) error {
	// Supermicro binaries must be provided locally
	if vc.VendorLocalBinariesPath == "" {
		return fmt.Errorf("vendor-local-binaries-path is required for Supermicro (binaries cannot be downloaded from vendor)")
	}

	var catalog *supermicroCatalog
	var err error

	// Try to read catalog if provided, otherwise scan the binaries folder
	if vc.VendorLocalCatalogPath != "" {
		logger.Get().Info().Msgf("Processing Supermicro catalog from local file: %s", vc.VendorLocalCatalogPath)
		catalog, err = vc.readSupermicroCatalog()
		if err != nil {
			return fmt.Errorf("failed to read Supermicro catalog: %v", err)
		}
	} else {
		logger.Get().Info().Msgf("No catalog file provided, scanning binaries folder: %s", vc.VendorLocalBinariesPath)
		catalog, err = vc.scanSupermicroBinariesFolder()
		if err != nil {
			return fmt.Errorf("failed to scan Supermicro binaries folder: %v", err)
		}
	}

	if catalog == nil || len(catalog.Components) == 0 {
		return fmt.Errorf("no components found in Supermicro catalog")
	}

	logger.Get().Info().Msgf("Found %d components in Supermicro catalog", len(catalog.Components))

	// For Supermicro, we don't filter by vendor systems because:
	// - The firmware applies to motherboard models (e.g., H13SSF-1D16)
	// - The vendor SKU from servers is different (e.g., ASG-1115S-NE316R)
	// - User explicitly provides the binaries they want to upload
	logger.Get().Debug().Msgf("VendorSystemsFilter: %v (will not be used for Supermicro)", vc.VendorSystemsFilter)

	// Process each component
	for _, component := range catalog.Components {
		componentId := fmt.Sprintf("%s-%s", component.ComponentType, component.Model)
		logger.Get().Debug().Msgf("Processing Supermicro component: %s version %s", componentId, component.Version)

		// Supermicro binaries must be provided locally, not downloaded from vendor site
		if vc.VendorLocalBinariesPath == "" {
			return fmt.Errorf("vendor-local-binaries-path is required for Supermicro firmware binaries")
		}

		// Verify the binary file exists in the local path
		binaryLocalPath := path.Join(vc.VendorLocalBinariesPath, component.FileName)
		if _, err := os.Stat(binaryLocalPath); os.IsNotExist(err) {
			logger.Get().Warn().Msgf("Binary file not found: %s - skipping component %s", binaryLocalPath, componentId)
			continue
		}

		// Supermicro doesn't provide direct download URLs, use placeholder
		// The actual cache download URL will be set during upload process
		binaryDownloadUrl := "NotAvailableForSupermicro"

		// Create supported devices list with id and model
		// id is used as the unique identifier for the device
		supportedDevices := []map[string]interface{}{
			{
				"id":    component.Model,
				"model": component.Model,
			},
		}

		// Prepare vendor configuration
		vendorConfiguration := map[string]interface{}{
			"componentType": component.ComponentType,
			"model":         component.Model,
		}

		if component.MD5 != "" {
			vendorConfiguration["md5"] = component.MD5
		}
		if component.SHA256 != "" {
			vendorConfiguration["sha256"] = component.SHA256
		}

		// Parse release date
		var timestamp time.Time
		if component.ReleaseDate != "" {
			timestamp, _ = time.Parse("2006-01-02", component.ReleaseDate)
		}

		// Extract the .bin file from the .zip archive
		binFileName, err := vc.extractBinFileFromZip(binaryLocalPath, component.FileName)
		if err != nil {
			logger.Get().Warn().Msgf("Failed to extract .bin file from %s: %v - skipping", component.FileName, err)
			continue
		}

		// Create firmware binary using the extracted .bin filename
		binary := &sdk.FirmwareBinary{
			ExternalId:             sdk.PtrString(binFileName),
			Name:                   fmt.Sprintf("%s - %s", component.ComponentType, component.Model),
			VendorInfoUrl:          nil,
			VendorDownloadUrl:      binaryDownloadUrl,
			CacheDownloadUrl:       nil,
			PackageId:              sdk.PtrString(component.FileName),
			PackageVersion:         sdk.PtrString(component.Version),
			RebootRequired:         false,
			UpdateSeverity:         sdk.FIRMWAREBINARYUPDATESEVERITY_UNKNOWN,
			VendorSupportedDevices: supportedDevices,
			VendorSupportedSystems: []map[string]interface{}{},
			VendorReleaseTimestamp: sdk.PtrString(timestamp.Format(time.RFC3339)),
			Vendor:                 vendorConfiguration,
		}

		vc.Binaries = append(vc.Binaries, binary)
	}

	vc.CatalogInfo.VendorId = sdk.PtrString(catalog.CatalogVersion)
	vc.CatalogInfo.VendorReleaseTimestamp = sdk.PtrString(catalog.ReleaseDate)

	logger.Get().Info().Msgf("Processed %d Supermicro firmware binaries", len(vc.Binaries))

	return nil
}

func (vc *VendorCatalog) readSupermicroCatalog() (*supermicroCatalog, error) {
	// Supermicro only supports local catalog files
	logger.Get().Info().Msgf("Reading Supermicro catalog from local file: %s", vc.VendorLocalCatalogPath)
	catalogData, err := os.ReadFile(vc.VendorLocalCatalogPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local catalog file: %v", err)
	}

	var catalog supermicroCatalog
	if err := json.Unmarshal(catalogData, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse catalog JSON: %v", err)
	}

	return &catalog, nil
}

func (vc *VendorCatalog) scanSupermicroBinariesFolder() (*supermicroCatalog, error) {
	logger.Get().Info().Msgf("Scanning Supermicro binaries folder: %s", vc.VendorLocalBinariesPath)

	if _, err := os.Stat(vc.VendorLocalBinariesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("binaries folder does not exist: %s", vc.VendorLocalBinariesPath)
	}

	var components []supermicroFirmwareComponent

	err := filepath.Walk(vc.VendorLocalBinariesPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		lowerName := strings.ToLower(fileName)

		if !strings.HasSuffix(lowerName, ".zip") {
			logger.Get().Debug().Msgf("Skipping non-zip file: %s", fileName)
			return nil
		}

		// Expected patterns:
		// BIOS_<model>_<date>_<version>_<type>.zip
		// BMC_<model>_<date>_<version>_<type>.zip
		var componentType, model, version string

		if strings.HasPrefix(fileName, "BIOS_") {
			componentType = "BIOS"
			// Extract model: BIOS_H13SSF-1D16_20251009_3.7a_STDsp.zip -> H13SSF-1D16
			parts := strings.Split(fileName, "_")
			if len(parts) >= 4 {
				model = parts[1]
				version = parts[3]
			}
		} else if strings.HasPrefix(fileName, "BMC_") {
			componentType = "BMC"
			// Extract model: BMC_H13AST2600-ROT20-E401MS_20250919_01.05.11_STDsp.zip -> H13AST2600-ROT20-E401MS
			parts := strings.Split(fileName, "_")
			if len(parts) >= 4 {
				model = parts[1]
				version = parts[3]
			}
		} else {
			// Unknown firmware type, try to extract from filename
			logger.Get().Debug().Msgf("Unknown firmware type for file: %s", fileName)
			return nil
		}

		if model == "" || version == "" {
			logger.Get().Warn().Msgf("Could not parse model/version from filename: %s", fileName)
			return nil
		}

		// Get relative path from binaries folder
		relPath, err := filepath.Rel(vc.VendorLocalBinariesPath, filePath)
		if err != nil {
			relPath = fileName
		}

		logger.Get().Info().Msgf("Found %s firmware: %s version %s", componentType, model, version)

		component := supermicroFirmwareComponent{
			ComponentType: componentType,
			Model:         model,
			Version:       version,
			FileName:      relPath,
			FileSize:      info.Size(),
			Description:   fmt.Sprintf("%s firmware for %s", componentType, model),
			ReleaseDate:   info.ModTime().Format("2006-01-02"),
		}

		components = append(components, component)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning binaries folder: %v", err)
	}

	if len(components) == 0 {
		return nil, fmt.Errorf("no firmware binaries found in folder: %s", vc.VendorLocalBinariesPath)
	}

	catalog := &supermicroCatalog{
		CatalogVersion: "auto-generated",
		ReleaseDate:    time.Now().Format(time.RFC3339),
		Components:     components,
	}

	return catalog, nil
}

// extractBinFileFromZip extracts the .bin file from a Supermicro .zip archive
// and returns the path to the extracted .bin file
func (vc *VendorCatalog) extractBinFileFromZip(zipPath string, zipFileName string) (string, error) {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %v", err)
	}
	defer zipReader.Close()

	var binFile *zip.File
	for _, file := range zipReader.File {
		if strings.HasSuffix(strings.ToLower(file.Name), ".bin") {
			binFile = file
			break
		}
	}

	if binFile == nil {
		return "", fmt.Errorf("no .bin file found in zip archive")
	}

	extractDir := filepath.Join(vc.VendorLocalBinariesPath, ".extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create extraction directory: %v", err)
	}

	binFilePath := filepath.Join(extractDir, binFile.Name)

	if _, err := os.Stat(binFilePath); err == nil {
		logger.Get().Debug().Msgf("Using cached extracted file: %s", binFilePath)
		return binFile.Name, nil
	}

	logger.Get().Info().Msgf("Extracting %s from %s", binFile.Name, zipFileName)

	srcFile, err := binFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file in zip: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(binFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to extract file: %v", err)
	}

	logger.Get().Info().Msgf("Extracted %s (%d bytes)", binFile.Name, binFile.UncompressedSize64)

	return binFile.Name, nil
}
