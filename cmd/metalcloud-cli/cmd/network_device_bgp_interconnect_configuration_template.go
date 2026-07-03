package cmd

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device_bgp_interconnect_configuration_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	networkDeviceBGPInterconnectConfigurationTemplateFlags = struct {
		configSource              string
		filterId                  []string
		filterNetworkDeviceDriver []string
	}{}

	networkDeviceBGPInterconnectConfigurationTemplateCmd = &cobra.Command{
		Use:     "bgp-interconnect-template [command]",
		Aliases: []string{"bit"},
		Short:   "Manage network devices BGP interconnect configuration templates",
		Long: `Network device BGP interconnect configuration template commands.

Network device BGP interconnect configuration templates are used to deploy BGP interconnect configurations to network devices
Available commands:
  list                List all available Network device BGP interconnect configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device BGP interconnect configuration template from JSON configuration
  update              Update an existing Network device BGP interconnect configuration template
  delete              Delete a Network device BGP interconnect configuration template`,
	}

	networkDeviceBGPInterconnectConfigurationTemplateListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network device BGP interconnect configuration templates with optional filtering",
		Long: `List all network device BGP interconnect configuration templates with optional filtering.

This command displays all network device BGP interconnect configuration templates that are registered in the system.
You can filter the results by network device driver to focus on specific groups of templates.
Flags:
  --filter-network-device-driver   Filter templates by network device driver

Examples:
  # List all network device BGP interconnect configuration templates (default)
  metalcloud-cli network-configuration bgp-interconnect-template list

  # List templates for a specific network device driver
  metalcloud-cli network-configuration bgp-interconnect-template list --filter-network-device-driver junos`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateList(cmd.Context(), networkDeviceBGPInterconnectConfigurationTemplateFlags.filterId, networkDeviceBGPInterconnectConfigurationTemplateFlags.filterNetworkDeviceDriver)
		},
	}

	networkDeviceBGPInterconnectConfigurationTemplateGetCmd = &cobra.Command{
		Use:     "get <network_device_bgp_interconnect_configuration_template_id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific network device BGP interconnect configuration template",
		Long: `Display detailed information about a specific network device BGP interconnect configuration template.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template

Examples:
  # Get details for network device BGP interconnect configuration template with ID 12345
  metalcloud-cli network-configuration bgp-interconnect-template get 12345
  # Using alias
  metalcloud-cli network-configuration bgp-interconnect-template show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateGet(cmd.Context(), args[0])
		},
	}

	networkDeviceBGPInterconnectConfigurationTemplateConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Generate example configuration template for network device BGP interconnect configuration template",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateConfigExample(cmd.Context())
		},
	}

	networkDeviceBGPInterconnectConfigurationTemplateCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new network device BGP interconnect configuration template with specified configuration",
		Long: `Create a new network device BGP interconnect configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "label": "string",
    "name": "string",
    "networkDeviceDriver": "junos",
    "executionType": "cli",
    "addGlobalConfig": "string",
    "removeGlobalConfig": "string",
    "addNeighbor": "string",
    "removeNeighbor": "string"
  }

Note: addGlobalConfig, removeGlobalConfig, addNeighbor and removeNeighbor fields need to be base64 encoded when submitted.

Examples:
  # Create template from JSON file
  metalcloud-cli network-configuration bgp-interconnect-template create --config-source template.json

  # Create template from pipe input
  cat template.json | metalcloud-cli network-configuration bgp-interconnect-template create --config-source pipe

  # Create template with inline JSON
  echo '{"label":"l","name":"n","networkDeviceDriver":"junos","executionType":"cli"}' | metalcloud-cli nc bit create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceBGPInterconnectConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateCreate(cmd.Context(), config)
		},
	}

	networkDeviceBGPInterconnectConfigurationTemplateUpdateCmd = &cobra.Command{
		Use:     "update <network_device_bgp_interconnect_configuration_template_id>",
		Aliases: []string{"modify"},
		Short:   "Update configuration of an existing network device BGP interconnect configuration template",
		Long: `Update the configuration of an existing network device BGP interconnect configuration template using JSON configuration
provided via file or pipe. Only the specified fields will be updated; other
configuration will remain unchanged.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template to update

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Update template from JSON file
  metalcloud-cli network-configuration bgp-interconnect-template update 12345 --config-source updates.json

  # Update template from pipe input
  cat updates.json | metalcloud-cli network-configuration bgp-interconnect-template update 12345 --config-source pipe

  # Update specific field
  echo '{"name":"new name"}' | metalcloud-cli nc bit update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceBGPInterconnectConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceBGPInterconnectConfigurationTemplateDeleteCmd = &cobra.Command{
		Use:     "delete <network_device_bgp_interconnect_configuration_template_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a network device BGP interconnect configuration template from the system",
		Long: `Delete a network device BGP interconnect configuration template from the system.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template to delete

Examples:
  # Delete network device BGP interconnect configuration template
  metalcloud-cli network-configuration bgp-interconnect-template delete 12345

  # Using alias
  metalcloud-cli nc bit rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRIC_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_bgp_interconnect_configuration_template.NetworkDeviceBGPInterconnectConfigurationTemplateDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	networkDeviceBGPInterconnectConfigurationTemplateConfigExampleCmd.Long = fmt.Sprintf(`Generate an example JSON configuration template that can be used to create
or update network device BGP interconnect configuration templates.

The addGlobalConfig, removeGlobalConfig, addNeighbor and removeNeighbor fields need to be base64 encoded when submitted.

Accepted field values:
  networkDeviceDriver: %s
  executionType:       %s

Examples:
  # Display example configuration
  metalcloud-cli network-configuration bgp-interconnect-template config-example -f json

  # Save example to file
  metalcloud-cli network-configuration bgp-interconnect-template config-example -f json > template.json`,
		strings.Join(network_device_bgp_interconnect_configuration_template.ValidNetworkDeviceDrivers, ", "),
		strings.Join(network_device_bgp_interconnect_configuration_template.ValidExecutionTypes, ", "),
	)

	networkConfigurationCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateCmd)

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateListCmd)
	networkDeviceBGPInterconnectConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceBGPInterconnectConfigurationTemplateFlags.filterId, "filter-id", nil, "Filter by template ID.")
	networkDeviceBGPInterconnectConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceBGPInterconnectConfigurationTemplateFlags.filterNetworkDeviceDriver, "filter-network-device-driver", nil, "Filter by network device driver.")

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateGetCmd)

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateConfigExampleCmd)

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateCreateCmd)
	networkDeviceBGPInterconnectConfigurationTemplateCreateCmd.Flags().StringVar(&networkDeviceBGPInterconnectConfigurationTemplateFlags.configSource, "config-source", "", "Source of the new network device BGP interconnect configuration template. Can be 'pipe' or path to a JSON file.")
	networkDeviceBGPInterconnectConfigurationTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateUpdateCmd)
	networkDeviceBGPInterconnectConfigurationTemplateUpdateCmd.Flags().StringVar(&networkDeviceBGPInterconnectConfigurationTemplateFlags.configSource, "config-source", "", "Source of the network device BGP interconnect configuration template updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceBGPInterconnectConfigurationTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceBGPInterconnectConfigurationTemplateCmd.AddCommand(networkDeviceBGPInterconnectConfigurationTemplateDeleteCmd)
}
