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
				"config_format":          c.FlagSet.String("config-format", command.NilDefaultStr, "The format of the config file. Supported values are 'json' and 'yaml'."),
				"raw_config":             c.FlagSet.String("raw-config", command.NilDefaultStr, "The path to the config file."),
				"download_binaries":      c.FlagSet.Bool("download-binaries", false, colors.Yellow("(Optional)")+" Download firmware binaries from the catalog to the local filesystem."),
				"skip_upload_to_repo":    c.FlagSet.Bool("skip-upload-to-repo", false, colors.Yellow("(Optional)")+" Skip firmware binaries upload to the HTTP repository."),
				"skip_host_key_checking": c.FlagSet.Bool("skip-host-key-checking", false, colors.Yellow("(Optional)")+" Skip check when adding a host key to the known_hosts file in the firmware binary upload process."),
				"replace_if_exists":      c.FlagSet.Bool("replace-if-exists", false, colors.Yellow("(Optional)")+" Replaces firmware binaries if the already exist in the HTTP repository."),
				"filter_server_types":	  c.FlagSet.String("filter-server-types", command.NilDefaultStr, colors.Yellow("(Optional)")+" Comma separated list of server types to filter the firmware catalog by. Defaults to all supported server types."),
				"debug":                  c.FlagSet.Bool("debug", false, colors.Green("(Flag)")+" If set, increases log level."),
			}
		},
		ExecuteFunc: firmwareCatalogCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
	},
}

func firmwareCatalogCreateCmd(c *command.Command, client metalcloud.MetalCloudClient) (string, error) {
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

	downloadBinaries := false
	if command.GetBoolParam(c.Arguments["download_binaries"]) {
		downloadBinaries = true
	}

	configFile := rawConfigFile{}
	err := parseConfigFile(configFormat, rawConfigFileContents, &configFile, downloadBinaries)

	if err != nil {
		return "", err
	}

	uploadToRepo := true
	if command.GetBoolParam(c.Arguments["skip_upload_to_repo"]) {
		uploadToRepo = false
	}

	filterServerTypes := ""
	if filterValue, ok := command.GetStringParamOk(c.Arguments["filter_server_types"]); ok {
		filterServerTypes = filterValue
	}

	var catalog firmwareCatalog
	var binaryCollection []*firmwareBinary

	switch configFile.Vendor {
	case catalogVendorDell:
		catalog, binaryCollection, err = parseDellCatalog(client, configFile, filterServerTypes, uploadToRepo, downloadBinaries)

		if err != nil {
			return "", err
		}

	case catalogVendorLenovo:
		catalog, binaryCollection, err = parseLenovoCatalog(configFile, client, filterServerTypes, uploadToRepo, downloadBinaries)

		if err != nil {
			return "", err
		}

	case catalogVendorHp:
		return "", fmt.Errorf("vendor '%s' is not yet supported", catalogVendorHp)

	default:
		validVendors := []string{catalogVendorDell, catalogVendorLenovo, catalogVendorHp}
		return "", fmt.Errorf("invalid vendor '%s' found in the raw-config file. Supported vendors are %v", configFile.Vendor, validVendors)
	}

	if downloadBinaries {
		err := downloadBinariesFromCatalog(binaryCollection)

		if err != nil {
			return "", err
		}
	}

	replaceIfExists := false
	if command.GetBoolParam(c.Arguments["replace_if_exists"]) {
		replaceIfExists = true
	}

	skipHostKeyChecking := false
	if command.GetBoolParam(c.Arguments["skip_host_key_checking"]) {
		skipHostKeyChecking = true
	}

	if uploadToRepo {
		err := uploadBinariesToRepository(binaryCollection, replaceIfExists, skipHostKeyChecking)

		if err != nil {
			return "", err
		}
	}

	err = sendCatalog(catalog)
	if err != nil {
		return "", err
	}

	err = sendBinaries(binaryCollection)
	if err != nil {
		return "", err
	}

	return "", nil
}
