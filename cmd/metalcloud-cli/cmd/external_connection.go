package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/external_connection"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	externalConnectionFlags = struct {
		configSource   string
		filterId       []string
		filterFabricId []string
		filterLabel    []string
		filterName     []string
		label          string
		name           string
		fabric         string
		interfaceIds   []string
	}{}

	externalConnectionCmd = &cobra.Command{
		Use:     "external-connection [command]",
		Aliases: []string{"ext-conn", "external-connections"},
		Short:   "External connection management",
		Long: `Manage external connections: the network device interfaces of a fabric that
carry traffic towards networks outside the MetalSoft managed infrastructure.

Command categories:
  Lifecycle:        list, get, create, update, delete, config-example
  Interfaces:       interface list|get|add|update|remove
  Logical networks: logical-network list|get|add|remove
  Discovery:        get-network-device-interfaces`,
	}

	externalConnectionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List external connections",
		Long: `List all external connections.

Optional Flags:
  --filter-id strings          Filter by external connection ID. Repeatable or comma-separated.
  --filter-fabric-id strings   Filter by fabric ID. Repeatable or comma-separated.
  --filter-label strings       Filter by label. Repeatable or comma-separated.
  --filter-name strings        Filter by name. Repeatable or comma-separated.

Examples:
  metalcloud external-connection list
  metalcloud ext-conn list --filter-fabric-id 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionList(cmd.Context(), external_connection.ExternalConnectionListFilters{
				Id:       externalConnectionFlags.filterId,
				FabricId: externalConnectionFlags.filterFabricId,
				Label:    externalConnectionFlags.filterLabel,
				Name:     externalConnectionFlags.filterName,
			})
		},
	}

	externalConnectionGetCmd = &cobra.Command{
		Use:     "get external_connection_id_or_label",
		Aliases: []string{"show"},
		Short:   "Get external connection details",
		Long: `Get the details of an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection get 12
  metalcloud ext-conn get dc1-external-connection`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionGet(cmd.Context(), args[0])
		},
	}

	externalConnectionConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example external connection create configuration",
		Long: `Print an example configuration that can be edited and passed to
'external-connection create --config-source'.

