package cmd

import (
	"encoding/json"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_catalog"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	firmwareCatalogFlags = struct {
		configSource            string
		name                    string
		description             string
		vendor                  string
		updateType              string
		vendorUrl               string
		vendorToken             string
		vendorLocalCatalogPath  string
		vendorLocalBinariesPath string
		serverTypes             []string
		vendorSystems           []string
		downloadBinaries        bool
		uploadBinaries          bool
		repoBaseUrl             string
		repoSshHost             string
		repoSshUser             string
		userPrivateKeyPath      string
		knownHostsPath          string
		ignoreHostKeyCheck      bool
	}{}

	firmwareCatalogCmd = &cobra.Command{
		Use:     "firmware-catalog [command]",
		Aliases: []string{"fw-catalog", "firmware"},
		Short:   "Firmware catalog management",
		Long:    `Firmware catalog management commands.`,
	}

	firmwareCatalogListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all firmware catalogs.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogList(cmd.Context())
		},
	}

	firmwareCatalogGetCmd = &cobra.Command{
		Use:          "get firmware_catalog_id",
		Aliases:      []string{"show"},
		Short:        "Get firmware catalog details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogGet(cmd.Context(), args[0])
		},
	}

	firmwareCatalogCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new firmware catalog.",
		Long:    `Create a new firmware catalog.`,
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
localFirmwarePath: ./hp_downloads
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var firmwareCatalogOptions firmware_catalog.FirmwareCatalogCreateOptions

			// If config source is provided, use it
			if firmwareCatalogFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
				if err != nil {
					return err
				}

				err = json.Unmarshal(config, &firmwareCatalogOptions)
				if err != nil {
					err = yaml.Unmarshal(config, &firmwareCatalogOptions)
					if err != nil {
						return err
					}
				}
			} else {
				// Otherwise build config from command line parameters
				firmwareCatalogOptions = firmware_catalog.FirmwareCatalogCreateOptions{
					Name:                    firmwareCatalogFlags.name,
					Description:             firmwareCatalogFlags.description,
					Vendor:                  firmwareCatalogFlags.vendor,
					UpdateType:              firmwareCatalogFlags.updateType,
					VendorUrl:               firmwareCatalogFlags.vendorUrl,
					VendorToken:             firmwareCatalogFlags.vendorToken,
					ServerTypesFilter:       firmwareCatalogFlags.serverTypes,
					VendorSystemsFilter:     firmwareCatalogFlags.vendorSystems,
					VendorLocalCatalogPath:  firmwareCatalogFlags.vendorLocalCatalogPath,
					VendorLocalBinariesPath: firmwareCatalogFlags.vendorLocalBinariesPath,
					DownloadBinaries:        firmwareCatalogFlags.downloadBinaries,
					UploadBinaries:          firmwareCatalogFlags.uploadBinaries,
					RepoBaseUrl:             firmwareCatalogFlags.repoBaseUrl,
					RepoSshHost:             firmwareCatalogFlags.repoSshHost,
					RepoSshUser:             firmwareCatalogFlags.repoSshUser,
					UserPrivateKeyPath:      firmwareCatalogFlags.userPrivateKeyPath,
					KnownHostsPath:          firmwareCatalogFlags.knownHostsPath,
					IgnoreHostKeyCheck:      firmwareCatalogFlags.ignoreHostKeyCheck,
				}
			}

			return firmware_catalog.FirmwareCatalogCreate(cmd.Context(), firmwareCatalogOptions)
		},
	}

	firmwareCatalogUpdateCmd = &cobra.Command{
		Use:          "update firmware_catalog_id",
		Aliases:      []string{"edit"},
		Short:        "Update a firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
			if err != nil {
				return err
			}

			firmwareCatalogId := ""
			if len(args) > 0 {
				firmwareCatalogId = args[0]
			}

			return firmware_catalog.FirmwareCatalogUpdate(cmd.Context(), firmwareCatalogId, config)
		},
	}

	firmwareCatalogDeleteCmd = &cobra.Command{
		Use:          "delete firmware_catalog_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(firmwareCatalogCmd)

	firmwareCatalogCmd.AddCommand(firmwareCatalogListCmd)
	firmwareCatalogCmd.AddCommand(firmwareCatalogGetCmd)

	firmwareCatalogCmd.AddCommand(firmwareCatalogCreateCmd)
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.configSource, "config-source", "", "Source of the new firmware catalog configuration. Can be 'pipe' or path to a JSON file.")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.name, "name", "", "Name of the firmware catalog")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.description, "description", "", "Description of the firmware catalog")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendor, "vendor", "", "Vendor type (e.g., 'dell', 'hp')")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.updateType, "update-type", "", "Update type (e.g., 'online', 'offline')")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorUrl, "vendor-url", "", "URL of the online vendor catalog")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorToken, "vendor-token", "", "Token for accessing the online vendor catalog")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorLocalCatalogPath, "vendor-local-catalog-path", "", "Path to the local catalog file")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorLocalBinariesPath, "vendor-local-binaries-path", "", "Path to the local binaries directory")
	firmwareCatalogCreateCmd.Flags().StringSliceVar(&firmwareCatalogFlags.serverTypes, "server-types", []string{}, "List of supported Metalsoft server types (comma-separated)")
	firmwareCatalogCreateCmd.Flags().StringSliceVar(&firmwareCatalogFlags.vendorSystems, "vendor-systems", []string{}, "List of supported vendor systems (comma-separated)")
	firmwareCatalogCreateCmd.Flags().BoolVar(&firmwareCatalogFlags.downloadBinaries, "download-binaries", false, "Download binaries from the vendor catalog")
	firmwareCatalogCreateCmd.Flags().BoolVar(&firmwareCatalogFlags.uploadBinaries, "upload-binaries", false, "Upload binaries to the offline repository")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.repoBaseUrl, "repo-base-url", "", "Base URL of the offline repository")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.repoSshHost, "repo-ssh-host", "", "SSH host with port of the offline repository")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.repoSshUser, "repo-ssh-user", "", "SSH user for the offline repository")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.userPrivateKeyPath, "user-private-key-path", "~/.ssh/id_rsa", "Path to the user's private SSH key")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.knownHostsPath, "known-hosts-path", "~/.ssh/known_hosts", "Path to the known hosts file for SSH connections")
	firmwareCatalogCreateCmd.Flags().BoolVar(&firmwareCatalogFlags.ignoreHostKeyCheck, "ignore-host-key-check", false, "Ignore host key check for SSH connections")
	firmwareCatalogCreateCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	firmwareCatalogCreateCmd.MarkFlagsRequiredTogether("name", "vendor", "update-type")
	firmwareCatalogCreateCmd.MarkFlagsMutuallyExclusive("vendor-url", "vendor-local-catalog-path")
	firmwareCatalogCreateCmd.MarkFlagsRequiredTogether("upload-binaries", "repo-ssh-host", "repo-ssh-user")
	firmwareCatalogCreateCmd.MarkFlagsMutuallyExclusive("known-hosts-path", "ignore-host-key-check")

	firmwareCatalogCmd.AddCommand(firmwareCatalogUpdateCmd)
	firmwareCatalogUpdateCmd.Flags().StringVar(&firmwareCatalogFlags.configSource, "config-source", "", "Source of the firmware catalog configuration updates. Can be 'pipe' or path to a JSON file.")
	firmwareCatalogUpdateCmd.MarkFlagsOneRequired("config-source")

	firmwareCatalogCmd.AddCommand(firmwareCatalogDeleteCmd)
}
