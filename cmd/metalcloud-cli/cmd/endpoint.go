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
		filterSite               []string
		filterExternalId         []string
		configSource             string
		siteId                   int
		name                     string
		label                    string
		externalId               string
		networkDeviceInterfaceId int
		macAddress               string
	}{}

	endpointCmd = &cobra.Command{
		Use:     "endpoint [command]",
		Aliases: []string{"ep", "endpoints"},
		Short:   "Endpoint management",
		Long: `Endpoint management commands.

An endpoint represents a device attached to the fabric. Its interfaces bind the
endpoint to individual switch ports.

Available Commands:
  list, get, create, create-bulk, update, delete
  interfaces                      List the interfaces of an endpoint
  interface                       Get, add, update or remove one interface
  get-network-device-interfaces   List a switch's ports and their endpoints`,
	}

	endpointListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List endpoints",
		Long: `List all endpoints in MetalSoft with optional filtering.

This command displays a table of all endpoints available in the system. You can filter the results 
by site or external ID to narrow down the output.

Flags:
  --filter-site strings           Filter results by site name(s). Can be specified multiple times.
  --filter-external-id strings    Filter results by external ID(s). Can be specified multiple times.

Examples:
  metalcloud-cli endpoint list
  metalcloud-cli endpoint ls --filter-site "site1" --filter-site "site2"
  metalcloud-cli endpoint list --filter-external-id "ext-001"
  metalcloud-cli endpoint list --filter-site "production" --filter-external-id "api-endpoint"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointList(cmd.Context(),
				endpointFlags.filterSite,
				endpointFlags.filterExternalId)
		},
	}

	endpointGetCmd = &cobra.Command{
		Use:     "get endpoint_id",
		Aliases: []string{"show"},
		Short:   "Display detailed information about a specific endpoint",
		Long: `Display detailed information about a specific endpoint in MetalSoft.

This command retrieves and displays comprehensive information about a single endpoint, 
including its configuration, status, and associated details.

Arguments:
  endpoint_id    The unique identifier of the endpoint to retrieve (required)

Examples:
  metalcloud-cli endpoint get 123
  metalcloud-cli endpoint show 456
  metalcloud-cli ep get endpoint-uuid-123`,
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
		Short:   "Create a new endpoint",
		Long: `Create a new endpoint in MetalSoft.

You can specify the endpoint configuration either by providing individual flags (--site-id, --name, --label, --external-id) 
or by supplying a configuration file or piped JSON/YAML using --config-source. 
When using --config-source, the file or piped content must contain a valid endpoint configuration in JSON or YAML format.

Required flags (when not using --config-source):
  --site-id     Site ID where the endpoint will be created
  --name        Name of the endpoint
  --label       Label of the endpoint

Optional flags:
  --external-id string       External ID of the endpoint
  --config-source string     Source of configuration (file path or 'pipe')

Flag dependencies:
  - When using --config-source, all other flags are ignored
  - When not using --config-source, --site-id, --name, and --label are required together

Examples:
  metalcloud-cli endpoint create --site-id 1 --name "api-endpoint" --label "API Endpoint"
  metalcloud-cli endpoint create --site-id 1 --name "api-endpoint" --label "API Endpoint" --external-id "ext-001"
  metalcloud-cli endpoint create --config-source ./endpoint.json
  cat endpoint.yaml | metalcloud-cli endpoint create --config-source pipe`,
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
					SiteId: int64(endpointFlags.siteId),
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

	endpointCreateBulkCmd = &cobra.Command{
		Use:     "create-bulk",
		Aliases: []string{"import"},
		Short:   "Create multiple endpoints in one call",
		Long: `Create multiple endpoints in a single bulk call.

The configuration must be a JSON/YAML list of endpoint definitions. Each entry
requires at least siteId, name, and label; endpointInterfaces and other fields
are optional.

Each endpoint interface can be specified in one of two ways:
  - by numeric id:    { "networkDeviceInterfaceId": 12345 }
  - by label:         { "networkDevice": "leaf-01", "interface": "swp9s0" }

When specified by label, the network device is resolved by numeric id or by its
identifierString (switch hostname), and the interface by its name (e.g.
"swp9s0"); each device's ports are fetched once and cached. An optional
macAddress may be set on either form.

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file
                    containing a list of endpoints.

