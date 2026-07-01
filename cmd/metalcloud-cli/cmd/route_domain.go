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
these commands to list, create, update, and delete route domains.`,
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
}
