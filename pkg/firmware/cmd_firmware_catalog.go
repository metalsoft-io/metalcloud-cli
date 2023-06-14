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

const configFormatJSON = "json"
const configFormatYAML = "yaml"

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
				"config-format": c.FlagSet.String("config-format", "", "The format of the config file. Supported values are 'json' and 'yaml'."),
				"raw-config":    c.FlagSet.String("raw-config", "*", "The path to the config file."),
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
	if configFormatValue, ok := command.GetStringParamOk(c.Arguments["config-format"]); !ok {
		return "", fmt.Errorf("the 'config-format' parameter must be specified when creating a firmware catalog")
	} else {
		validFormats := []string{configFormatJSON, configFormatYAML}
		if !slices.Contains[string](validFormats, configFormatValue) {
			return "", fmt.Errorf("the 'config-format' parameter must be one of %v", validFormats)
		}
		configFormat = configFormatValue
	}

	var rawConfigFileContents string
	if rawConfigFilePathValue, ok := command.GetStringParamOk(c.Arguments["raw-config"]); !ok {
		return "", fmt.Errorf("the 'raw-config' parameter must be specified when creating a firmware catalog")
	} else {
		fileContents, err := os.ReadFile(rawConfigFilePathValue)

		if err != nil {
			return "", fmt.Errorf("error reading file %s: %s", rawConfigFilePathValue, err.Error())
		}

		rawConfigFileContents = string(fileContents)
	}

	fmt.Printf("Creating firmware catalog with label '%s' and config format '%s' and contents \n%s\n", label, configFormat, rawConfigFileContents)

	return "", nil
}
