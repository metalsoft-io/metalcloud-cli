package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_endpoint_group"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	networkEndpointGroupFlags = struct {
		configSource            string
		filterId                []string
		filterName              []string
		filterSiteId            []string
		name                    string
		siteId                  int64
		accessMode              string
		tagged                  bool
		mtu                     int32
		providesDefaultRoute    bool
		disableAutoIpAllocation bool
	}{}

	networkEndpointGroupCmd = &cobra.Command{
		Use:     "network-endpoint-group [command]",
		Aliases: []string{"neg"},
		Short:   "Network endpoint group management",
		Long: `Manage network endpoint groups: named sets of logical network connections that
server instance groups and VM instance groups attach to.

Command categories:
  Lifecycle:        list, get, create, update, delete, config-example
  Logical networks: logical-network list|get|add|update|remove`,
	}

	networkEndpointGroupListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List network endpoint groups",
		Long: `List all network endpoint groups.

Optional Flags:
  --filter-id strings        Filter by network endpoint group ID. Repeatable or comma-separated.
  --filter-name strings      Filter by name. Repeatable or comma-separated.
  --filter-site-id strings   Filter by site ID. Repeatable or comma-separated.

Examples:
  metalcloud network-endpoint-group list
  metalcloud neg list --filter-site-id 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupList(cmd.Context(), network_endpoint_group.NetworkEndpointGroupListFilters{
				Id:     networkEndpointGroupFlags.filterId,
				Name:   networkEndpointGroupFlags.filterName,
				SiteId: networkEndpointGroupFlags.filterSiteId,
			})
		},
	}

	networkEndpointGroupGetCmd = &cobra.Command{
		Use:     "get network_endpoint_group_id_or_name",
		Aliases: []string{"show"},
		Short:   "Get network endpoint group details",
		Long: `Get the details of a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Examples:
  metalcloud network-endpoint-group get 3
  metalcloud neg get dc1-endpoint-group`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupGet(cmd.Context(), args[0])
		},
	}

	networkEndpointGroupConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example network endpoint group create configuration",
		Long: `Print an example configuration that can be edited and passed to
'network-endpoint-group create --config-source'.

