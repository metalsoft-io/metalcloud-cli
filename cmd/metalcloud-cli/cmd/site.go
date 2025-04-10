package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/spf13/cobra"
)

var (
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteList(cmd.Context())
		},
	}

	siteGetCmd = &cobra.Command{
		Use:          "get site_id_or_name",
		Aliases:      []string{"show"},
		Short:        "Get site details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_WRITE}, // TODO: Use specific permission
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_WRITE}, // TODO: Use specific permission
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_WRITE}, // TODO: Use specific permission
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGetAgents(cmd.Context(), args[0])
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
}
