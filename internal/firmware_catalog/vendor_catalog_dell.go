package firmware_catalog

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const (
	componentTypeFirmware string = "FRMW"
)

type dellName struct {
	Display string `xml:"Display"`
}

type dellCriticality struct {
	Value string `xml:"value,attr"`
}

type dellImportantInfo struct {
	URL string `xml:"URL,attr"`
}

type dellCategory struct {
	Display string `xml:"Display"`
}

type dellComponentType struct {
	Value string `xml:"value,attr"`
}

type dellDescription struct {
	Display string `xml:"Display"`
}

type dellDevice struct {
	ComponentID string `xml:"componentID,attr"`
	Display     string `xml:"Display"`
	Prefix      string `xml:"prefix,attr"`
}

type dellModel struct {
	SystemID     string `xml:"systemID,attr"`
	SystemIDType string `xml:"systemIDType,attr"`
	Display      string `xml:"Display"`
}

type dellBrand struct {
	Models  []dellModel `xml:"Model"`
	Display string      `xml:"Display"`
	Prefix  string      `xml:"prefix,attr"`
}

type dellSupportedDevices struct {
	Devices []dellDevice `xml:"Device"`
}

type dellSupportedSystems struct {
	Brands []dellBrand `xml:"Brand"`
}

type dellSoftwareComponent struct {
	XMLName          xml.Name             `xml:"SoftwareComponent"`
	DateTime         string               `xml:"dateTime,attr"`
	Path             string               `xml:"path,attr"`
	PackageID        string               `xml:"packageID,attr"`
	VendorVersion    string               `xml:"vendorVersion,attr"`
	RebootRequired   string               `xml:"rebootRequired,attr"`
	UpdateSeverity   dellCriticality      `xml:"Criticality"`
	ReleaseDate      string               `xml:"releaseDate,attr"`
	Size             string               `xml:"size,attr"`
	HashMD5          string               `xml:"hashMD5,attr"`
	Category         dellCategory         `xml:"Category"`
	ReleaseID        string               `xml:"releaseID,attr"`
	PackageType      string               `xml:"packageType,attr"`
	ImportantInfo    dellImportantInfo    `xml:"ImportantInfo"`
	ComponentType    dellComponentType    `xml:"ComponentType"`
	Name             dellName             `xml:"Name"`
	Description      dellDescription      `xml:"Description"`
	SupportedDevices dellSupportedDevices `xml:"SupportedDevices"`
	SupportedSystems dellSupportedSystems `xml:"SupportedSystems"`
}

type dellSoftwareBundle struct {
	XMLName          xml.Name             `xml:"SoftwareBundle"`
	DateTime         string               `xml:"dateTime,attr"`
	Path             string               `xml:"path,attr"`
	RebootRequired   string               `xml:"rebootRequired,attr"`
	Name             dellName             `xml:"Name"`
	SupportedDevices dellSupportedDevices `xml:"SupportedDevices"`
	SupportedSystems dellSupportedSystems `xml:"SupportedSystems"`
}

type dellManifest struct {
	BaseLocation                string                  `xml:"baseLocation,attr"`
	BaseLocationAccessProtocols string                  `xml:"baseLocationAccessProtocols,attr"`
	DateTime                    string                  `xml:"dateTime,attr"`
	Identifier                  string                  `xml:"identifier,attr"`
	ReleaseID                   string                  `xml:"releaseID,attr"`
	Version                     string                  `xml:"version,attr"`
	PredecessorID               string                  `xml:"predecessorID,attr"`
	Software                    []dellSoftwareBundle    `xml:"SoftwareBundle"`
	Components                  []dellSoftwareComponent `xml:"SoftwareComponent"`
}

