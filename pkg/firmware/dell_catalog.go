package firmware

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html/charset"

	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/internal/networking"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

const (
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

func downloadDellCatalog(catalogURL string, catalogFilePath string) error {
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
	defer resp.Body.Close()

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

func parseDellCatalog(client metalcloud.MetalCloudClient, configFile rawConfigFile, serverTypesFilter string, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) (firmwareCatalog, []*firmwareBinary, error) {
	supportedServerTypes, _, err := retrieveSupportedServerTypes(client, serverTypesFilter)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	catalog, manifest, err := processDellCatalog(configFile)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	firmwareBinaryCollection, err := processDellBinaries(configFile, manifest, &catalog, supportedServerTypes, uploadToRepo, downloadBinaries, repoConfig)
	if err != nil {
		return firmwareCatalog{}, nil, err
	}

	return catalog, firmwareBinaryCollection, nil
}

func processDellCatalog(configFile rawConfigFile) (firmwareCatalog, manifest, error) {
	if configFile.CatalogUrl != "" {
		err := downloadDellCatalog(configFile.CatalogUrl, configFile.LocalCatalogPath)
		if err != nil {
			return firmwareCatalog{}, manifest{}, err
		}
	}

	xmlFile, err := os.Open(configFile.LocalCatalogPath)
	if err != nil {
		return firmwareCatalog{}, manifest{}, err
	}
	defer xmlFile.Close()

	if err != nil {
		return firmwareCatalog{}, manifest{}, err
	}

	reader, err := charset.NewReader(xmlFile, "utf-16")
	if err != nil {
		return firmwareCatalog{}, manifest{}, err
	}

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = BypassReader

	var manifest manifest

	err = decoder.Decode(&manifest)
	if err != nil {
		log.Fatal(err)
	}

	catalogConfiguration := map[string]any{
		"baseLocation":                manifest.BaseLocation,
		"baseLocationAccessProtocols": manifest.BaseLocationAccessProtocols,
		"dateTime":                    manifest.DateTime,
		"identifier":                  manifest.Identifier,
		"releaseID":                   manifest.ReleaseID,
		"version":                     manifest.Version,
		"predecessorID":               manifest.PredecessorID,
	}

	vendorId := configFile.Vendor
	err = checkStringSize(vendorId, 1, 255)
	if err != nil {
		return firmwareCatalog{}, manifest, err
	}

	catalog := firmwareCatalog{
		Name:                          configFile.Name,
		Description:                   configFile.Description,
		Vendor:                        configFile.Vendor,
		VendorID:                      vendorId,
		VendorURL:                     configFile.CatalogUrl,
		VendorReleaseTimestamp:        manifest.DateTime,
		UpdateType:                    getUpdateType(configFile),
		MetalSoftServerTypesSupported: []string{},
		ServerTypesSupported:          []string{},
		Configuration:                 catalogConfiguration,
		CreatedTimestamp:              time.Now().Format(time.RFC3339),
	}

	return catalog, manifest, nil
}

func processDellBinaries(configFile rawConfigFile, dellManifest manifest, catalog *firmwareCatalog, supportedServerTypes map[string][]string, uploadToRepo, downloadBinaries bool, repoConfig repoConfiguration) ([]*firmwareBinary, error) {
	metalsoftServerTypes := []string{}
	dellServerTypes := []string{}

	for dellServerType, msServerTypes := range supportedServerTypes {
		for _, msServerType := range msServerTypes {
			if !slices.Contains[string](metalsoftServerTypes, msServerType) {
				metalsoftServerTypes = append(metalsoftServerTypes, msServerType)
			}
		}
		dellServerTypes = append(dellServerTypes, dellServerType)
	}

	baseCatalogURL := ""

	if configFile.CatalogUrl != "" {
		url, err := url.Parse(configFile.CatalogUrl)

		if err != nil {
			return nil, fmt.Errorf("unable to parse catalog URL: %v", err)
		}

		baseCatalogURL = url.Scheme + "://" + url.Host
	}

	firmwareBinaryCollection := []*firmwareBinary{}
	repositoryURL := repoConfig.HttpUrl
	if repositoryURL == "" {
		var err error
		repositoryURL, err = configuration.GetFirmwareRepositoryURL()
		if uploadToRepo && err != nil {
			return nil, fmt.Errorf("error getting firmware repository URL: %v", err)
		}
	}

	for _, component := range dellManifest.Components {
		// We only check for components that are of type firmware
		if component.ComponentType.Value != componentTypeFirmware {
			continue
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

		var systemName string
		validBinary := false
		for _, supportedSystem := range supportedSystems {
			systemName = supportedSystem["brandName"] + " " + supportedSystem["modelName"]
			if !validBinary && slices.Contains[string](dellServerTypes, systemName) {
				validBinary = true
				break
			}
		}

		if !validBinary {
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

		componentVendorConfiguration := map[string]any{
			"path":          component.Path,
			"size":          component.Size,
			"category":      component.Category.Display,
			"datetime":      component.DateTime,
			"packageId":     component.PackageID,
			"releaseId":     component.ReleaseID,
			"packageType":   component.PackageType,
			"importantInfo": component.ImportantInfo.URL,
		}

		timestamp, err := time.Parse("January 02, 2006", component.ReleaseDate)

		if err != nil {
			return nil, fmt.Errorf("error parsing release date: %v", err)
		}

		severity, err := getSeverity(component.UpdateSeverity.Value)

		if err != nil {
			return nil, fmt.Errorf("error parsing severity: %v", err)
		}

		downloadURL := ""
		if configFile.CatalogUrl != "" {
			downloadURL, err = url.JoinPath(baseCatalogURL, component.Path)
			if err != nil {
				return nil, err
			}
		}

		componentPathArr := strings.Split(component.Path, "/")
		componentName := componentPathArr[len(componentPathArr)-1]
		componentRepoUrl, err := url.JoinPath(repositoryURL, componentName)
		if err != nil {
			return nil, err
		}

		localPath := ""
		if configFile.LocalFirmwarePath != "" && (downloadBinaries || configFile.CatalogUrl == "") {
			localPath, err = filepath.Abs(filepath.Join(configFile.LocalFirmwarePath, componentName))

			if err != nil {
				return nil, fmt.Errorf("error getting download binary absolute path: %v", err)
			}
		}

		firmwareBinary := firmwareBinary{
			ExternalId:             component.Path,
			Name:                   component.Name.Display,
			FileName:               componentName,
			Description:            component.Description.Display,
			PackageId:              component.PackageID,
			PackageVersion:         component.VendorVersion,
			RebootRequired:         rebootRequired,
			UpdateSeverity:         severity,
			Hash:                   component.HashMD5,
			HashingAlgorithm:       networking.HashingAlgorithmMD5,
			SupportedDevices:       supportedDevices,
			SupportedSystems:       supportedSystems,
			VendorProperties:       componentVendorConfiguration,
			VendorReleaseTimestamp: timestamp.Format(time.RFC3339),
			CreatedTimestamp:       time.Now().Format(time.RFC3339),
			DownloadURL:            downloadURL,
			RepoURL:                componentRepoUrl,
			LocalPath:              localPath,
		}

		firmwareBinaryCollection = append(firmwareBinaryCollection, &firmwareBinary)

		if !slices.Contains[string](catalog.ServerTypesSupported, systemName) {
			catalog.ServerTypesSupported = append(catalog.ServerTypesSupported, systemName)

			for _, supportedServerType := range supportedServerTypes[systemName] {
				if !slices.Contains[string](catalog.MetalSoftServerTypesSupported, supportedServerType) {
					catalog.MetalSoftServerTypesSupported = append(catalog.MetalSoftServerTypesSupported, supportedServerType)
				}
			}
		}
	}

	return firmwareBinaryCollection, nil
}
