package cmd

import (
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/storage"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	storageFlags = struct {
		filterTechnology []string
		configSource     string
		limit            string
		page             string
	}{}

	storageCmd = &cobra.Command{
		Use:     "storage [command]",
		Aliases: []string{"storage-pool"},
		Short:   "Manage storage pools and related resources",
		Long: `Manage storage pools and their associated resources in the MetalCloud infrastructure.

Storage pools are external storage systems that provide block, file, or object storage
to instances. This command group allows you to create, configure, and manage storage
pools, as well as access their drives, file shares, buckets, and network configurations.

Available commands:
  list             List all storage pools
  get              Get detailed information about a specific storage pool
  create           Create a new storage pool
  delete           Delete an existing storage pool
  config-example   Display a configuration template for creating storage pools
  credentials      Retrieve credentials for a storage pool
  drives           List drives available in a storage pool
  file-shares      List file shares in a storage pool
  buckets          List object storage buckets in a storage pool
  network-configs  List network device configurations for a storage pool

Use "metalcloud storage [command] --help" for more information about a command.`,
	}

	storageListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all storage pools",
		Long: `List all storage pools with optional filtering.

This command displays information about storage pools including their ID, site, driver,
technology, type, name, and status. The output can be filtered by storage technology.

Flags:
  --filter-technology strings   Filter results by storage technology (e.g., block, file, object)

Examples:
  # List all storage pools
  metalcloud storage list

  # List only block storage pools
  metalcloud storage list --filter-technology block

  # List multiple storage types
  metalcloud storage list --filter-technology block,file`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageList(cmd.Context(), storageFlags.filterTechnology)
		},
	}

	storageGetCmd = &cobra.Command{
		Use:     "get storage_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific storage pool",
		Long: `Get detailed information about a specific storage pool by its ID.

This command displays comprehensive information about a storage pool including its
configuration, status, driver details, technologies, and associated metadata.

Arguments:
  storage_id    The numeric ID of the storage pool to retrieve

Examples:
  # Get details for storage pool with ID 123
  metalcloud storage get 123

  # Using the show alias
  metalcloud storage show 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageGet(cmd.Context(), args[0])
		},
	}

	storageConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display a configuration template for creating storage pools",
		Long: `Display a configuration template for creating storage pools.

This command outputs a JSON template showing all available configuration options
for creating a storage pool. The template includes required fields, optional fields,
and example values to help you create valid storage configurations.

The output can be used as a starting point for creating storage pool configurations
that can be passed to the 'create' command via the --config-source flag.

Examples:
  # Display the configuration template
  metalcloud storage config-example

  # Save the template to a file for editing
  metalcloud storage config-example > storage-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageConfigExample(cmd.Context())
		},
	}

	storageCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new storage pool",
		Long: `Create a new storage pool using a configuration file or piped input.

This command creates a new storage pool in the MetalCloud infrastructure. The storage
configuration must be provided as JSON either through a file or piped input.

Required flags:
  --config-source    Source of the storage configuration. Can be 'pipe' for piped input
                     or a path to a JSON file containing the storage configuration.

The configuration must include required fields such as siteId, driver, technologies,
type, name, managementHost, username, password, and subnetType. Use the 'config-example'
command to see a complete template with all available options.

Examples:
  # Create storage from a JSON file
  metalcloud storage create --config-source ./storage-config.json

  # Create storage from piped input
  cat storage-config.json | metalcloud storage create --config-source pipe

  # Generate template, edit, and create
  metalcloud storage config-example > config.json
  # Edit config.json with your storage details
  metalcloud storage create --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(storageFlags.configSource)
			if err != nil {
				return err
			}

			return storage.StorageCreate(cmd.Context(), config)
		},
	}

	storageDeleteCmd = &cobra.Command{
		Use:   "delete storage_id",
		Short: "Delete an existing storage pool",
		Long: `Delete an existing storage pool by its ID.

This command permanently removes a storage pool from the MetalCloud infrastructure.
Warning: This action is irreversible and will remove all associated data.

Arguments:
  storage_id    The numeric ID of the storage pool to delete

