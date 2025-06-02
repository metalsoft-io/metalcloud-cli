package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	siteFlags = struct {
		configSource string
	}{}

	siteCmd = &cobra.Command{
		Use:     "site [command]",
		Aliases: []string{"datacenter", "dc"},
		Short:   "Site management",
		Long:    `Site management commands.`,
	}

	siteListCmd = &cobra.Command{
		Use:          "list [flags...]",
		Aliases:      []string{"ls"},
		Short:        "List all sites.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteList(cmd.Context())
		},
	}

	siteGetCmd = &cobra.Command{
		Use:          "get site_id_or_name",
		Aliases:      []string{"show"},
		Short:        "Get site details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGet(cmd.Context(), args[0])
		},
	}

	siteCreateCmd = &cobra.Command{
		Use:          "create name",
		Aliases:      []string{"new"},
		Short:        "Create new site.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteCreate(cmd.Context(), args[0])
		},
	}

	siteUpdateCmd = &cobra.Command{
		Use:          "update site_id_or_name [new_label]",
		Aliases:      []string{"edit"},
		Short:        "Update site configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			return site.SiteUpdate(cmd.Context(), args[0], name)
		},
	}

	siteDecommissionCmd = &cobra.Command{
		Use:          "decommission site_id_or_name",
		Aliases:      []string{"archive"},
		Short:        "Decommission site.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteDecommission(cmd.Context(), args[0])
		},
	}

	siteGetAgentsCmd = &cobra.Command{
		Use:          "agents site_id_or_name",
		Aliases:      []string{"get-agents", "list-agents"},
		Short:        "Get agents for a site.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGetAgents(cmd.Context(), args[0])
		},
	}

	siteGetConfigCmd = &cobra.Command{
		Use:          "get-config site_id_or_name",
		Aliases:      []string{"config", "show-config"},
		Short:        "Get site configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGetConfig(cmd.Context(), args[0])
		},
	}

	siteUpdateConfigCmd = &cobra.Command{
		Use:          "update-config site_id_or_name json_config",
		Aliases:      []string{"edit-config"},
		Short:        "Update site configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(siteFlags.configSource)
			if err != nil {
				return err
			}

			return site.SiteUpdateConfig(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(siteCmd)

	siteCmd.AddCommand(siteListCmd)

	siteCmd.AddCommand(siteGetCmd)

	siteCmd.AddCommand(siteCreateCmd)

	siteCmd.AddCommand(siteUpdateCmd)

	siteCmd.AddCommand(siteDecommissionCmd)

	siteCmd.AddCommand(siteGetAgentsCmd)

	siteCmd.AddCommand(siteGetConfigCmd)

	siteCmd.AddCommand(siteUpdateConfigCmd)
	siteUpdateConfigCmd.Flags().StringVar(&siteFlags.configSource, "config-source", "", "Source of the site configuration. Can be 'pipe' or path to a JSON file.")
}