Examples:
  metalcloud network-endpoint-group config-example > neg.json
  metalcloud neg config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupConfigExample(cmd.Context())
		},
	}

	networkEndpointGroupCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a network endpoint group",
		Long: `Create a new network endpoint group.

Required Flags (one of):
  --config-source string   Source of the configuration. Can be 'pipe' or path to a JSON/YAML file.
  --name string            Name of the new network endpoint group

Optional Flags (when not using --config-source):
  --site-id int            ID of the site the network endpoint group belongs to

Examples:
  metalcloud network-endpoint-group create --config-source neg.json
  metalcloud neg create --name dc1-endpoint-group --site-id 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateNetworkEndpointGroup

			if networkEndpointGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(networkEndpointGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.CreateNetworkEndpointGroup{
					Name: networkEndpointGroupFlags.name,
				}
				if networkEndpointGroupFlags.siteId != 0 {
					create.SiteId = sdk.PtrInt64(networkEndpointGroupFlags.siteId)
				}
			}

			return network_endpoint_group.NetworkEndpointGroupCreate(cmd.Context(), create)
		},
	}

	networkEndpointGroupUpdateCmd = &cobra.Command{
		Use:     "update network_endpoint_group_id_or_name",
		Aliases: []string{"edit"},
		Short:   "Update a network endpoint group",
		Long: `Update the name or site of a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Required Flags:
  --config-source string              Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud network-endpoint-group update 3 --config-source update.json
  echo '{"name":"new name"}' | metalcloud neg update dc1-endpoint-group --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(networkEndpointGroupFlags.configSource)
			if err != nil {
				return err
			}

			return network_endpoint_group.NetworkEndpointGroupUpdate(cmd.Context(), args[0], config)
		},
	}

	networkEndpointGroupDeleteCmd = &cobra.Command{
		Use:     "delete network_endpoint_group_id_or_name",
		Aliases: []string{"rm"},
		Short:   "Delete a network endpoint group",
		Long: `Delete a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Examples:
  metalcloud network-endpoint-group delete 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupDelete(cmd.Context(), args[0])
		},
	}

	networkEndpointGroupLogicalNetworkCmd = &cobra.Command{
		Use:     "logical-network [command]",
		Aliases: []string{"logical-networks"},
		Short:   "Network endpoint group logical network management",
		Long: `Manage the logical networks attached to a network endpoint group.

Commands: list, get, add, update, remove`,
	}

	networkEndpointGroupLogicalNetworkListCmd = &cobra.Command{
		Use:     "list network_endpoint_group_id_or_name",
		Aliases: []string{"ls"},
		Short:   "List the logical networks of a network endpoint group",
		Long: `List the logical networks attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Examples:
  metalcloud network-endpoint-group logical-network list 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupLogicalNetworkList(cmd.Context(), args[0])
		},
	}

	networkEndpointGroupLogicalNetworkGetCmd = &cobra.Command{
		Use:     "get network_endpoint_group_id_or_name logical_network_id",
		Aliases: []string{"show"},
		Short:   "Get one logical network of a network endpoint group",
		Long: `Get the settings of one logical network attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Examples:
  metalcloud network-endpoint-group logical-network get 3 44`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupLogicalNetworkGet(cmd.Context(), args[0], args[1])
		},
	}

	networkEndpointGroupLogicalNetworkAddCmd = &cobra.Command{
		Use:     "add network_endpoint_group_id_or_name logical_network_id",
		Aliases: []string{"new", "create"},
		Short:   "Attach a logical network to a network endpoint group",
		Long: `Attach a logical network to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Optional Flags:
  --config-source string        Source of the connection configuration. Can be 'pipe' or path to a JSON/YAML file.
                                Mutually exclusive with the flags below.
  --access-mode string          Access mode of the connection (default "l2")
  --tagged                      Attach the logical network tagged
  --mtu int                     MTU of the connection
  --provides-default-route      The logical network provides the default route
  --disable-auto-ip-allocation  Disable automatic IPv4 allocation on this connection

Examples:
  metalcloud network-endpoint-group logical-network add 3 44 --tagged
  metalcloud neg logical-network add dc1-endpoint-group 44 --config-source connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateNetworkEndpointGroupLogicalNetwork

			if networkEndpointGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(networkEndpointGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				accessMode, err := sdk.NewNetworkEndpointGroupAllowedAccessModeFromValue(networkEndpointGroupFlags.accessMode)
				if err != nil {
					return err
				}

				create = sdk.CreateNetworkEndpointGroupLogicalNetwork{
					AccessMode: *accessMode,
					Tagged:     networkEndpointGroupFlags.tagged,
				}
				if networkEndpointGroupFlags.mtu != 0 {
					create.Mtu = sdk.PtrInt32(networkEndpointGroupFlags.mtu)
				}
				if cmd.Flags().Changed("provides-default-route") {
					create.ProvidesDefaultRoute = sdk.PtrBool(networkEndpointGroupFlags.providesDefaultRoute)
				}
				if cmd.Flags().Changed("disable-auto-ip-allocation") {
					create.DisableAutoIpAllocation = sdk.PtrBool(networkEndpointGroupFlags.disableAutoIpAllocation)
				}
			}

			// The logical network is always taken from the positional argument
			// so that the command and the payload cannot disagree.
			create.LogicalNetworkId = args[1]

			return network_endpoint_group.NetworkEndpointGroupLogicalNetworkAdd(cmd.Context(), args[0], create)
		},
	}

	networkEndpointGroupLogicalNetworkUpdateCmd = &cobra.Command{
		Use:     "update network_endpoint_group_id_or_name logical_network_id",
		Aliases: []string{"edit"},
		Short:   "Update a logical network of a network endpoint group",
		Long: `Update the settings of one logical network attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Optional Flags:
  --config-source string        Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.
                                Mutually exclusive with the flags below.
  --access-mode string          Access mode of the connection
  --tagged                      Attach the logical network tagged
  --mtu int                     MTU of the connection
  --provides-default-route      The logical network provides the default route
  --disable-auto-ip-allocation  Disable automatic IPv4 allocation on this connection

Examples:
  metalcloud network-endpoint-group logical-network update 3 44 --mtu 9000
  metalcloud neg logical-network update dc1-endpoint-group 44 --config-source connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var update sdk.UpdateNetworkEndpointGroupLogicalNetwork

			if networkEndpointGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(networkEndpointGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &update); err != nil {
					return err
				}
			} else {
				if cmd.Flags().Changed("access-mode") {
					accessMode, err := sdk.NewNetworkEndpointGroupAllowedAccessModeFromValue(networkEndpointGroupFlags.accessMode)
					if err != nil {
						return err
					}
					update.AccessMode = accessMode
				}
				if cmd.Flags().Changed("tagged") {
					update.Tagged = sdk.PtrBool(networkEndpointGroupFlags.tagged)
				}
				if networkEndpointGroupFlags.mtu != 0 {
					update.Mtu = sdk.PtrInt32(networkEndpointGroupFlags.mtu)
				}
				if cmd.Flags().Changed("provides-default-route") {
					update.ProvidesDefaultRoute = sdk.PtrBool(networkEndpointGroupFlags.providesDefaultRoute)
				}
				if cmd.Flags().Changed("disable-auto-ip-allocation") {
					update.DisableAutoIpAllocation = sdk.PtrBool(networkEndpointGroupFlags.disableAutoIpAllocation)
				}
			}

			return network_endpoint_group.NetworkEndpointGroupLogicalNetworkUpdate(cmd.Context(), args[0], args[1], update)
		},
	}

	networkEndpointGroupLogicalNetworkRemoveCmd = &cobra.Command{
		Use:     "remove network_endpoint_group_id_or_name logical_network_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Detach a logical network from a network endpoint group",
		Long: `Detach a logical network from a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Examples:
  metalcloud network-endpoint-group logical-network remove 3 44`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_ENDPOINT_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return network_endpoint_group.NetworkEndpointGroupLogicalNetworkRemove(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(networkEndpointGroupCmd)

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupListCmd)
	networkEndpointGroupListCmd.Flags().StringSliceVar(&networkEndpointGroupFlags.filterId, "filter-id", nil, "Filter by network endpoint group ID.")
	networkEndpointGroupListCmd.Flags().StringSliceVar(&networkEndpointGroupFlags.filterName, "filter-name", nil, "Filter by name.")
	networkEndpointGroupListCmd.Flags().StringSliceVar(&networkEndpointGroupFlags.filterSiteId, "filter-site-id", nil, "Filter by site ID.")

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupGetCmd)
	networkEndpointGroupCmd.AddCommand(networkEndpointGroupConfigExampleCmd)

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupCreateCmd)
	networkEndpointGroupCreateCmd.Flags().StringVar(&networkEndpointGroupFlags.configSource, "config-source", "", "Source of the new network endpoint group configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkEndpointGroupCreateCmd.Flags().StringVar(&networkEndpointGroupFlags.name, "name", "", "Name of the new network endpoint group.")
	networkEndpointGroupCreateCmd.Flags().Int64Var(&networkEndpointGroupFlags.siteId, "site-id", 0, "ID of the site the network endpoint group belongs to.")
	networkEndpointGroupCreateCmd.MarkFlagsOneRequired("config-source", "name")
	networkEndpointGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "name")

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupUpdateCmd)
	networkEndpointGroupUpdateCmd.Flags().StringVar(&networkEndpointGroupFlags.configSource, "config-source", "", "Source of the updated network endpoint group configuration. Can be 'pipe' or path to a JSON/YAML file.")
	networkEndpointGroupUpdateCmd.MarkFlagsOneRequired("config-source")

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupDeleteCmd)

	networkEndpointGroupCmd.AddCommand(networkEndpointGroupLogicalNetworkCmd)
	networkEndpointGroupLogicalNetworkCmd.AddCommand(networkEndpointGroupLogicalNetworkListCmd)
	networkEndpointGroupLogicalNetworkCmd.AddCommand(networkEndpointGroupLogicalNetworkGetCmd)

	networkEndpointGroupLogicalNetworkCmd.AddCommand(networkEndpointGroupLogicalNetworkAddCmd)
	registerNetworkEndpointGroupLogicalNetworkFlags(networkEndpointGroupLogicalNetworkAddCmd, "l2")

	networkEndpointGroupLogicalNetworkCmd.AddCommand(networkEndpointGroupLogicalNetworkUpdateCmd)
	registerNetworkEndpointGroupLogicalNetworkFlags(networkEndpointGroupLogicalNetworkUpdateCmd, "l2")

	networkEndpointGroupLogicalNetworkCmd.AddCommand(networkEndpointGroupLogicalNetworkRemoveCmd)
}

