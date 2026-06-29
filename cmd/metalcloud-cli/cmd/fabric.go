package cmd

import (
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		dryRun               bool
	}{}

	// configureSwitchesFlags are the per-property alternatives to --config-source
	// for the configure-switches command. A section is enabled if its enable bool
	// is set or any of its sub-flags is provided (detected via Flags().Changed).
	configureSwitchesFlags = struct {
		ordering            string
		enablePhysicalPorts bool
		descriptionTemplate string

		hostname           bool
		hostnameLeaf       string
		hostnameSpine      string
		hostnameSuperSpine string
		hostnameSkip       []string

		asn                bool
		asnLeafStart       int64
		asnSpineStart      int64
		asnSuperSpineStart int64

		loopback       bool
		loopbackSubnet string

		topoLeafSpine          bool
		topoLeafSpineLPP       string
		topoSpineSuperSpine    bool
		topoSpineSuperSpineLPP string

		topoLeafHost            bool
		topoLeafHostNodeCount   int
		topoLeafHostNodes       []int
		topoLeafHostPortPattern string
		topoLeafHostNicNames    []string
		topoLeafHostDescription string

		p2p                    bool
		p2pPoolLeafSpine       string
		p2pPoolSpineSuperSpine string
		p2pPoolLeafHost        string
		p2pMtu                 int32
	}{}

	// configureSwitchesDetailFlags is every per-property flag; each is marked
	// mutually exclusive with --config-source.
	configureSwitchesDetailFlags = []string{
		"ordering", "enable-physical-ports", "description-template",
		"hostname", "hostname-leaf", "hostname-spine", "hostname-super-spine", "hostname-skip",
		"asn", "asn-leaf-start", "asn-spine-start", "asn-super-spine-start",
		"loopback", "loopback-subnet",
		"topology-leaf-spine", "topology-leaf-spine-links-per-pair",
		"topology-spine-super-spine", "topology-spine-super-spine-links-per-pair",
		"topology-leaf-host", "topology-leaf-host-node-count", "topology-leaf-host-nodes",
		"topology-leaf-host-port-pattern", "topology-leaf-host-nic-names",
		"topology-leaf-host-description-template",
		"p2p", "p2p-pool-leaf-spine", "p2p-pool-spine-super-spine", "p2p-pool-leaf-host", "p2p-mtu",
	}

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

	fabricConfigureSwitchesCmd = &cobra.Command{
		Use:     "configure-switches fabric_id",
		Aliases: []string{"configure-switch"},
		Short:   "Configure all switches of a fabric from a declarative YAML/JSON",
		Long: `Configure every network device attached to a fabric from one declarative
configuration: hostnames (identifierString), ASNs, loopback IPs, physical-port
enable + interface descriptions, and point-to-point links with deterministic
/31 IPAM subnets.

Each feature section is optional - omit one to skip that step. Every step is
idempotent: current state is read first and only differences are written. Use
--dry-run to compute and preview the full plan without making any changes.

Arguments:
  fabric_id    The ID or label of the fabric to configure

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a YAML/JSON config file.

Optional Flags:
  --dry-run         Compute the plan and report what would change, without writing.

Config sections: hostname, asn, loopback, topology (leafSpine/spineSuperSpine/
leafHost), p2p, descriptionTemplate, enablePhysicalPorts, ordering.

Examples:
  metalcloud-cli fabric configure-switches 5 --config-source fabric-config.yaml --dry-run
  metalcloud-cli fabric configure-switches my-fabric --config-source fabric-config.yaml
  cat fabric-config.yaml | metalcloud-cli fabric configure-switches 5 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var config []byte
			var err error
			if fabricFlags.configSource != "" {
				config, err = utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
			} else {
				config, err = buildSwitchConfigFromFlags(cmd)
			}
			if err != nil {
				return err
			}
			return fabric.FabricConfigureSwitches(cmd.Context(), args[0], config, fabricFlags.dryRun)
		},
	}

	fabricConfigureSwitchesExampleCmd = &cobra.Command{
		Use:     "configure-switches-example",
		Aliases: []string{"configure-switches-config-example"},
		Short:   "Show an example switch configuration for configure-switches",
		Long: `Print a commented, ready-to-edit example of the configuration accepted by
'fabric configure-switches'. The output is valid YAML; redirect it to a file or
pipe it straight into the command.

Examples:
  metalcloud-cli fabric configure-switches-example
  metalcloud-cli fabric configure-switches-example > fabric-config.yaml
  metalcloud-cli fabric configure-switches-example | metalcloud-cli fabric configure-switches 5 --config-source pipe --dry-run`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricConfigureSwitchesExample(cmd.Context())
		},
	}
)

// buildSwitchConfigFromFlags assembles a configure-switches YAML document from
// the individual --... flags that were set, then hands it to the same loader as
// --config-source. A section appears only if its enable flag or any of its
// sub-flags was provided (so presence semantics match the file path exactly).
func buildSwitchConfigFromFlags(cmd *cobra.Command) ([]byte, error) {
	f := cmd.Flags()
	cs := &configureSwitchesFlags
	doc := map[string]interface{}{}

	if f.Changed("ordering") {
		doc["ordering"] = cs.ordering
	}
	if f.Changed("enable-physical-ports") {
		doc["enablePhysicalPorts"] = cs.enablePhysicalPorts
	}
	if f.Changed("description-template") {
		doc["descriptionTemplate"] = cs.descriptionTemplate
	}

	// hostname
	hostname := map[string]interface{}{}
	hostnamePresent := f.Changed("hostname") && cs.hostname
	if f.Changed("hostname-leaf") {
		hostname["leaf"] = cs.hostnameLeaf
		hostnamePresent = true
	}
	if f.Changed("hostname-spine") {
		hostname["spine"] = cs.hostnameSpine
		hostnamePresent = true
	}
	if f.Changed("hostname-super-spine") {
		hostname["super_spine"] = cs.hostnameSuperSpine
		hostnamePresent = true
	}
	for _, position := range cs.hostnameSkip {
		hostname[position] = nil // explicit null => skip the position
		hostnamePresent = true
	}
	if hostnamePresent {
		doc["hostname"] = hostname
	}

	// asn
	asn := map[string]interface{}{}
	asnPresent := f.Changed("asn") && cs.asn
	if f.Changed("asn-leaf-start") {
		asn["leafStart"] = cs.asnLeafStart
		asnPresent = true
	}
	if f.Changed("asn-spine-start") {
		asn["spineStart"] = cs.asnSpineStart
		asnPresent = true
	}
	if f.Changed("asn-super-spine-start") {
		asn["superSpineStart"] = cs.asnSuperSpineStart
		asnPresent = true
	}
	if asnPresent {
		doc["asn"] = asn
	}

	// loopback
	loopback := map[string]interface{}{}
	loopbackPresent := f.Changed("loopback") && cs.loopback
	if f.Changed("loopback-subnet") {
		loopback["subnet"] = cs.loopbackSubnet
		loopbackPresent = true
	}
	if loopbackPresent {
		doc["loopback"] = loopback
	}

	// topology
	topology := map[string]interface{}{}
	if leafSpine, present, err := buildLayerFlags(f, "topology-leaf-spine", cs.topoLeafSpine, "topology-leaf-spine-links-per-pair", cs.topoLeafSpineLPP); err != nil {
		return nil, err
	} else if present {
		topology["leafSpine"] = leafSpine
	}
	if spineSsp, present, err := buildLayerFlags(f, "topology-spine-super-spine", cs.topoSpineSuperSpine, "topology-spine-super-spine-links-per-pair", cs.topoSpineSuperSpineLPP); err != nil {
		return nil, err
	} else if present {
		topology["spineSuperSpine"] = spineSsp
	}
	leafHost := map[string]interface{}{}
	leafHostPresent := f.Changed("topology-leaf-host") && cs.topoLeafHost
	if f.Changed("topology-leaf-host-node-count") {
		leafHost["nodeCount"] = cs.topoLeafHostNodeCount
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-nodes") {
		leafHost["nodes"] = cs.topoLeafHostNodes
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-port-pattern") {
		leafHost["portPattern"] = cs.topoLeafHostPortPattern
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-nic-names") {
		leafHost["nicNames"] = cs.topoLeafHostNicNames
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-description-template") {
		leafHost["descriptionTemplate"] = cs.topoLeafHostDescription
		leafHostPresent = true
	}
	if leafHostPresent {
		topology["leafHost"] = leafHost
	}
	if len(topology) > 0 {
		doc["topology"] = topology
	}

	// p2p
	p2p := map[string]interface{}{}
	p2pPresent := f.Changed("p2p") && cs.p2p
	pools := map[string]interface{}{}
	if f.Changed("p2p-pool-leaf-spine") {
		pools["leafSpine"] = cs.p2pPoolLeafSpine
	}
	if f.Changed("p2p-pool-spine-super-spine") {
		pools["spineSuperSpine"] = cs.p2pPoolSpineSuperSpine
	}
	if f.Changed("p2p-pool-leaf-host") {
		pools["leafHost"] = cs.p2pPoolLeafHost
	}
	if len(pools) > 0 {
		p2p["pools"] = pools
		p2pPresent = true
	}
	if f.Changed("p2p-mtu") {
		p2p["mtu"] = cs.p2pMtu
		p2pPresent = true
	}
	if p2pPresent {
		doc["p2p"] = p2p
	}

	if len(doc) == 0 {
		return nil, fmt.Errorf("specify --config-source or at least one configuration flag (see 'fabric configure-switches --help')")
	}
	return yaml.Marshal(doc)
}

// buildLayerFlags builds a topology fabric-layer (leafSpine / spineSuperSpine)
// section from its enable flag and its links-per-pair flag.
func buildLayerFlags(f interface{ Changed(string) bool }, enableName string, enabled bool, lppName, lpp string) (map[string]interface{}, bool, error) {
	layer := map[string]interface{}{}
	present := f.Changed(enableName) && enabled
	if f.Changed(lppName) {
		value, err := parseLinksPerPair(lpp)
		if err != nil {
			return nil, false, err
		}
		layer["linksPerPair"] = value
		present = true
	}
	return layer, present, nil
}

// parseLinksPerPair maps the CLI string to what the loader expects: the literal
// "auto", or an integer.
func parseLinksPerPair(s string) (interface{}, error) {
	if s == "" || s == "auto" {
		return "auto", nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("links-per-pair must be 'auto' or an integer, got %q", s)
	}
	return n, nil
}

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

	fabricCmd.AddCommand(fabricConfigureSwitchesCmd)
	csCmd := fabricConfigureSwitchesCmd
	csCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the switch configuration. Can be 'pipe' or path to a YAML/JSON file. Mutually exclusive with the per-property flags below.")
	csCmd.Flags().BoolVar(&fabricFlags.dryRun, "dry-run", false, "Compute and preview the plan without making any changes.")

	cs := &configureSwitchesFlags
	csCmd.Flags().StringVar(&cs.ordering, "ordering", "managementAddress", "Device ordering: managementAddress | identifierString | id.")
	csCmd.Flags().BoolVar(&cs.enablePhysicalPorts, "enable-physical-ports", true, "Enable every physical port's staged config.")
	csCmd.Flags().StringVar(&cs.descriptionTemplate, "description-template", "", "Interface description template (placeholders {peerHostname}, {peerPort}). Requires a topology section.")

	csCmd.Flags().BoolVar(&cs.hostname, "hostname", false, "Enable hostname computation using the built-in reference templates.")
	csCmd.Flags().StringVar(&cs.hostnameLeaf, "hostname-leaf", "", "Hostname template for leaf devices.")
	csCmd.Flags().StringVar(&cs.hostnameSpine, "hostname-spine", "", "Hostname template for spine devices.")
	csCmd.Flags().StringVar(&cs.hostnameSuperSpine, "hostname-super-spine", "", "Hostname template for super_spine devices.")
	csCmd.Flags().StringSliceVar(&cs.hostnameSkip, "hostname-skip", nil, "Positions to skip (set to null), e.g. spine.")

	csCmd.Flags().BoolVar(&cs.asn, "asn", false, "Enable ASN assignment using the default starts.")
	csCmd.Flags().Int64Var(&cs.asnLeafStart, "asn-leaf-start", 0, "Starting ASN for leaves.")
	csCmd.Flags().Int64Var(&cs.asnSpineStart, "asn-spine-start", 0, "Starting ASN for spine groups.")
	csCmd.Flags().Int64Var(&cs.asnSuperSpineStart, "asn-super-spine-start", 0, "Shared ASN for superspines.")

	csCmd.Flags().BoolVar(&cs.loopback, "loopback", false, "Enable loopback IP allocation using the default subnet.")
	csCmd.Flags().StringVar(&cs.loopbackSubnet, "loopback-subnet", "", "Pool the loopback /32s are carved from.")

	csCmd.Flags().BoolVar(&cs.topoLeafSpine, "topology-leaf-spine", false, "Enable leaf<->spine pairing.")
	csCmd.Flags().StringVar(&cs.topoLeafSpineLPP, "topology-leaf-spine-links-per-pair", "", "Leaf<->spine links per pair: 'auto' or an integer.")
	csCmd.Flags().BoolVar(&cs.topoSpineSuperSpine, "topology-spine-super-spine", false, "Enable spine<->superspine pairing (3-tier only).")
	csCmd.Flags().StringVar(&cs.topoSpineSuperSpineLPP, "topology-spine-super-spine-links-per-pair", "", "Spine<->superspine links per pair: 'auto' or an integer.")

	csCmd.Flags().BoolVar(&cs.topoLeafHost, "topology-leaf-host", false, "Enable leaf->host downlinks.")
	csCmd.Flags().IntVar(&cs.topoLeafHostNodeCount, "topology-leaf-host-node-count", 0, "Number of host port-pairs per leaf.")
	csCmd.Flags().IntSliceVar(&cs.topoLeafHostNodes, "topology-leaf-host-nodes", nil, "Exact 0-based node indices (mutually exclusive with node-count).")
	csCmd.Flags().StringVar(&cs.topoLeafHostPortPattern, "topology-leaf-host-port-pattern", "", "Leaf host port pattern, e.g. swp{port}s{sub}.")
	csCmd.Flags().StringSliceVar(&cs.topoLeafHostNicNames, "topology-leaf-host-nic-names", nil, "Remote host NIC names (even count).")
	csCmd.Flags().StringVar(&cs.topoLeafHostDescription, "topology-leaf-host-description-template", "", "Leaf->host description template.")

	csCmd.Flags().BoolVar(&cs.p2p, "p2p", false, "Enable point-to-point link creation with reference default pools.")
	csCmd.Flags().StringVar(&cs.p2pPoolLeafSpine, "p2p-pool-leaf-spine", "", "Leaf<->spine /31 pool.")
	csCmd.Flags().StringVar(&cs.p2pPoolSpineSuperSpine, "p2p-pool-spine-super-spine", "", "Spine<->superspine /31 pool.")
	csCmd.Flags().StringVar(&cs.p2pPoolLeafHost, "p2p-pool-leaf-host", "", "Leaf->host /31 pool.")
	csCmd.Flags().Int32Var(&cs.p2pMtu, "p2p-mtu", 0, "MTU applied to created links.")

	// --config-source is mutually exclusive with each per-property flag; the
	// per-property flags can be combined freely with one another.
	for _, name := range configureSwitchesDetailFlags {
		csCmd.MarkFlagsMutuallyExclusive("config-source", name)
	}

	fabricCmd.AddCommand(fabricConfigureSwitchesExampleCmd)
}