Examples:
  metalcloud-cli endpoint create-bulk --config-source endpoints.yaml
  cat endpoints.json | metalcloud-cli endpoint create-bulk --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointFlags.configSource)
			if err != nil {
				return err
			}
			return endpoint.EndpointCreateBulk(cmd.Context(), config)
		},
	}

	endpointUpdateCmd = &cobra.Command{
		Use:     "update endpoint_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing endpoint",
		Long: `Update an existing endpoint in MetalSoft.

You can update the endpoint by specifying new values for its name, label, or external ID using flags, 
or by providing a configuration file or piped JSON/YAML with --config-source. 
When using --config-source, the file or piped content must contain the fields to update in JSON or YAML format.

Arguments:
  endpoint_id    The unique identifier of the endpoint to update (required)

Optional flags:
  --name string              New name for the endpoint
  --label string             New label for the endpoint
  --external-id string       New external ID for the endpoint
  --config-source string     Source of configuration (file path or 'pipe')

Flag dependencies:
  - Flags are mutually exclusive with --config-source
  - At least one flag must be provided to update the endpoint
  - Only the fields provided will be updated

Examples:
  metalcloud-cli endpoint update 123 --name "new-name"
  metalcloud-cli endpoint update 123 --label "New Label" --external-id "new-ext-001"
  metalcloud-cli endpoint update 123 --config-source ./update.json
  cat update.yaml | metalcloud-cli endpoint update 123 --config-source pipe`,
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
		Use:     "delete endpoint_id",
		Aliases: []string{"rm", "del"},
		Short:   "Delete an endpoint",
		Long: `Delete an endpoint from MetalSoft.

This command permanently removes an endpoint from the system. This action cannot be undone.

Arguments:
  endpoint_id    The unique identifier of the endpoint to delete (required)

Examples:
  metalcloud-cli endpoint delete 123
  metalcloud-cli endpoint rm 456
  metalcloud-cli ep del endpoint-uuid-123

Warning: This operation is irreversible. Make sure you have the correct endpoint ID before proceeding.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointDelete(cmd.Context(), args[0])
		},
	}

	endpointInterfaceListCmd = &cobra.Command{
		Use:     "interfaces endpoint_id",
		Aliases: []string{"ifaces", "ifs"},
		Short:   "List interfaces of an endpoint",
		Long: `List all network interfaces of a specific endpoint in MetalSoft.

This command displays detailed information about all network interfaces associated 
with the specified endpoint, including their configuration and status.

Arguments:
  endpoint_id    The unique identifier of the endpoint whose interfaces to list (required)

Examples:
  metalcloud-cli endpoint interfaces 123
  metalcloud-cli endpoint ifaces 456
  metalcloud-cli ep ifs endpoint-uuid-123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointInterfaceList(cmd.Context(), args[0])
		},
	}

	endpointInterfaceCmd = &cobra.Command{
		Use:     "interface [command]",
		Aliases: []string{"iface", "if"},
		Short:   "Endpoint interface management",
		Long: `Manage the individual network interfaces of an endpoint.

Each endpoint interface binds the endpoint to one switch port
(networkDeviceInterfaceId), optionally with the MAC address seen on it.

Available Commands:
  get      Get one endpoint interface
  add      Add an interface to an endpoint
  update   Update an existing endpoint interface
  remove   Remove an interface from an endpoint

Use 'endpoint interfaces endpoint_id' to list all interfaces of an endpoint.`,
	}

	endpointInterfaceGetCmd = &cobra.Command{
		Use:     "get endpoint_id interface_id",
		Aliases: []string{"show"},
		Short:   "Get one interface of an endpoint",
		Long: `Display the details of a single endpoint interface.

Required Arguments:
  endpoint_id    The ID of the endpoint
  interface_id   The ID of the endpoint interface

Examples:
  metalcloud-cli endpoint interface get 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointInterfaceGet(cmd.Context(), args[0], args[1])
		},
	}

	endpointInterfaceAddCmd = &cobra.Command{
		Use:     "add endpoint_id",
		Aliases: []string{"create", "new"},
		Short:   "Add an interface to an endpoint",
		Long: `Add a new interface to an existing endpoint.

The interface can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

Required Arguments:
  endpoint_id  The ID of the endpoint

Required Flags (one of):
  --network-device-interface-id  The numeric ID of the switch port to attach
  --config-source                'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --mac-address  The MAC address seen on the interface

Examples:
  metalcloud-cli endpoint interface add 12 --network-device-interface-id 4567
  metalcloud-cli endpoint interface add 12 --network-device-interface-id 4567 --mac-address AA:BB:CC:DD:EE:FF
  metalcloud-cli endpoint interface add 12 --config-source interface.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var createConfig sdk.CreateEndpointInterface

			if endpointFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &createConfig); err != nil {
					return err
				}
			} else {
				createConfig.NetworkDeviceInterfaceId = int64(endpointFlags.networkDeviceInterfaceId)
				if endpointFlags.macAddress != "" {
					createConfig.MacAddress = &endpointFlags.macAddress
				}
			}

			return endpoint.EndpointInterfaceCreate(cmd.Context(), args[0], createConfig)
		},
	}

	endpointInterfaceUpdateCmd = &cobra.Command{
		Use:     "update endpoint_id interface_id",
		Aliases: []string{"edit"},
		Short:   "Update an interface of an endpoint",
		Long: `Update an existing endpoint interface.

The updates can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source. The interface's current revision is
sent as the If-Match entity tag.

Required Arguments:
  endpoint_id    The ID of the endpoint
  interface_id   The ID of the endpoint interface

Optional Flags:
  --network-device-interface-id  Move the interface to another switch port
  --mac-address                  The new MAC address of the interface
  --config-source                'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli endpoint interface update 12 3 --mac-address AA:BB:CC:DD:EE:FF
  metalcloud-cli endpoint interface update 12 3 --config-source interface.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var updateConfig sdk.UpdateEndpointInterface

			if endpointFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &updateConfig); err != nil {
					return err
				}
			} else {
				if endpointFlags.networkDeviceInterfaceId != 0 {
					interfaceId := int64(endpointFlags.networkDeviceInterfaceId)
					updateConfig.NetworkDeviceInterfaceId = &interfaceId
				}
				if endpointFlags.macAddress != "" {
					updateConfig.MacAddress = &endpointFlags.macAddress
				}
			}

			return endpoint.EndpointInterfaceUpdate(cmd.Context(), args[0], args[1], updateConfig)
		},
	}

	endpointInterfaceRemoveCmd = &cobra.Command{
		Use:     "remove endpoint_id interface_id",
		Aliases: []string{"delete", "rm", "del"},
		Short:   "Remove an interface from an endpoint",
		Long: `Remove an interface from an endpoint.

The interface's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  endpoint_id    The ID of the endpoint
  interface_id   The ID of the endpoint interface

Examples:
  metalcloud-cli endpoint interface remove 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointInterfaceDelete(cmd.Context(), args[0], args[1])
		},
	}

	endpointNetworkDeviceInterfacesCmd = &cobra.Command{
		Use:     "get-network-device-interfaces network_device_id",
		Aliases: []string{"device-interfaces"},
		Short:   "List a network device's interfaces and their endpoints",
		Long: `List every interface of a network device together with the endpoint attached to it.

Interfaces with no endpoint are listed with empty endpoint columns, so this is
the command to use when looking for a free switch port.

Required Arguments:
  network_device_id  The ID or label (identifier string) of the network device

Examples:
  metalcloud-cli endpoint get-network-device-interfaces 45
  metalcloud-cli endpoint get-network-device-interfaces leaf-01`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint.EndpointNetworkDeviceInterfaces(cmd.Context(), args[0])
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

	endpointCmd.AddCommand(endpointCreateBulkCmd)
	endpointCreateBulkCmd.Flags().StringVar(&endpointFlags.configSource, "config-source", "", "Source of the endpoints list. Can be 'pipe' or path to a JSON/YAML file.")
	endpointCreateBulkCmd.MarkFlagRequired("config-source")

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

	endpointCmd.AddCommand(endpointInterfaceCmd)

	endpointInterfaceCmd.AddCommand(endpointInterfaceGetCmd)

	endpointInterfaceCmd.AddCommand(endpointInterfaceAddCmd)
	endpointInterfaceAddCmd.Flags().StringVar(&endpointFlags.configSource, "config-source", "", "Source of the new endpoint interface configuration. Can be 'pipe' or path to a JSON file.")
	endpointInterfaceAddCmd.Flags().IntVar(&endpointFlags.networkDeviceInterfaceId, "network-device-interface-id", 0, "The network device interface (switch port) ID to attach.")
	endpointInterfaceAddCmd.Flags().StringVar(&endpointFlags.macAddress, "mac-address", "", "The MAC address of the endpoint interface.")
	endpointInterfaceAddCmd.MarkFlagsOneRequired("config-source", "network-device-interface-id")
	endpointInterfaceAddCmd.MarkFlagsMutuallyExclusive("config-source", "network-device-interface-id")
	endpointInterfaceAddCmd.MarkFlagsMutuallyExclusive("config-source", "mac-address")

	endpointInterfaceCmd.AddCommand(endpointInterfaceUpdateCmd)
	endpointInterfaceUpdateCmd.Flags().StringVar(&endpointFlags.configSource, "config-source", "", "Source of the endpoint interface updates. Can be 'pipe' or path to a JSON file.")
	endpointInterfaceUpdateCmd.Flags().IntVar(&endpointFlags.networkDeviceInterfaceId, "network-device-interface-id", 0, "The new network device interface (switch port) ID.")
	endpointInterfaceUpdateCmd.Flags().StringVar(&endpointFlags.macAddress, "mac-address", "", "The new MAC address of the endpoint interface.")
	endpointInterfaceUpdateCmd.MarkFlagsOneRequired("config-source", "network-device-interface-id", "mac-address")
	endpointInterfaceUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "network-device-interface-id")
	endpointInterfaceUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "mac-address")

	endpointInterfaceCmd.AddCommand(endpointInterfaceRemoveCmd)

	endpointCmd.AddCommand(endpointNetworkDeviceInterfacesCmd)
}
