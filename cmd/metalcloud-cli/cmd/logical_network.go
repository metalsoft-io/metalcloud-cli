package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	logicalNetworkFlags = struct {
		configSource           string
		filterId               []string
		filterLabel            []string
		filterFabricId         []string
		filterInfrastructureId []string
		filterKind             []string
		sortBy                 []string
		page                   int
		limit                  int
		profileId              int
		label                  string
		name                   string
		infrastructureId       int
		mtu                    int
	}{}

	logicalNetworkCmd = &cobra.Command{
		Use:     "logical-network [command]",
		Aliases: []string{"ln", "network", "logical_network"},
		Short:   "Manage logical networks within fabrics",
		Long: `Manage logical networks within fabrics for network segmentation and isolation.

Logical networks provide Layer 2 network isolation within a fabric, allowing you to create
separate broadcast domains for different applications or tenants. Each logical network is
associated with a fabric and can have specific configurations based on its kind (vlan, vxlan, etc.).

Available Commands:
  list                 List logical networks with optional filtering
  get                  Get detailed information about a specific logical network
  create               Create a new logical network from configuration
  create-from-profile  Create a logical network from a logical network profile
  update               Update an existing logical network
  delete               Delete a logical network
  config-example       Get example configuration for a specific network kind
  get-config           Get the config sub-resource of a logical network
  update-config        Update the global settings of a logical network config
  apply-profiles       Apply a logical network profile to a logical network config
  allocation-strategy  Manage the config's allocation strategies
  get-external-connections                  List the attached external connections
  get-external-connection-logical-networks  List the external connection attachments
  get-interconnects                         List the attached logical network interconnects
  detach-external-connection                Detach an external connection

Examples:
  # List all logical networks
  metalcloud-cli logical-network list

  # List logical networks in a specific fabric
  metalcloud-cli logical-network list fabric-1

  # Create a VLAN logical network
  metalcloud-cli logical-network create vlan --config-source config.json`,
	}

	logicalNetworkListCmd = &cobra.Command{
		Use:     "list [fabric_id_or_label]",
		Aliases: []string{"ls"},
		Short:   "List logical networks with optional filtering and sorting",
		Long: `List all logical networks, optionally filtered by fabric and other criteria.

This command displays logical networks in a tabular format. You can optionally provide
a fabric ID or label to filter results to networks within that specific fabric.

Arguments:
  fabric_id_or_label  Optional fabric identifier to filter networks (can be ID or label)

Flags:
  --filter-id                Filter results by logical network ID(s) (can be used multiple times)
  --filter-label             Filter results by logical network label(s) (can be used multiple times) 
  --filter-fabric-id         Filter results by fabric ID(s) (can be used multiple times)
  --filter-infrastructure-id Filter results by infrastructure ID(s) (can be used multiple times). Use 'null' to filter public logical networks
  --filter-kind              Filter results by network kind(s) like 'vlan', 'vxlan' (can be used multiple times)
  --sort-by                  Sort results by field(s) with direction (e.g., id:ASC, name:DESC)
  --page                     Page number to retrieve (default: all records)
  --limit                    Number of records per page (default: all records)

Examples:
  # List all logical networks
  metalcloud-cli logical-network list

  # List networks in a specific fabric
  metalcloud-cli logical-network list fabric-production

  # Filter by network kind
  metalcloud-cli logical-network list --filter-kind vlan

  # Filter by multiple criteria
  metalcloud-cli logical-network list --filter-kind vlan --filter-label test

  # Sort by name descending
  metalcloud-cli logical-network list --sort-by name:DESC

  # Paginate results (get page 2 with 50 records per page)
  metalcloud-cli logical-network list --page 2 --limit 50

  # Combine fabric filter with additional filters and pagination
  metalcloud-cli logical-network list fabric-1 --filter-kind vxlan --sort-by id:ASC --page 1 --limit 10`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fabricIdOrLabel := ""
			if len(args) > 0 {
				fabricIdOrLabel = args[0]
			}

			return logical_network.LogicalNetworkList(cmd.Context(), fabricIdOrLabel, logical_network.ListFlags{
				FilterId:               logicalNetworkFlags.filterId,
				FilterLabel:            logicalNetworkFlags.filterLabel,
				FilterFabricId:         logicalNetworkFlags.filterFabricId,
				FilterInfrastructureId: logicalNetworkFlags.filterInfrastructureId,
				FilterKind:             logicalNetworkFlags.filterKind,
				SortBy:                 logicalNetworkFlags.sortBy,
				Page:                   logicalNetworkFlags.page,
				Limit:                  logicalNetworkFlags.limit,
			})
		},
	}

	logicalNetworkGetCmd = &cobra.Command{
		Use:     "get logical_network_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a logical network",
		Long: `Display detailed information about a specific logical network including its configuration,
associated fabric, network kind, and other properties.

Arguments:
  logical_network_id  The unique identifier of the logical network to retrieve (required)

The command shows comprehensive details including:
- Network identification (ID, label)
- Associated fabric and infrastructure
- Network kind and configuration
- Creation and modification timestamps
- Current status and operational state

Examples:
  # Get details of a logical network by ID
  metalcloud-cli logical-network get 12345

  # Get details using the 'show' alias
  metalcloud-cli logical-network show network-production-vlan

  # Use with pipe or redirect for further processing
  metalcloud-cli logical-network get 12345 | jq .`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkConfigExampleCmd = &cobra.Command{
		Use:     "config-example kind",
		Aliases: []string{"example"},
		Short:   "Generate example configuration for a logical network kind",
		Long: `Generate example configuration templates for different logical network kinds.

This command provides sample JSON configurations that can be used as templates when
creating logical networks. The configuration examples show the structure and required
fields for each network kind.

Arguments:
  kind  The type of logical network for which to generate example configuration
        Supported kinds include: vlan, vxlan, flat, and others

The generated configuration can be used with the 'create' command by saving it to a file
and using the --config-source flag, or by piping it directly.

Examples:
  # Get example configuration for a VLAN network
  metalcloud-cli logical-network config-example vlan

  # Save example to file for editing
  metalcloud-cli logical-network config-example vxlan > network-config.json

  # Use with create command via pipe
  metalcloud-cli logical-network config-example vlan | metalcloud-cli logical-network create vlan --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkConfigExample(cmd.Context(), args[0])
		},
	}

	logicalNetworkCreateCmd = &cobra.Command{
		Use:     "create kind",
		Aliases: []string{"new"},
		Short:   "Create a new logical network from configuration",
		Long: `Create a new logical network of the specified kind using configuration from a file or pipe.

This command creates a logical network with the provided configuration. The configuration
must be in JSON format and contain all required fields for the specified network kind.

Arguments:
  kind  The type of logical network to create (e.g., vlan, vxlan, flat)

Required Flags:
  --config-source  Source of the logical network configuration (required)
                   Can be 'pipe' to read from stdin or path to a JSON file

Configuration Format:
The configuration file must contain a JSON object with the network specification.
Use 'config-example' command to see the expected structure for each kind.

Examples:
  # Create from a configuration file
  metalcloud-cli logical-network create vlan --config-source network.json

  # Create using pipe input
  cat network.json | metalcloud-cli logical-network create vlan --config-source pipe

  # Create using generated example (edit as needed)
  metalcloud-cli logical-network config-example vlan > config.json
  # Edit config.json with your values
  metalcloud-cli logical-network create vlan --config-source config.json

  # One-liner with example and pipe
  metalcloud-cli logical-network config-example vxlan | \
    jq '.label = "my-network"' | \
    metalcloud-cli logical-network create vxlan --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkCreate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkUpdateCmd = &cobra.Command{
		Use:     "update logical_network_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing logical network configuration",
		Long: `Update an existing logical network with new configuration from a file or pipe.

This command updates the configuration of an existing logical network. The configuration
must be in JSON format and can contain partial updates or complete new configuration.

Arguments:
  logical_network_id  The unique identifier of the logical network to update (required)

Required Flags:
  --config-source  Source of the logical network configuration updates (required)
                   Can be 'pipe' to read from stdin or path to a JSON file

Configuration Format:
The configuration file must contain a JSON object with the network specification updates.
You can provide partial updates (only the fields you want to change) or complete configuration.

Examples:
  # Update from a configuration file
  metalcloud-cli logical-network update 12345 --config-source updates.json

  # Update using pipe input
  cat updates.json | metalcloud-cli logical-network update 12345 --config-source pipe

  # Update specific field using jq and pipe
  echo '{"label": "new-network-name"}' | metalcloud-cli logical-network update 12345 --config-source pipe

  # Get current config, edit, and update
  metalcloud-cli logical-network get 12345 --output json > current.json
  # Edit current.json with your changes
  metalcloud-cli logical-network update 12345 --config-source current.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkUpdate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkDeleteCmd = &cobra.Command{
		Use:     "delete logical_network_id",
		Aliases: []string{"rm"},
		Short:   "Delete a logical network",
		Long: `Delete a specific logical network by its unique identifier.

This command permanently removes a logical network from the system. The deletion is
irreversible, so use with caution. Make sure the logical network is not in use by
any resources before attempting to delete it.

Arguments:
  logical_network_id  The unique identifier of the logical network to delete (required)

Warning:
- This operation is irreversible
- Ensure the logical network is not referenced by other resources
- Any dependent configurations may need to be updated after deletion

Examples:
  # Delete a logical network by ID
  metalcloud-cli logical-network delete 12345

  # Delete using the 'rm' alias
  metalcloud-cli logical-network rm network-test-vlan

  # Confirm deletion with output redirection
  metalcloud-cli logical-network delete 12345 2>&1 | tee delete.log`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkDelete(cmd.Context(), args[0])
		},
	}

	logicalNetworkGetConfigCmd = &cobra.Command{
		Use:     "get-config logical_network_id",
		Aliases: []string{"config", "show-config"},
		Short:   "Get the config of a logical network",
		Long: `Display the config object of a logical network.

The config is a separate sub-resource holding the desired state of the network:
its kind, MTU, deploy type and status, and the allocation strategies that the
'logical-network allocation-strategy' commands manage. It carries its own
revision, distinct from the logical network's.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-config 12
  metalcloud-cli logical-network get-config 12 -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkConfigGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkUpdateConfigCmd = &cobra.Command{
		Use:     "update-config logical_network_id",
		Aliases: []string{"edit-config"},
		Short:   "Update the global settings of a logical network config",
		Long: `Update the global settings of a logical network's config.

Only the config's global settings are updated here (for example the MTU and the
VXLAN properties); the allocation strategies are managed with the
'logical-network allocation-strategy' commands. The config's own revision is
sent as the If-Match entity tag, so a concurrent change is rejected instead of
being overwritten.

Required Arguments:
  logical_network_id  The ID of the logical network

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli logical-network update-config 12 --config-source settings.json
  echo '{"mtu":9000}' | metalcloud-cli logical-network update-config 12 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkConfigUpdate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkApplyProfilesCmd = &cobra.Command{
		Use:     "apply-profiles logical_network_id profile_id",
		Aliases: []string{"apply-profile"},
		Short:   "Apply a logical network profile to a logical network config",
		Long: `Apply a logical network profile onto the config of an existing logical network.

The profile's allocation strategies replace the ones currently held by the
config. The config's own revision is sent as the If-Match entity tag.

Required Arguments:
  logical_network_id  The ID of the logical network
  profile_id          The ID of the logical network profile to apply

Examples:
  metalcloud-cli logical-network apply-profiles 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkApplyProfiles(cmd.Context(), args[0], args[1])
		},
	}

	logicalNetworkCreateFromProfileCmd = &cobra.Command{
		Use:     "create-from-profile",
		Aliases: []string{"new-from-profile"},
		Short:   "Create a logical network from a logical network profile",
		Long: `Create a logical network whose configuration is taken from a logical network profile.

The network can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

Required Flags (one of):
  --profile-id       The ID of the logical network profile to instantiate
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --label              The label of the new logical network
  --name               The name of the new logical network
  --infrastructure-id  The infrastructure the network belongs to
  --mtu                Maximum Transmission Unit in bytes

Examples:
  metalcloud-cli logical-network create-from-profile --profile-id 3 --label my-network
  metalcloud-cli logical-network create-from-profile --profile-id 3 --infrastructure-id 7 --mtu 9000
  metalcloud-cli logical-network create-from-profile --config-source network.json
  echo '{"logicalNetworkProfileId":3,"label":"my-network"}' | \
    metalcloud-cli logical-network create-from-profile --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var createConfig sdk.CreateLogicalNetworkFromProfile

			if logicalNetworkFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &createConfig); err != nil {
					return err
				}
			} else {
				createConfig.LogicalNetworkProfileId = int64(logicalNetworkFlags.profileId)

				if logicalNetworkFlags.label != "" {
					createConfig.Label = &logicalNetworkFlags.label
				}
				if logicalNetworkFlags.name != "" {
					createConfig.Name = &logicalNetworkFlags.name
				}
				if logicalNetworkFlags.infrastructureId != 0 {
					createConfig.InfrastructureId = *sdk.NewNullableInt64(sdk.PtrInt64(int64(logicalNetworkFlags.infrastructureId)))
				}
				if logicalNetworkFlags.mtu != 0 {
					createConfig.Mtu = *sdk.NewNullableInt32(sdk.PtrInt32(int32(logicalNetworkFlags.mtu)))
				}
			}

			return logical_network.LogicalNetworkCreateFromProfile(cmd.Context(), createConfig)
		},
	}

	logicalNetworkGetExternalConnectionsCmd = &cobra.Command{
		Use:     "get-external-connections logical_network_id",
		Aliases: []string{"external-connections"},
		Short:   "List the external connections attached to a logical network",
		Long: `List the external connections attached to a logical network.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-external-connections 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkExternalConnections(cmd.Context(), args[0])
		},
	}

	logicalNetworkGetExternalConnectionLogicalNetworksCmd = &cobra.Command{
		Use:     "get-external-connection-logical-networks logical_network_id",
		Aliases: []string{"external-connection-logical-networks"},
		Short:   "List the external connection attachments of a logical network",
		Long: `List the external connection logical networks of a logical network.

Each record is one attachment between the logical network and an external
connection, with the status of that attachment.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-external-connection-logical-networks 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkExternalConnectionLogicalNetworks(cmd.Context(), args[0])
		},
	}

	logicalNetworkGetInterconnectsCmd = &cobra.Command{
		Use:     "get-interconnects logical_network_id",
		Aliases: []string{"interconnects"},
		Short:   "List the logical network interconnects of a logical network",
		Long: `List the logical network interconnects attached to a logical network.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-interconnects 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkInterconnects(cmd.Context(), args[0])
		},
	}

	logicalNetworkDetachExternalConnectionCmd = &cobra.Command{
		Use:     "detach-external-connection logical_network_id external_connection_id",
		Aliases: []string{"detach-external"},
		Short:   "Detach an external connection from a logical network",
		Long: `Detach an external connection from a logical network.

Required Arguments:
  logical_network_id      The ID of the logical network
  external_connection_id  The ID of the external connection to detach

Examples:
  metalcloud-cli logical-network detach-external-connection 12 4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkDetachExternalConnection(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(logicalNetworkCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkListCmd)
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterId, "filter-id", nil, "Filter by logical network ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterLabel, "filter-label", nil, "Filter by logical network label.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterFabricId, "filter-fabric-id", nil, "Filter by fabric ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID. Use 'null' to filter public logical networks.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterKind, "filter-kind", nil, "Filter by logical network kind.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., id:ASC, name:DESC).")
	logicalNetworkListCmd.Flags().IntVar(&logicalNetworkFlags.page, "page", 0, "Page number to retrieve (default: return all records).")
	logicalNetworkListCmd.Flags().IntVar(&logicalNetworkFlags.limit, "limit", 0, "Number of records per page (default: return all records).")

	logicalNetworkCmd.AddCommand(logicalNetworkGetCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkConfigExampleCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkCreateCmd)
	logicalNetworkCreateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the new logical network configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkCreateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkUpdateCmd)
	logicalNetworkUpdateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the logical network updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkUpdateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkDeleteCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkGetConfigCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkUpdateConfigCmd)
	logicalNetworkUpdateConfigCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the logical network config updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkApplyProfilesCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkCreateFromProfileCmd)
	logicalNetworkCreateFromProfileCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the create-from-profile configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkCreateFromProfileCmd.Flags().IntVar(&logicalNetworkFlags.profileId, "profile-id", 0, "The ID of the logical network profile to instantiate.")
	logicalNetworkCreateFromProfileCmd.Flags().StringVar(&logicalNetworkFlags.label, "label", "", "The label of the new logical network.")
	logicalNetworkCreateFromProfileCmd.Flags().StringVar(&logicalNetworkFlags.name, "name", "", "The name of the new logical network.")
	logicalNetworkCreateFromProfileCmd.Flags().IntVar(&logicalNetworkFlags.infrastructureId, "infrastructure-id", 0, "The infrastructure the new logical network belongs to.")
	logicalNetworkCreateFromProfileCmd.Flags().IntVar(&logicalNetworkFlags.mtu, "mtu", 0, "Maximum Transmission Unit (MTU) in bytes.")
	logicalNetworkCreateFromProfileCmd.MarkFlagsOneRequired("config-source", "profile-id")
	logicalNetworkCreateFromProfileCmd.MarkFlagsMutuallyExclusive("config-source", "profile-id")
	logicalNetworkCreateFromProfileCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	logicalNetworkCreateFromProfileCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	logicalNetworkCreateFromProfileCmd.MarkFlagsMutuallyExclusive("config-source", "infrastructure-id")
	logicalNetworkCreateFromProfileCmd.MarkFlagsMutuallyExclusive("config-source", "mtu")

	logicalNetworkCmd.AddCommand(logicalNetworkGetExternalConnectionsCmd)
	logicalNetworkCmd.AddCommand(logicalNetworkGetExternalConnectionLogicalNetworksCmd)
	logicalNetworkCmd.AddCommand(logicalNetworkGetInterconnectsCmd)
	logicalNetworkCmd.AddCommand(logicalNetworkDetachExternalConnectionCmd)
}
