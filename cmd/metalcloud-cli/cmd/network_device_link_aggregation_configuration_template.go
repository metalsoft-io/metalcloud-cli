package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device_link_aggregation_configuration_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	networkDeviceLinkAggregationConfigurationTemplateFlags = struct {
		configSource       string
		filterId           []string
		filterLibraryLabel []string
	}{}

	networkDeviceLinkAggregationConfigurationTemplateCmd = &cobra.Command{
		Use:     "link-aggregation-template [command]",
		Aliases: []string{"lat"},
		Short:   "Manage network devices link aggregation configuration templates",
		Long: `Network device link aggregation configuration template commands.

Network device link aggregation configuration templates are used to deploy link aggregation configurations to network devices
Available commands:
  list                List all available Network device link aggregation configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device link aggregation configuration template from JSON configuration
  update              Update an existing Network device link aggregation configuration template
  delete              Delete a Network device link aggregation configuration template`,
	}

	networkDeviceLinkAggregationConfigurationTemplateListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network device link aggregation configuration templates with optional filtering",
		Long: `List all network device link aggregation configuration templates with optional filtering.

This command displays all network device link aggregation configuration templates that are registered in the system.
You can filter the results by library label to focus on specific groups of templates.
Flags:
  --filter-library-label   Filter templates by library label

Examples:
  # List all network device link aggregation configuration templates (default)
  metalcloud-cli network-configuration link-aggregation-template list

  # List templates with a specific library label
  metalcloud-cli network-configuration link-aggregation-template list --filter-library-label example-label`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateList(cmd.Context(), networkDeviceLinkAggregationConfigurationTemplateFlags.filterId, networkDeviceLinkAggregationConfigurationTemplateFlags.filterLibraryLabel)
		},
	}

	networkDeviceLinkAggregationConfigurationTemplateGetCmd = &cobra.Command{
		Use:     "get <network_device_link_aggregation_configuration_template_id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific network device link aggregation configuration template",
		Long: `Display detailed information about a specific network device link aggregation configuration template.

Arguments:
  network_device_link_aggregation_configuration_template_id   The unique identifier of the network device link aggregation configuration template

Examples:
  # Get details for network device link aggregation configuration template with ID 12345
  metalcloud-cli network-configuration link-aggregation-template get 12345
  # Using alias
  metalcloud-cli network-configuration link-aggregation-template show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateGet(cmd.Context(), args[0])
		},
	}

	networkDeviceLinkAggregationConfigurationTemplateConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Generate example configuration template for network device link aggregation configuration template",
		Long: `Generate an example JSON configuration template that can be used to create
or update network device link aggregation configuration templates. This template includes all available configuration
options with example values and documentation.

Preparation and configuration fields need to be base64 encoded when submitted.

The generated template can be saved to a file and modified as needed for actual
template configuration.

Examples:
  # Display example configuration
  metalcloud-cli network-configuration link-aggregation-template config-example -f json

  # Save example to file
  metalcloud-cli network-configuration link-aggregation-template config-example -f json > template.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateConfigExample(cmd.Context())
		},
	}

	networkDeviceLinkAggregationConfigurationTemplateCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new network device link aggregation configuration template with specified configuration",
		Long: `Create a new network device link aggregation configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "action": "create",
    "aggregationType": "lacp",
    "networkDeviceDriver": "junos",
    "executionType": "cli",
    "libraryLabel": "string",
    "preparation": "string",
    "configuration": "string"
  }

Note: Preparation and configuration fields need to be base64 encoded when submitted.

Examples:
  # Create template from JSON file
  metalcloud-cli network-configuration link-aggregation-template create --config-source template.json

  # Create template from pipe input
  cat template.json | metalcloud-cli network-configuration link-aggregation-template create --config-source pipe

  # Create template with inline JSON
  echo '{"action":"create","aggregationType":"lacp","networkDeviceDriver":"junos","executionType":"cli","libraryLabel":"label","configuration":"string"}' | metalcloud-cli nc lat create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceLinkAggregationConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateCreate(cmd.Context(), config)
		},
	}

	networkDeviceLinkAggregationConfigurationTemplateUpdateCmd = &cobra.Command{
		Use:     "update <network_device_link_aggregation_configuration_template_id>",
		Aliases: []string{"modify"},
		Short:   "Update configuration of an existing network device link aggregation configuration template",
		Long: `Update the configuration of an existing network device link aggregation configuration template using JSON configuration
provided via file or pipe. Only the specified fields will be updated; other
configuration will remain unchanged.

Arguments:
  network_device_link_aggregation_configuration_template_id   The unique identifier of the network device link aggregation configuration template to update

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Update template from JSON file
  metalcloud-cli network-configuration link-aggregation-template update 12345 --config-source updates.json

  # Update template from pipe input
  cat updates.json | metalcloud-cli network-configuration link-aggregation-template update 12345 --config-source pipe

  # Update specific field
  echo '{"aggregationType":"static"}' | metalcloud-cli nc lat update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkDeviceLinkAggregationConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	networkDeviceLinkAggregationConfigurationTemplateDeleteCmd = &cobra.Command{
		Use:     "delete <network_device_link_aggregation_configuration_template_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a network device link aggregation configuration template from the system",
		Long: `Delete a network device link aggregation configuration template from the system.

Arguments:
  network_device_link_aggregation_configuration_template_id   The unique identifier of the network device link aggregation configuration template to delete

Examples:
  # Delete network device link aggregation configuration template
  metalcloud-cli network-configuration link-aggregation-template delete 12345

  # Using alias
  metalcloud-cli nc lat rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_LINK_AGGREGATION_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_device_link_aggregation_configuration_template.NetworkDeviceLinkAggregationConfigurationTemplateDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	networkConfigurationCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateCmd)

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateListCmd)
	networkDeviceLinkAggregationConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceLinkAggregationConfigurationTemplateFlags.filterId, "filter-id", nil, "Filter by template ID.")
	networkDeviceLinkAggregationConfigurationTemplateListCmd.Flags().StringSliceVar(&networkDeviceLinkAggregationConfigurationTemplateFlags.filterLibraryLabel, "filter-library-label", nil, "Filter by template library label.")

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateGetCmd)

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateConfigExampleCmd)

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateCreateCmd)
	networkDeviceLinkAggregationConfigurationTemplateCreateCmd.Flags().StringVar(&networkDeviceLinkAggregationConfigurationTemplateFlags.configSource, "config-source", "", "Source of the new network device link aggregation configuration template. Can be 'pipe' or path to a JSON file.")
	networkDeviceLinkAggregationConfigurationTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateUpdateCmd)
	networkDeviceLinkAggregationConfigurationTemplateUpdateCmd.Flags().StringVar(&networkDeviceLinkAggregationConfigurationTemplateFlags.configSource, "config-source", "", "Source of the network device link aggregation configuration template updates. Can be 'pipe' or path to a JSON file.")
	networkDeviceLinkAggregationConfigurationTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	networkDeviceLinkAggregationConfigurationTemplateCmd.AddCommand(networkDeviceLinkAggregationConfigurationTemplateDeleteCmd)
}
