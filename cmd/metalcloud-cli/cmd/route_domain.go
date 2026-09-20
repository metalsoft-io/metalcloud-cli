package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/route_domain"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	routeDomainFlags = struct {
		configSource string
	}{}

	routeDomainCmd = &cobra.Command{
		Use:     "route-domain [command]",
		Aliases: []string{"route-domains", "rd"},
		Short:   "Manage route domains (tenant VRFs)",
		Long: `Manage route domains in the MetalCloud infrastructure.

A route domain is a tenant VRF: an EVPN-L3VPN / VRF-Lite routing instance that L3
logical networks attach to (via a logical network profile's routeDomainId). Use
these commands to list, create, update, and delete route domains.

Available Commands:
  list, get, create, update, delete, config-example
  get-config           Get the config sub-resource of a route domain
  update-config        Update the global settings of a route domain config
  allocation-strategy  Manage the config's allocation strategies`,
	}

	routeDomainListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all route domains",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return route_domain.RouteDomainList(cmd.Context())
		},
	}

	routeDomainGetCmd = &cobra.Command{
		Use:          "get route_domain_id",
		Aliases:      []string{"show"},
		Short:        "Get details about a specific route domain",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return route_domain.RouteDomainGet(cmd.Context(), args[0])
		},
	}

	routeDomainConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Display a route domain configuration example",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return route_domain.RouteDomainConfigExample(cmd.Context())
		},
	}

	routeDomainCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new route domain",
		Long: `Create a new route domain (tenant VRF) from a JSON/YAML configuration.

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

The configuration must include the route domain kind (evpn_l3vpn | mpls_l3vpn |
vrf_lite) and its VRF allocation strategy; an l3evpn tenant VRF also carries an
L3VNI allocation strategy. Run 'route-domain config-example' for a template.

Examples:
  metalcloud-cli route-domain config-example > route-domain.yaml
  metalcloud-cli route-domain create --config-source route-domain.yaml
  cat route-domain.yaml | metalcloud-cli route-domain create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(routeDomainFlags.configSource)
			if err != nil {
				return err
			}
			return route_domain.RouteDomainCreate(cmd.Context(), config)
		},
	}

	routeDomainUpdateCmd = &cobra.Command{
		Use:          "update route_domain_id",
		Aliases:      []string{"modify"},
		Short:        "Update an existing route domain",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(routeDomainFlags.configSource)
			if err != nil {
				return err
			}
			return route_domain.RouteDomainUpdate(cmd.Context(), args[0], config)
		},
	}

	routeDomainDeleteCmd = &cobra.Command{
		Use:          "delete route_domain_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a route domain",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return route_domain.RouteDomainDelete(cmd.Context(), args[0])
		},
	}

	routeDomainGetConfigCmd = &cobra.Command{
		Use:     "get-config route_domain_id",
		Aliases: []string{"config", "show-config"},
		Short:   "Get the config of a route domain",
		Long: `Display the config object of a route domain.

The config is a separate sub-resource holding the desired state of the route
domain: its kind, the auto route distinguisher / route target settings, and the
allocation strategies that the 'route-domain allocation-strategy' commands
manage. It carries its own revision, distinct from the route domain's.

Required Arguments:
  route_domain_id  The ID of the route domain

Examples:
  metalcloud-cli route-domain get-config 2
  metalcloud-cli route-domain get-config 2 -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return route_domain.RouteDomainConfigGet(cmd.Context(), args[0])
		},
	}

	routeDomainUpdateConfigCmd = &cobra.Command{
		Use:     "update-config route_domain_id",
		Aliases: []string{"edit-config"},
		Short:   "Update the global settings of a route domain config",
		Long: `Update the global settings of a route domain's config.

Only the config's global settings are updated here (autoRouteDistinguisher and
autoRouteTarget); the allocation strategies are managed with the 'route-domain
allocation-strategy' commands. The config's own revision is sent as the
If-Match entity tag, so a concurrent change is rejected instead of being
overwritten.

Required Arguments:
  route_domain_id  The ID of the route domain

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli route-domain update-config 2 --config-source settings.json
  echo '{"autoRouteTarget":true}' | metalcloud-cli route-domain update-config 2 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(routeDomainFlags.configSource)
			if err != nil {
				return err
			}
			return route_domain.RouteDomainConfigUpdate(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(routeDomainCmd)

	routeDomainCmd.AddCommand(routeDomainListCmd)
	routeDomainCmd.AddCommand(routeDomainGetCmd)
	routeDomainCmd.AddCommand(routeDomainConfigExampleCmd)

	routeDomainCmd.AddCommand(routeDomainCreateCmd)
	routeDomainCreateCmd.Flags().StringVar(&routeDomainFlags.configSource, "config-source", "", "Source of the new route domain configuration. Can be 'pipe' or path to a JSON/YAML file.")
	routeDomainCreateCmd.MarkFlagRequired("config-source")

	routeDomainCmd.AddCommand(routeDomainUpdateCmd)
	routeDomainUpdateCmd.Flags().StringVar(&routeDomainFlags.configSource, "config-source", "", "Source of the route domain configuration updates. Can be 'pipe' or path to a JSON/YAML file.")
	routeDomainUpdateCmd.MarkFlagRequired("config-source")

	routeDomainCmd.AddCommand(routeDomainDeleteCmd)

	routeDomainCmd.AddCommand(routeDomainGetConfigCmd)

	routeDomainCmd.AddCommand(routeDomainUpdateConfigCmd)
	routeDomainUpdateConfigCmd.Flags().StringVar(&routeDomainFlags.configSource, "config-source", "", "Source of the route domain config updates. Can be 'pipe' or path to a JSON/YAML file.")
	routeDomainUpdateConfigCmd.MarkFlagsOneRequired("config-source")
}
