package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/file_share"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	fileShareFlags = struct {
		configSource string
		filterStatus []string
	}{}

	fileShareCmd = &cobra.Command{
		Use:     "file-share [command]",
		Aliases: []string{"fs"},
		Short:   "Manage file shares for infrastructure resources",
		Long: `Manage file shares for infrastructure resources including creating, updating, deleting, and configuring shared storage.

File shares provide shared storage capabilities across multiple instances within an infrastructure.
Use these commands to manage the lifecycle and configuration of file shares.

Available Commands:
  list           List all file shares for an infrastructure
  get            Get detailed information about a specific file share
  create         Create a new file share with specified configuration
  delete         Delete an existing file share
  update-config  Update file share configuration
  update-meta    Update file share metadata
  get-hosts      Get hosts configured for a file share
  update-hosts   Update hosts configuration for a file share
  config-info    Get configuration information for a file share

Examples:
  # List all file shares for an infrastructure
  metalcloud-cli file-share list my-infrastructure

  # Get details of a specific file share
  metalcloud-cli file-share get my-infrastructure 12345

  # Create a new file share from a configuration file
  metalcloud-cli file-share create my-infrastructure --config-source config.json`,
	}

	fileShareListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all file shares for an infrastructure",
		Long: `List all file shares associated with the specified infrastructure.

This command displays file shares with their basic information including ID, name, 
status, and other key attributes. Results can be filtered by status.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label to list file shares for

Optional Flags:
  --filter-status               Filter results by file share status (can be used multiple times)
                               Common statuses: active, inactive, creating, deleting, error

Examples:
  # List all file shares for an infrastructure
  metalcloud-cli file-share list my-infrastructure

  # List file shares with ID
  metalcloud-cli file-share list 12345

  # Filter by status
  metalcloud-cli file-share list my-infrastructure --filter-status active
  
  # Filter by multiple statuses
  metalcloud-cli file-share list my-infrastructure --filter-status active --filter-status creating`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareList(cmd.Context(), args[0], fileShareFlags.filterStatus)
		},
	}

	fileShareGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label file_share_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific file share",
		Long: `Get detailed information about a specific file share including configuration, 
status, hosts, and metadata.

This command displays comprehensive information about a file share such as its
current configuration, operational status, associated hosts, and metadata.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to retrieve information for

Examples:
  # Get details of a file share by infrastructure label and file share ID
  metalcloud-cli file-share get my-infrastructure 12345

  # Get details using infrastructure ID
  metalcloud-cli file-share get 100 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGet(cmd.Context(), args[0], args[1])
		},
	}

	fileShareCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new file share with specified configuration",
		Long: `Create a new file share within the specified infrastructure using configuration 
from a JSON file or pipe input.

This command creates a new file share resource with the specified configuration
including storage size, access permissions, and other settings.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label to create the file share in

Required Flags:
  --config-source               Source of the new file share configuration
                               Can be 'pipe' for stdin input or path to a JSON file

Configuration Format:
The configuration should be a JSON object containing file share properties such as:
- name: File share name
- size: Storage size
- type: File share type
- access_mode: Access permissions
- mount_path: Mount path for clients

Examples:
  # Create a file share from a JSON configuration file
  metalcloud-cli file-share create my-infrastructure --config-source config.json

  # Create a file share using pipe input
  echo '{"name":"shared-storage","size":"100GB"}' | metalcloud-cli file-share create my-infrastructure --config-source pipe

  # Create with infrastructure ID
  metalcloud-cli file-share create 12345 --config-source /path/to/config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareCreate(cmd.Context(), args[0], config)
		},
	}

	fileShareDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label file_share_id",
		Aliases: []string{"rm"},
		Short:   "Delete an existing file share",
		Long: `Delete an existing file share from the specified infrastructure.

This command permanently removes a file share and all its associated data.
Use with caution as this operation cannot be undone.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to delete

Examples:
  # Delete a file share by infrastructure label and file share ID
  metalcloud-cli file-share delete my-infrastructure 12345

  # Delete using infrastructure ID
  metalcloud-cli file-share delete 100 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareDelete(cmd.Context(), args[0], args[1])
		},
	}

	fileShareUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label file_share_id",
		Aliases: []string{"config-update"},
		Short:   "Update file share configuration",
		Long: `Update the configuration of an existing file share with new settings.

This command allows you to modify various configuration parameters of a file share
such as storage size, access permissions, mount options, and other settings.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to update

Required Flags:
  --config-source               Source of the file share configuration updates
                               Can be 'pipe' for stdin input or path to a JSON file

Configuration Format:
The configuration should be a JSON object containing the file share properties to update:
- name: File share name
- size: Storage size (if supported for expansion)
- access_mode: Access permissions
- mount_options: Mount configuration options
- description: File share description

