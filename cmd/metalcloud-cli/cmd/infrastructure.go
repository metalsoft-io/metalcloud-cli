package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/spf13/cobra"
)

var (
	showOwnOnly     bool
	showOrdered     bool
	showDeleted     bool
	customVariables string

	infrastructureCmd = &cobra.Command{
		Use:     "infrastructure [command]",
		Aliases: []string{"infra"},
		Short:   "Infrastructure management",
		Long:    `Infrastructure management commands.`,
	}

	infrastructureListCmd = &cobra.Command{
		Use:          "list [flags...]",
		Aliases:      []string{"ls"},
		Short:        "List all infrastructures.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureList(cmd.Context(), showOwnOnly, showOrdered, showDeleted)
		},
	}

	infrastructureGetCmd = &cobra.Command{
		Use:          "get infrastructure_id_or_label",
		Aliases:      []string{"show"},
		Short:        "Get infrastructure details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGet(cmd.Context(), args[0])
		},
	}

	infrastructureCreateCmd = &cobra.Command{
		Use:          "create site_id label",
		Aliases:      []string{"new"},
		Short:        "Create new infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_WRITE}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureCreate(cmd.Context(), args[0], args[1])
		},
	}

	infrastructureUpdateCmd = &cobra.Command{
		Use:          "update infrastructure_id_or_label [new_label]",
		Aliases:      []string{"edit"},
		Short:        "Update infrastructure configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_WRITE}, // TODO: Use specific permission
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label := ""
			if len(args) > 1 {
				label = args[1]
			}

			return infrastructure.InfrastructureUpdate(cmd.Context(), args[0], label, customVariables)
		},
	}
)

func init() {
	rootCmd.AddCommand(infrastructureCmd)

	infrastructureCmd.AddCommand(infrastructureListCmd)
	infrastructureListCmd.Flags().BoolVar(&showOwnOnly, "show-own-only", true, "If set will return only infrastructures owned by the current user.")
	infrastructureListCmd.Flags().BoolVar(&showOrdered, "show-ordered", false, "If set will also return ordered (created but not deployed) infrastructures.")
	infrastructureListCmd.Flags().BoolVar(&showDeleted, "show-deleted", false, "If set will also return deleted infrastructures.")

	infrastructureCmd.AddCommand(infrastructureGetCmd)

	infrastructureCmd.AddCommand(infrastructureCreateCmd)

	infrastructureCmd.AddCommand(infrastructureUpdateCmd)
	infrastructureListCmd.Flags().StringVar(&customVariables, "custom-variables", "", "Set of infrastructure custom variables.")
}
