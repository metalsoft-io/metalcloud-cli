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
		Short:   "Storage management",
		Long:    `Storage commands.`,
	}

	storageListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all storages.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageList(cmd.Context(), storageFlags.filterTechnology)
		},
	}

	storageGetCmd = &cobra.Command{
		Use:          "get storage_id",
		Aliases:      []string{"show"},
		Short:        "Get storage details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageGet(cmd.Context(), args[0])
		},
	}

	storageConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get storage configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageConfigExample(cmd.Context())
		},
	}

	storageCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a new storage.",
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
		Use:          "delete storage_id",
		Short:        "Delete a storage.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageDelete(cmd.Context(), args[0])
		},
	}

	storageGetCredentialsCmd = &cobra.Command{
		Use:          "credentials storage_id",
		Short:        "Get storage credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.StorageGetCredentials(cmd.Context(), args[0])
		},
	}

	storageGetDrivesCmd = &cobra.Command{
		Use:          "drives storage_id",
		Short:        "Get drives for a storage.",
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
		Use:          "file-shares storage_id",
		Short:        "Get file shares for a storage.",
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
		Use:          "buckets storage_id",
		Short:        "Get buckets for a storage.",
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
		Use:          "network-configs storage_id",
		Short:        "Get network device configurations for a storage.",
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
