package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	fabricFlags = struct {
		configSource         string
		networkDeviceA       string
		interfaceA           string
		networkDeviceB       string
		interfaceB           string
		linkType             string
		bgpNumbering         string
		bgpLinkConfiguration string
		customVariables      []string
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
  deploy         Deploy a fabric
  config-example Show configuration example
  get-devices    List fabric devices
  add-device     Add devices to fabric
  remove-device  Remove device from fabric
  get-links      List fabric links
  add-link       Add fabric link
  remove-link    Remove fabric link`,
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

	fabricDeployCmd = &cobra.Command{
		Use:   "deploy fabric_id",
		Short: "Deploy a fabric",
		Long: `Deploy a network fabric underlay.

This command deploys fabric underlay using the configured links and templates.

Arguments:
  fabric_id    The ID or label of the fabric to deploy

Examples:
  # Deploy fabric by ID
  metalcloud fabric deploy 12345
  
  # Deploy fabric by label
  metalcloud fabric deploy my-fabric-label`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDeploy(cmd.Context(), args[0])
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

	fabricLinksGetCmd = &cobra.Command{
		Use:     "get-links fabric_id",
		Aliases: []string{"show-links", "list-links"},
		Short:   "List links in a fabric",
		Long: `List all network fabric links in a specific fabric.

This command displays a table showing all links that are part of the specified fabric,
including their link information, status, and connection details.

Arguments:
  fabric_id    The ID or label of the fabric to list links for

Examples:
  # List links in fabric by ID
  metalcloud fabric get-links 12345
  
  # List links in fabric by label
  metalcloud fabric get-links my-fabric-label
  
  # Using alias
  metalcloud fabric list-links 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricLinksGet(cmd.Context(), args[0])
		},
	}

	fabricLinkAddCmd = &cobra.Command{
		Use:     "add-link fabric_id",
		Aliases: []string{"create-link"},
		Short:   "Add a network fabric link",
		Long: `Add a new network fabric link to an existing fabric.

This command creates a new link in the fabric using the configuration provided through
the --config-source flag. The configuration must be a JSON file or piped input containing
the link details such as source and destination network devices and interfaces.

Arguments:
  fabric_id     The ID or label of the fabric to add the link to

Required Flags when using raw configuration:
  --config-source string   Source of the link configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the link configuration

Required Flags when using individual flags:
  --networkDeviceA string  Identifier string of network device A
  --InterfaceA     string  Name of the interface A
  --networkDeviceB string  Identifier string of network device B
  --InterfaceB     string  Name of the interface B
  --linkType       string  Link type: point-to-point, broadcast

Optional Flags when using individual flags:
  --bgpNumbering         string   inherited, numbered, unnumbered
  --bgpLinkConfiguration string   disabled, active, passive
  --customVariables

Examples:
  # Add link with configuration from file
  metalcloud fabric add-link my-fabric --config-source link-config.json
  
  # Add link with piped configuration
  cat link-config.json | metalcloud fabric add-link 12345 --config-source pipe
  
  # Using alias
  metalcloud fabric create-link my-fabric --config-source link.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fabricFlags.configSource != "" {
				// Create the link using raw source format
				config, err := utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
				if err != nil {
					return err
				}

				var createLink sdk.CreateNetworkFabricLink
				err = utils.UnmarshalContent(config, &createLink)
				if err != nil {
					return err
				}

				return fabric.FabricLinkAdd(cmd.Context(), args[0], createLink)
			}

			// Create the link using flags
			return fabric.FabricLinkAddEx(cmd.Context(), args[0],
				fabricFlags.networkDeviceA,
				fabricFlags.interfaceA,
				fabricFlags.networkDeviceB,
				fabricFlags.interfaceB,
				fabricFlags.linkType,
				fabricFlags.bgpNumbering,
				fabricFlags.bgpLinkConfiguration,
				fabricFlags.customVariables,
			)
		},
	}

	fabricLinkRemoveCmd = &cobra.Command{
		Use:     "remove-link fabric_id link_id",
		Aliases: []string{"delete-link"},
		Short:   "Remove a network fabric link",
		Long: `Remove a network fabric link from an existing fabric.

This command removes a link from the fabric, disconnecting the associated network devices.

Arguments:
  fabric_id    The ID or label of the fabric to remove the link from
  link_id      The ID of the link to remove from the fabric

Examples:
  # Remove link from fabric by IDs
  metalcloud fabric remove-link 12345 67890
  
  # Remove link using fabric label
  metalcloud fabric remove-link my-fabric 67890
  
  # Using alias
  metalcloud fabric delete-link my-fabric 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricLinkRemove(cmd.Context(), args[0], args[1])
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
	fabricCmd.AddCommand(fabricDeployCmd)

	fabricCmd.AddCommand(fabricDevicesGetCmd)
	fabricCmd.AddCommand(fabricDevicesAddCmd)
	fabricCmd.AddCommand(fabricDevicesRemoveCmd)

	fabricCmd.AddCommand(fabricLinksGetCmd)

	fabricCmd.AddCommand(fabricLinkAddCmd)
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the link configuration. Can be 'pipe' or path to a JSON file.")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.networkDeviceA, "network-device-a", "", "Identifier of the network device A")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.interfaceA, "interface-a", "", "Name of the interface A")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.networkDeviceB, "network-device-b", "", "Identifier of the network device B")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.interfaceB, "interface-b", "", "Name of the interface B")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.linkType, "link-type", "", "Type of the link")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.bgpNumbering, "bgp-numbering", "inherited", "BGP numbering")
	fabricLinkAddCmd.Flags().StringVar(&fabricFlags.bgpLinkConfiguration, "bgp-link-configuration", "disabled", "BGP configuration")
	fabricLinkAddCmd.Flags().StringArrayVar(&fabricFlags.customVariables, "custom-variable", []string{}, "Custom variable")
	fabricLinkAddCmd.MarkFlagsOneRequired("config-source", "network-device-a")
	fabricLinkAddCmd.MarkFlagsRequiredTogether("network-device-a", "interface-a", "network-device-b", "interface-b", "link-type")

	fabricCmd.AddCommand(fabricLinkRemoveCmd)
}
