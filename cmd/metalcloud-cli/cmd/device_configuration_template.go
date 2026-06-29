package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/device_configuration_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	deviceConfigurationTemplateFlags = struct {
		configSource          string
		filterId              []string
		filterLabel           []string
		filterName            []string
		filterTemplateId      []string
		filterNetworkDeviceId []string
		filterNetworkFabricId []string
	}{}

	deviceConfigurationTemplateCmd = &cobra.Command{
		Use:     "device-config-template [command]",
		Aliases: []string{"dct"},
		Short:   "Manage device configuration templates and profiles",
		Long: `Device configuration template commands.

Device configuration templates hold renderable configuration content for network
devices, and profiles bind those templates to specific devices or fabrics.

Available commands:
  list                List device configuration templates
  get                 Show details about a specific template
  create              Create a new template from JSON configuration
  update              Update an existing template
  delete              Delete a template
  render              Render arbitrary template content
  render-saved        Render a saved template by ID
  config-example      Generate an example template configuration
  profile             Manage device configuration template profiles`,
	}

	deviceConfigurationTemplateListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List device configuration templates with optional filtering",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateList(
				cmd.Context(),
				deviceConfigurationTemplateFlags.filterId,
				deviceConfigurationTemplateFlags.filterLabel,
				deviceConfigurationTemplateFlags.filterName,
			)
		},
	}

	deviceConfigurationTemplateGetCmd = &cobra.Command{
		Use:          "get <device_configuration_template_id>",
		Aliases:      []string{"show"},
		Short:        "Get detailed information about a specific device configuration template",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateGet(cmd.Context(), args[0])
		},
	}

	deviceConfigurationTemplateCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new device configuration template",
		Long: `Create a new device configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration.

Examples:
  metalcloud-cli device-config-template create --config-source template.json
  cat template.json | metalcloud-cli dct create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateCreate(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateUpdateCmd = &cobra.Command{
		Use:          "update <device_configuration_template_id>",
		Aliases:      []string{"modify"},
		Short:        "Update an existing device configuration template",
		Long: `Update an existing device configuration template using JSON configuration provided via
