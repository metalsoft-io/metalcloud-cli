package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/custom_iso"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	customIsoFlags = struct {
		configSource string
	}{}

	customIsoCmd = &cobra.Command{
		Use:     "custom-iso [command]",
		Aliases: []string{"iso", "isos"},
		Short:   "Manage custom ISO images for server provisioning",
		Long: `Manage custom ISO images that can be used for server provisioning and booting.

Custom ISOs allow you to create bootable images with custom operating systems, 
tools, or configurations that can be mounted and used to boot servers in your infrastructure.

Available commands:
  list           List all custom ISOs
  get            Get details of a specific custom ISO
  config-example Show configuration example for creating custom ISOs
  create         Create a new custom ISO from configuration
  update         Update an existing custom ISO
  delete         Delete a custom ISO
  make-public    Make a custom ISO available to all users
  boot-server    Boot a server using a custom ISO

Examples:
  # List all custom ISOs
  metalcloud-cli custom-iso list

  # Get details of a specific custom ISO
  metalcloud-cli custom-iso get 12345

  # Create a new custom ISO from a JSON file
  metalcloud-cli custom-iso create --config-source config.json

  # Boot a server with a custom ISO
  metalcloud-cli custom-iso boot-server 12345 67890`,
	}

	customIsoListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available custom ISO images",
		Long: `List all custom ISO images available in your account.

This command displays a table of all custom ISOs showing their ID, name, 
description, size, creation date, and availability status.

Required permissions:
  - custom_iso:read

Examples:
  # List all custom ISOs
  metalcloud-cli custom-iso list
  
  # List custom ISOs with shorter alias
  metalcloud-cli iso ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoList(cmd.Context())
		},
	}

	customIsoGetCmd = &cobra.Command{
		Use:     "get custom_iso_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific custom ISO",
		Long: `Get detailed information about a specific custom ISO including its configuration,
metadata, size, creation date, and current status.

Arguments:
  custom_iso_id   ID of the custom ISO to retrieve (required)

Required permissions:
  - custom_iso:read

Examples:
  # Get details of custom ISO with ID 12345
  metalcloud-cli custom-iso get 12345
  
  # Show custom ISO details using alias
  metalcloud-cli iso show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoGet(cmd.Context(), args[0])
		},
	}

	customIsoConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show configuration example for creating custom ISOs",
		Long: `Display a JSON configuration example that can be used as a template for creating
or updating custom ISO images.

The configuration example shows all available fields, their data types, and
expected values. You can save this output to a file and modify it according
to your requirements.

Required permissions:
  - custom_iso:write

Examples:
  # Show configuration example
  metalcloud-cli custom-iso config-example
  
  # Save example to file for editing
  metalcloud-cli custom-iso config-example > custom-iso-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoConfigExample(cmd.Context())
		},
	}

	customIsoCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new custom ISO from configuration",
		Long: `Create a new custom ISO image using a JSON configuration file or piped input.

The configuration must include all required fields such as name, description,
and ISO source. Use the config-example command to see the expected format.

Required flags:
  --config-source   Source of the configuration (required)
                    Can be 'pipe' to read from stdin or path to a JSON file

Required permissions:
  - custom_iso:write

Dependencies:
  - Valid JSON configuration matching the expected schema
  - Accessible ISO source if specified in configuration

