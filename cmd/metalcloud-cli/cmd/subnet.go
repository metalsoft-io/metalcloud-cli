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
		Short:   "Subnet management",
		Long:    `Subnet management commands.`,
	}

	subnetListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all subnets.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetList(cmd.Context())
		},
	}

	subnetGetCmd = &cobra.Command{
		Use:          "get subnet_id",
		Aliases:      []string{"show"},
		Short:        "Get subnet details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetGet(cmd.Context(), args[0])
		},
	}

	subnetConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get subnet configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SUBNETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return subnet.SubnetConfigExample(cmd.Context())
		},
	}

	subnetCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a subnet.",
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
		Use:          "update subnet_id",
		Aliases:      []string{"modify"},
		Short:        "Update a subnet.",
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
		Use:          "delete subnet_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a subnet.",
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