Examples:
  metalcloud external-connection config-example > external-connection.json
  metalcloud ext-conn config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionConfigExample(cmd.Context())
		},
	}

	externalConnectionCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create an external connection",
		Long: `Create a new external connection.

The external connection can be described either with a configuration file
(--config-source) or with individual flags (--label, --name, --fabric).

Required Flags (one of):
  --config-source string   Source of the external connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Label of the new external connection (used together with the flags below)

Optional Flags (when not using --config-source):
  --name string            Name of the new external connection (required with --label)
  --fabric string          ID or name of the fabric the connection belongs to (required with --label)
  --interface-ids strings  Network device interface IDs to attach. Repeatable or comma-separated.

Examples:
  metalcloud external-connection create --config-source external-connection.json
  cat ext-conn.yaml | metalcloud ext-conn create --config-source pipe
  metalcloud ext-conn create --label dc1-ext --name "DC1 external" --fabric dc1-fabric --interface-ids 101,102`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateExternalConnection

			if externalConnectionFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(externalConnectionFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				fabricId, err := fabric.ResolveFabricNumericId(cmd.Context(), externalConnectionFlags.fabric)
				if err != nil {
					return err
				}

				create = sdk.CreateExternalConnection{
					Label:    externalConnectionFlags.label,
					Name:     externalConnectionFlags.name,
					FabricId: fabricId,
				}

				interfaceIds, err := utils.GetInt64SliceFromStrings(externalConnectionFlags.interfaceIds)
				if err != nil {
					return err
				}
				for _, interfaceId := range interfaceIds {
					create.ExternalConnectionInterfaces = append(create.ExternalConnectionInterfaces,
						sdk.CreateExternalConnectionInterface{NetworkDeviceInterfaceId: interfaceId})
				}
			}

			return external_connection.ExternalConnectionCreate(cmd.Context(), create)
		},
	}

	externalConnectionUpdateCmd = &cobra.Command{
		Use:     "update external_connection_id_or_label",
		Aliases: []string{"edit"},
		Short:   "Update an external connection",
		Long: `Update the label or name of an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud external-connection update 12 --config-source update.json
  echo '{"name":"new name"}' | metalcloud ext-conn update dc1-ext --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(externalConnectionFlags.configSource)
			if err != nil {
				return err
			}

			return external_connection.ExternalConnectionUpdate(cmd.Context(), args[0], config)
		},
	}

	externalConnectionDeleteCmd = &cobra.Command{
		Use:     "delete external_connection_id_or_label",
		Aliases: []string{"rm"},
		Short:   "Delete an external connection",
		Long: `Delete an external connection together with its interfaces.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection delete 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionDelete(cmd.Context(), args[0])
		},
	}

	externalConnectionInterfaceCmd = &cobra.Command{
		Use:     "interface [command]",
		Aliases: []string{"interfaces"},
		Short:   "External connection interface management",
		Long: `Manage the network device interfaces that make up an external connection.

Commands: list, get, add, update, remove`,
	}

	externalConnectionInterfaceListCmd = &cobra.Command{
		Use:     "list external_connection_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the interfaces of an external connection",
		Long: `List the network device interfaces attached to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection interface list 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionInterfaceList(cmd.Context(), args[0])
		},
	}

	externalConnectionInterfaceGetCmd = &cobra.Command{
		Use:     "get external_connection_id_or_label interface_id",
		Aliases: []string{"show"},
		Short:   "Get one interface of an external connection",
		Long: `Get the details of a single external connection interface.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  interface_id                      The ID of the external connection interface

Examples:
  metalcloud external-connection interface get 12 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionInterfaceGet(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionInterfaceAddCmd = &cobra.Command{
		Use:     "add external_connection_id_or_label network_device_interface_id",
		Aliases: []string{"new", "create"},
		Short:   "Add an interface to an external connection",
		Long: `Attach a network device interface to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  network_device_interface_id       The ID of the network device interface (see 'get-network-device-interfaces')

Examples:
  metalcloud external-connection interface add 12 101`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionInterfaceAdd(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionInterfaceUpdateCmd = &cobra.Command{
		Use:     "update external_connection_id_or_label interface_id network_device_interface_id",
		Aliases: []string{"edit"},
		Short:   "Update an interface of an external connection",
		Long: `Point an external connection interface at a different network device interface.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  interface_id                      The ID of the external connection interface
  network_device_interface_id       The new network device interface ID

Examples:
  metalcloud external-connection interface update 12 5 102`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionInterfaceUpdate(cmd.Context(), args[0], args[1], args[2])
		},
	}

	externalConnectionInterfaceRemoveCmd = &cobra.Command{
		Use:     "remove external_connection_id_or_label interface_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an interface from an external connection",
		Long: `Detach a network device interface from an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  interface_id                      The ID of the external connection interface

Examples:
  metalcloud external-connection interface remove 12 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionInterfaceRemove(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionLogicalNetworkCmd = &cobra.Command{
		Use:     "logical-network [command]",
		Aliases: []string{"logical-networks"},
		Short:   "External connection logical network management",
		Long: `Manage the logical networks attached to an external connection.

Commands: list, get, add, remove`,
	}

	externalConnectionLogicalNetworkListCmd = &cobra.Command{
		Use:     "list external_connection_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the logical networks of an external connection",
		Long: `List the logical networks attached to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection logical-network list 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionLogicalNetworkList(cmd.Context(), args[0])
		},
	}

	externalConnectionLogicalNetworkGetCmd = &cobra.Command{
		Use:     "get external_connection_id_or_label id",
		Aliases: []string{"show"},
		Short:   "Get one logical network attachment of an external connection",
		Long: `Get the details of a single external connection logical network attachment.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  id                                The ID of the attachment (not the logical network ID)

Examples:
  metalcloud external-connection logical-network get 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionLogicalNetworkGet(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionLogicalNetworkAddCmd = &cobra.Command{
		Use:     "add external_connection_id_or_label logical_network_id",
		Aliases: []string{"new", "create"},
		Short:   "Attach a logical network to an external connection",
		Long: `Attach a logical network to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  logical_network_id                The numeric ID of the logical network

Examples:
  metalcloud external-connection logical-network add 12 44`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionLogicalNetworkAdd(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionLogicalNetworkRemoveCmd = &cobra.Command{
		Use:     "remove external_connection_id_or_label id",
		Aliases: []string{"rm", "delete"},
		Short:   "Detach a logical network from an external connection",
		Long: `Detach a logical network from an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  id                                The ID of the attachment (not the logical network ID)

Examples:
  metalcloud external-connection logical-network remove 12 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.ExternalConnectionLogicalNetworkRemove(cmd.Context(), args[0], args[1])
		},
	}

	externalConnectionNetworkDeviceInterfacesCmd = &cobra.Command{
		Use:     "get-network-device-interfaces network_device_id",
		Aliases: []string{"network-device-interfaces"},
		Short:   "List a network device's interfaces and their external connections",
		Long: `List all interfaces of a network device together with the external connection
each interface belongs to, if any. Use it to find the network device interface IDs
accepted by 'external-connection interface add'.

Required Arguments:
  network_device_id   The ID or identifier of the network device

Examples:
  metalcloud external-connection get-network-device-interfaces 45
  metalcloud ext-conn get-network-device-interfaces border-leaf-01`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTERNAL_CONNECTIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return external_connection.NetworkDeviceInterfacesGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(externalConnectionCmd)

	externalConnectionCmd.AddCommand(externalConnectionListCmd)
	externalConnectionListCmd.Flags().StringSliceVar(&externalConnectionFlags.filterId, "filter-id", nil, "Filter by external connection ID.")
	externalConnectionListCmd.Flags().StringSliceVar(&externalConnectionFlags.filterFabricId, "filter-fabric-id", nil, "Filter by fabric ID.")
	externalConnectionListCmd.Flags().StringSliceVar(&externalConnectionFlags.filterLabel, "filter-label", nil, "Filter by label.")
	externalConnectionListCmd.Flags().StringSliceVar(&externalConnectionFlags.filterName, "filter-name", nil, "Filter by name.")

	externalConnectionCmd.AddCommand(externalConnectionGetCmd)
	externalConnectionCmd.AddCommand(externalConnectionConfigExampleCmd)

	externalConnectionCmd.AddCommand(externalConnectionCreateCmd)
	externalConnectionCreateCmd.Flags().StringVar(&externalConnectionFlags.configSource, "config-source", "", "Source of the new external connection configuration. Can be 'pipe' or path to a JSON/YAML file.")
	externalConnectionCreateCmd.Flags().StringVar(&externalConnectionFlags.label, "label", "", "Label of the new external connection.")
	externalConnectionCreateCmd.Flags().StringVar(&externalConnectionFlags.name, "name", "", "Name of the new external connection.")
	externalConnectionCreateCmd.Flags().StringVar(&externalConnectionFlags.fabric, "fabric", "", "ID or name of the fabric the external connection belongs to.")
	externalConnectionCreateCmd.Flags().StringSliceVar(&externalConnectionFlags.interfaceIds, "interface-ids", nil, "Network device interface IDs to attach to the new external connection.")
	externalConnectionCreateCmd.MarkFlagsOneRequired("config-source", "label")
	externalConnectionCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	externalConnectionCreateCmd.MarkFlagsRequiredTogether("label", "name", "fabric")

	externalConnectionCmd.AddCommand(externalConnectionUpdateCmd)
	externalConnectionUpdateCmd.Flags().StringVar(&externalConnectionFlags.configSource, "config-source", "", "Source of the updated external connection configuration. Can be 'pipe' or path to a JSON/YAML file.")
	externalConnectionUpdateCmd.MarkFlagsOneRequired("config-source")

	externalConnectionCmd.AddCommand(externalConnectionDeleteCmd)
	externalConnectionCmd.AddCommand(externalConnectionNetworkDeviceInterfacesCmd)

	externalConnectionCmd.AddCommand(externalConnectionInterfaceCmd)
	externalConnectionInterfaceCmd.AddCommand(externalConnectionInterfaceListCmd)
	externalConnectionInterfaceCmd.AddCommand(externalConnectionInterfaceGetCmd)
	externalConnectionInterfaceCmd.AddCommand(externalConnectionInterfaceAddCmd)
	externalConnectionInterfaceCmd.AddCommand(externalConnectionInterfaceUpdateCmd)
	externalConnectionInterfaceCmd.AddCommand(externalConnectionInterfaceRemoveCmd)

	externalConnectionCmd.AddCommand(externalConnectionLogicalNetworkCmd)
	externalConnectionLogicalNetworkCmd.AddCommand(externalConnectionLogicalNetworkListCmd)
	externalConnectionLogicalNetworkCmd.AddCommand(externalConnectionLogicalNetworkGetCmd)
	externalConnectionLogicalNetworkCmd.AddCommand(externalConnectionLogicalNetworkAddCmd)
	externalConnectionLogicalNetworkCmd.AddCommand(externalConnectionLogicalNetworkRemoveCmd)
}
