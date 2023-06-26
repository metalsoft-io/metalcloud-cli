package firmware

import (
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

const (
	configFormatJSON          = "json"
	configFormatYAML          = "yaml"
	catalogVendorDell         = "Dell"
	catalogVendorLenovo       = "Lenovo"
	catalogVendorHp           = "Hp"
	catalogUpdateTypeOffline  = "offline"
	catalogUpdateTypeOnline   = "online"
	stringMinimumSize         = 1
	stringMaximumSize         = 255
	updateSeverityUnknown     = "unknown"
	updateSeverityRecommended = "recommended"
	updateSeverityCritical    = "critical"
	updateSeverityOptional    = "optional"
)

type serverInfo struct {
	MachineType  string `json:"machineType" yaml:"machine_type"`
	SerialNumber string `json:"serialNumber" yaml:"serial_number"`
}

type rawConfigFile struct {
	Name            string       `json:"name" yaml:"name"`
	Description     string       `json:"description" yaml:"description"`
	Vendor          string       `json:"vendor" yaml:"vendor"`
	CatalogUrl      string       `json:"catalogUrl" yaml:"catalog_url"`
	DownloadCatalog bool         `json:"downloadCatalog" yaml:"download_catalog"`
	CatalogPath     string       `json:"catalogPath" yaml:"catalog_path"`
	ServersList     []serverInfo `json:"serversList" yaml:"servers_list"`
}

type catalog struct {
	Name                   string
	Description            string
	Vendor                 string
	VendorID               string
	VendorURL              string
	VendorReleaseTimestamp string
	UpdateType             string
	ServerTypesSupported   []string
	Configuration          map[string]string
	CreatedTimestamp       string
}

type firmwareBinary struct {
	ExternalId             string
	Name                   string
	PackageId              string
	PackageVersion         string
	RebootRequired         bool
	UpdateSeverity         string
	SupportedDevices       []map[string]string
	SupportedSystems       []map[string]string
	VendorProperties       map[string]string
	VendorReleaseTimestamp string
	CreatedTimestamp       string
}

func parseConfigFile(configFormat string, rawConfigFileContents []byte, configFile *rawConfigFile) error {
	switch configFormat {
	case configFormatJSON:
		err := json.Unmarshal(rawConfigFileContents, &configFile)

		if err != nil {
			return fmt.Errorf("error parsing json file %s: %s", rawConfigFileContents, err.Error())
		}

	case configFormatYAML:
		err := yaml.Unmarshal(rawConfigFileContents, &configFile)

		if err != nil {
			return fmt.Errorf("error parsing yaml file %s: %s", rawConfigFileContents, err.Error())
		}

	default:
		validFormats := []string{configFormatJSON, configFormatYAML}
		return fmt.Errorf("the 'config-format' parameter must be one of %v", validFormats)
	}

	structValue := reflect.ValueOf(configFile).Elem()
	fieldNum := structValue.NumField()

	optionalFields := []string{"CatalogUrl", "ServersList"}

	for i := 0; i < fieldNum; i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if !isSet && !slices.Contains[string](optionalFields, fieldName) {
			return fmt.Errorf("the '%s' field must be set in the raw-config file", fieldName)
		}
	}

	if configFile.DownloadCatalog && configFile.CatalogUrl == "" {
		return fmt.Errorf("the 'catalogUrl' field must be set in the raw-config file when 'downloadCatalog' is set to true")
	}

	if configFile.Vendor == catalogVendorLenovo && configFile.ServersList == nil {
		return fmt.Errorf("the 'serversList' field must be set in the raw-config file when 'vendor' is set to '%s'", catalogVendorLenovo)
	}

	checkStringSize(configFile.Name)
	checkStringSize(configFile.Description)
	checkStringSize(configFile.CatalogUrl)

	return nil
}

func checkStringSize(str string) error {
	if len(str) < stringMinimumSize {
		return fmt.Errorf("the '%s' field must be at least %d characters", str, stringMinimumSize)
	}

	if len(str) > stringMaximumSize {
		return fmt.Errorf("the '%s' field must be less than %d characters", str, stringMaximumSize)
	}

	return nil
}

func getUpdateType(rawConfigFile rawConfigFile) string {
	if rawConfigFile.DownloadCatalog {
		return catalogUpdateTypeOnline
	}

	return catalogUpdateTypeOffline
}

func getSeverity(input string) (string, error) {
	switch input {
	case "0":
		return updateSeverityUnknown, nil
	case "1":
		return updateSeverityRecommended, nil
	case "2":
		return updateSeverityCritical, nil
	case "3":
		return updateSeverityOptional, nil
	default:
		return "", fmt.Errorf("invalid severity value: %s", input)
	}
}