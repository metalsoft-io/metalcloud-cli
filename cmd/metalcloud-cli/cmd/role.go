package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/role"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	roleFlags = struct {
		configSource string
	}{}

	roleCmd = &cobra.Command{
		Use:     "role [command]",
		Aliases: []string{"roles"},
		Short:   "Role management",
		Long:    `Role management commands.`,
	}

	roleListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all roles.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.List(cmd.Context())
		},
	}

	roleGetCmd = &cobra.Command{
		Use:          "get role_name",
		Aliases:      []string{"show"},
		Short:        "Get role details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.Get(cmd.Context(), args[0])
		},
	}

	roleCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new role.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(roleFlags.configSource)
			if err != nil {
				return err
			}
			return role.Create(cmd.Context(), config)
		},
	}

	roleDeleteCmd = &cobra.Command{
		Use:          "delete role_name",
		Aliases:      []string{"remove"},
		Short:        "Delete a role.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.Delete(cmd.Context(), args[0])
		},
	}

	roleUpdateCmd = &cobra.Command{
		Use:          "update role_name",
		Aliases:      []string{"edit"},
		Short:        "Update a role.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(roleFlags.configSource)
			if err != nil {
				return err
			}
			return role.Update(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(roleCmd)

	// Role list
	roleCmd.AddCommand(roleListCmd)

	// Role get
	roleCmd.AddCommand(roleGetCmd)

	// Role create
	roleCmd.AddCommand(roleCreateCmd)
	roleCreateCmd.Flags().StringVar(&roleFlags.configSource, "config-source", "", "Source of the role configuration. Can be 'pipe' or path to a JSON file.")
	roleCreateCmd.MarkFlagRequired("config-source")

	// Role delete
	roleCmd.AddCommand(roleDeleteCmd)

	// Role update
	roleCmd.AddCommand(roleUpdateCmd)
	roleUpdateCmd.Flags().StringVar(&roleFlags.configSource, "config-source", "", "Source of the role configuration. Can be 'pipe' or path to a JSON file.")
	roleUpdateCmd.MarkFlagRequired("config-source")
}
