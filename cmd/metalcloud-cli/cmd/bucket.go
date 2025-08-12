package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/bucket"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	bucketFlags = struct {
		configSource string
		filterStatus []string
	}{}

	bucketCmd = &cobra.Command{
		Use:     "bucket [command]",
		Aliases: []string{"buckets", "s3"},
		Short:   "Manage S3-compatible object storage buckets",
		Long: `Comprehensive bucket management commands for S3-compatible object storage.

This command group provides complete lifecycle management for buckets including:
- Creating and deleting buckets
- Listing buckets with filtering capabilities
- Retrieving bucket details and configuration
- Updating bucket configuration and metadata
- Managing bucket credentials

All commands require appropriate permissions and operate within the context of an infrastructure.`,
	}

	bucketListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all buckets in an infrastructure",
		Long: `List all buckets within a specified infrastructure with optional filtering capabilities.

This command displays all buckets associated with the given infrastructure, showing their 
status, configuration, and metadata. Results can be filtered by bucket status to focus 
on specific states.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure

Optional Flags:
  --filter-status strings      Filter buckets by their status (can be repeated for multiple statuses)

Examples:
  # List all buckets in infrastructure with ID 100
  metalcloud-cli bucket list 100

  # List all buckets in infrastructure labeled "production"
  metalcloud-cli bucket list production

  # List only active buckets
  metalcloud-cli bucket list 100 --filter-status active

  # List buckets with multiple status filters
  metalcloud-cli bucket list production --filter-status active --filter-status pending`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketList(cmd.Context(), args[0], bucketFlags.filterStatus)
		},
	}

	bucketGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label bucket_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific bucket",
		Long: `Retrieve comprehensive details about a specific bucket within an infrastructure.

This command displays detailed information about a bucket including its configuration,
status, metadata, and other properties. The output provides complete visibility into
the bucket's current state and settings.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  bucket_id                     The unique identifier of the bucket

Examples:
  # Get details for bucket with ID 42 in infrastructure 100
  metalcloud-cli bucket get 100 42

  # Get details using infrastructure label
  metalcloud-cli bucket get production bucket-abc123

  # Display bucket information
  metalcloud-cli bucket show staging my-bucket-id`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketGet(cmd.Context(), args[0], args[1])
		},
	}

	bucketCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new bucket with specified configuration",
		Long: `Create a new bucket within the specified infrastructure using provided configuration.

This command creates a new bucket with the configuration specified through either a JSON file
or piped input. The configuration must include all required bucket parameters such as name,
region, and any specific bucket settings.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure where the bucket will be created

Required Flags:
  --config-source string       Source of the new bucket configuration
                               Accepts either 'pipe' for piped JSON input or a path to a JSON file

Examples:
  # Create bucket using configuration from a JSON file
  metalcloud-cli bucket create 100 --config-source bucket-config.json

  # Create bucket using piped configuration
  echo '{"name": "my-bucket", "region": "us-east-1"}' | metalcloud-cli bucket create production --config-source pipe

  # Create bucket with configuration file in different directory
  metalcloud-cli bucket create staging --config-source /path/to/configs/bucket.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(bucketFlags.configSource)
			if err != nil {
				return err
			}

			return bucket.BucketCreate(cmd.Context(), args[0], config)
		},
	}

	bucketDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label bucket_id",
		Aliases: []string{"rm"},
		Short:   "Delete a bucket and all its contents",
		Long: `Permanently delete a bucket and all its contents from the specified infrastructure.

This command removes a bucket and all objects stored within it. This operation is 
irreversible and will permanently destroy all data in the bucket. Use with caution 
in production environments.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket to delete

Examples:
  # Delete bucket with ID 42 from infrastructure 100
  metalcloud-cli bucket delete 100 42

  # Delete bucket using infrastructure label
  metalcloud-cli bucket delete production bucket-abc123

  # Remove bucket using alias
  metalcloud-cli bucket rm staging my-bucket-id`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketDelete(cmd.Context(), args[0], args[1])
		},
	}

	bucketUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label bucket_id",
		Aliases: []string{"config-update"},
		Short:   "Update bucket configuration with new settings",
		Long: `Update the configuration of an existing bucket with new settings or modifications.

This command allows you to modify bucket configuration parameters such as access policies,
storage settings, or other configurable properties. The new configuration is provided
through either a JSON file or piped input containing the updated settings.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket to update

Required Flags:
  --config-source string       Source of the bucket configuration updates
                               Accepts either 'pipe' for piped JSON input or a path to a JSON file

