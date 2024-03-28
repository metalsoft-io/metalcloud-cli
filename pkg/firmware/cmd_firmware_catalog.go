package firmware

import (
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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
				"config_format":             c.FlagSet.String("config-format", command.NilDefaultStr, colors.Red("(Required)")+" The format of the config file. Supported values are 'json' and 'yaml'."),
				"raw_config":                c.FlagSet.String("raw-config", command.NilDefaultStr, colors.Red("(Required)")+" The path to the config file."),
				"download_binaries":         c.FlagSet.Bool("download-binaries", false, colors.Yellow("(Optional)")+" Download firmware binaries from the catalog to the local filesystem."),
				"skip_upload_to_repo":       c.FlagSet.Bool("skip-upload-to-repo", false, colors.Yellow("(Optional)")+" Skip firmware binaries upload to the HTTP repository."),
				"skip_host_key_checking":    c.FlagSet.Bool("skip-host-key-checking", false, colors.Yellow("(Optional)")+" Skip check when adding a host key to the known_hosts file in the firmware binary upload process."),
				"replace_if_exists":         c.FlagSet.Bool("replace-if-exists", false, colors.Yellow("(Optional)")+" Replaces firmware binaries if the already exist in the HTTP repository."),
				"filter_server_types":       c.FlagSet.String("filter-server-types", command.NilDefaultStr, colors.Yellow("(Optional)")+" Comma separated list of server types to filter the firmware catalog by. Defaults to all supported server types."),
				"repo_http_url":             c.FlagSet.String("repo-http-url", command.NilDefaultStr, colors.Yellow("(Optional)")+" The HTTP URL of the firmware repository. Replaces the value of the METALCLOUD_FIRMWARE_REPOSITORY_URL environment variable."),
				"repo_ssh_path":             c.FlagSet.String("repo-ssh-path", command.NilDefaultStr, colors.Yellow("(Optional)")+" The SSH path of the firmware repository. Replaces the value of the METALCLOUD_FIRMWARE_REPOSITORY_SSH_PATH environment variable."),
				"repo_ssh_port":             c.FlagSet.String("repo-ssh-port", command.NilDefaultStr, colors.Yellow("(Optional)")+" The SSH port of the firmware repository. Replaces the value of the METALCLOUD_FIRMWARE_REPOSITORY_SSH_PORT environment variable."),
				"repo_ssh_user":             c.FlagSet.String("repo-ssh-user", command.NilDefaultStr, colors.Yellow("(Optional)")+" The SSH user of the firmware repository. Replaces the value of the METALCLOUD_FIRMWARE_REPOSITORY_SSH_USER environment variable."),
				"user_private_ssh_key_path": c.FlagSet.String("user-private-ssh-key-path", command.NilDefaultStr, colors.Yellow("(Optional)")+" The path to the private SSH key of the user. Replaces the value of the METALCLOUD_USER_PRIVATE_OPENSSH_KEY_PATH environment variable."),
				"debug":                     c.FlagSet.Bool("debug", false, colors.Green("(Flag)")+" If set, increases log level."),
			}
		},
		ExecuteFunc: firmwareCatalogCreateCmd,
		Endpoint:    configuration.DeveloperEndpoint,
		Example: `
Dell example:
metalcloud-cli firmware-catalog create --config-format yaml --raw-config .\example_dell.yaml --download-binaries --filter-server-types "M.32.32.2, M.32.64.2, M.40.32.1.v2" --repo-http-url http://176.223.226.61/repo/firmware/ --repo-ssh-path /home/repo/firmware --repo-ssh-port 22 --repo-ssh-user root --user-private-ssh-key-path ~/.ssh/id_rsa

example_dell.yaml:

name: test-dell
description: test
vendor: dell
catalogUrl: https://downloads.dell.com/FOLDER04655306M/1/ESXi_Catalog.xml.gz
localCatalogPath: ./ESXi_Catalog.xml
localFirmwarePath: ./downloads

Lenovo example:
metalcloud-cli firmware-catalog create --config-format json --raw-config .\example_lenovo.json --filter-server-types "M.8.8.2.v5"

example_lenovo.json:

{
	"name": "test-lenovo",
	"description": "lenovo test",
	"vendor": "lenovo",
	"localCatalogPath": "./lenovo_catalogs",
	"overwriteCatalogs": false,
	"localFirmwarePath": "./lenovo_downloads",
	"serversList": [
		{
			"machineType": "7Y51",
			"serialNumber": "J10227CF"
		}
	]
}

For Lenovo servers, we can filter by server type like for the other vendors, but we can also specify a list of servers to filter by machine type and serial number.
The list of servers specified in the config serversList parameter takes precedence over the filter-server-types parameter, which will be ignored.

HP example:
metalcloud-cli firmware-catalog create --config-format yaml --raw-config .\example_hp_gen_11.yaml

example_hp_gen_11.yaml:

name: test-hp
description: test
vendor: hp
catalogUrl: https://downloads.linux.hpe.com/SDR/repo/fwpp-gen11/current/fwrepodata/fwrepo.json
localCatalogPath: ./fwrepo.json
localFirmwarePath: ./hp_downloads`,
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

	repoConfig := getRepoConfiguration(c)

	err = sendHealthCheck()
	if err != nil {
		return "", err
	}

	var catalog firmwareCatalog
	var binaryCollection []*firmwareBinary
	downloadUser, downloadPassword := "", ""

	switch configFile.Vendor {
	case catalogVendorDell:
		catalog, binaryCollection, err = parseDellCatalog(client, configFile, filterServerTypes, uploadToRepo, downloadBinaries, repoConfig)

		if err != nil {
			return "", err
		}

	case catalogVendorLenovo:
		catalog, binaryCollection, err = parseLenovoCatalog(configFile, client, filterServerTypes, uploadToRepo, downloadBinaries, repoConfig)

		if err != nil {
			return "", err
		}

	case catalogVendorHp:
		catalog, binaryCollection, err = parseHpCatalog(configFile, client, filterServerTypes, uploadToRepo, downloadBinaries, repoConfig)

		hpSupportToken := os.Getenv("METALCLOUD_HP_SUPPORT_TOKEN")
		if hpSupportToken != "" {
			downloadUser = hpSupportToken
			downloadPassword = "null"
		}

		if err != nil {
			return "", err
		}

	default:
		validVendors := []string{catalogVendorDell, catalogVendorLenovo, catalogVendorHp}
		return "", fmt.Errorf("invalid vendor '%s' found in the raw-config file. Supported vendors are %v", configFile.Vendor, validVendors)
	}

	if downloadBinaries {
		err := downloadBinariesFromCatalog(binaryCollection, downloadUser, downloadPassword)

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
		err := uploadBinariesToRepository(binaryCollection, replaceIfExists, skipHostKeyChecking, downloadUser, downloadPassword, repoConfig)

		if err != nil {
			return "", err
		}
	}

	catalogObject, err := sendCatalog(catalog)
	if err != nil {
		return "", err
	}

	if catalogObject.ServerFirmwareCatalogId == 0 {
		return "", fmt.Errorf("received invalid firmware catalog ID. Catalog might already exist or was not created.")
	}

	err = sendBinaries(binaryCollection, catalogObject.ServerFirmwareCatalogId)
	if err != nil {
		return "", err
	}

	return "", nil
}
