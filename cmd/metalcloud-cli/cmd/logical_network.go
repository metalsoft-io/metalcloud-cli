package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
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
  list         List logical networks with optional filtering
  get          Get detailed information about a specific logical network
  create       Create a new logical network from configuration
  update       Update an existing logical network
  delete       Delete a logical network
  config-example  Get example configuration for a specific network kind

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
  --filter-infrastructure-id Filter results by infrastructure ID(s) (can be used multiple times)
  --filter-kind              Filter results by network kind(s) like 'vlan', 'vxlan' (can be used multiple times)
  --sort-by                  Sort results by field(s) with direction (e.g., id:ASC, name:DESC)

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

  # Combine fabric filter with additional filters
  metalcloud-cli logical-network list fabric-1 --filter-kind vxlan --sort-by id:ASC`,
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
)

func init() {
	rootCmd.AddCommand(logicalNetworkCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkListCmd)
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterId, "filter-id", nil, "Filter by logical network ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterLabel, "filter-label", nil, "Filter by logical network label.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterFabricId, "filter-fabric-id", nil, "Filter by fabric ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterKind, "filter-kind", nil, "Filter by logical network kind.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., id:ASC, name:DESC).")

	logicalNetworkCmd.AddCommand(logicalNetworkGetCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkConfigExampleCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkCreateCmd)
	logicalNetworkCreateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the new logical network configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkCreateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkUpdateCmd)
	logicalNetworkUpdateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the logical network updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkUpdateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkDeleteCmd)
}