Examples:
  # Update bucket configuration from a JSON file
  metalcloud-cli bucket update-config 100 42 --config-source new-config.json

  # Update configuration using piped input
  echo '{"policy": "public-read"}' | metalcloud-cli bucket update-config production bucket-123 --config-source pipe

  # Update bucket with configuration file
  metalcloud-cli bucket config-update staging my-bucket --config-source /configs/update.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(bucketFlags.configSource)
			if err != nil {
				return err
			}

			return bucket.BucketUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	bucketUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label bucket_id",
		Aliases: []string{"meta-update"},
		Short:   "Update bucket metadata and custom properties",
		Long: `Update the metadata and custom properties of an existing bucket.

This command allows you to modify bucket metadata such as labels, descriptions, 
custom tags, and other non-configuration properties. The metadata updates are 
provided through either a JSON file or piped input.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket to update

Required Flags:
  --config-source string       Source of the bucket metadata updates
                               Accepts either 'pipe' for piped JSON input or a path to a JSON file

Examples:
  # Update bucket metadata from a JSON file
  metalcloud-cli bucket update-meta 100 42 --config-source metadata.json

  # Update metadata using piped input
  echo '{"label": "production-storage", "description": "Main storage bucket"}' | metalcloud-cli bucket update-meta production bucket-123 --config-source pipe

  # Update bucket metadata with file
  metalcloud-cli bucket meta-update staging my-bucket --config-source /configs/meta.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(bucketFlags.configSource)
			if err != nil {
				return err
			}

			return bucket.BucketUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	bucketGetConfigInfoCmd = &cobra.Command{
		Use:     "config-info infrastructure_id_or_label bucket_id",
		Aliases: []string{"get-config-info"},
		Short:   "Get configuration information for a bucket",
		Long: `Retrieve detailed configuration information for a specific bucket.

This command displays the current configuration settings of a bucket including
access policies, storage parameters, and other configuration details. This is
useful for reviewing current settings before making updates.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket

Examples:
  # Get configuration info for bucket with ID 42
  metalcloud-cli bucket config-info 100 42

  # Get config info using infrastructure label
  metalcloud-cli bucket config-info production bucket-abc123

  # Display configuration information using alias
  metalcloud-cli bucket get-config-info staging my-bucket-id`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}

	bucketGetCredentialsCmd = &cobra.Command{
		Use:     "get-credentials infrastructure_id_or_label bucket_id",
		Aliases: []string{"credentials"},
		Short:   "Get access credentials for a bucket",
		Long: `Retrieve access credentials for a specific bucket within an infrastructure.

This command displays the credentials required to access the bucket programmatically,
including access keys, secrets, and endpoint information. These credentials can be
used with S3-compatible tools and SDKs to interact with the bucket.

Required Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure containing the bucket
  bucket_id                     The unique identifier of the bucket

Examples:
  # Get credentials for bucket with ID 42
  metalcloud-cli bucket get-credentials 100 42

  # Get credentials using infrastructure label
  metalcloud-cli bucket get-credentials production bucket-abc123

  # Display credentials using alias
  metalcloud-cli bucket credentials staging my-bucket-id`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_BUCKETS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketGetCredentials(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(bucketCmd)

	bucketCmd.AddCommand(bucketListCmd)
	bucketListCmd.Flags().StringSliceVar(&bucketFlags.filterStatus, "filter-status", nil, "Filter the result by bucket status.")

	bucketCmd.AddCommand(bucketGetCmd)

	bucketCmd.AddCommand(bucketCreateCmd)
	bucketCreateCmd.Flags().StringVar(&bucketFlags.configSource, "config-source", "", "Source of the new bucket configuration. Can be 'pipe' or path to a JSON file.")
	bucketCreateCmd.MarkFlagsOneRequired("config-source")

	bucketCmd.AddCommand(bucketDeleteCmd)

	bucketCmd.AddCommand(bucketUpdateConfigCmd)
	bucketUpdateConfigCmd.Flags().StringVar(&bucketFlags.configSource, "config-source", "", "Source of the bucket configuration updates. Can be 'pipe' or path to a JSON file.")
	bucketUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	bucketCmd.AddCommand(bucketUpdateMetaCmd)
	bucketUpdateMetaCmd.Flags().StringVar(&bucketFlags.configSource, "config-source", "", "Source of the bucket metadata updates. Can be 'pipe' or path to a JSON file.")
	bucketUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	bucketCmd.AddCommand(bucketGetConfigInfoCmd)

	bucketCmd.AddCommand(bucketGetCredentialsCmd)
}
