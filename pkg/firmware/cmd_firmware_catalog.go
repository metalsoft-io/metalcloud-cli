package firmware

import (
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"

	"flag"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

var FirmwareCatalogCmds = []command.Command{
	{
		Description:  "Creates a firmware catalog.",
		Subject:      "firmware-catalog",
		AltSubject:   "catalog",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create firmware catalog", flag.ExitOnError),
		InitFunc: func(c *command.Command) {
			c.Arguments = map[string]interface{}{
				"label":         c.FlagSet.String("label", command.NilDefaultStr, colors.Red("(Required)")+" Firmware catalog's label"),
				"config-format": c.FlagSet.String("config-format", command.NilDefaultStr, "The format of the config file. Supported values are 'json' and 'yaml'."),
				"raw-config":    c.FlagSet.String("raw-config", command.NilDefaultStr, "The path to the config file."),
			}
		},
		ExecuteFunc: firmwareCatalogCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func firmwareCatalogCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
	var label string
	if labelValue, ok := command.GetStringParamOk(c.Arguments["label"]); !ok {
		return "", fmt.Errorf("the 'label' parameter must be specified when creating a firmware catalog")
	} else {
		label = labelValue
	}

	var configFormat string
	validFormats := []string{configFormatJSON, configFormatYAML}

	if configFormatValue, ok := command.GetStringParamOk(c.Arguments["config-format"]); !ok {
		return "", fmt.Errorf("the 'config-format' parameter must be specified when creating a firmware catalog")
	} else {
		if !slices.Contains[string](validFormats, configFormatValue) {
			return "", fmt.Errorf("the 'config-format' parameter must be one of %v", validFormats)
		}
		configFormat = configFormatValue
	}

	var rawConfigFileContents []byte
	if rawConfigFilePathValue, ok := command.GetStringParamOk(c.Arguments["raw-config"]); !ok {
		return "", fmt.Errorf("the 'raw-config' parameter must be specified when creating a firmware catalog")
	} else {
		fileContents, err := os.ReadFile(rawConfigFilePathValue)

		if err != nil {
			return "", fmt.Errorf("error reading file %s: %s", rawConfigFilePathValue, err.Error())
		}

		rawConfigFileContents = fileContents
	}

	fmt.Printf("Creating firmware catalog with label '%s' and config format '%s' and contents \n%s\n", label, configFormat, string(rawConfigFileContents))

	configFile := rawConfigFile{}
	err := parseConfigFile(configFormat, rawConfigFileContents, &configFile)

	if err != nil {
		return "", err
	}

	fmt.Printf("Parsed config file: %+v\n", configFile)

	switch configFile.Vendor {
	case catalogVendorDell:
		err := parseDellCatalog(configFile)

		if err != nil {
			return "", err
		}

	case catalogVendorLenovo:
		parseLenovoCatalog(configFile)

	case catalogVendorHp:
		return "", fmt.Errorf("vendor '%s' is not yet supported", catalogVendorHp)

	default:
		validVendors := []string{catalogVendorDell, catalogVendorLenovo, catalogVendorHp}
		return "", fmt.Errorf("invalid vendor '%s' found in the raw-config file. Supported vendors are %v", configFile.Vendor, validVendors)
	}

	return "", nil
}