// registerNetworkEndpointGroupLogicalNetworkFlags wires the connection setting
// flags shared by the 'logical-network add' and 'logical-network update'
// sub-commands.
func registerNetworkEndpointGroupLogicalNetworkFlags(cmd *cobra.Command, defaultAccessMode string) {
	cmd.Flags().StringVar(&networkEndpointGroupFlags.configSource, "config-source", "", "Source of the connection configuration. Can be 'pipe' or path to a JSON/YAML file.")
	cmd.Flags().StringVar(&networkEndpointGroupFlags.accessMode, "access-mode", defaultAccessMode, "Access mode of the connection.")
	cmd.Flags().BoolVar(&networkEndpointGroupFlags.tagged, "tagged", false, "Attach the logical network tagged.")
	cmd.Flags().Int32Var(&networkEndpointGroupFlags.mtu, "mtu", 0, "MTU of the connection.")
	cmd.Flags().BoolVar(&networkEndpointGroupFlags.providesDefaultRoute, "provides-default-route", false, "The logical network provides the default route.")
	cmd.Flags().BoolVar(&networkEndpointGroupFlags.disableAutoIpAllocation, "disable-auto-ip-allocation", false, "Disable automatic IPv4 allocation on this connection.")
	cmd.MarkFlagsMutuallyExclusive("config-source", "access-mode")
	cmd.MarkFlagsMutuallyExclusive("config-source", "tagged")
	cmd.MarkFlagsMutuallyExclusive("config-source", "mtu")
	cmd.MarkFlagsMutuallyExclusive("config-source", "provides-default-route")
	cmd.MarkFlagsMutuallyExclusive("config-source", "disable-auto-ip-allocation")
}
