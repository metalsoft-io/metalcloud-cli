package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/endpoint"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	endpointFlags = struct {
		filterSite       []string
		filterExternalId []string
		configSource     string
		siteId           int
		name             string
		label            string
		externalId       string
	}{}

	endpointCmd = &cobra.Command{
		Use:     "endpoint [command]",
		Aliases: []string{"ep", "endpoints"},
		Short:   "Endpoint management",
		Long:    `Endpoint management commands.`,
	}

	endpointListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all endpoints.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointList(cmd.Context(),
				endpointFlags.filterSite,
				endpointFlags.filterExternalId)
		},
	}

	endpointGetCmd = &cobra.Command{
		Use:          "get endpoint_id",
		Aliases:      []string{"show"},
		Short:        "Get endpoint info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointGet(cmd.Context(), args[0])
		},
	}

	endpointCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"add", "new"},
		Short:   "Create a new endpoint.",
		Long: `Create a new endpoint in MetalSoft.

You can specify the endpoint configuration either by providing individual flags (--site-id, --name, --label, --external-id) 
or by supplying a configuration file or piped JSON/YAML using --config-source. 
When using --config-source, the file or piped content must contain a valid endpoint configuration in JSON or YAML format.

Examples:
  metalcloud-cli endpoint create --site-id 1 --name "api-endpoint" --label "API Endpoint"
  metalcloud-cli endpoint create --config-source ./endpoint.json
  cat endpoint.yaml | metalcloud-cli endpoint create --config-source pipe

Note: --site-id, --name, and --label are required unless --config-source is used. 
Flags are mutually exclusive with --config-source.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var endpointConfig sdk.CreateEndpoint

			if endpointFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointFlags.configSource)
				if err != nil {
					return err
				}
				err = utils.UnmarshalContent(config, &endpointConfig)
				if err != nil {
					return err
				}
			} else {
				endpointConfig = sdk.CreateEndpoint{
					SiteId: float32(endpointFlags.siteId),
					Name:   endpointFlags.name,
					Label:  endpointFlags.label,
				}

				if endpointFlags.externalId != "" {
					endpointConfig.ExternalId = &endpointFlags.externalId
				}
			}

			return endpoint.EndpointCreate(cmd.Context(), endpointConfig)
		},
	}

	endpointUpdateCmd = &cobra.Command{
		Use:     "update endpoint_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing endpoint.",
		Long: `Update an existing endpoint in MetalSoft.

You can update the endpoint by specifying new values for its name, label, or external ID using flags, 
or by providing a configuration file or piped JSON/YAML with --config-source. 
When using --config-source, the file or piped content must contain the fields to update in JSON or YAML format.

Examples:
  metalcloud-cli endpoint update 123 --name "new-name"
  metalcloud-cli endpoint update 123 --config-source ./update.json
  cat update.yaml | metalcloud-cli endpoint update 123 --config-source pipe

Note: Flags are mutually exclusive with --config-source. 
Only the fields provided will be updated.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var endpointUpdates sdk.UpdateEndpoint

			if endpointFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointFlags.configSource)
				if err != nil {
					return err
				}
				err = utils.UnmarshalContent(config, &endpointUpdates)
				if err != nil {
					return err
				}
			} else {
				if endpointFlags.name != "" {
					endpointUpdates.Name = &endpointFlags.name
				}
				if endpointFlags.label != "" {
					endpointUpdates.Label = &endpointFlags.label
				}
				if endpointFlags.externalId != "" {
					endpointUpdates.ExternalId = &endpointFlags.externalId
				}
			}

			return endpoint.EndpointUpdate(cmd.Context(), args[0], endpointUpdates)
		},
	}

	endpointDeleteCmd = &cobra.Command{
		Use:          "delete endpoint_id",
		Aliases:      []string{"rm", "del"},
		Short:        "Delete an endpoint.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointDelete(cmd.Context(), args[0])
		},
	}

	endpointInterfaceListCmd = &cobra.Command{
		Use:          "interfaces endpoint_id",
		Aliases:      []string{"ifaces", "ifs"},
		Short:        "List interfaces of an endpoint.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointInterfaceList(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(endpointCmd)

	endpointCmd.AddCommand(endpointListCmd)
	endpointListCmd.Flags().StringSliceVar(&endpointFlags.filterSite, "filter-site", nil, "Filter the result by site.")
	endpointListCmd.Flags().StringSliceVar(&endpointFlags.filterExternalId, "filter-external-id", nil, "Filter the result by endpoint external Id.")

	endpointCmd.AddCommand(endpointGetCmd)

	endpointCmd.AddCommand(endpointCreateCmd)
	endpointCreateCmd.Flags().StringVar(&endpointFlags.configSource, "config-source", "", "Source of the new endpoint configuration. Can be 'pipe' or path to a JSON file.")
	endpointCreateCmd.Flags().IntVar(&endpointFlags.siteId, "site-id", 0, "The site ID where the endpoint will be created.")
	endpointCreateCmd.Flags().StringVar(&endpointFlags.name, "name", "", "The name of the endpoint.")
	endpointCreateCmd.Flags().StringVar(&endpointFlags.label, "label", "", "The label of the endpoint.")
	endpointCreateCmd.Flags().StringVar(&endpointFlags.externalId, "external-id", "", "The external ID of the endpoint.")
	endpointCreateCmd.MarkFlagsMutuallyExclusive("config-source", "site-id")
	endpointCreateCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	endpointCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	endpointCreateCmd.MarkFlagsMutuallyExclusive("config-source", "external-id")
	endpointCreateCmd.MarkFlagsRequiredTogether("site-id", "name", "label")

	endpointCmd.AddCommand(endpointUpdateCmd)
	endpointUpdateCmd.Flags().StringVar(&endpointFlags.configSource, "config-source", "", "Source of the endpoint configuration to update. Can be 'pipe' or path to a JSON file.")
	endpointUpdateCmd.Flags().StringVar(&endpointFlags.name, "name", "", "The new name of the endpoint.")
	endpointUpdateCmd.Flags().StringVar(&endpointFlags.label, "label", "", "The new label of the endpoint.")
	endpointUpdateCmd.Flags().StringVar(&endpointFlags.externalId, "external-id", "", "The new external ID of the endpoint.")
	endpointUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	endpointUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	endpointUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "external-id")

	endpointCmd.AddCommand(endpointDeleteCmd)

	endpointCmd.AddCommand(endpointInterfaceListCmd)
}
