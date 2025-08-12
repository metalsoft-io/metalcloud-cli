package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_catalog"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
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
		repoSshPath             string
		repoSshUser             string
		userPrivateKeyPath      string
		knownHostsPath          string
		ignoreHostKeyCheck      bool
	}{}

	firmwareCatalogCmd = &cobra.Command{
		Use:     "firmware-catalog [command]",
		Aliases: []string{"fw-catalog", "firmware"},
		Short:   "Manage firmware catalogs for server hardware updates",
		Long: `Manage firmware catalogs for server hardware updates.

Firmware catalogs contain collections of firmware packages and updates for different server models
and hardware components. They can be sourced from vendor catalogs (Dell, HP, Lenovo) and used to
maintain up-to-date firmware across your infrastructure.

Supported vendors:
  • Dell - Using Dell Repository Manager (DRM) XML catalogs
  • HP/HPE - Using Service Pack for ProLiant (SPP) JSON repositories  
  • Lenovo - Using XClarity Administrator (XCA) catalogs

Update types:
  • online - Direct updates from vendor repositories
  • offline - Updates using local mirror repositories`,
	}

	firmwareCatalogListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all firmware catalogs",
		Long: `List all firmware catalogs in the system.

This command displays all available firmware catalogs with their basic information including:
- Catalog ID and name
- Vendor type (Dell, HP, Lenovo)
- Update type (online/offline)
- Creation date and status

No additional flags are required for this command.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogList(cmd.Context())
		},
	}

	firmwareCatalogGetCmd = &cobra.Command{
		Use:     "get firmware_catalog_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific firmware catalog",
		Long: `Get detailed information about a specific firmware catalog.

This command displays comprehensive information about a firmware catalog including:
- Basic metadata (name, description, vendor, creation date)
- Configuration details (update type, server types, vendor systems)
- Repository information for offline catalogs
- Available firmware packages and their versions
- Download/upload status and statistics

Arguments:
  firmware_catalog_id    The ID of the firmware catalog to retrieve

Examples:
  metalcloud-cli firmware-catalog get 12345
  metalcloud-cli fw-catalog show dell-r640-catalog`,
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
		Short:   "Create a new firmware catalog from vendor sources",
		Long: `Create a new firmware catalog from vendor sources.

This command creates a firmware catalog by downloading and processing firmware information
from vendor repositories. The catalog can be configured for online or offline updates.

Configuration Methods:
1. Command-line flags - Specify all options using individual flags
2. Configuration file - Use --config-source to load settings from JSON/YAML file

Required Flags (when not using --config-source):
  --name            Name of the firmware catalog
  --vendor          Vendor type: 'dell', 'hp', or 'lenovo'
  --update-type     Update method: 'online' or 'offline'

Source Configuration (mutually exclusive):
  --vendor-url                 URL of the online vendor catalog
  --vendor-local-catalog-path  Path to a local catalog file

Optional Flags:
  --description                   Description of the firmware catalog
  --vendor-token                  Authentication token for vendor API access
  --server-types                  Comma-separated list of Metalsoft server types to filter
  --vendor-systems                Comma-separated list of vendor system models to filter
  --vendor-local-binaries-path    Local directory for downloaded firmware binaries
  --download-binaries             Download firmware binaries locally
  --upload-binaries               Upload binaries to offline repository

Offline Repository Configuration (required when --upload-binaries is used):
  --repo-base-url        Base URL of the offline repository
  --repo-ssh-host        SSH hostname:port for repository upload
  --repo-ssh-user        SSH username for repository access
  --repo-ssh-path        Target directory path on SSH server

SSH Configuration (mutually exclusive):
  --user-private-key-path    Path to SSH private key (default: ~/.ssh/id_rsa)
  --known-hosts-path         Path to SSH known hosts file (default: ~/.ssh/known_hosts)
  --ignore-host-key-check    Skip SSH host key verification`,
		Example: `
