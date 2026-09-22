package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network_interconnect"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	logicalNetworkInterconnectFlags = struct {
		configSource               string
		filterId                   []string
		filterLabel                []string
		filterName                 []string
		filterKind                 []string
		filterStatus               []string
		filterFabricInterconnectId []string
		filterLogicalNetworkId     []string
		label                      string
		name                       string
		kind                       string
		fabricInterconnectId       int64
	}{}

	logicalNetworkInterconnectCmd = &cobra.Command{
		Use:     "logical-network-interconnect [command]",
		Aliases: []string{"ln-interconnect", "lni"},
		Short:   "Logical network interconnect management",
		Long: `Manage logical network interconnects: the per-logical-network stretch of a
network fabric interconnect, extending a logical network across the interconnected fabrics.

Command categories:
  Lifecycle:  list, get, create, update, delete, config-example
  Links:      link list|get|add|remove`,
	}

	logicalNetworkInterconnectListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List logical network interconnects",
		Long: `List all logical network interconnects.

Optional Flags:
  --filter-id strings                      Filter by ID. Repeatable or comma-separated.
  --filter-label strings                   Filter by label. Repeatable or comma-separated.
  --filter-name strings                    Filter by name. Repeatable or comma-separated.
  --filter-kind strings                    Filter by kind (e.g. dci-evpn). Repeatable or comma-separated.
  --filter-status strings                  Filter by status. Repeatable or comma-separated.
  --filter-fabric-interconnect-id strings  Filter by network fabric interconnect ID. Repeatable or comma-separated.

Examples:
  metalcloud logical-network-interconnect list
  metalcloud lni list --filter-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectList(cmd.Context(), logical_network_interconnect.LogicalNetworkInterconnectListFilters{
				Id:                   logicalNetworkInterconnectFlags.filterId,
				Label:                logicalNetworkInterconnectFlags.filterLabel,
				Name:                 logicalNetworkInterconnectFlags.filterName,
				Kind:                 logicalNetworkInterconnectFlags.filterKind,
				Status:               logicalNetworkInterconnectFlags.filterStatus,
				FabricInterconnectId: logicalNetworkInterconnectFlags.filterFabricInterconnectId,
			})
		},
	}

	logicalNetworkInterconnectGetCmd = &cobra.Command{
		Use:     "get interconnect_id_or_label",
		Aliases: []string{"show"},
		Short:   "Get logical network interconnect details",
		Long: `Get the details of a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Examples:
  metalcloud logical-network-interconnect get 4
  metalcloud lni get dc1-dc2-logical-interconnect`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkInterconnectConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example logical network interconnect create configuration",
		Long: `Print an example configuration that can be edited and passed to
'logical-network-interconnect create --config-source'.

