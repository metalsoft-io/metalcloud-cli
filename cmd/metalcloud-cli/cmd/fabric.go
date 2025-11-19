package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	fabricFlags = struct {
		configSource string
	}{}

	fabricCmd = &cobra.Command{
		Use:     "fabric [command]",
		Aliases: []string{"fc", "fabrics"},
		Short:   "Manage network fabrics",
		Long: `Manage network fabrics in MetalCloud.

Fabrics are logical network constructs that group network devices and define how they are interconnected.
This command provides operations to create, configure, activate, and manage fabric devices.

Available Commands:
  list           List all fabrics
  get            Get fabric details
  create         Create a new fabric
  update         Update fabric configuration
  activate       Activate a fabric
  config-example Show configuration example
  get-devices    List fabric devices
  add-device     Add devices to fabric
  remove-device  Remove device from fabric`,
	}

	fabricListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all network fabrics",
		Long: `List all network fabrics in your MetalCloud infrastructure.

This command displays a table showing all fabrics with their details including:
- Fabric ID and name
- Fabric type and status
- Associated site
- Creation timestamp

Examples:
  # List all fabrics
  metalcloud fabric list
  
  # Using alias
  metalcloud fc ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricList(cmd.Context())
		},
	}

	fabricGetCmd = &cobra.Command{
		Use:     "get fabric_id",
		Aliases: []string{"show"},
		Short:   "Get detailed fabric information",
		Long: `Get detailed information about a specific network fabric.

This command displays comprehensive fabric details including:
- Basic fabric properties (ID, name, type, status)
- Associated site information
- Network configuration
- Device associations
- Creation and modification timestamps

Arguments:
  fabric_id    The ID or label of the fabric to retrieve

