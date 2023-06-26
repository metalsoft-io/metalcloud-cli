package firmware

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	STOP_AFTER int = 1
)

type Manifest struct {
	BaseLocation                string              `xml:"baseLocation,attr"`
	BaseLocationAccessProtocols string              `xml:"baseLocationAccessProtocols,attr"`
	DateTime                    string              `xml:"dateTime,attr"`
	Identifier                  string              `xml:"identifier,attr"`
	ReleaseID                   string              `xml:"releaseID,attr"`
	Version                     string              `xml:"version,attr"`
	PredecessorID               string              `xml:"predecessorID,attr"`
	Software                    []SoftwareBundle    `xml:"SoftwareBundle"`
	Components                  []SoftwareComponent `xml:"SoftwareComponent"`
}

type SoftwareBundle struct {
	XMLName          xml.Name         `xml:"SoftwareBundle"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	Name             Name             `xml:"Name"`
	SupportedDevices SupportedDevices `xml:"SupportedDevices"`
	SupportedSystems SupportedSystems `xml:"SupportedSystems"`
}

type Name struct {
	Display string `xml:"Display"`
}

type Criticality struct {
	Value string `xml:"value,attr"`
}

type ImportantInfo struct {
	URL string `xml:"URL,attr"`
}

type Category struct {
	Display string `xml:"Display"`
}

type SupportedDevices struct {
	Devices []Device `xml:"Device"`
}

type Device struct {
	ComponentID string `xml:"componentID,attr"`
	Display     string `xml:"Display"`
	Prefix      string `xml:"prefix,attr"`
}

type SupportedSystems struct {
	Brands []Brand `xml:"Brand"`
}

type Brand struct {
	Models  []Model `xml:"Model"`
	Display string  `xml:"Display"`
	Prefix  string  `xml:"prefix,attr"`
}

type Model struct {
	SystemID     string `xml:"systemID,attr"`
	SystemIDType string `xml:"systemIDType,attr"`
	Display      string `xml:"Display"`
}

type SoftwareComponent struct {
	XMLName          xml.Name         `xml:"SoftwareComponent"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	PackageID        string           `xml:"packageID,attr"`
	VendorVersion    string           `xml:"vendorVersion,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	UpdateSeverity   Criticality      `xml:"Criticality"`
	ReleaseDate      string           `xml:"releaseDate,attr"`
	Size             string           `xml:"size,attr"`
	Hashmd5          string           `xml:"hashMD5,attr"`
	Category         Category         `xml:"Category"`
	ReleaseID        string           `xml:"releaseID,attr"`
	PackageType      string           `xml:"packageType,attr"`
	ImportantInfo    ImportantInfo    `xml:"ImportantInfo"`
	Name             Name             `xml:"Name"`
	SupportedDevices SupportedDevices `xml:"SupportedDevices"`
	SupportedSystems SupportedSystems `xml:"SupportedSystems"`
}

func downloadCatalog(catalogURL string, catalogFilePath string) error {
	req, err := http.NewRequest(http.MethodGet, catalogURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(catalogFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, gzReader)
	if err != nil {
		return err
	}

	return nil
}

func BypassReader(label string, input io.Reader) (io.Reader, error) {
	return input, nil
}

func parseDellCatalog(configFile rawConfigFile) error {
	if configFile.DownloadCatalog {
		err := downloadCatalog(configFile.CatalogUrl, configFile.CatalogPath)
		if err != nil {
			return err
		}
	}

	xmlFile, err := os.Open(configFile.CatalogPath)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	if err != nil {
		return err
	}

	reader, err := charset.NewReader(xmlFile, "utf-16")
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = BypassReader

	var manifest Manifest

	err = decoder.Decode(&manifest)
	if err != nil {
		log.Fatal(err)
	}

	catalogConfiguration := map[string]string{
		"baseLocation":                manifest.BaseLocation,
		"baseLocationAccessProtocols": manifest.BaseLocationAccessProtocols,
		"dateTime":                    manifest.DateTime,
		"identifier":                  manifest.Identifier,
		"releaseID":                   manifest.ReleaseID,
		"version":                     manifest.Version,
		"predecessorID":               manifest.PredecessorID,
	}

	vendorId := configFile.Vendor //TODO: What is this???
	checkStringSize(vendorId)

	catalog := catalog{
		Name:                   configFile.Name,
		Description:            configFile.Description,
		Vendor:                 configFile.Vendor,
		VendorID:               vendorId,
		VendorURL:              configFile.CatalogUrl,
		VendorReleaseTimestamp: manifest.DateTime,
		UpdateType:             getUpdateType(configFile),
		ServerTypesSupported:   []string{}, //TODO: this needs to be updated after parsing the firmware binaries
		Configuration:          catalogConfiguration,
		CreatedTimestamp:       time.Now().Format(time.RFC3339),
	}

	catalogJSON, err := json.MarshalIndent(catalog, "", "  ")
	fmt.Printf("Created catalog for %s: %+v\n", configFile.Name, string(catalogJSON))

	/**
	(
		server_firmware_binary_id -> 1,
		server_firmware_binary_catalog_id -> 1,
		server_firmware_binary_external_id -> 'FOLDER04177723M/1/Serial-ATA_Firmware_H09VC_LN_MA8F_A00.BIN',
		server_firmware_binary_name -> 'Seagate MA8F for model number(s) ST6000NM0024-1US17Z..',
		server_firmware_binary_package_id- > 'H09VC',
		server_firmware_binary_package_version -> 'MA8F',
		server_firmware_binary_reboot_required -> 'no',
		server_firmware_binary_update_severity -> 'recommended',
		server_firmware_binary_vendor_supported_devices_json -> '[{\"id\": \"103733\", \"name\": \"Makara SATA 512e\"}]',
		server_firmware_binary_vendor_supported_systems_json -> '[{\"id\": \"0723\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"DSS1500\", \"brandPrefix\": \"PE\"}, {\"id\": \"0722\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"DSS1510\", \"brandPrefix\": \"PE\"}, {\"id\": \"0724\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"DSS2500\", \"brandPrefix\": \"PE\"}, {\"id\": \"06A5\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R230\", \"brandPrefix\": \"PE\"}, {\"id\": \"0627\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R730xd\", \"brandPrefix\": \"PE\"}, {\"id\": \"0639\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R430\", \"brandPrefix\": \"PE\"}, {\"id\": \"063A\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R530\", \"brandPrefix\": \"PE\"}, {\"id\": \"06A6\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R330\", \"brandPrefix\": \"PE\"}, {\"id\": \"063B\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"T430\", \"brandPrefix\": \"PE\"}, {\"id\": \"0600\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"R730\", \"brandPrefix\": \"PE\"}, {\"id\": \"06A7\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"T330\", \"brandPrefix\": \"PE\"}, {\"id\": \"0602\", \"idType\": \"BIOS\", \"brandName\": \"PowerEdge\", \"modelName\": \"T630\", \"brandPrefix\": \"PE\"}]',
		server_firmware_binary_vendor_json -> '{\"path\": \"FOLDER04177723M/1/Serial-ATA_Firmware_H09VC_LN_MA8F_A00.BIN\", \"size\": \"40821221\", \"hashmd5\": \"3979c65df3c67a5342d707af89923de5\", \"category\": \"Serial ATA\", \"datetime\": \"2017-03-24T21:05:18+05:30\", \"packageId\": \"H09VC\", \"releaseId\": \"H09VC\", \"packageType\": \"LLXP\"}',
		server_firmware_binary_vendor_release_timestamp -> '2017-03-23 22:00:00',
		server_firmware_binary_created_timestamp -> '2023-06-23 12:46:59'
	),
	*/

	for idx, component := range manifest.Components {
		rebootRequired := false
		if component.RebootRequired == "true" {
			rebootRequired = true
		}

		supportedDevices := []map[string]string{}
		for _, device := range component.SupportedDevices.Devices {
			deviceInfo := map[string]string{
				"id":   device.ComponentID,
				"name": device.Display,
			}
			supportedDevices = append(supportedDevices, deviceInfo)
		}

		supportedSystems := []map[string]string{}
		for _, brand := range component.SupportedSystems.Brands {
			for _, model := range brand.Models {
				systemInfo := map[string]string{
					"brandName":   brand.Display,
					"brandPrefix": brand.Prefix,
					"id":          model.SystemID,
					"idType":      model.SystemIDType,
					"modelName":   model.Display,
				}
				supportedSystems = append(supportedSystems, systemInfo)
			}
		}

		componentVendorConfiguration := map[string]string{
			"path":          component.Path,
			"size":          component.Size,
			"hashmd5":       component.Hashmd5,
			"category":      component.Category.Display,
			"datetime":      component.DateTime,
			"packageId":     component.PackageID,
			"releaseId":     component.ReleaseID,
			"packageType":   component.PackageType,
			"importantInfo": component.ImportantInfo.URL,
		}

		fmt.Println(component.ReleaseDate)
		timestamp, err := time.Parse("January 02, 2006", component.ReleaseDate)

		if err != nil {
			return fmt.Errorf("Error parsing release date: %v", err)
		}

		severity, err := getSeverity(component.UpdateSeverity.Value)

		if err != nil {
			return fmt.Errorf("Error parsing severity: %v", err)
		}

		firmwareBinary := firmwareBinary{
			ExternalId:             component.Path,
			Name:                   component.Name.Display,
			PackageId:              component.PackageID,
			PackageVersion:         component.VendorVersion,
			RebootRequired:         rebootRequired,
			UpdateSeverity:         severity,
			SupportedDevices:       supportedDevices,
			SupportedSystems:       supportedSystems,
			VendorProperties:       componentVendorConfiguration,
			VendorReleaseTimestamp: timestamp.Format(time.RFC3339),
			CreatedTimestamp:       time.Now().Format(time.RFC3339),
		}

		firmwareBinaryJson, _ := json.MarshalIndent(firmwareBinary, "", "  ")
		fmt.Printf("Created firmware binary: %v\n", string(firmwareBinaryJson))

		if idx > STOP_AFTER {
			break
		}
	}

	return nil
}
