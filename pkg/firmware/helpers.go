package firmware

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

const (
	configFormatJSON = "json"
	configFormatYAML = "yaml"
	catalogVendorDell = "Dell"
	catalogVendorLenovo = "Lenovo"
	catalogVendorHp = "Hp"
)

type rawConfigFile struct {
	Name string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Vendor string `json:"vendor" yaml:"vendor"`
	CatalogUrl string `json:"catalogUrl" yaml:"catalog_url"`
	DownloadCatalog bool `json:"downloadCatalog" yaml:"download_catalog"`
	CatalogPath string `json:"catalogPath" yaml:"catalog_path"`
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

	for i := 0; i < fieldNum; i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if !isSet && fieldName != "CatalogUrl" {
			return fmt.Errorf("the '%s' field must be set in the raw-config file", fieldName)
		}
	}

	if configFile.DownloadCatalog && configFile.CatalogUrl == "" {
		return fmt.Errorf("the 'catalogUrl' field must be set in the raw-config file when 'downloadCatalog' is set to true")
	}

	return nil
}
