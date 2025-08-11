package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/subnet"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	subnetFlags = struct {
		configSource string
	}{}

	subnetCmd = &cobra.Command{
		Use:     "subnet [command]",
		Aliases: []string{"subnets", "net"},
		Short:   "Manage network subnets and IP address pools",
		Long: `Manage network subnets and IP address pools in the MetalCloud infrastructure.

Subnets define network segments with specific IP address ranges and can be configured as:
- Regular subnets: Fixed network segments with defined address ranges
- IP pools: Dynamic address pools for automatic IP allocation

Available commands allow you to list, create, update, delete subnets and view configuration examples.`,
	}

	subnetListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all subnets and IP pools",
		Long: `List all subnets and IP address pools in the MetalCloud infrastructure.

This command displays a tabular view of all subnets with key information including:
- Subnet ID and name
- IP version (IPv4/IPv6)  
- Network address and prefix length
- Netmask
- Pool status (whether it's configured as an IP pool)
- Creation timestamp

Examples:
  metalcloud-cli subnet list
  metalcloud-cli subnets ls
  metalcloud-cli net list`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetList(cmd.Context())
		},
	}

	subnetGetCmd = &cobra.Command{
		Use:     "get subnet_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific subnet",
		Long: `Get detailed information about a specific subnet including all its configuration properties.

The command shows comprehensive subnet details including:
- Basic subnet information (ID, name, label)
- Network configuration (address, prefix, netmask, gateway)
- IP pool configuration status
- Allocation denylist and rules
- Creation and modification timestamps
- Associated tags and annotations

Arguments:
  subnet_id    The ID of the subnet to retrieve

Examples:
  metalcloud-cli subnet get 123
  metalcloud-cli subnets show 456
  metalcloud-cli net get 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetGet(cmd.Context(), args[0])
		},
	}

	subnetConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display subnet configuration example",
		Long: `Display a complete example of subnet configuration in JSON format.

This command shows a sample configuration that can be used as a template for creating subnets.
The example includes all available fields with descriptions of their purpose and common values.

The output can be saved to a file and modified for your specific requirements:

Examples:
  metalcloud-cli subnet config-example
  metalcloud-cli subnet config-example > subnet-config.json
  metalcloud-cli subnets config-example | jq . > my-subnet.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetConfigExample(cmd.Context())
		},
	}

	subnetCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new subnet",
		Long: `Create a new subnet in the MetalCloud infrastructure.

This command creates a subnet with the configuration provided through JSON input.
The configuration must include required fields like network address, prefix length,
and whether the subnet should be configured as an IP pool.

Required Flags:
  --config-source    Source of the new subnet configuration. Can be 'pipe' to read
                     from stdin, or a path to a JSON file containing the configuration.

Configuration Format:
The JSON configuration should contain the following fields:
- networkAddress (required): IP network address (e.g., "192.168.1.0")
- prefixLength (required): Network prefix length (e.g., 24)
- isPool (required): Whether this subnet is an IP pool (true/false)
- label (optional): Human-readable label for the subnet
- name (optional): Subnet name
- defaultGatewayAddress (optional): Gateway IP address
- parentSubnetId (optional): ID of parent subnet
- allocationDenylist (optional): List of IP ranges to exclude from allocation
- childOverlapAllowRules (optional): Rules for allowing child subnet overlaps
- tags (optional): Key-value pairs for tagging
- annotations (optional): Additional metadata

Examples:
  # Create subnet from stdin
  echo '{"networkAddress":"10.0.1.0","prefixLength":24,"isPool":false}' | metalcloud-cli subnet create --config-source pipe
  
  # Create subnet from file
  metalcloud-cli subnet create --config-source subnet-config.json
  
  # Show example configuration first
  metalcloud-cli subnet config-example > example.json
  # Edit example.json and then create
  metalcloud-cli subnet create --config-source example.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(subnetFlags.configSource)
			if err != nil {
				return err
			}

			return subnet.SubnetCreate(cmd.Context(), config)
		},
	}

	subnetUpdateCmd = &cobra.Command{
		Use:     "update subnet_id",
		Aliases: []string{"modify"},
		Short:   "Update an existing subnet",
		Long: `Update an existing subnet in the MetalCloud infrastructure.

This command updates a subnet with the configuration provided through JSON input.
Only the fields included in the configuration will be updated, other fields remain unchanged.

Arguments:
  subnet_id    The ID of the subnet to update

Required Flags:
  --config-source    Source of the subnet configuration updates. Can be 'pipe' to read
                     from stdin, or a path to a JSON file containing the configuration.

Configuration Format:
The JSON configuration can contain the following fields (all optional for updates):
- label: Human-readable label for the subnet
- name: Subnet name
- defaultGatewayAddress: Gateway IP address
- isPool: Whether this subnet is an IP pool (true/false)
- allocationDenylist: List of IP ranges to exclude from allocation
- childOverlapAllowRules: Rules for allowing child subnet overlaps
- tags: Key-value pairs for tagging
- annotations: Additional metadata

Note: Core network settings (networkAddress, prefixLength) typically cannot be modified
after subnet creation due to infrastructure constraints.

Examples:
  # Update subnet from stdin
  echo '{"label":"updated-subnet","isPool":true}' | metalcloud-cli subnet update 123 --config-source pipe
  
  # Update subnet from file
  metalcloud-cli subnet update 456 --config-source updates.json
  
  # Update only tags
  echo '{"tags":{"environment":"production","team":"networking"}}' | metalcloud-cli subnet update 789 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(subnetFlags.configSource)
			if err != nil {
				return err
			}

			return subnet.SubnetUpdate(cmd.Context(), args[0], config)
		},
	}

	subnetDeleteCmd = &cobra.Command{
		Use:     "delete subnet_id",
		Aliases: []string{"rm"},
		Short:   "Delete a subnet",
		Long: `Delete a subnet from the MetalCloud infrastructure.

This command permanently removes a subnet and all its associated configuration.
The subnet must not be in use by any resources before it can be deleted.

Arguments:
  subnet_id    The ID of the subnet to delete

Warning: This operation is irreversible. Ensure the subnet is not being used by any
infrastructure components before deletion.

Examples:
  metalcloud-cli subnet delete 123
  metalcloud-cli subnets rm 456
  metalcloud-cli net delete 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(subnetCmd)

	subnetCmd.AddCommand(subnetListCmd)
	subnetCmd.AddCommand(subnetGetCmd)
	subnetCmd.AddCommand(subnetConfigExampleCmd)

	subnetCmd.AddCommand(subnetCreateCmd)
	subnetCreateCmd.Flags().StringVar(&subnetFlags.configSource, "config-source", "", "Source of the new subnet configuration. Can be 'pipe' or path to a JSON file.")
	subnetCreateCmd.MarkFlagsOneRequired("config-source")

	subnetCmd.AddCommand(subnetUpdateCmd)
	subnetUpdateCmd.Flags().StringVar(&subnetFlags.configSource, "config-source", "", "Source of the subnet configuration updates. Can be 'pipe' or path to a JSON file.")
	subnetUpdateCmd.MarkFlagsOneRequired("config-source")

	subnetCmd.AddCommand(subnetDeleteCmd)
}