func (vc *VendorCatalog) processDellCatalog(ctx context.Context) error {
	catalogUrl := ""
	if vc.CatalogInfo.VendorUrl != nil {
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
		err = downloadGzipCatalog(catalogUrl, localPath)
		if err != nil {
			return fmt.Errorf("failed to download catalog: %v", err)
		}
	}

	var manifest dellManifest

	err := readXmlDocument(localPath, &manifest)
	if err != nil {
		return fmt.Errorf("error reading XML document: %v", err)
	}

	// Update vendor catalog info
	vc.CatalogInfo.VendorId = sdk.PtrString(manifest.Identifier)
	vc.CatalogInfo.VendorReleaseTimestamp = sdk.PtrString(manifest.DateTime)
	vc.CatalogInfo.VendorConfiguration = map[string]any{
		"baseLocation":                manifest.BaseLocation,
		"baseLocationAccessProtocols": manifest.BaseLocationAccessProtocols,
		"releaseID":                   manifest.ReleaseID,
		"version":                     manifest.Version,
		"predecessorID":               manifest.PredecessorID,
	}

	baseVendorDownloadUrl, err := url.Parse(strings.ToLower(strings.Split(manifest.BaseLocationAccessProtocols, ",")[0]) + "://" + manifest.BaseLocation)
	if err != nil {
		return fmt.Errorf("error parsing base catalog URL: %v", err)
	}

	// Process catalog components
	for _, component := range manifest.Components {
		// We only check for components that are of type firmware
		if component.ComponentType.Value != componentTypeFirmware {
			continue
		}

		includedBinary := false
		supportedSystems := []map[string]interface{}{}
		for _, brand := range component.SupportedSystems.Brands {
			for _, model := range brand.Models {
				systemInfo := map[string]interface{}{
					"id": model.SystemID,
					// "idType":      model.SystemIDType,
					// "brandName":   brand.Display,
					// "brandPrefix": brand.Prefix,
					// "modelName":   model.Display,
				}

				supportedSystems = append(supportedSystems, systemInfo)

				if !includedBinary {
					systemName := brand.Display + " " + model.Display
					if len(vc.VendorSystemsFilter) == 0 || slices.Contains(vc.VendorSystemsFilter, systemName) {
						includedBinary = true
					}
				}
			}
		}

		if !includedBinary {
			continue
		}

		rebootRequired := false
		if component.RebootRequired == "true" {
			rebootRequired = true
		}

		supportedDevices := []map[string]interface{}{}
		for _, device := range component.SupportedDevices.Devices {
			deviceInfo := map[string]interface{}{
				"id":    device.ComponentID,
				"model": device.Display,
			}

			supportedDevices = append(supportedDevices, deviceInfo)
		}

		timestamp, err := time.Parse("January 02, 2006", component.ReleaseDate)
		if err != nil {
			return fmt.Errorf("error parsing release date: %v", err)
		}

		severity := parseDellUpdateSeverity(component.UpdateSeverity.Value)

		firmwareBinary := sdk.FirmwareBinary{
			ExternalId:             sdk.PtrString(component.Path),
			Name:                   component.Name.Display,
			VendorInfoUrl:          sdk.PtrString(component.ImportantInfo.URL),
			VendorDownloadUrl:      baseVendorDownloadUrl.JoinPath(component.Path).String(),
			CacheDownloadUrl:       nil, //	Will be set after the binary is downloaded
			PackageId:              sdk.PtrString(component.PackageID),
			PackageVersion:         sdk.PtrString(component.VendorVersion),
			RebootRequired:         rebootRequired,
			UpdateSeverity:         severity,
			VendorSupportedDevices: supportedDevices,
			VendorSupportedSystems: supportedSystems,
			VendorReleaseTimestamp: sdk.PtrString(timestamp.Format(time.RFC3339)),
			Vendor: map[string]any{
				"path":        component.Path,
				"size":        component.Size,
				"category":    component.Category.Display,
				"datetime":    component.DateTime,
				"releaseId":   component.ReleaseID,
				"packageType": component.PackageType,
			},
		}

		vc.Binaries = append(vc.Binaries, &firmwareBinary)
	}

	return nil
}

func parseDellUpdateSeverity(dellCriticality string) sdk.FirmwareBinaryUpdateSeverity {
	switch dellCriticality {
	case "2":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_CRITICAL
	case "Urgent":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_CRITICAL
	case "1":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_RECOMMENDED
	case "Recommended":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_RECOMMENDED
	case "3":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_OPTIONAL
	case "Optional":
		return sdk.FIRMWAREBINARYUPDATESEVERITY_OPTIONAL
	default:
		return sdk.FIRMWAREBINARYUPDATESEVERITY_UNKNOWN
	}
}