file or pipe. Only the specified fields will be updated; other configuration remains unchanged.

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  metalcloud-cli device-config-template update 12345 --config-source updates.json
  cat updates.json | metalcloud-cli dct update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	deviceConfigurationTemplateDeleteCmd = &cobra.Command{
		Use:          "delete <device_configuration_template_id>",
		Aliases:      []string{"rm"},
		Short:        "Delete a device configuration template",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateDelete(cmd.Context(), args[0])
		},
	}

	deviceConfigurationTemplateRenderCmd = &cobra.Command{
		Use:          "render",
		Short:        "Render arbitrary device configuration template content",
		Long: `Render device configuration template content provided inline (not saved) with the given variables.

Required Flags:
  --config-source   Source of the render request (required)
                   Values: 'pipe' for stdin input, or path to JSON file

The request body accepts: templateContent (required), variables, debug.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateRender(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateRenderSavedCmd = &cobra.Command{
		Use:          "render-saved <device_configuration_template_id>",
		Short:        "Render a saved device configuration template by ID",
		Long: `Render a previously saved device configuration template, identified by ID, with the given variables.

Required Flags:
  --config-source   Source of the render request (required)
                   Values: 'pipe' for stdin input, or path to JSON file

The request body accepts: variables, debug.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateRenderSaved(cmd.Context(), args[0], config)
		},
	}

	deviceConfigurationTemplateConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Generate an example device configuration template",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateConfigExample(cmd.Context())
		},
	}

	// Profile sub-commands -----------------------------------------------------

	deviceConfigurationTemplateProfileCmd = &cobra.Command{
		Use:     "profile [command]",
		Aliases: []string{"p"},
		Short:   "Manage device configuration template profiles",
		Long: `Device configuration template profile commands.

Profiles bind a device configuration template to a specific network device or fabric,
with variables and lifecycle/apply settings.

Available commands:
  list                List profiles
  get                 Show details about a specific profile
  create              Create a new profile from JSON configuration
  update              Update an existing profile
  delete              Delete a profile
  render              Render a profile for a given device
  find-applicable     Find profiles applicable to a device/fabric
  render-applicable   Render profiles applicable to a device/fabric
  bulk-assign         Bulk-assign a template to multiple devices
  config-example      Generate an example profile configuration`,
	}

	deviceConfigurationTemplateProfileListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List device configuration template profiles with optional filtering",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateProfileList(
				cmd.Context(),
				deviceConfigurationTemplateFlags.filterId,
				deviceConfigurationTemplateFlags.filterTemplateId,
				deviceConfigurationTemplateFlags.filterNetworkDeviceId,
				deviceConfigurationTemplateFlags.filterNetworkFabricId,
			)
		},
	}

	deviceConfigurationTemplateProfileGetCmd = &cobra.Command{
		Use:          "get <profile_id>",
		Aliases:      []string{"show"},
		Short:        "Get detailed information about a specific profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateProfileGet(cmd.Context(), args[0])
		},
	}

	deviceConfigurationTemplateProfileCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new device configuration template profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileCreate(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateProfileUpdateCmd = &cobra.Command{
		Use:          "update <profile_id>",
		Aliases:      []string{"modify"},
		Short:        "Update an existing device configuration template profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileUpdate(cmd.Context(), args[0], config)
		},
	}

	deviceConfigurationTemplateProfileDeleteCmd = &cobra.Command{
		Use:          "delete <profile_id>",
		Aliases:      []string{"rm"},
		Short:        "Delete a device configuration template profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateProfileDelete(cmd.Context(), args[0])
		},
	}

	deviceConfigurationTemplateProfileRenderCmd = &cobra.Command{
		Use:          "render <profile_id>",
		Short:        "Render a device configuration template profile for a device",
		Long: `Render a profile for a given network device.

Required Flags:
  --config-source   Source of the render request (required)

The request body accepts: networkDeviceId (required), extraVariables, debug.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileRender(cmd.Context(), args[0], config)
		},
	}

	deviceConfigurationTemplateProfileFindApplicableCmd = &cobra.Command{
		Use:          "find-applicable",
		Short:        "Find device configuration template profiles applicable to a device or fabric",
		Long: `Find profiles applicable to a network device or fabric.

Required Flags:
  --config-source   Source of the request (required)

The request body accepts: networkDeviceId, networkFabricId, lifecycleStage, includeDisabled.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileFindApplicable(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateProfileRenderApplicableCmd = &cobra.Command{
		Use:          "render-applicable",
		Short:        "Render device configuration template profiles applicable to a device or fabric",
		Long: `Render profiles applicable to a network device or fabric.

Required Flags:
  --config-source   Source of the request (required)

The request body accepts: networkDeviceId, networkFabricId, lifecycleStage, includeDisabled, extraVariables, debug.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileRenderApplicable(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateProfileBulkAssignCmd = &cobra.Command{
		Use:          "bulk-assign",
		Short:        "Bulk-assign a device configuration template to multiple devices",
		Long: `Bulk-assign a device configuration template to multiple network devices as profiles.

Required Flags:
  --config-source   Source of the request (required)

The request body accepts: deviceConfigurationTemplateId (required), networkFabricId,
networkDeviceIds, lifecycleStage, variables, isEnabled, priority, applyMode, annotations, tags.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(deviceConfigurationTemplateFlags.configSource)
			if err != nil {
				return err
			}
			return device_configuration_template.DeviceConfigurationTemplateProfileBulkAssign(cmd.Context(), config)
		},
	}

	deviceConfigurationTemplateProfileConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Generate an example device configuration template profile",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_DEVICE_CONFIGURATION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return device_configuration_template.DeviceConfigurationTemplateProfileConfigExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(deviceConfigurationTemplateCmd)

	// Template (config) commands
	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateListCmd)
	deviceConfigurationTemplateListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterId, "filter-id", nil, "Filter by template ID.")
	deviceConfigurationTemplateListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterLabel, "filter-label", nil, "Filter by template label.")
	deviceConfigurationTemplateListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterName, "filter-name", nil, "Filter by template name.")

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateGetCmd)

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateCreateCmd)
	deviceConfigurationTemplateCreateCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the new template configuration. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateUpdateCmd)
	deviceConfigurationTemplateUpdateCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the template updates. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateDeleteCmd)

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateRenderCmd)
	deviceConfigurationTemplateRenderCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the render request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateRenderCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateRenderSavedCmd)
	deviceConfigurationTemplateRenderSavedCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the render request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateRenderSavedCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateConfigExampleCmd)

	// Profile commands
	deviceConfigurationTemplateCmd.AddCommand(deviceConfigurationTemplateProfileCmd)

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileListCmd)
	deviceConfigurationTemplateProfileListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterId, "filter-id", nil, "Filter by profile ID.")
	deviceConfigurationTemplateProfileListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterTemplateId, "filter-template-id", nil, "Filter by device configuration template ID.")
	deviceConfigurationTemplateProfileListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterNetworkDeviceId, "filter-network-device-id", nil, "Filter by network device ID.")
	deviceConfigurationTemplateProfileListCmd.Flags().StringSliceVar(&deviceConfigurationTemplateFlags.filterNetworkFabricId, "filter-network-fabric-id", nil, "Filter by network fabric ID.")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileGetCmd)

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileCreateCmd)
	deviceConfigurationTemplateProfileCreateCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the new profile configuration. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileCreateCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileUpdateCmd)
	deviceConfigurationTemplateProfileUpdateCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the profile updates. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileUpdateCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileDeleteCmd)

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileRenderCmd)
	deviceConfigurationTemplateProfileRenderCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the render request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileRenderCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileFindApplicableCmd)
	deviceConfigurationTemplateProfileFindApplicableCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileFindApplicableCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileRenderApplicableCmd)
	deviceConfigurationTemplateProfileRenderApplicableCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileRenderApplicableCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileBulkAssignCmd)
	deviceConfigurationTemplateProfileBulkAssignCmd.Flags().StringVar(&deviceConfigurationTemplateFlags.configSource, "config-source", "", "Source of the request. Can be 'pipe' or path to a JSON file.")
	deviceConfigurationTemplateProfileBulkAssignCmd.MarkFlagsOneRequired("config-source")

	deviceConfigurationTemplateProfileCmd.AddCommand(deviceConfigurationTemplateProfileConfigExampleCmd)
}