Examples:
  # Delete storage pool with ID 123
  metalcloud storage delete 123

  # Delete storage pool with confirmation
  metalcloud storage delete 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageDelete(cmd.Context(), args[0])
		},
	}

	storageGetCredentialsCmd = &cobra.Command{
		Use:   "credentials storage_id",
		Short: "Retrieve credentials for a storage pool",
		Long: `Retrieve authentication credentials for a specific storage pool.

This command returns the credentials required to connect to and manage the storage
pool. The credentials typically include connection information, authentication tokens,
or access keys depending on the storage driver type.

Arguments:
  storage_id    The numeric ID of the storage pool

Examples:
  # Get credentials for storage pool with ID 123
  metalcloud storage credentials 123

  # Get credentials and save to file
  metalcloud storage credentials 456 > storage-creds.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageGetCredentials(cmd.Context(), args[0])
		},
	}

	storageGetDrivesCmd = &cobra.Command{
		Use:   "drives storage_id",
		Short: "List drives available in a storage pool",
		Long: `List drives available in a specific storage pool.

This command retrieves all drives associated with a storage pool, showing their
configuration, status, and specifications. Results can be paginated using the
limit and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all drives for storage pool 123
  metalcloud storage drives 123

  # List first 10 drives
  metalcloud storage drives 123 --limit 10

  # List second page with 10 drives per page
  metalcloud storage drives 123 --limit 10 --page 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if storageFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(storageFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if storageFlags.page != "" {
				pageVal, err := strconv.ParseFloat(storageFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return storage.StorageGetDrives(cmd.Context(), args[0], limit, page)
		},
	}

	storageGetFileSharesCmd = &cobra.Command{
		Use:   "file-shares storage_id",
		Short: "List file shares in a storage pool",
		Long: `List file shares available in a specific storage pool.

This command retrieves all file shares associated with a storage pool, showing their
configuration, status, and access information. File shares are typically used for
NFS or CIFS/SMB file storage. Results can be paginated using the limit and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all file shares for storage pool 123
  metalcloud storage file-shares 123

  # List first 5 file shares
  metalcloud storage file-shares 123 --limit 5

  # List third page with 5 file shares per page
  metalcloud storage file-shares 123 --limit 5 --page 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if storageFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(storageFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if storageFlags.page != "" {
				pageVal, err := strconv.ParseFloat(storageFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return storage.StorageGetFileShares(cmd.Context(), args[0], limit, page)
		},
	}

	storageGetBucketsCmd = &cobra.Command{
		Use:   "buckets storage_id",
		Short: "List object storage buckets in a storage pool",
		Long: `List object storage buckets available in a specific storage pool.

This command retrieves all object storage buckets associated with a storage pool,
showing their configuration, status, and access information. Buckets are typically
used for S3-compatible object storage. Results can be paginated using the limit
and page flags.

Arguments:
  storage_id    The numeric ID of the storage pool

Optional flags:
  --limit       Number of records per page (default: all records)
  --page        Page number for pagination (requires --limit)

Examples:
  # List all buckets for storage pool 123
  metalcloud storage buckets 123

  # List first 10 buckets
  metalcloud storage buckets 123 --limit 10

  # List second page with 10 buckets per page
  metalcloud storage buckets 123 --limit 10 --page 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if storageFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(storageFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if storageFlags.page != "" {
				pageVal, err := strconv.ParseFloat(storageFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return storage.StorageGetBuckets(cmd.Context(), args[0], limit, page)
		},
	}

	storageGetNetworkDeviceConfigurationsCmd = &cobra.Command{
		Use:   "network-configs storage_id",
		Short: "List network device configurations for a storage pool",
		Long: `List network device configurations for a specific storage pool.

This command retrieves network device configurations associated with a storage pool,
showing network settings, device mappings, and connectivity information.

Arguments:
  storage_id    The numeric ID of the storage pool

Examples:
  # List network configurations for storage pool 123
  metalcloud storage network-configs 123

  # Save network configurations to file
  metalcloud storage network-configs 456 > network-configs.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageGetNetworkDeviceConfigurations(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(storageCmd)

	// Add existing commands
	storageCmd.AddCommand(storageListCmd)
	storageListCmd.Flags().StringSliceVar(&storageFlags.filterTechnology, "filter-technology", nil, "Filter the result by storage technology.")

	storageCmd.AddCommand(storageGetCmd)

	storageCmd.AddCommand(storageConfigExampleCmd)

	// Add new commands
	storageCmd.AddCommand(storageCreateCmd)
	storageCreateCmd.Flags().StringVar(&storageFlags.configSource, "config-source", "", "Source of the new storage configuration. Can be 'pipe' or path to a JSON file.")
	storageCreateCmd.MarkFlagsOneRequired("config-source")

	storageCmd.AddCommand(storageDeleteCmd)

	storageCmd.AddCommand(storageGetCredentialsCmd)

	// Add commands for retrieving related resources
	storageCmd.AddCommand(storageGetDrivesCmd)
	storageGetDrivesCmd.Flags().StringVar(&storageFlags.limit, "limit", "", "Number of records per page")
	storageGetDrivesCmd.Flags().StringVar(&storageFlags.page, "page", "", "Page number")

	storageCmd.AddCommand(storageGetFileSharesCmd)
	storageGetFileSharesCmd.Flags().StringVar(&storageFlags.limit, "limit", "", "Number of records per page")
	storageGetFileSharesCmd.Flags().StringVar(&storageFlags.page, "page", "", "Page number")

	storageCmd.AddCommand(storageGetBucketsCmd)
	storageGetBucketsCmd.Flags().StringVar(&storageFlags.limit, "limit", "", "Number of records per page")
	storageGetBucketsCmd.Flags().StringVar(&storageFlags.page, "page", "", "Page number")

	storageCmd.AddCommand(storageGetNetworkDeviceConfigurationsCmd)
}
