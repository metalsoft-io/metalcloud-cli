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
		Short:   "Bucket management",
		Long:    `Bucket management commands for S3-compatible object storage.`,
	}

	bucketListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all buckets for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketList(cmd.Context(), args[0], bucketFlags.filterStatus)
		},
	}

	bucketGetCmd = &cobra.Command{
		Use:          "get infrastructure_id_or_label bucket_id",
		Aliases:      []string{"show"},
		Short:        "Get bucket details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketGet(cmd.Context(), args[0], args[1])
		},
	}

	bucketCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label",
		Aliases:      []string{"new"},
		Short:        "Create a new bucket.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
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
		Use:          "delete infrastructure_id_or_label bucket_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a bucket.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketDelete(cmd.Context(), args[0], args[1])
		},
	}

	bucketUpdateConfigCmd = &cobra.Command{
		Use:          "update-config infrastructure_id_or_label bucket_id",
		Aliases:      []string{"config-update"},
		Short:        "Update bucket configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
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
		Use:          "update-meta infrastructure_id_or_label bucket_id",
		Aliases:      []string{"meta-update"},
		Short:        "Update bucket metadata.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_WRITE},
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
		Use:          "config-info infrastructure_id_or_label bucket_id",
		Aliases:      []string{"get-config-info"},
		Short:        "Get configuration information for a bucket.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return bucket.BucketGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}

	bucketGetCredentialsCmd = &cobra.Command{
		Use:          "get-credentials infrastructure_id_or_label bucket_id",
		Aliases:      []string{"credentials"},
		Short:        "Get credentials for a bucket.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
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
