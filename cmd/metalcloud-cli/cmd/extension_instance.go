package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/extension_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	extensionInstanceFlags = struct {
		configSource string
	}{}

	extensionInstanceCmd = &cobra.Command{
		Use:     "extension-instance [command]",
		Aliases: []string{"ext-inst"},
		Short:   "Extension instance management",
		Long:    `Extension instance management commands.`,
	}

	extensionInstanceListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all extension instances in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceList(cmd.Context(), args[0])
		},
	}

	extensionInstanceGetCmd = &cobra.Command{
		Use:          "get extension_instance_id",
		Aliases:      []string{"show"},
		Short:        "Get extension instance details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceGet(cmd.Context(), args[0])
		},
	}

	extensionInstanceCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label",
		Aliases:      []string{"new"},
		Short:        "Create new extension instance in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(extensionInstanceFlags.configSource)
			if err != nil {
				return err
			}
			return extension_instance.ExtensionInstanceCreate(cmd.Context(), args[0], config)
		},
	}

	extensionInstanceUpdateCmd = &cobra.Command{
		Use:          "update extension_instance_id",
		Aliases:      []string{"edit"},
		Short:        "Update extension instance configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(extensionInstanceFlags.configSource)
			if err != nil {
				return err
			}
			return extension_instance.ExtensionInstanceUpdate(cmd.Context(), args[0], config)
		},
	}

	extensionInstanceDeleteCmd = &cobra.Command{
		Use:          "delete extension_instance_id",
		Aliases:      []string{"rm"},
		Short:        "Delete extension instance.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSION_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension_instance.ExtensionInstanceDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(extensionInstanceCmd)

	extensionInstanceCmd.AddCommand(extensionInstanceListCmd)
	extensionInstanceCmd.AddCommand(extensionInstanceGetCmd)

	extensionInstanceCmd.AddCommand(extensionInstanceCreateCmd)
	extensionInstanceCreateCmd.Flags().StringVar(&extensionInstanceFlags.configSource, "config-source", "", "Source of the new extension instance configuration. Can be 'pipe' or path to a JSON file.")
	extensionInstanceCreateCmd.MarkFlagsOneRequired("config-source")

	extensionInstanceCmd.AddCommand(extensionInstanceUpdateCmd)
	extensionInstanceUpdateCmd.Flags().StringVar(&extensionInstanceFlags.configSource, "config-source", "", "Source of the extension instance configuration updates. Can be 'pipe' or path to a JSON file.")
	extensionInstanceUpdateCmd.MarkFlagsOneRequired("config-source")

	extensionInstanceCmd.AddCommand(extensionInstanceDeleteCmd)
}
