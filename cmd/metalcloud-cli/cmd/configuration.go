package cmd

import (
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	configurationFlags = struct {
		configSource string
	}{}

	// configurationServices is the list of service sections accepted as the
	// filter argument, rendered for the help texts.
	configurationServices = strings.Join(configuration.KnownServices, ", ")

	configurationCmd = &cobra.Command{
		// "config" is not used as an alias: the root command already owns the
		// persistent --config/-c flag for the CLI's own config file, and a
		// command by that name reads as if it managed that file.
		Use:     "configuration [command]",
		Aliases: []string{"system-config", "platform-config"},
		Short:   "Global platform configuration management",
		Long: `Read and change the GLOBAL configuration of the MetalSoft platform.

This is not the configuration of the CLI (see the --config flag for that): these
commands read and write the platform-wide settings served by /api/v2/config, which
affect every user and every infrastructure of the installation.

The configuration is split into service sections, addressed by their name:
  ` + configurationServices + `

Command categories:
  Read:    get
  Write:   replace (full section), update (partial section)`,
	}

	configurationGetCmd = &cobra.Command{
		Use:     "get [service]",
		Aliases: []string{"show"},
		Short:   "Get the global platform configuration",
		Long: `Show the global platform configuration, optionally restricted to one service section.

Optional Arguments:
  service   One of: ` + configurationServices + `. When omitted, every section is returned.

The document is deeply nested, so text output is rendered as YAML. Use -f json or
-f yaml for machine-readable output.

Examples:
  metalcloud configuration get
  metalcloud configuration get platform
  metalcloud configuration get auth -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONFIGURATION_READ},
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			service := ""
			if len(args) > 0 {
				service = args[0]
			}

			return configuration.ConfigurationGet(cmd.Context(), service)
		},
	}

	configurationReplaceCmd = &cobra.Command{
		Use:   "replace service",
		Short: "Replace a section of the global platform configuration",
		Long: `Replace the ENTIRE configuration of one service section of the platform (HTTP PUT).

Any setting missing from the supplied document is reset by the platform: use
'configuration update' to change individual settings. This affects every user of
the installation - take a copy with 'configuration get <service>' first.

Required Arguments:
  service   The configuration section to replace. One of: ` + configurationServices + `

Required Flags:
  --config-source string   Source of the new section configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud configuration get platform -f json > platform.json
  metalcloud configuration replace platform --config-source platform.json
  cat platform.yaml | metalcloud configuration replace platform --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONFIGURATION_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(configurationFlags.configSource)
			if err != nil {
				return err
			}

			return configuration.ConfigurationReplace(cmd.Context(), args[0], config)
		},
	}

	configurationUpdateCmd = &cobra.Command{
		Use:     "update service",
		Aliases: []string{"edit", "patch"},
		Short:   "Update part of a section of the global platform configuration",
		Long: `Merge a partial document into one service section of the global platform
configuration (HTTP PATCH). Settings absent from the document are left untouched.

This changes platform-wide settings that affect every user of the installation.

Required Arguments:
  service   The configuration section to update. One of: ` + configurationServices + `

Required Flags:
  --config-source string   Source of the partial configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud configuration update notification --config-source notification-patch.json
  echo '{"smtp":{"port":587}}' | metalcloud configuration update notification --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONFIGURATION_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(configurationFlags.configSource)
			if err != nil {
				return err
			}

			return configuration.ConfigurationUpdate(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(configurationCmd)

	configurationCmd.AddCommand(configurationGetCmd)

	configurationCmd.AddCommand(configurationReplaceCmd)
	configurationReplaceCmd.Flags().StringVar(&configurationFlags.configSource, "config-source", "", "Source of the new section configuration. Can be 'pipe' or path to a JSON/YAML file.")
	configurationReplaceCmd.MarkFlagRequired("config-source")

	configurationCmd.AddCommand(configurationUpdateCmd)
	configurationUpdateCmd.Flags().StringVar(&configurationFlags.configSource, "config-source", "", "Source of the partial configuration. Can be 'pipe' or path to a JSON/YAML file.")
	configurationUpdateCmd.MarkFlagRequired("config-source")
}
