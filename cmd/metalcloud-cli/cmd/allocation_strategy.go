package cmd

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/allocation_strategy"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// allocationStrategyFlags is shared by every registered allocation-strategy
// group; only one command runs per invocation.
var allocationStrategyFlags = struct {
	configSource string
}{}

// registerAllocationStrategyCommands attaches an "allocation-strategy" group
// with list/get/config-example/add/replace/remove sub-commands to parentCmd.
// The same generic implementation serves logical networks, logical network
// profiles, route domains and point-to-point links, so the group is built
// once here instead of being repeated in each resource's command file.
func registerAllocationStrategyCommands(parentCmd *cobra.Command, parent allocation_strategy.Parent, parentArg string, readPermission string, writePermission string) {
	families := strings.Join(parent.FamilyNames(), ", ")
	cli := parentCmd.Name()

	groupCmd := &cobra.Command{
		Use:     "allocation-strategy [command]",
		Aliases: []string{"strategy", "strategies"},
		Short:   fmt.Sprintf("Manage %s allocation strategies", parent.Name),
		Long: fmt.Sprintf(`Manage the allocation strategies of a %s.

Each strategy family is a separate collection. Supported families for this resource:
  %s

Commands:
  list, get, config-example, add, replace, remove`, parent.Name, families),
	}

	listCmd := &cobra.Command{
		Use:     fmt.Sprintf("list %s family", parentArg),
		Aliases: []string{"ls"},
		Short:   "List allocation strategies of one family",
		Long: fmt.Sprintf(`List the allocation strategies of one family on a %s.

Required Arguments:
  %-12s The ID of the %s
  family       One of: %s

Examples:
  metalcloud %s allocation-strategy list 12 %s`, parent.Name, parentArg, parent.Name, families, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: readPermission},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return allocation_strategy.List(cmd.Context(), parent, args[0], args[1])
		},
	}

	getCmd := &cobra.Command{
		Use:     fmt.Sprintf("get %s family strategy_id", parentArg),
		Aliases: []string{"show"},
		Short:   "Get one allocation strategy",
		Long: fmt.Sprintf(`Get the details of one allocation strategy of a %s.

Required Arguments:
  %-12s The ID of the %s
  family       One of: %s
  strategy_id  The ID of the allocation strategy

Examples:
  metalcloud %s allocation-strategy get 12 %s 3`, parent.Name, parentArg, parent.Name, families, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: readPermission},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return allocation_strategy.Get(cmd.Context(), parent, args[0], args[1], args[2])
		},
	}

	configExampleCmd := &cobra.Command{
		Use:   "config-example family",
		Short: "Print example allocation strategy configurations",
		Long: fmt.Sprintf(`Print one example create body per strategy kind (auto, manual, ...) for a family.
Edit one of them and pass it to 'add' or 'replace' via --config-source.

Required Arguments:
  family       One of: %s

Examples:
  metalcloud %s allocation-strategy config-example %s -f yaml`, families, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: readPermission},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return allocation_strategy.ConfigExample(parent, args[0])
		},
	}

	addCmd := &cobra.Command{
		Use:     fmt.Sprintf("add %s family", parentArg),
		Aliases: []string{"create", "new"},
		Short:   "Add an allocation strategy",
		Long: fmt.Sprintf(`Add an allocation strategy to a %s.

Required Arguments:
  %-12s The ID of the %s
  family       One of: %s

Required Flags:
  --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud %s allocation-strategy add 12 %s --config-source strategy.json
  echo '{"kind":"auto","scope":{"kind":"global"}}' | metalcloud %s allocation-strategy add 12 %s --config-source pipe`,
			parent.Name, parentArg, parent.Name, families, cli, parent.Families[0].Name, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: writePermission},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(allocationStrategyFlags.configSource)
			if err != nil {
				return err
			}
			return allocation_strategy.Create(cmd.Context(), parent, args[0], args[1], config)
		},
	}

	replaceCmd := &cobra.Command{
		Use:     fmt.Sprintf("replace %s family strategy_id", parentArg),
		Aliases: []string{"update"},
		Short:   "Replace an allocation strategy",
		Long: fmt.Sprintf(`Replace the whole configuration of one allocation strategy of a %s.

Required Arguments:
  %-12s The ID of the %s
  family       One of: %s
  strategy_id  The ID of the allocation strategy

Required Flags:
  --config-source string   Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud %s allocation-strategy replace 12 %s 3 --config-source strategy.json`,
			parent.Name, parentArg, parent.Name, families, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: writePermission},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(allocationStrategyFlags.configSource)
			if err != nil {
				return err
			}
			return allocation_strategy.Replace(cmd.Context(), parent, args[0], args[1], args[2], config)
		},
	}

	removeCmd := &cobra.Command{
		Use:     fmt.Sprintf("remove %s family strategy_id", parentArg),
		Aliases: []string{"delete", "rm"},
		Short:   "Remove an allocation strategy",
		Long: fmt.Sprintf(`Remove one allocation strategy from a %s.

Required Arguments:
  %-12s The ID of the %s
  family       One of: %s
  strategy_id  The ID of the allocation strategy

Examples:
  metalcloud %s allocation-strategy remove 12 %s 3`, parent.Name, parentArg, parent.Name, families, cli, parent.Families[0].Name),
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: writePermission},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return allocation_strategy.Delete(cmd.Context(), parent, args[0], args[1], args[2])
		},
	}

	parentCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(listCmd, getCmd, configExampleCmd)

	groupCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&allocationStrategyFlags.configSource, "config-source", "", "Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.")
	addCmd.MarkFlagsOneRequired("config-source")

	groupCmd.AddCommand(replaceCmd)
	replaceCmd.Flags().StringVar(&allocationStrategyFlags.configSource, "config-source", "", "Source of the strategy configuration. Can be 'pipe' or path to a JSON/YAML file.")
	replaceCmd.MarkFlagsOneRequired("config-source")

	groupCmd.AddCommand(removeCmd)
}

func init() {
	registerAllocationStrategyCommands(logicalNetworkCmd, allocation_strategy.LogicalNetwork, "logical_network_id",
		system.PERMISSION_NETWORK_PROFILES_READ, system.PERMISSION_NETWORK_PROFILES_WRITE)
	registerAllocationStrategyCommands(logicalNetworkProfileCmd, allocation_strategy.LogicalNetworkProfile, "logical_network_profile_id",
		system.PERMISSION_NETWORK_PROFILES_READ, system.PERMISSION_NETWORK_PROFILES_WRITE)
	registerAllocationStrategyCommands(routeDomainCmd, allocation_strategy.RouteDomain, "route_domain_id",
		system.PERMISSION_NETWORK_PROFILES_READ, system.PERMISSION_NETWORK_PROFILES_WRITE)
	registerAllocationStrategyCommands(pointToPointLinkCmd, allocation_strategy.PointToPointLink, "link_id",
		system.PERMISSION_NETWORK_FABRICS_READ, system.PERMISSION_NETWORK_FABRICS_WRITE)
}
