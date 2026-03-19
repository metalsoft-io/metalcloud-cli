package cmd

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device_configuration_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	networkDeviceConfigurationTemplateFlags = struct {
		configSource       string
		filterId           []string
		filterLibraryLabel []string
	}{}

	networkDeviceConfigurationTemplateCmd = &cobra.Command{
		Use:     "device-template [command]",
		Aliases: []string{"dt"},
		Short:   "Manage network devices configuration templates",
		Long: `Network device configuration template commands.

Network device configuration templates are used to deploy configurations to network devices
Available commands:
  list                List all available Network device configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device configuration template from JSON configuration
  update              Update an existing Network device configuration template
  delete              Delete a Network device configuration template`,
	}

	networkDeviceConfigurationTemplateListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network device configuration templates with optional filtering",
		Long: `List all network device configuration templates with optional filtering.

This command displays all network device configuration templates that are registered in the system.
You can filter the results by library label to focus on specific groups of templates.
Flags:
  --filter-library-label   Filter templates by library label

Examples:
  # List all network device configuration templates (default)
  metalcloud-cli network-configuration device-template list

  # List templates with a specific library label
  metalcloud-cli network-configuration device-template list --filter-library-label example-label`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_configuration_template.NetworkDeviceConfigurationTemplateList(cmd.Context(), networkDeviceConfigurationTemplateFlags.filterId, networkDeviceConfigurationTemplateFlags.filterLibraryLabel)
		},
	}

	networkDeviceConfigurationTemplateGetCmd = &cobra.Command{
		Use:     "get <network_device_configuration_template_id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific network device configuration template",
		Long: `Display detailed information about a specific network device configuration template.

Arguments:
  network_device_configuration_template_id   The unique identifier of the network device configuration template

Examples:
  # Get details for network device configuration template with ID 12345
  metalcloud-cli network-configuration device-template get 12345
  # Using alias
  metalcloud-cli network-configuration device-template show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_configuration_template.NetworkDeviceConfigurationTemplateGet(cmd.Context(), args[0])
		},
	}

	networkDeviceConfigurationTemplateConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Generate example configuration template for network device configuration template",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_configuration_template.NetworkDeviceConfigurationTemplateConfigExample(cmd.Context())
		},
	}

	networkDeviceConfigurationTemplateCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new network device configuration template with specified configuration",
		Long: `Create a new network device configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "action": "add-global-config",
    "networkType": "underlay",
    "networkDeviceDriver": "junos",
    "networkDevicePosition": "all",
    "remoteNetworkDevicePosition": "all",
    "bgpNumbering": "numbered",
    "bgpLinkConfiguration": "disabled",
    "executionType": "cli",
    "libraryLabel": "string",
    "preparation": "string",
    "configuration": "string"
  }

Note: Preparation and configuration fields need to be base64 encoded when submitted.

Examples:
  # Create template from JSON file
  metalcloud-cli network-configuration device-template create --config-source template.json

  # Create template from pipe input
  cat template.json | metalcloud-cli network-configuration device-template create --config-source pipe

  # Create template with inline JSON
  echo '{"networkDevicePosition":"all","remoteNetworkDevicePosition":"all"}' | metalcloud-cli nc dt create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_configuration_template.NetworkDeviceConfigurationTemplateCreate(cmd.Context(), config)
		},
	}

	networkDeviceConfigurationTemplateUpdateCmd = &cobra.Command{
		Use:     "update <network_device_configuration_template_id>",
		Aliases: []string{"modify"},
		Short:   "Update configuration of an existing network device configuration template",
		Long: `Update the configuration of an existing network device configuration template using JSON configuration
provided via file or pipe. Only the specified fields will be updated; other
configuration will remain unchanged.

Arguments:
  network_device_configuration_template_id   The unique identifier of the network device configuration template to update

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Update template from JSON file
  metalcloud-cli network-configuration device-template update 12345 --config-source updates.json

  # Update template from pipe input
  cat updates.json | metalcloud-cli network-configuration device-template update 12345 --config-source pipe

  # Update specific field
  echo '{"networkDevicePosition":"all"}' | metalcloud-cli nc dt update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_configuration_template.NetworkDeviceConfigurationTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceConfigurationTemplateDeleteCmd = &cobra.Command{
		Use:     "delete <network_device_configuration_template_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a network device configuration template from the system",
		Long: `Delete a network device configuration template from the system.

Arguments:
  network_device_configuration_template_id   The unique identifier of the network device configuration template to delete

Examples:
  # Delete network device configuration template
  metalcloud-cli network-configuration device-template delete 12345

  # Using alias
  metalcloud-cli nc dt rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_configuration_template.NetworkDeviceConfigurationTemplateDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	networkDeviceConfigurationTemplateConfigExampleCmd.Long = fmt.Sprintf(`Generate an example JSON configuration template that can be used to create
or update network device configuration templates.

Preparation and configuration fields need to be base64 encoded when submitted.

Accepted field values:
  action:               %s
  networkType:          %s
  networkDeviceDriver:  %s
  networkDevicePosition / remoteNetworkDevicePosition:
                        %s
  bgpNumbering:         %s
  bgpLinkConfiguration: %s

Examples:
  # Display example configuration
  metalcloud-cli network-configuration device-template config-example -f json

  # Save example to file
  metalcloud-cli network-configuration device-template config-example -f json > template.json`,
		strings.Join(network_device_configuration_template.ValidDeviceTemplateActions, ", "),
		strings.Join(network_device_configuration_template.ValidDeviceTemplateNetworkTypes, ", "),
		strings.Join(network_device_configuration_template.ValidNetworkDeviceDrivers, ", "),
		strings.Join(network_device_configuration_template.ValidNetworkDevicePositions, ", "),
		strings.Join(network_device_configuration_template.ValidBgpNumberings, ", "),
		strings.Join(network_device_configuration_template.ValidBgpLinkConfigurations, ", "),
	)

	networkConfigurationCmd.AddCommand(networkDeviceConfigurationTemplateCmd)

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateListCmd)
	networkDeviceConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceConfigurationTemplateFlags.filterId, "filter-id", nil, "Filter by template ID.")
	networkDeviceConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceConfigurationTemplateFlags.filterLibraryLabel, "filter-library-label", nil, "Filter by template library label.")

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateGetCmd)

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateConfigExampleCmd)

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateCreateCmd)
	networkDeviceConfigurationTemplateCreateCmd.Flags().StringVar(&networkDeviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the new network device configuration template. Can be 'pipe' or path to a JSON file.")
	networkDeviceConfigurationTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateUpdateCmd)
	networkDeviceConfigurationTemplateUpdateCmd.Flags().StringVar(&networkDeviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the network device configuration template updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceConfigurationTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceConfigurationTemplateCmd.AddCommand(networkDeviceConfigurationTemplateDeleteCmd)
}