Examples:
  metalcloud logical-network-interconnect config-example > lni.json
  metalcloud lni config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectConfigExample(cmd.Context())
		},
	}

	logicalNetworkInterconnectCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a logical network interconnect",
		Long: `Create a new logical network interconnect on top of an existing network fabric
interconnect.

Required Flags (one of):
  --config-source string        Source of the configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string                Label of the new logical network interconnect (used together with the flags below)

Required Flags (when not using --config-source):
  --name string                 Name of the new logical network interconnect
  --fabric-interconnect-id int  ID of the network fabric interconnect it runs on

Optional Flags (when not using --config-source):
  --kind string                 Interconnect kind (default "dci-evpn")

Examples:
  metalcloud logical-network-interconnect create --config-source lni.json
  metalcloud lni create --label dc1-dc2-ln --name "DC1 to DC2 network" --fabric-interconnect-id 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateLogicalNetworkInterconnect

			if logicalNetworkInterconnectFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkInterconnectFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				kind, err := sdk.NewLogicalNetworkInterconnectKindFromValue(logicalNetworkInterconnectFlags.kind)
				if err != nil {
					return err
				}

				create = sdk.CreateLogicalNetworkInterconnect{
					Label:                logicalNetworkInterconnectFlags.label,
					Name:                 logicalNetworkInterconnectFlags.name,
					Kind:                 kind,
					FabricInterconnectId: logicalNetworkInterconnectFlags.fabricInterconnectId,
				}
			}

			return logical_network_interconnect.InterconnectCreate(cmd.Context(), create)
		},
	}

	logicalNetworkInterconnectUpdateCmd = &cobra.Command{
		Use:     "update interconnect_id_or_label",
		Aliases: []string{"edit"},
		Short:   "Update a logical network interconnect",
		Long: `Update the label, name, kind or annotations of a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Required Flags:
  --config-source string     Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud logical-network-interconnect update 4 --config-source update.json
  echo '{"name":"new name"}' | metalcloud lni update dc1-dc2-ln --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkInterconnectFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network_interconnect.InterconnectUpdate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkInterconnectDeleteCmd = &cobra.Command{
		Use:     "delete interconnect_id_or_label",
		Aliases: []string{"rm"},
		Short:   "Delete a logical network interconnect",
		Long: `Delete a logical network interconnect together with its links.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Examples:
  metalcloud logical-network-interconnect delete 4`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectDelete(cmd.Context(), args[0])
		},
	}

	logicalNetworkInterconnectLinkCmd = &cobra.Command{
		Use:     "link [command]",
		Aliases: []string{"links"},
		Short:   "Logical network interconnect link management",
		Long: `Manage the logical networks linked to a logical network interconnect.

Commands: list, get, add, remove`,
	}

	logicalNetworkInterconnectLinkListCmd = &cobra.Command{
		Use:     "list interconnect_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the links of a logical network interconnect",
		Long: `List the logical networks linked to a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Optional Flags:
  --filter-logical-network-id strings   Filter by logical network ID. Repeatable or comma-separated.
  --filter-status strings               Filter by status. Repeatable or comma-separated.

Examples:
  metalcloud logical-network-interconnect link list 4
  metalcloud lni link list dc1-dc2-ln --filter-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectLinkList(cmd.Context(), args[0],
				logicalNetworkInterconnectFlags.filterLogicalNetworkId,
				logicalNetworkInterconnectFlags.filterStatus)
		},
	}

	logicalNetworkInterconnectLinkGetCmd = &cobra.Command{
		Use:     "get interconnect_id_or_label link_id",
		Aliases: []string{"show"},
		Short:   "Get one link of a logical network interconnect",
		Long: `Get the details of a single logical network interconnect link.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  link_id                    The ID of the link

Examples:
  metalcloud logical-network-interconnect link get 4 9`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectLinkGet(cmd.Context(), args[0], args[1])
		},
	}

	logicalNetworkInterconnectLinkAddCmd = &cobra.Command{
		Use:     "add interconnect_id_or_label logical_network_id",
		Aliases: []string{"new", "create"},
		Short:   "Link a logical network to a logical network interconnect",
		Long: `Link a logical network to a logical network interconnect so that it is stretched
across the interconnected fabrics.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  logical_network_id         The numeric ID of the logical network

Examples:
  metalcloud logical-network-interconnect link add 4 44`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectLinkAdd(cmd.Context(), args[0], args[1])
		},
	}

	logicalNetworkInterconnectLinkRemoveCmd = &cobra.Command{
		Use:     "remove interconnect_id_or_label link_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a link from a logical network interconnect",
		Long: `Remove a logical network link from a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect
  link_id                    The ID of the link to remove

Examples:
  metalcloud logical-network-interconnect link remove 4 9`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_LOGICAL_NETWORK_INTERCONNECTS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_interconnect.InterconnectLinkRemove(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(logicalNetworkInterconnectCmd)

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectListCmd)
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterId, "filter-id", nil, "Filter by logical network interconnect ID.")
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterLabel, "filter-label", nil, "Filter by label.")
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterName, "filter-name", nil, "Filter by name.")
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterKind, "filter-kind", nil, "Filter by kind (e.g. dci-evpn).")
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterStatus, "filter-status", nil, "Filter by status.")
	logicalNetworkInterconnectListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterFabricInterconnectId, "filter-fabric-interconnect-id", nil, "Filter by network fabric interconnect ID.")

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectGetCmd)
	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectConfigExampleCmd)

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectCreateCmd)
	logicalNetworkInterconnectCreateCmd.Flags().StringVar(&logicalNetworkInterconnectFlags.configSource, "config-source", "", "Source of the new logical network interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.")
	logicalNetworkInterconnectCreateCmd.Flags().StringVar(&logicalNetworkInterconnectFlags.label, "label", "", "Label of the new logical network interconnect.")
	logicalNetworkInterconnectCreateCmd.Flags().StringVar(&logicalNetworkInterconnectFlags.name, "name", "", "Name of the new logical network interconnect.")
	logicalNetworkInterconnectCreateCmd.Flags().StringVar(&logicalNetworkInterconnectFlags.kind, "kind", string(sdk.LOGICALNETWORKINTERCONNECTKIND_DCI_EVPN), "Kind of the logical network interconnect.")
	logicalNetworkInterconnectCreateCmd.Flags().Int64Var(&logicalNetworkInterconnectFlags.fabricInterconnectId, "fabric-interconnect-id", 0, "ID of the network fabric interconnect it runs on.")
	logicalNetworkInterconnectCreateCmd.MarkFlagsOneRequired("config-source", "label")
	logicalNetworkInterconnectCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	logicalNetworkInterconnectCreateCmd.MarkFlagsRequiredTogether("label", "name", "fabric-interconnect-id")

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectUpdateCmd)
	logicalNetworkInterconnectUpdateCmd.Flags().StringVar(&logicalNetworkInterconnectFlags.configSource, "config-source", "", "Source of the updated logical network interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.")
	logicalNetworkInterconnectUpdateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectDeleteCmd)

	logicalNetworkInterconnectCmd.AddCommand(logicalNetworkInterconnectLinkCmd)
	logicalNetworkInterconnectLinkCmd.AddCommand(logicalNetworkInterconnectLinkListCmd)
	logicalNetworkInterconnectLinkListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterLogicalNetworkId, "filter-logical-network-id", nil, "Filter by logical network ID.")
	logicalNetworkInterconnectLinkListCmd.Flags().StringSliceVar(&logicalNetworkInterconnectFlags.filterStatus, "filter-status", nil, "Filter by status.")

	logicalNetworkInterconnectLinkCmd.AddCommand(logicalNetworkInterconnectLinkGetCmd)
	logicalNetworkInterconnectLinkCmd.AddCommand(logicalNetworkInterconnectLinkAddCmd)
	logicalNetworkInterconnectLinkCmd.AddCommand(logicalNetworkInterconnectLinkRemoveCmd)
}
