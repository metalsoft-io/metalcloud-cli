package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/point_to_point_link"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	pointToPointLinkFlags = struct {
		configSource    string
		interfaceId     string
		routeDomainId   string
		subnetId        int64
		binding         string
		scopeKind       string
		scopeResourceId int64
	}{}

	pointToPointLinkCmd = &cobra.Command{
		Use:     "point-to-point-link [command]",
		Aliases: []string{"p2p-link", "p2p"},
		Short:   "Manage point-to-point links between network interfaces",
		Long: `Manage point-to-point links between network device (and server) interfaces.

A point-to-point link connects two interfaces (or a single interface, for a
half-connected link) and can carry IPv4/IPv6 subnet allocation strategies that
assign the link's addresses. Links can be created fully staged (interfaces plus
a manual /31 strategy) in one call via the create command's config source.`,
	}

	pointToPointLinkListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List point-to-point links",
		Long: `List point-to-point links, optionally filtered by a referenced interface id
or route domain id.

Examples:
  metalcloud-cli point-to-point-link list
  metalcloud-cli p2p ls --interface-id 1001
  metalcloud-cli p2p ls --route-domain-id 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return point_to_point_link.PointToPointLinkList(cmd.Context(), pointToPointLinkFlags.interfaceId, pointToPointLinkFlags.routeDomainId)
		},
	}

	pointToPointLinkGetCmd = &cobra.Command{
		Use:          "get link_id",
		Aliases:      []string{"show"},
		Short:        "Get details about a specific point-to-point link",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return point_to_point_link.PointToPointLinkGet(cmd.Context(), args[0])
		},
	}

	pointToPointLinkConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Display a point-to-point link configuration example",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return point_to_point_link.PointToPointLinkConfigExample(cmd.Context())
		},
	}

	pointToPointLinkCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a point-to-point link",
		Long: `Create a point-to-point link from a JSON/YAML configuration.

The configuration may include:
- interfaceA / interfaceB: { type: "network_equipment_interface", interfaceId: <id> }
  (omit interfaceB for a half-connected link)
- description, mtu, routingActivation ("default" or "while_transporting_logical_network")
- ipv4.subnetAllocationStrategies: one or more strategies staged on create
  (e.g. a manual strategy with a subnetId and interfaceABinding)

Required Flags:
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli point-to-point-link config-example > link.json
  metalcloud-cli point-to-point-link create --config-source link.json
  cat link.json | metalcloud-cli p2p create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(pointToPointLinkFlags.configSource)
			if err != nil {
				return err
			}

			return point_to_point_link.PointToPointLinkCreate(cmd.Context(), config)
		},
	}

	pointToPointLinkAddIpv4StrategyCmd = &cobra.Command{
		Use:   "add-ipv4-strategy link_id",
		Short: "Attach a manual IPv4 subnet allocation strategy to a link",
		Long: `Attach a manual IPv4 subnet allocation strategy to an existing point-to-point
link. This is the repair path for a link that was created without its strategy;
new links should stage the strategy on create instead.

Arguments:
  link_id            The ID of the point-to-point link

Required Flags:
  --subnet-id        ID of the IPAM subnet to allocate from
  --binding          Which interface gets the first address: a_first, b_first, or auto

Optional Flags:
  --scope-kind       Allocation scope kind (default: global)
  --scope-resource-id   Resource id for non-global scopes (default: 0)

Examples:
  metalcloud-cli p2p add-ipv4-strategy 42 --subnet-id 12345 --binding a_first
  metalcloud-cli p2p add-ipv4-strategy 42 --subnet-id 12345 --binding b_first`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return point_to_point_link.PointToPointLinkAddIpv4Strategy(
				cmd.Context(),
				args[0],
				pointToPointLinkFlags.subnetId,
				pointToPointLinkFlags.binding,
				pointToPointLinkFlags.scopeKind,
				pointToPointLinkFlags.scopeResourceId,
			)
		},
	}

	pointToPointLinkDeleteCmd = &cobra.Command{
		Use:          "delete link_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a point-to-point link",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return point_to_point_link.PointToPointLinkDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(pointToPointLinkCmd)

	pointToPointLinkCmd.AddCommand(pointToPointLinkListCmd)
	pointToPointLinkListCmd.Flags().StringVar(&pointToPointLinkFlags.interfaceId, "interface-id", "", "Filter links by a referenced interface id.")
	pointToPointLinkListCmd.Flags().StringVar(&pointToPointLinkFlags.routeDomainId, "route-domain-id", "", "Filter links by route domain id.")

	pointToPointLinkCmd.AddCommand(pointToPointLinkGetCmd)

	pointToPointLinkCmd.AddCommand(pointToPointLinkConfigExampleCmd)

	pointToPointLinkCmd.AddCommand(pointToPointLinkCreateCmd)
	pointToPointLinkCreateCmd.Flags().StringVar(&pointToPointLinkFlags.configSource, "config-source", "", "Source of the new link configuration. Can be 'pipe' or path to a JSON/YAML file.")
	pointToPointLinkCreateCmd.MarkFlagRequired("config-source")

	pointToPointLinkCmd.AddCommand(pointToPointLinkAddIpv4StrategyCmd)
	pointToPointLinkAddIpv4StrategyCmd.Flags().Int64Var(&pointToPointLinkFlags.subnetId, "subnet-id", 0, "ID of the IPAM subnet to allocate from.")
	pointToPointLinkAddIpv4StrategyCmd.Flags().StringVar(&pointToPointLinkFlags.binding, "binding", "a_first", "Interface A binding: a_first, b_first, or auto.")
	pointToPointLinkAddIpv4StrategyCmd.Flags().StringVar(&pointToPointLinkFlags.scopeKind, "scope-kind", "global", "Allocation scope kind.")
	pointToPointLinkAddIpv4StrategyCmd.Flags().Int64Var(&pointToPointLinkFlags.scopeResourceId, "scope-resource-id", 0, "Resource id for non-global scopes.")
	pointToPointLinkAddIpv4StrategyCmd.MarkFlagRequired("subnet-id")

	pointToPointLinkCmd.AddCommand(pointToPointLinkDeleteCmd)
}
