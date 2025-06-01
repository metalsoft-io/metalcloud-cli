package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/custom_iso"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	customIsoFlags = struct {
		configSource string
	}{}

	customIsoCmd = &cobra.Command{
		Use:     "custom-iso [command]",
		Aliases: []string{"iso", "isos"},
		Short:   "Custom ISO management",
		Long:    `Custom ISO management commands.`,
	}

	customIsoListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all custom ISOs.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoList(cmd.Context())
		},
	}

	customIsoGetCmd = &cobra.Command{
		Use:          "get custom_iso_id",
		Aliases:      []string{"show"},
		Short:        "Get custom ISO details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoGet(cmd.Context(), args[0])
		},
	}

	customIsoConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get custom ISO configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoConfigExample(cmd.Context())
		},
	}

	customIsoCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new custom ISO.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(customIsoFlags.configSource)
			if err != nil {
				return err
			}

			return custom_iso.CustomIsoCreate(cmd.Context(), config)
		},
	}

	customIsoUpdateCmd = &cobra.Command{
		Use:          "update custom_iso_id",
		Aliases:      []string{"edit"},
		Short:        "Update a custom ISO.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(customIsoFlags.configSource)
			if err != nil {
				return err
			}

			return custom_iso.CustomIsoUpdate(cmd.Context(), args[0], config)
		},
	}

	customIsoDeleteCmd = &cobra.Command{
		Use:          "delete custom_iso_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a custom ISO.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoDelete(cmd.Context(), args[0])
		},
	}

	customIsoMakePublicCmd = &cobra.Command{
		Use:          "make-public custom_iso_id",
		Short:        "Make a custom ISO public.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CUSTOM_ISO_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoMakePublic(cmd.Context(), args[0])
		},
	}

	customIsoBootServerCmd = &cobra.Command{
		Use:          "boot-server custom_iso_id server_id",
		Aliases:      []string{"boot"},
		Short:        "Boot a server using a custom ISO.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return custom_iso.CustomIsoBootServer(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(customIsoCmd)

	customIsoCmd.AddCommand(customIsoListCmd)
	customIsoCmd.AddCommand(customIsoGetCmd)
	customIsoCmd.AddCommand(customIsoConfigExampleCmd)

	customIsoCmd.AddCommand(customIsoCreateCmd)
	customIsoCreateCmd.Flags().StringVar(&customIsoFlags.configSource, "config-source", "", "Source of the new custom ISO configuration. Can be 'pipe' or path to a JSON file.")
	customIsoCreateCmd.MarkFlagsOneRequired("config-source")

	customIsoCmd.AddCommand(customIsoUpdateCmd)
	customIsoUpdateCmd.Flags().StringVar(&customIsoFlags.configSource, "config-source", "", "Source of the custom ISO configuration updates. Can be 'pipe' or path to a JSON file.")
	customIsoUpdateCmd.MarkFlagsOneRequired("config-source")

	customIsoCmd.AddCommand(customIsoDeleteCmd)
	customIsoCmd.AddCommand(customIsoMakePublicCmd)
	customIsoCmd.AddCommand(customIsoBootServerCmd)
}