Dell example (online):
metalcloud-cli firmware-catalog create \
  --name "Dell R640 Catalog" \
  --description "Dell PowerEdge R640 firmware catalog" \
  --vendor dell \
  --vendor-url https://downloads.dell.com/FOLDER06417267M/1/ESXi_Catalog.xml.gz \
  --vendor-systems "R640" \
  --server-types "M.24.64.2,M.32.64.2" \
  --update-type online

Dell example (offline with upload):
metalcloud-cli firmware-catalog create \
  --name "Dell Offline Catalog" \
  --vendor dell \
  --vendor-url https://downloads.dell.com/catalog.xml.gz \
  --download-binaries \
  --vendor-local-binaries-path ./downloads \
  --upload-binaries \
  --repo-base-url http://repo.mycloud.com/dell \
  --repo-ssh-host repo.mycloud.com:22 \
  --repo-ssh-user admin \
  --repo-ssh-path /var/www/html/dell \
  --update-type offline

HP example:
metalcloud-cli firmware-catalog create \
  --name "HP Gen11 Catalog" \
  --vendor hp \
  --vendor-url https://downloads.linux.hpe.com/SDR/repo/fwpp-gen11/current/fwrepodata/fwrepo.json \
  --server-types "M.8.8.2.v5" \
  --update-type online

Lenovo example using config file:
metalcloud-cli firmware-catalog create --config-source ./lenovo-config.json

Configuration file examples:

lenovo-config.json:
{
  "name": "Lenovo Catalog",
  "description": "Lenovo server firmware catalog",
  "vendor": "lenovo",
  "update_type": "offline",
  "vendor_local_catalog_path": "./lenovo_catalogs",
  "vendor_local_binaries_path": "./lenovo_downloads",
  "server_types_filter": ["M.8.8.2.v5"],
  "vendor_systems_filter": ["7Y51"],
  "download_binaries": true
}

hp-config.yaml:
name: HP Gen11 Catalog
description: HP ProLiant Gen11 firmware
vendor: hp
update_type: online
vendor_url: https://downloads.linux.hpe.com/SDR/repo/fwpp-gen11/current/fwrepodata/fwrepo.json
vendor_local_catalog_path: ./fwrepo.json
vendor_local_binaries_path: ./hp_downloads
server_types_filter:
  - M.8.8.2.v5`,
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

				err = utils.UnmarshalContent(config, &firmwareCatalogOptions)
				if err != nil {
					return err
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
					RepoSshPath:             firmwareCatalogFlags.repoSshPath,
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
		Use:     "update firmware_catalog_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing firmware catalog",
		Long: `Update an existing firmware catalog.

This command allows you to update the configuration of an existing firmware catalog.
Updates are provided through a configuration file (JSON or YAML format) that contains
the new settings to apply.

The configuration file should contain only the fields you want to update. The catalog
will be refreshed with the new configuration, potentially downloading new firmware
information from vendor sources.

Required Flags:
  --config-source    Source of the configuration updates (JSON/YAML file path or 'pipe')

Arguments:
  firmware_catalog_id    The ID of the firmware catalog to update (optional if provided in config)

Examples:
  metalcloud-cli firmware-catalog update 12345 --config-source ./update-config.json
  cat update-config.json | metalcloud-cli firmware-catalog update --config-source pipe

Configuration file example (update-config.json):
{
  "description": "Updated description",
  "vendor_url": "https://downloads.dell.com/new-catalog.xml.gz",
  "server_types_filter": ["M.24.64.2", "M.32.128.2"],
  "download_binaries": true
}`,
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
		Use:     "delete firmware_catalog_id",
		Aliases: []string{"rm"},
		Short:   "Delete a firmware catalog",
		Long: `Delete a firmware catalog permanently.

This command removes a firmware catalog from the system. This action is irreversible
and will delete all associated firmware package information.

Arguments:
  firmware_catalog_id    The ID of the firmware catalog to delete

Examples:
  metalcloud-cli firmware-catalog delete 12345
  metalcloud-cli fw-catalog rm dell-r640-catalog`,
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
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.repoSshPath, "repo-ssh-path", "", "The path to the target folder in the SSH repository")
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