Examples:
  # Update file share configuration from a JSON file
  metalcloud-cli file-share update-config my-infrastructure 12345 --config-source config.json

  # Update using pipe input
  echo '{"description":"Updated shared storage"}' | metalcloud-cli file-share update-config my-infrastructure 12345 --config-source pipe

  # Update with infrastructure ID
  metalcloud-cli file-share update-config 100 12345 --config-source /path/to/config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label file_share_id",
		Aliases: []string{"meta-update"},
		Short:   "Update file share metadata",
		Long: `Update the metadata of an existing file share with new information.

This command allows you to modify metadata properties of a file share such as
labels, tags, descriptions, and other custom metadata fields.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to update

Required Flags:
  --config-source               Source of the file share metadata updates
                               Can be 'pipe' for stdin input or path to a JSON file

Metadata Format:
The configuration should be a JSON object containing the metadata properties to update:
- labels: Key-value pairs for labeling
- tags: Array of tags for categorization
- description: Detailed description of the file share
- custom_fields: Custom metadata fields

Examples:
  # Update file share metadata from a JSON file
  metalcloud-cli file-share update-meta my-infrastructure 12345 --config-source metadata.json

  # Update using pipe input
  echo '{"labels":{"env":"production","team":"devops"}}' | metalcloud-cli file-share update-meta my-infrastructure 12345 --config-source pipe

  # Update with infrastructure ID
  metalcloud-cli file-share update-meta 100 12345 --config-source /path/to/metadata.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareGetHostsCmd = &cobra.Command{
		Use:     "get-hosts infrastructure_id_or_label file_share_id",
		Aliases: []string{"hosts"},
		Short:   "Get hosts configured for a file share",
		Long: `Get the list of hosts that are configured to access the specified file share.

This command displays all hosts that have been granted access to the file share,
including their mount configurations and access permissions.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to get hosts for

Examples:
  # Get hosts for a file share by infrastructure label and file share ID
  metalcloud-cli file-share get-hosts my-infrastructure 12345

  # Get hosts using infrastructure ID
  metalcloud-cli file-share get-hosts 100 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGetHosts(cmd.Context(), args[0], args[1])
		},
	}

	fileShareUpdateHostsCmd = &cobra.Command{
		Use:     "update-hosts infrastructure_id_or_label file_share_id",
		Aliases: []string{"hosts-update"},
		Short:   "Update hosts configuration for a file share",
		Long: `Update the hosts configuration for an existing file share.

This command allows you to modify which hosts have access to the file share
and their mount configurations. You can add new hosts, remove existing ones,
or update their access permissions and mount settings.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to update hosts for

Required Flags:
  --config-source               Source of the file share hosts configuration
                               Can be 'pipe' for stdin input or path to a JSON file

Hosts Configuration Format:
The configuration should be a JSON object containing the hosts configuration:
- hosts: Array of host objects with their access settings
- mount_path: Default mount path for hosts
- access_permissions: Default access permissions (read, write, read-write)

Examples:
  # Update hosts configuration from a JSON file
  metalcloud-cli file-share update-hosts my-infrastructure 12345 --config-source hosts.json

  # Update using pipe input
  echo '{"hosts":[{"id":"host1","access":"read-write"}]}' | metalcloud-cli file-share update-hosts my-infrastructure 12345 --config-source pipe

  # Update with infrastructure ID
  metalcloud-cli file-share update-hosts 100 12345 --config-source /path/to/hosts.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateHosts(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareGetConfigInfoCmd = &cobra.Command{
		Use:     "config-info infrastructure_id_or_label file_share_id",
		Aliases: []string{"get-config-info"},
		Short:   "Get configuration information for a file share",
		Long: `Get configuration information for a specific file share including technical details,
settings, and current configuration state.

This command displays the complete configuration profile of a file share including
storage configuration, networking settings, access control, and other technical
parameters that may be needed for troubleshooting or integration purposes.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to get configuration information for

Examples:
  # Get configuration information for a file share
  metalcloud-cli file-share config-info my-infrastructure 12345

  # Get configuration info using infrastructure ID
  metalcloud-cli file-share config-info 100 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(fileShareCmd)

	fileShareCmd.AddCommand(fileShareListCmd)
	fileShareListCmd.Flags().StringSliceVar(&fileShareFlags.filterStatus, "filter-status", nil, "Filter the result by file share status.")

	fileShareCmd.AddCommand(fileShareGetCmd)

	fileShareCmd.AddCommand(fileShareCreateCmd)
	fileShareCreateCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the new file share configuration. Can be 'pipe' or path to a JSON file.")
	fileShareCreateCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareDeleteCmd)

	fileShareCmd.AddCommand(fileShareUpdateConfigCmd)
	fileShareUpdateConfigCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share configuration updates. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareUpdateMetaCmd)
	fileShareUpdateMetaCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share metadata updates. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareGetHostsCmd)

	fileShareCmd.AddCommand(fileShareUpdateHostsCmd)
	fileShareUpdateHostsCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share hosts configuration. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateHostsCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareGetConfigInfoCmd)
}