Examples:
  # Get fabric by ID
  metalcloud fabric get 12345
  
  # Get fabric by label
  metalcloud fabric get my-fabric-label
  
  # Using alias
  metalcloud fc show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricGet(cmd.Context(), args[0])
		},
	}

	fabricConfigExampleCmd = &cobra.Command{
		Use:   "config-example fabric_type",
		Short: "Show example fabric configuration",
		Long: `Show example configuration for the specified fabric type.

This command returns a JSON configuration template that can be used as a starting point
for creating or updating fabrics. The configuration includes all available options
and their expected formats for the specified fabric type.

Arguments:
  fabric_type    The type of fabric to get configuration example for
                 (e.g., "spine-leaf", "collapsed-core", "hybrid")

Examples:
  # Get configuration example for spine-leaf fabric
  metalcloud fabric config-example spine-leaf
  
  # Get example and save to file
  metalcloud fabric config-example collapsed-core > fabric-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricConfigExample(cmd.Context(), args[0])
		},
	}

	fabricCreateCmd = &cobra.Command{
		Use:     "create site_id_or_label fabric_name fabric_type [fabric_description]",
		Aliases: []string{"new"},
		Short:   "Create a new fabric",
		Long: `Create a new network fabric in MetalCloud.

This command creates a new fabric with the specified configuration. The fabric configuration
must be provided through the --config-source flag, which can be a JSON file or piped input.

Arguments:
  site_id_or_label       The ID or label of the site where the fabric will be created
  fabric_name           The name for the new fabric
  fabric_type           The type of fabric to create (e.g., "ethernet", "infiniband")
  fabric_description    Optional description for the fabric (defaults to fabric_name if not provided)

Required Flags:
  --config-source string   Source of the fabric configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the fabric configuration

Examples:
  # Create fabric with configuration from file
  metalcloud fabric create site1 my-fabric ethernet "Production fabric" --config-source fabric-config.json
  
  # Create fabric with piped configuration
  cat fabric-config.json | metalcloud fabric create site1 my-fabric ethernet --config-source pipe
  
  # Get example config and create fabric
  metalcloud fabric config-example ethernet > config.json
  # Edit config.json as needed
  metalcloud fabric create site1 my-fabric ethernet --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[1]
			if len(args) > 3 {
				description = args[3]
			}

			config, err := utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
			if err != nil {
				return err
			}

			return fabric.FabricCreate(cmd.Context(), args[0], args[1], args[2], description, config)
		},
	}

	fabricUpdateCmd = &cobra.Command{
		Use:     "update fabric_id [fabric_name [fabric_description]]",
		Aliases: []string{"edit"},
		Short:   "Update fabric configuration",
		Long: `Update the configuration, name, or description of an existing fabric.

This command allows you to modify fabric properties and configuration. The fabric
configuration can be updated by providing a new configuration through the --config-source flag.
The name and description are optional and will only be updated if provided.

Arguments:
  fabric_id            The ID or label of the fabric to update
  fabric_name          Optional new name for the fabric
  fabric_description   Optional new description for the fabric

Required Flags:
  --config-source string   Source of the updated fabric configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the updated configuration

Examples:
  # Update fabric configuration from file
  metalcloud fabric update 12345 --config-source updated-config.json
  
  # Update name, description and configuration
  metalcloud fabric update my-fabric "New Name" "New Description" --config-source config.json
  
  # Update with piped configuration
  cat new-config.json | metalcloud fabric update 12345 --config-source pipe
  
  # Update only configuration, keeping existing name and description
  metalcloud fabric update my-fabric --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			description := ""
			if len(args) > 2 {
				description = args[2]
			}

			config, err := utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
			if err != nil {
				return err
			}

			return fabric.FabricUpdate(cmd.Context(), args[0], name, description, config)
		},
	}

	fabricActivateCmd = &cobra.Command{
		Use:     "activate fabric_id",
		Aliases: []string{"start"},
		Short:   "Activate a fabric",
		Long: `Activate a network fabric to make it operational.

This command activates a fabric that has been created and configured. Once activated,
the fabric will begin managing the network connectivity according to its configuration.
Only fabrics in an inactive state can be activated.

Arguments:
  fabric_id    The ID or label of the fabric to activate

Examples:
  # Activate fabric by ID
  metalcloud fabric activate 12345
  
  # Activate fabric by label
  metalcloud fabric activate my-fabric-label
  
  # Using alias
  metalcloud fc start 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricActivate(cmd.Context(), args[0])
		},
	}

	fabricDevicesGetCmd = &cobra.Command{
		Use:     "get-devices fabric_id",
		Aliases: []string{"show-devices"},
		Short:   "List devices in a fabric",
		Long: `List all network devices associated with a specific fabric.

This command displays a table showing all devices that are part of the specified fabric,
including their device information, status, and role within the fabric configuration.

Arguments:
  fabric_id    The ID or label of the fabric to list devices for

Examples:
  # List devices in fabric by ID
  metalcloud fabric get-devices 12345
  
  # List devices in fabric by label
  metalcloud fabric get-devices my-fabric-label
  
  # Using alias
  metalcloud fabric show-devices 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesGet(cmd.Context(), args[0])
		},
	}

	fabricDevicesAddCmd = &cobra.Command{
		Use:     "add-device fabric_id device_id...",
		Aliases: []string{"join-device"},
		Short:   "Add network device(s) to a fabric",
		Long: `Add one or more network devices to an existing fabric.

This command associates network devices with a fabric, making them part of the fabric's
network topology. Multiple devices can be added in a single command by specifying
multiple device IDs.

Arguments:
  fabric_id     The ID or label of the fabric to add devices to
  device_id...  One or more device IDs or labels to add to the fabric

Examples:
  # Add a single device to fabric
  metalcloud fabric add-device my-fabric device123
  
  # Add multiple devices to fabric
  metalcloud fabric add-device 12345 device1 device2 device3
  
  # Using alias
  metalcloud fabric join-device my-fabric switch-01`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesAdd(cmd.Context(), args[0], args[1:])
		},
	}

	fabricDevicesRemoveCmd = &cobra.Command{
		Use:     "remove-device fabric_id device_id",
		Aliases: []string{"delete-device"},
		Short:   "Remove network device from a fabric",
		Long: `Remove a network device from an existing fabric.

This command disassociates a network device from a fabric, removing it from the fabric's
network topology. The device will no longer be managed by the fabric configuration.

Arguments:
  fabric_id    The ID or label of the fabric to remove the device from
  device_id    The ID or label of the device to remove from the fabric

Examples:
  # Remove device from fabric by IDs
  metalcloud fabric remove-device 12345 device123
  
  # Remove device using labels
  metalcloud fabric remove-device my-fabric switch-01
  
  # Using alias
  metalcloud fabric delete-device my-fabric device123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesRemove(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(fabricCmd)

	fabricCmd.AddCommand(fabricListCmd)

	fabricCmd.AddCommand(fabricGetCmd)

	fabricCmd.AddCommand(fabricConfigExampleCmd)

	fabricCmd.AddCommand(fabricCreateCmd)
	fabricCreateCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the new fabric configuration. Can be 'pipe' or path to a JSON file.")
	fabricCreateCmd.MarkFlagsOneRequired("config-source")

	fabricCmd.AddCommand(fabricUpdateCmd)
	fabricUpdateCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the updated fabric configuration. Can be 'pipe' or path to a JSON file.")
	fabricUpdateCmd.MarkFlagsOneRequired("config-source")

	fabricCmd.AddCommand(fabricActivateCmd)

	fabricCmd.AddCommand(fabricDevicesGetCmd)
	fabricCmd.AddCommand(fabricDevicesAddCmd)
	fabricCmd.AddCommand(fabricDevicesRemoveCmd)
}