Examples:
  # Create custom ISO from a JSON file
  metalcloud-cli custom-iso create --config-source config.json
  
  # Create custom ISO from piped JSON
  cat config.json | metalcloud-cli custom-iso create --config-source pipe
  
  # Create using shorter alias
  metalcloud-cli iso new --config-source my-iso-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(customIsoFlags.configSource)
			if err != nil {
				return err
			}

			return custom_iso.CustomIsoCreate(cmd.Context(), config)
		},
	}

	customIsoUpdateCmd = &cobra.Command{
		Use:     "update custom_iso_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing custom ISO with new configuration",
		Long: `Update an existing custom ISO image using a JSON configuration file or piped input.

This command allows you to modify properties of an existing custom ISO such as
name, description, or other metadata. The configuration should contain only the
fields you want to update.

Arguments:
  custom_iso_id   ID of the custom ISO to update (required)

Required flags:
  --config-source   Source of the configuration updates (required)
                    Can be 'pipe' to read from stdin or path to a JSON file

Required permissions:
  - custom_iso:write

Dependencies:
  - Valid JSON configuration with fields to update
  - Custom ISO must exist and be accessible

Examples:
  # Update custom ISO from a JSON file
  metalcloud-cli custom-iso update 12345 --config-source updates.json
  
  # Update custom ISO from piped JSON
  echo '{"name":"New Name"}' | metalcloud-cli custom-iso update 12345 --config-source pipe
  
  # Update using shorter alias
  metalcloud-cli iso edit 12345 --config-source config-updates.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(customIsoFlags.configSource)
			if err != nil {
				return err
			}

			return custom_iso.CustomIsoUpdate(cmd.Context(), args[0], config)
		},
	}

	customIsoDeleteCmd = &cobra.Command{
		Use:     "delete custom_iso_id",
		Aliases: []string{"rm"},
		Short:   "Delete a custom ISO permanently",
		Long: `Delete a custom ISO image permanently from your account.

This action cannot be undone. The custom ISO will be removed from all servers
where it might be mounted and will no longer be available for provisioning.

Arguments:
  custom_iso_id   ID of the custom ISO to delete (required)

Required permissions:
  - custom_iso:write

Dependencies:
  - Custom ISO must exist and be accessible
  - Custom ISO should not be actively used by running servers

Examples:
  # Delete custom ISO with ID 12345
  metalcloud-cli custom-iso delete 12345
  
  # Delete using shorter alias
  metalcloud-cli iso rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoDelete(cmd.Context(), args[0])
		},
	}

	customIsoMakePublicCmd = &cobra.Command{
		Use:   "make-public custom_iso_id",
		Short: "Make a custom ISO available to all users",
		Long: `Make a custom ISO image publicly available to all users in the organization.

This command changes the visibility of a custom ISO from private (accessible only
to the owner) to public (accessible to all users with appropriate permissions).

Arguments:
  custom_iso_id   ID of the custom ISO to make public (required)

Required permissions:
  - custom_iso:write

Dependencies:
  - Custom ISO must exist and be accessible
  - User must be the owner of the custom ISO or have admin privileges

Examples:
  # Make custom ISO with ID 12345 public
  metalcloud-cli custom-iso make-public 12345
  
  # Make ISO public using shorter alias
  metalcloud-cli iso make-public 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoMakePublic(cmd.Context(), args[0])
		},
	}

	customIsoBootServerCmd = &cobra.Command{
		Use:     "boot-server custom_iso_id server_id",
		Aliases: []string{"boot"},
		Short:   "Boot a server using a custom ISO",
		Long: `Boot a server using a custom ISO image.

This command mounts the specified custom ISO to the server and initiates a boot
process. The server will boot from the custom ISO, allowing you to use custom
operating systems or recovery tools.

Arguments:
  custom_iso_id   ID of the custom ISO to use for booting (required)
  server_id       ID of the server to boot (required)

Required permissions:
  - server_instances:write

Dependencies:
  - Custom ISO must exist and be accessible
  - Server must exist and be accessible
  - Server should be in a powered-off state for optimal results

Examples:
  # Boot server 67890 using custom ISO 12345
  metalcloud-cli custom-iso boot-server 12345 67890
  
  # Boot using shorter alias
  metalcloud-cli iso boot 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoBootServer(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(customIsoCmd)

	customIsoCmd.AddCommand(customIsoListCmd)
	customIsoCmd.AddCommand(customIsoGetCmd)
	customIsoCmd.AddCommand(customIsoConfigExampleCmd)

	customIsoCmd.AddCommand(customIsoCreateCmd)
	customIsoCreateCmd.Flags().StringVar(&customIsoFlags.configSource, "config-source", "", "Source of the new custom ISO configuration. Can be 'pipe' or path to a JSON file.")
	customIsoCreateCmd.MarkFlagsOneRequired("config-source")

	customIsoCmd.AddCommand(customIsoUpdateCmd)
	customIsoUpdateCmd.Flags().StringVar(&customIsoFlags.configSource, "config-source", "", "Source of the custom ISO configuration updates. Can be 'pipe' or path to a JSON file.")
	customIsoUpdateCmd.MarkFlagsOneRequired("config-source")

	customIsoCmd.AddCommand(customIsoDeleteCmd)
	customIsoCmd.AddCommand(customIsoMakePublicCmd)
	customIsoCmd.AddCommand(customIsoBootServerCmd)
}
