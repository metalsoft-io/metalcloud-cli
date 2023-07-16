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
				"label":                    c.FlagSet.String("label", command.NilDefaultStr, colors.Red("(Required)")+" Firmware catalog's label"),
				"config_format":            c.FlagSet.String("config-format", command.NilDefaultStr, "The format of the config file. Supported values are 'json' and 'yaml'."),
				"raw_config":               c.FlagSet.String("raw-config", command.NilDefaultStr, "The path to the config file."),
				"download_binaries":        c.FlagSet.Bool("download-binaries", false, colors.Yellow("(Optional)")+"Download firmware binaries from the catalog to the local filesystem."),
				"skip_upload_to_repo":      c.FlagSet.Bool("skip-upload-to-repo", false, colors.Yellow("(Optional)")+"Skip firmware binaries upload to the HTTP repository."),
				"strict_host_key_checking": c.FlagSet.Bool("strict-host-key-checking", false, colors.Yellow("(Optional)")+"Do a strict check when adding a host key to the known_hosts file in the firmware binary upload process."),
				"replace_if_exists":        c.FlagSet.Bool("replace-if-exists", false, colors.Yellow("(Optional)")+"Replaces firmware binaries if the already exist in the HTTP repository."),
				"debug":                    c.FlagSet.Bool("debug", false, colors.Green("(Flag)")+"If set, increases log level."),
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

	if configFormatValue, ok := command.GetStringParamOk(c.Arguments["config_format"]); !ok {
		return "", fmt.Errorf("the 'config-format' parameter must be specified when creating a firmware catalog")
	} else {
		if !slices.Contains[string](validFormats, configFormatValue) {
			return "", fmt.Errorf("the 'config-format' parameter must be one of %v", validFormats)
		}
		configFormat = configFormatValue
	}

	var rawConfigFileContents []byte
	if rawConfigFilePathValue, ok := command.GetStringParamOk(c.Arguments["raw_config"]); !ok {
		return "", fmt.Errorf("the 'raw-config' parameter must be specified when creating a firmware catalog")
	} else {
		fileContents, err := os.ReadFile(rawConfigFilePathValue)

		if err != nil {
			return "", fmt.Errorf("error reading file %s: %s", rawConfigFilePathValue, err.Error())
		}

		rawConfigFileContents = fileContents
	}

	fmt.Printf("Creating firmware catalog with label '%s' and config format '%s' and contents \n%s\n", label, configFormat, string(rawConfigFileContents))

	downloadBinaries := false
	if command.GetBoolParam(c.Arguments["download_binaries"]) {
		downloadBinaries = true
	}

	fmt.Printf("Download binaries is set to: %v\n", downloadBinaries)

	configFile := rawConfigFile{}
	err := parseConfigFile(configFormat, rawConfigFileContents, &configFile, downloadBinaries)

	if err != nil {
		return "", err
	}

	fmt.Printf("Parsed config file: %+v\n", configFile)

	uploadToRepo := true
	if command.GetBoolParam(c.Arguments["skip_upload_to_repo"]) {
		uploadToRepo = false
	}

	fmt.Printf("Uploading to repo is set to: %v\n", uploadToRepo)

	var catalog firmwareCatalog
	var binaryCollection []firmwareBinary

	switch configFile.Vendor {
	case catalogVendorDell:
		catalog, binaryCollection, err = parseDellCatalog(configFile, client, []string{}, uploadToRepo)

		if err != nil {
			return "", err
		}

	case catalogVendorLenovo:
		catalog, binaryCollection, err = parseLenovoCatalog(configFile, client, "*", uploadToRepo)

		if err != nil {
			return "", err
		}

	case catalogVendorHp:
		return "", fmt.Errorf("vendor '%s' is not yet supported", catalogVendorHp)

	default:
		validVendors := []string{catalogVendorDell, catalogVendorLenovo, catalogVendorHp}
		return "", fmt.Errorf("invalid vendor '%s' found in the raw-config file. Supported vendors are %v", configFile.Vendor, validVendors)
	}

	if downloadBinaries	{
		err := downloadBinariesFromCatalog(binaryCollection)

		if err != nil {
			return "", err
		}
	}
	
	replaceIfExists := false
	if command.GetBoolParam(c.Arguments["replace_if_exists"]) {
		replaceIfExists = true
	}

	strictHostKeyChecking := false
	if command.GetBoolParam(c.Arguments["strict_host_key_checking"]) {
		strictHostKeyChecking = true
	}

	if uploadToRepo {
		err := uploadBinariesToRepository(binaryCollection, replaceIfExists, strictHostKeyChecking, downloadBinaries)

		if err != nil {
			return "", err
		}
	}

	sendCatalog(catalog)
	// sendBinaries(binaryCollection)

	return "", nil
}
