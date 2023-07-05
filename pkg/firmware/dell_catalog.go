package firmware

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	STOP_AFTER            int    = 1
	componentTypeFirmware string = "FRMW"
)

type manifest struct {
	BaseLocation                string              `xml:"baseLocation,attr"`
	BaseLocationAccessProtocols string              `xml:"baseLocationAccessProtocols,attr"`
	DateTime                    string              `xml:"dateTime,attr"`
	Identifier                  string              `xml:"identifier,attr"`
	ReleaseID                   string              `xml:"releaseID,attr"`
	Version                     string              `xml:"version,attr"`
	PredecessorID               string              `xml:"predecessorID,attr"`
	Software                    []softwareBundle    `xml:"SoftwareBundle"`
	Components                  []softwareComponent `xml:"SoftwareComponent"`
}

type softwareBundle struct {
	XMLName          xml.Name         `xml:"SoftwareBundle"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	Name             name             `xml:"Name"`
	SupportedDevices supportedDevices `xml:"SupportedDevices"`
	SupportedSystems supportedSystems `xml:"SupportedSystems"`
}

type name struct {
	Display string `xml:"Display"`
}

type criticality struct {
	Value string `xml:"value,attr"`
}

type importantInfo struct {
	URL string `xml:"URL,attr"`
}

type category struct {
	Display string `xml:"Display"`
}

type supportedDevices struct {
	Devices []device `xml:"Device"`
}

type componentType struct {
	Value string `xml:"value,attr"`
}

type description struct {
	Display string `xml:"Display"`
}

type device struct {
	ComponentID string `xml:"componentID,attr"`
	Display     string `xml:"Display"`
	Prefix      string `xml:"prefix,attr"`
}

type supportedSystems struct {
	Brands []brand `xml:"Brand"`
}

type brand struct {
	Models  []model `xml:"Model"`
	Display string  `xml:"Display"`
	Prefix  string  `xml:"prefix,attr"`
}

type model struct {
	SystemID     string `xml:"systemID,attr"`
	SystemIDType string `xml:"systemIDType,attr"`
	Display      string `xml:"Display"`
}

type softwareComponent struct {
	XMLName          xml.Name         `xml:"SoftwareComponent"`
	DateTime         string           `xml:"dateTime,attr"`
	Path             string           `xml:"path,attr"`
	PackageID        string           `xml:"packageID,attr"`
	VendorVersion    string           `xml:"vendorVersion,attr"`
	RebootRequired   string           `xml:"rebootRequired,attr"`
	UpdateSeverity   criticality      `xml:"Criticality"`
	ReleaseDate      string           `xml:"releaseDate,attr"`
	Size             string           `xml:"size,attr"`
	HashMD5          string           `xml:"hashMD5,attr"`
	Category         category         `xml:"Category"`
	ReleaseID        string           `xml:"releaseID,attr"`
	PackageType      string           `xml:"packageType,attr"`
	ImportantInfo    importantInfo    `xml:"ImportantInfo"`
	ComponentType    componentType    `xml:"ComponentType"`
	Name             name             `xml:"Name"`
	Description      description      `xml:"Description"`
	SupportedDevices supportedDevices `xml:"SupportedDevices"`
	SupportedSystems supportedSystems `xml:"SupportedSystems"`
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

func parseDellCatalog(configFile rawConfigFile) (firmwareCatalog, []firmwareBinary, error) {
	if configFile.DownloadCatalog {
		err := downloadCatalog(configFile.CatalogUrl, configFile.CatalogPath)
		if err != nil {
			return firmwareCatalog{}, nil, err
		}
	}

	xmlFile, err := os.Open(configFile.CatalogPath)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}
	defer xmlFile.Close()

	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	reader, err := charset.NewReader(xmlFile, "utf-16")
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = BypassReader

	var manifest manifest

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

	catalog := firmwareCatalog{
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

	baseCatalogURL := ""

	if configFile.DownloadCatalog {
		url, err := url.Parse(configFile.CatalogUrl)

		if err != nil {
			return firmwareCatalog{}, nil, fmt.Errorf("Unable to parse catalog URL: %v", err)
		}

		baseCatalogURL = url.Scheme + "://" + url.Host
	}

	firmwareBinaryCollection := []firmwareBinary{}
	repositoryURL := getFirmwareRepositoryURL()

	for idx, component := range manifest.Components {
		// We only check for components that are of type firmware
		if component.ComponentType.Value != "FRMW" {
			continue
		}

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
			"hashmd5":       component.HashMD5,
			"category":      component.Category.Display,
			"datetime":      component.DateTime,
			"packageId":     component.PackageID,
			"releaseId":     component.ReleaseID,
			"packageType":   component.PackageType,
			"importantInfo": component.ImportantInfo.URL,
		}

		timestamp, err := time.Parse("January 02, 2006", component.ReleaseDate)

		if err != nil {
			return firmwareCatalog{}, nil, fmt.Errorf("Error parsing release date: %v", err)
		}

		severity, err := getSeverity(component.UpdateSeverity.Value)

		if err != nil {
			return firmwareCatalog{}, nil, fmt.Errorf("Error parsing severity: %v", err)
		}

		downloadURL := ""
		if configFile.DownloadCatalog {
			downloadURL = baseCatalogURL + "/" + component.Path
		}

		firmwareBinary := firmwareBinary{
			ExternalId:             component.Path,
			Name:                   component.Name.Display,
			Description:            component.Description.Display,
			PackageId:              component.PackageID,
			PackageVersion:         component.VendorVersion,
			RebootRequired:         rebootRequired,
			UpdateSeverity:         severity,
			SupportedDevices:       supportedDevices,
			SupportedSystems:       supportedSystems,
			VendorProperties:       componentVendorConfiguration,
			VendorReleaseTimestamp: timestamp.Format(time.RFC3339),
			CreatedTimestamp:       time.Now().Format(time.RFC3339),
			DownloadURL:            downloadURL,
			RepoURL:                repositoryURL + "/" + component.Path,
		}

		firmwareBinaryCollection = append(firmwareBinaryCollection, firmwareBinary)

		firmwareBinaryJson, _ := json.MarshalIndent(firmwareBinary, "", "  ")
		fmt.Printf("Created firmware binary: %v\n", string(firmwareBinaryJson))

		if idx > STOP_AFTER {
			break
		}
	}

	return catalog, firmwareBinaryCollection, nil
}
