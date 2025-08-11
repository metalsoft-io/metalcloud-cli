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
		Short:   "Manage sites (datacenters) and their configurations",
		Long: `Manage sites (datacenters) including creation, configuration updates, agent management, and decommissioning.

Sites represent physical datacenters or locations where infrastructure is deployed. Each site can contain
multiple servers, agents, and has its own configuration parameters.

Available Commands:
  list           List all sites with their basic information
  get            Retrieve detailed information about a specific site
  create         Create a new site with the specified name
  update         Update site properties like label/name
  decommission   Archive a site and mark it as inactive
  agents         List all agents deployed in a specific site
  get-config     Retrieve the configuration settings for a site
  update-config  Update site configuration using JSON input

Examples:
  # List all sites
  metalcloud-cli site list

  # Get details for a specific site
  metalcloud-cli site get "site-01"
  metalcloud-cli site get 12345

  # Create a new site
  metalcloud-cli site create "new-datacenter"

  # Update site configuration from file
  metalcloud-cli site update-config site-01 --config-source config.json`,
	}

	siteListCmd = &cobra.Command{
		Use:     "list [flags...]",
		Aliases: []string{"ls"},
		Short:   "List all sites with their basic information",
		Long: `List all sites (datacenters) available in the system with their basic information.

This command displays a table containing site details including ID, name, label, status,
and creation date. Sites are physical datacenters or locations where infrastructure
is deployed.

Required Permissions:
  sites:read - Permission to view site information

Optional Flags:
  Common output flags are available (--format, --output, etc.)

Examples:
  # List all sites in table format
  metalcloud-cli site list

  # List sites with JSON output
  metalcloud-cli site list --format json

  # List sites with custom output format
  metalcloud-cli site list --output "{{.ID}} {{.Label}}"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteList(cmd.Context())
		},
	}

	siteGetCmd = &cobra.Command{
		Use:     "get site_id_or_name",
		Aliases: []string{"show"},
		Short:   "Retrieve detailed information about a specific site",
		Long: `Retrieve detailed information about a specific site (datacenter) including configuration, 
status, and metadata.

This command fetches comprehensive information about a site including its ID, name, label,
creation date, status, and any associated configuration parameters. You can specify the
site by either its ID (numeric) or name (string).

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to retrieve information for

Required Permissions:
  sites:read - Permission to view site information

Examples:
  # Get site details by name
  metalcloud-cli site get "datacenter-01"

  # Get site details by ID
  metalcloud-cli site get 12345

  # Get site details with JSON output
  metalcloud-cli site get "datacenter-01" --format json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGet(cmd.Context(), args[0])
		},
	}

	siteCreateCmd = &cobra.Command{
		Use:     "create name",
		Aliases: []string{"new"},
		Short:   "Create a new site with the specified name",
		Long: `Create a new site (datacenter) with the specified name in the system.

This command creates a new site that can be used to host infrastructure components.
The site name must be unique within the system and will serve as the identifier
for the new datacenter location.

Required Arguments:
  name    The name for the new site (must be unique)

Required Permissions:
  sites:write - Permission to create sites

Examples:
  # Create a new site
  metalcloud-cli site create "datacenter-west"

  # Create a site with a descriptive name
  metalcloud-cli site create "production-datacenter-01"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteCreate(cmd.Context(), args[0])
		},
	}

	siteUpdateCmd = &cobra.Command{
		Use:     "update site_id_or_name [new_label]",
		Aliases: []string{"edit"},
		Short:   "Update site properties like label/name",
		Long: `Update site properties including the label/name of an existing site.

This command allows you to modify the properties of an existing site. Currently,
you can update the site's label (display name). The site is identified by either
its ID (numeric) or current name (string).

Required Arguments:
  site_id_or_name    Site identifier (ID or current name) to update

Optional Arguments:
  new_label          New label/name for the site (if not provided, only other properties are updated)

Required Permissions:
  sites:write - Permission to modify sites

Examples:
  # Update site label by name
  metalcloud-cli site update "old-datacenter" "new-datacenter-name"

  # Update site label by ID
  metalcloud-cli site update 12345 "production-datacenter"

  # Update site properties without changing label
  metalcloud-cli site update "datacenter-01"`,
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
		Use:     "decommission site_id_or_name",
		Aliases: []string{"archive"},
		Short:   "Archive a site and mark it as inactive",
		Long: `Decommission (archive) a site and mark it as inactive in the system.

This command permanently decommissions a site, making it unavailable for new deployments
while preserving historical data. Once decommissioned, a site cannot be reactivated
and any existing infrastructure should be migrated to other sites before decommissioning.

Warning: This operation is irreversible. Ensure all infrastructure has been properly
migrated before decommissioning a site.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to decommission

Required Permissions:
  sites:write - Permission to modify sites

Examples:
  # Decommission a site by name
  metalcloud-cli site decommission "old-datacenter"

  # Decommission a site by ID
  metalcloud-cli site decommission 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteDecommission(cmd.Context(), args[0])
		},
	}

	siteGetAgentsCmd = &cobra.Command{
		Use:     "agents site_id_or_name",
		Aliases: []string{"get-agents", "list-agents"},
		Short:   "List all agents deployed in a specific site",
		Long: `List all agents deployed in a specific site (datacenter) including their status and configuration.

This command retrieves information about all agents that are deployed within the specified site.
Agents are software components that manage and monitor infrastructure within a datacenter.
The output includes agent details such as ID, name, status, and deployment information.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to list agents for

Required Permissions:
  sites:read - Permission to view site information

Examples:
  # List agents in a site by name
  metalcloud-cli site agents "datacenter-01"

  # List agents in a site by ID
  metalcloud-cli site agents 12345

  # List agents with JSON output
  metalcloud-cli site agents "datacenter-01" --format json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGetAgents(cmd.Context(), args[0])
		},
	}

	siteGetConfigCmd = &cobra.Command{
		Use:     "get-config site_id_or_name",
		Aliases: []string{"config", "show-config"},
		Short:   "Retrieve the configuration settings for a site",
		Long: `Retrieve the configuration settings for a specific site (datacenter) in JSON format.

This command fetches the complete configuration settings for a site including
infrastructure parameters, deployment options, and other site-specific settings.
The configuration is returned in JSON format for easy parsing and modification.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to retrieve configuration for

Required Permissions:
  sites:read - Permission to view site information

Optional Flags:
  Common output flags are available (--format, --output, etc.)

Examples:
  # Get site configuration by name
  metalcloud-cli site get-config "datacenter-01"

  # Get site configuration by ID
  metalcloud-cli site get-config 12345

  # Save configuration to file
  metalcloud-cli site get-config "datacenter-01" > site-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SITES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return site.SiteGetConfig(cmd.Context(), args[0])
		},
	}

	siteUpdateConfigCmd = &cobra.Command{
		Use:     "update-config site_id_or_name",
		Aliases: []string{"edit-config"},
		Short:   "Update site configuration using JSON input",
		Long: `Update the configuration settings for a specific site (datacenter) using JSON input.

This command allows you to modify the configuration settings of an existing site.
The configuration can be provided through a file or piped from standard input.
The configuration must be in valid JSON format and contain the appropriate
site configuration parameters.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to update configuration for

Required Flags:
  --config-source    Source of the site configuration. Can be 'pipe' for stdin input
                     or path to a JSON file containing the configuration

Required Permissions:
  sites:write - Permission to modify sites

Dependencies:
  The --config-source flag is mandatory and must specify either:
  - 'pipe' to read JSON configuration from standard input
  - Path to a valid JSON file containing site configuration

Examples:
  # Update site configuration from a file
  metalcloud-cli site update-config "datacenter-01" --config-source config.json

  # Update site configuration from standard input
  cat config.json | metalcloud-cli site update-config 12345 --config-source pipe

  # Update site configuration with inline JSON
  echo '{"key": "value"}' | metalcloud-cli site update-config "datacenter-01" --config-source pipe`,
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
	siteUpdateConfigCmd.Flags().StringVar(&siteFlags.configSource, "config-source", "", "Source of the site configuration. Can be 'pipe' for stdin or path to a JSON file (required).")
	siteUpdateConfigCmd.MarkFlagRequired("config-source")
}
