package cmd

import (
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/endpoint_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

// Endpoint Instance management commands.
var (
	endpointInstanceFlags = struct {
		configSource             string
		infrastructure           string
		filterInfrastructureId   []string
		filterGroupId            []string
		filterEndpointId         []string
		filterServiceStatus      []string
		filterConfigEndpointId   []string
		filterConfigDeployStatus []string
		filterConfigDeployType   []string
		label                    string
		groupId                  int64
		endpointId               int64
		tags                     []string
	}{}

	endpointInstanceCmd = &cobra.Command{
		Use:     "endpoint-instance [command]",
		Aliases: []string{"ei"},
		Short:   "Endpoint instance management",
		Long: `Manage endpoint instances.

An endpoint instance is the deployment of an endpoint inside an infrastructure. It is
created within an infrastructure but is addressed by its own ID afterwards.

Command categories:
  Lifecycle:      list, get, create, delete
  Configuration:  config, update-config, update-meta, config-example

Use "metalcloud-cli endpoint-instance [command] --help" for detailed information about each command.`,
	}

	endpointInstanceListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List endpoint instances",
		Long: `List endpoint instances.

Without flags every endpoint instance visible to the user is listed. Pass
--infrastructure to restrict the listing to a single infrastructure; the
infrastructure is resolved by ID or by label.

Optional Flags:
  --infrastructure string               List only the endpoint instances of this infrastructure (ID or label).
  --filter-infrastructure-id strings    Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-group-id strings             Filter by endpoint instance group ID.
  --filter-endpoint-id strings          Filter by endpoint ID.
  --filter-service-status strings       Filter by service status.
  --filter-config-endpoint-id strings   Filter by the endpoint ID of the pending configuration.
  --filter-config-deploy-status strings Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings   Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli endpoint-instance list
  metalcloud-cli endpoint-instance list --infrastructure prod-env
  metalcloud-cli ei ls --filter-service-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceList(cmd.Context(), endpointInstanceFlags.infrastructure, endpoint_instance.EndpointInstanceFilters{
				InfrastructureId:   endpointInstanceFlags.filterInfrastructureId,
				GroupId:            endpointInstanceFlags.filterGroupId,
				EndpointId:         endpointInstanceFlags.filterEndpointId,
				ServiceStatus:      endpointInstanceFlags.filterServiceStatus,
				ConfigEndpointId:   endpointInstanceFlags.filterConfigEndpointId,
				ConfigDeployStatus: endpointInstanceFlags.filterConfigDeployStatus,
				ConfigDeployType:   endpointInstanceFlags.filterConfigDeployType,
			})
		},
	}

	endpointInstanceGetCmd = &cobra.Command{
		Use:     "get endpoint_instance_id",
		Aliases: []string{"show"},
		Short:   "Get endpoint instance details",
		Long: `Get the details of an endpoint instance.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Examples:
  metalcloud-cli endpoint-instance get 42
  metalcloud-cli ei show 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGet(cmd.Context(), args[0])
		},
	}

	endpointInstanceConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example endpoint instance create configuration",
		Long: `Print an example configuration that can be edited and passed to
'endpoint-instance create --config-source'.

Examples:
  metalcloud-cli endpoint-instance config-example > endpoint-instance.json
  metalcloud-cli ei config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceConfigExample(cmd.Context())
		},
	}

	endpointInstanceCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create an endpoint instance in an infrastructure",
		Long: `Create a new endpoint instance in an infrastructure.

The instance can be described either with a configuration file (--config-source) or
with individual flags (--endpoint-id, --label, ...).

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Required Flags (one of):
  --config-source string   Source of the new endpoint instance configuration. Can be 'pipe' or path to a JSON/YAML file.
  --endpoint-id int        ID of the endpoint deployed by this instance

Optional Flags (when not using --config-source):
  --label string      Label of the new endpoint instance
  --group-id int      ID of the endpoint instance group the instance belongs to
  --tags strings      Tags of the new endpoint instance

Examples:
  metalcloud-cli endpoint-instance create 1234 --endpoint-id 7
  metalcloud-cli endpoint-instance create prod-env --config-source endpoint-instance.json
  cat endpoint-instance.yaml | metalcloud-cli ei new prod-env --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.EndpointInstanceCreate

			if endpointInstanceFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.EndpointInstanceCreate{
					EndpointId: endpointInstanceFlags.endpointId,
				}
				if endpointInstanceFlags.label != "" {
					create.Label = sdk.PtrString(endpointInstanceFlags.label)
				}
				if endpointInstanceFlags.groupId != 0 {
					create.GroupId = sdk.PtrInt64(endpointInstanceFlags.groupId)
				}
				if len(endpointInstanceFlags.tags) > 0 {
					create.Tags = endpointInstanceFlags.tags
				}
			}

			return endpoint_instance.EndpointInstanceCreate(cmd.Context(), args[0], create)
		},
	}

	endpointInstanceDeleteCmd = &cobra.Command{
		Use:     "delete endpoint_instance_id",
		Aliases: []string{"rm"},
		Short:   "Delete an endpoint instance",
		Long: `Delete an endpoint instance.

The current revision of the instance is fetched first and sent as If-Match, so the
delete fails if the instance changed in the meantime.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Examples:
  metalcloud-cli endpoint-instance delete 42
  metalcloud-cli ei rm 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceDelete(cmd.Context(), args[0])
		},
	}

	endpointInstanceConfigCmd = &cobra.Command{
		Use:     "config endpoint_instance_id",
		Aliases: []string{"get-config"},
		Short:   "Get the pending configuration of an endpoint instance",
		Long: `Get the pending (not yet deployed) configuration of an endpoint instance.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Examples:
  metalcloud-cli endpoint-instance config 42
  metalcloud-cli ei get-config 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceConfigGet(cmd.Context(), args[0])
		},
	}

	endpointInstanceUpdateConfigCmd = &cobra.Command{
		Use:     "update-config endpoint_instance_id",
		Aliases: []string{"update"},
		Short:   "Update the configuration of an endpoint instance",
		Long: `Update the pending configuration of an endpoint instance.

The current configuration is fetched first and its revision is sent as If-Match,
because writes below the /config path are guarded by the configuration revision.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance update-config 42 --config-source update.json
  echo '{"label":"new-label"}' | metalcloud-cli ei update-config 42 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceConfigUpdate(cmd.Context(), args[0], config)
		},
	}

	endpointInstanceUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta endpoint_instance_id",
		Aliases: []string{"meta"},
		Short:   "Update the metadata of an endpoint instance",
		Long: `Update the metadata (tags) of an endpoint instance.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Required Flags:
  --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance update-meta 42 --config-source meta.json
  echo '{"tags":["prod"]}' | metalcloud-cli ei update-meta 42 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceMetaUpdate(cmd.Context(), args[0], config)
		},
	}
)

// Endpoint Instance Group management commands.
var (
	endpointInstanceGroupFlags = struct {
		configSource             string
		filterExtensionInstance  []string
		filterServiceStatus      []string
		filterConfigDeployStatus []string
		filterConfigDeployType   []string
		label                    string
		endpointGroupName        string
		extensionInstanceId      int64
		hostname                 string
		resourcePoolId           int64
		tags                     []string
		logicalNetworkId         string
		connectAccessMode        string
		accessMode               string
		tagged                   string
		mtu                      int32
		redundancy               string
		providesDefaultRoute     string
		ruleType                 string
		direction                string
		sequence                 int32
		forwardingAction         string
		enforcementPoint         string
		networkProtocol          string
		sourceAddress            string
		destinationAddress       string
		sourcePort               string
		destinationPort          string
	}{}

	endpointInstanceGroupCmd = &cobra.Command{
		Use:     "endpoint-instance-group [command]",
		Aliases: []string{"eig"},
		Short:   "Endpoint instance group management",
		Long: `Manage endpoint instance groups.

An endpoint instance group collects the endpoint instances of an infrastructure that
share the same configuration and network connections. Groups are created inside an
infrastructure but are addressed by their own ID afterwards.

Command categories:
  Lifecycle:      list, get, create, delete
  Configuration:  config, update-config, update-meta, config-example
  Members:        instances
  Networking:     network (connections and their security rules)

Use "metalcloud-cli endpoint-instance-group [command] --help" for detailed information about each command.`,
	}

	endpointInstanceGroupListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the endpoint instance groups of an infrastructure",
		Long: `List all endpoint instance groups of an infrastructure.

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Optional Flags:
  --filter-extension-instance-id strings  Filter by extension instance ID. Repeatable or comma-separated.
  --filter-service-status strings         Filter by service status.
  --filter-config-deploy-status strings   Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings     Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli endpoint-instance-group list 1234
  metalcloud-cli eig ls prod-env --filter-service-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupList(cmd.Context(), args[0], endpoint_instance.EndpointInstanceGroupFilters{
				ExtensionInstanceId: endpointInstanceGroupFlags.filterExtensionInstance,
				ServiceStatus:       endpointInstanceGroupFlags.filterServiceStatus,
				ConfigDeployStatus:  endpointInstanceGroupFlags.filterConfigDeployStatus,
				ConfigDeployType:    endpointInstanceGroupFlags.filterConfigDeployType,
			})
		},
	}

	endpointInstanceGroupGetCmd = &cobra.Command{
		Use:     "get endpoint_instance_group_id",
		Aliases: []string{"show"},
		Short:   "Get endpoint instance group details",
		Long: `Get the details of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group get 12
  metalcloud-cli eig show 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupGet(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example endpoint instance group create configuration",
		Long: `Print an example configuration that can be edited and passed to
'endpoint-instance-group create --config-source'.

Examples:
  metalcloud-cli endpoint-instance-group config-example > group.json
  metalcloud-cli eig config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupConfigExample(cmd.Context())
		},
	}

	endpointInstanceGroupCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create an endpoint instance group in an infrastructure",
		Long: `Create a new endpoint instance group in an infrastructure. The group starts empty.

The group can be described either with a configuration file (--config-source) or with
individual flags (--label, --endpoint-group-name, ...).

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Required Flags (one of):
  --config-source string        Source of the new group configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string                Label of the new endpoint instance group

Optional Flags (when not using --config-source):
  --endpoint-group-name string  Name of the endpoint group deployed by this group
  --extension-instance-id int   ID of the extension instance backing the group
  --hostname string             Custom hostname for the DNS load balancing record
  --resource-pool-id int        ID of the resource pool assigned to the group
  --tags strings                Tags of the new group

Examples:
  metalcloud-cli endpoint-instance-group create 1234 --label web-endpoints
  metalcloud-cli endpoint-instance-group create prod-env --config-source group.json
  cat group.yaml | metalcloud-cli eig new prod-env --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.EndpointInstanceGroupCreate

			if endpointInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.EndpointInstanceGroupCreate{
					Label: sdk.PtrString(endpointInstanceGroupFlags.label),
				}
				if endpointInstanceGroupFlags.endpointGroupName != "" {
					create.EndpointGroupName = sdk.PtrString(endpointInstanceGroupFlags.endpointGroupName)
				}
				if endpointInstanceGroupFlags.extensionInstanceId != 0 {
					create.ExtensionInstanceId = sdk.PtrInt64(endpointInstanceGroupFlags.extensionInstanceId)
				}
				if endpointInstanceGroupFlags.hostname != "" {
					create.Hostname = sdk.PtrString(endpointInstanceGroupFlags.hostname)
				}
				if endpointInstanceGroupFlags.resourcePoolId != 0 {
					create.ResourcePoolId = sdk.PtrInt64(endpointInstanceGroupFlags.resourcePoolId)
				}
				if len(endpointInstanceGroupFlags.tags) > 0 {
					create.Tags = endpointInstanceGroupFlags.tags
				}
			}

			return endpoint_instance.EndpointInstanceGroupCreate(cmd.Context(), args[0], create)
		},
	}

	endpointInstanceGroupDeleteCmd = &cobra.Command{
		Use:     "delete endpoint_instance_group_id",
		Aliases: []string{"rm"},
		Short:   "Delete an endpoint instance group",
		Long: `Delete an endpoint instance group.

The current revision of the group is fetched first and sent as If-Match, so the
delete fails if the group changed in the meantime.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group delete 12
  metalcloud-cli eig rm 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupDelete(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupConfigCmd = &cobra.Command{
		Use:     "config endpoint_instance_group_id",
		Aliases: []string{"get-config"},
		Short:   "Get the pending configuration of an endpoint instance group",
		Long: `Get the pending (not yet deployed) configuration of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group config 12
  metalcloud-cli eig get-config 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupConfigGet(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupUpdateConfigCmd = &cobra.Command{
		Use:     "update-config endpoint_instance_group_id",
		Aliases: []string{"update"},
		Short:   "Update the configuration of an endpoint instance group",
		Long: `Update the pending configuration of an endpoint instance group.

The current configuration is fetched first and its revision is sent as If-Match,
because writes below the /config path are guarded by the configuration revision.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group update-config 12 --config-source update.json
  echo '{"label":"new-label"}' | metalcloud-cli eig update-config 12 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceGroupConfigUpdate(cmd.Context(), args[0], config)
		},
	}

	endpointInstanceGroupUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta endpoint_instance_group_id",
		Aliases: []string{"meta"},
		Short:   "Update the metadata of an endpoint instance group",
		Long: `Update the metadata (tags) of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Required Flags:
  --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group update-meta 12 --config-source meta.json
  echo '{"tags":["prod"]}' | metalcloud-cli eig update-meta 12 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceGroupMetaUpdate(cmd.Context(), args[0], config)
		},
	}

	endpointInstanceGroupInstancesCmd = &cobra.Command{
		Use:     "instances endpoint_instance_group_id",
		Aliases: []string{"members"},
		Short:   "List the endpoint instances of a group",
		Long: `List all endpoint instances that belong to an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group instances 12
  metalcloud-cli eig members 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupInstances(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupNetworkCmd = &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"net"},
		Short:   "Manage the network configuration of an endpoint instance group",
		Long: `Manage the network configuration of an endpoint instance group.

Command categories:
  Configuration:  list, replace, config-example
  Connections:    connections, get, connect, update, disconnect
  Security:       acl (security rules of a connection)

Use "metalcloud-cli endpoint-instance-group network [command] --help" for detailed information about each command.`,
	}

	endpointInstanceGroupNetworkListCmd = &cobra.Command{
		Use:     "list endpoint_instance_group_id",
		Aliases: []string{"ls", "show"},
		Short:   "Get the network configuration of an endpoint instance group",
		Long: `Get the network endpoint group that holds the network configuration of an
endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group network list 12
  metalcloud-cli eig net ls 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupNetworkList(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupNetworkReplaceCmd = &cobra.Command{
		Use:   "replace endpoint_instance_group_id",
		Short: "Create or replace the network configuration of an endpoint instance group",
		Long: `Create or replace (PUT) the network configuration of an endpoint instance group.

Without --config-source the operation is sent with no body, which creates the network
configuration if it does not exist yet and returns it otherwise. With --config-source
the supplied document is sent as the request body together with the If-Match of the
group configuration.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Optional Flags:
  --config-source string   Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group network replace 12
  metalcloud-cli endpoint-instance-group network replace 12 --config-source networking.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceGroupNetworkReplace(cmd.Context(), args[0], config)
		},
	}

	endpointInstanceGroupNetworkConnectionsCmd = &cobra.Command{
		Use:     "connections endpoint_instance_group_id",
		Aliases: []string{"conns"},
		Short:   "List the network connections of an endpoint instance group",
		Long: `List all network connections of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group network connections 12
  metalcloud-cli eig net conns 12`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupNetworkConnections(cmd.Context(), args[0])
		},
	}

	endpointInstanceGroupNetworkGetCmd = &cobra.Command{
		Use:     "get endpoint_instance_group_id connection_id",
		Aliases: []string{"show"},
		Short:   "Get a network connection of an endpoint instance group",
		Long: `Get the details of one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Examples:
  metalcloud-cli endpoint-instance-group network get 12 5
  metalcloud-cli eig net show 12 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupNetworkGet(cmd.Context(), args[0], args[1])
		},
	}

	endpointInstanceGroupNetworkConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example network connection configuration",
		Long: `Print an example configuration that can be edited and passed to
'endpoint-instance-group network connect --config-source'.

Examples:
  metalcloud-cli endpoint-instance-group network config-example > connection.json
  metalcloud-cli eig net config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupNetworkConnectionConfigExample(cmd.Context())
		},
	}

	endpointInstanceGroupNetworkConnectCmd = &cobra.Command{
		Use:     "connect endpoint_instance_group_id",
		Aliases: []string{"add"},
		Short:   "Connect an endpoint instance group to a logical network",
		Long: `Create a network connection between an endpoint instance group and a logical network.

The connection can be described either with a configuration file (--config-source) or
with individual flags (--logical-network-id, --access-mode, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Required Flags (one of):
  --config-source string       Source of the new connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  --logical-network-id string  ID of the logical network to connect to

Optional Flags (when not using --config-source):
  --access-mode string            Access mode of the connection (default "l2")
  --tagged string                 Whether the logical network is tagged (true/false)
  --mtu int                       MTU of the connection
  --redundancy string             Redundancy mode of the connection
  --provides-default-route string Whether the connection provides the default route (true/false)

Examples:
  metalcloud-cli endpoint-instance-group network connect 12 --logical-network-id 7 --tagged true
  metalcloud-cli endpoint-instance-group network connect 12 --config-source connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateEndpointInstanceGroupNetworkConnection

			if endpointInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				tagged, err := eiParseOptionalBoolFlag(endpointInstanceGroupFlags.tagged, "tagged")
				if err != nil {
					return err
				}

				create = sdk.CreateEndpointInstanceGroupNetworkConnection{
					LogicalNetworkId: endpointInstanceGroupFlags.logicalNetworkId,
					AccessMode:       sdk.NetworkEndpointGroupAllowedAccessMode(endpointInstanceGroupFlags.connectAccessMode),
				}
				if tagged != nil {
					create.Tagged = *tagged
				}
				if endpointInstanceGroupFlags.mtu != 0 {
					create.Mtu = sdk.PtrInt32(endpointInstanceGroupFlags.mtu)
				}
				if endpointInstanceGroupFlags.providesDefaultRoute != "" {
					providesDefaultRoute, err := eiParseOptionalBoolFlag(endpointInstanceGroupFlags.providesDefaultRoute, "provides-default-route")
					if err != nil {
						return err
					}
					create.ProvidesDefaultRoute = providesDefaultRoute
				}
				if endpointInstanceGroupFlags.redundancy != "" {
					create.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
						Mode: sdk.NetworkEndpointGroupRedundancyMode(endpointInstanceGroupFlags.redundancy),
					})
				}
			}

			return endpoint_instance.EndpointInstanceGroupNetworkConnect(cmd.Context(), args[0], create)
		},
	}

	endpointInstanceGroupNetworkUpdateCmd = &cobra.Command{
		Use:     "update endpoint_instance_group_id connection_id",
		Aliases: []string{"edit"},
		Short:   "Update a network connection of an endpoint instance group",
		Long: `Update an existing network connection of an endpoint instance group.

The change can be described either with a configuration file (--config-source) or with
individual flags (--access-mode, --tagged, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Required Flags (one of):
  --config-source string          Source of the updated connection. Can be 'pipe' or path to a JSON/YAML file.
  --access-mode string            New access mode of the connection
  --tagged string                 Whether the logical network is tagged (true/false)
  --mtu int                       New MTU of the connection
  --redundancy string             New redundancy mode of the connection
  --provides-default-route string Whether the connection provides the default route (true/false)

Examples:
  metalcloud-cli endpoint-instance-group network update 12 5 --access-mode l2
  metalcloud-cli eig net edit 12 5 --config-source connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var update sdk.UpdateNetworkEndpointGroupLogicalNetwork

			if endpointInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &update); err != nil {
					return err
				}
			} else {
				if endpointInstanceGroupFlags.accessMode != "" {
					accessMode := sdk.NetworkEndpointGroupAllowedAccessMode(endpointInstanceGroupFlags.accessMode)
					update.AccessMode = &accessMode
				}
				if endpointInstanceGroupFlags.tagged != "" {
					tagged, err := eiParseOptionalBoolFlag(endpointInstanceGroupFlags.tagged, "tagged")
					if err != nil {
						return err
					}
					update.Tagged = tagged
				}
				if endpointInstanceGroupFlags.mtu != 0 {
					update.Mtu = sdk.PtrInt32(endpointInstanceGroupFlags.mtu)
				}
				if endpointInstanceGroupFlags.providesDefaultRoute != "" {
					providesDefaultRoute, err := eiParseOptionalBoolFlag(endpointInstanceGroupFlags.providesDefaultRoute, "provides-default-route")
					if err != nil {
						return err
					}
					update.ProvidesDefaultRoute = providesDefaultRoute
				}
				if endpointInstanceGroupFlags.redundancy != "" {
					update.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
						Mode: sdk.NetworkEndpointGroupRedundancyMode(endpointInstanceGroupFlags.redundancy),
					})
				}
			}

			return endpoint_instance.EndpointInstanceGroupNetworkUpdate(cmd.Context(), args[0], args[1], update)
		},
	}

	endpointInstanceGroupNetworkDisconnectCmd = &cobra.Command{
		Use:     "disconnect endpoint_instance_group_id connection_id",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove a network connection from an endpoint instance group",
		Long: `Remove a network connection from an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Examples:
  metalcloud-cli endpoint-instance-group network disconnect 12 5
  metalcloud-cli eig net rm 12 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupNetworkDisconnect(cmd.Context(), args[0], args[1])
		},
	}

	endpointInstanceGroupACLCmd = &cobra.Command{
		Use:     "acl [command]",
		Aliases: []string{"security"},
		Short:   "Manage the security rules of a network connection",
		Long: `Manage the security rules (ACLs) of one network connection of an endpoint
instance group.

Command categories:
  Rules:   list, get, add, update, remove, config-example

Use "metalcloud-cli endpoint-instance-group network acl [command] --help" for detailed information about each command.`,
	}

	endpointInstanceGroupACLListCmd = &cobra.Command{
		Use:     "list endpoint_instance_group_id connection_id",
		Aliases: []string{"ls"},
		Short:   "List the security rules of a network connection",
		Long: `List the security rules of one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Examples:
  metalcloud-cli endpoint-instance-group network acl list 12 5
  metalcloud-cli eig net acl ls 12 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupACLList(cmd.Context(), args[0], args[1])
		},
	}

	endpointInstanceGroupACLGetCmd = &cobra.Command{
		Use:     "get endpoint_instance_group_id connection_id rule_id",
		Aliases: []string{"show"},
		Short:   "Get a security rule of a network connection",
		Long: `Get the details of one security rule of a network connection.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection
  rule_id                      The numeric ID of the security rule

Examples:
  metalcloud-cli endpoint-instance-group network acl get 12 5 3
  metalcloud-cli eig net acl show 12 5 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupACLGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	endpointInstanceGroupACLConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example security rule configuration",
		Long: `Print an example configuration that can be edited and passed to
'endpoint-instance-group network acl add --config-source'.

Examples:
  metalcloud-cli endpoint-instance-group network acl config-example > rule.json
  metalcloud-cli eig net acl config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupACLConfigExample(cmd.Context())
		},
	}

	endpointInstanceGroupACLAddCmd = &cobra.Command{
		Use:     "add endpoint_instance_group_id connection_id",
		Aliases: []string{"create", "new"},
		Short:   "Add a security rule to a network connection",
		Long: `Add a security rule to one network connection of an endpoint instance group.

The rule can be described either with a configuration file (--config-source) or with
individual flags (--rule-type, --direction, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Required Flags (one of):
  --config-source string       Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.
  --rule-type string           Rule type: ipv4, ipv6 or mac

Optional Flags (when not using --config-source):
  --direction string           Rule direction: in or out (default "in")
  --sequence int               Evaluation order of the rule
  --forwarding-action string   Forwarding action: allow, deny, transit or discard (default "allow")
  --enforcement-point string   Enforcement point of the rule (default "svi")
  --network-protocol string    Network protocol of the rule
  --source-address string      Source address of the rule
  --destination-address string Destination address of the rule
  --source-port string         Source port of the rule
  --destination-port string    Destination port of the rule

Examples:
  metalcloud-cli endpoint-instance-group network acl add 12 5 --rule-type ipv4 --sequence 10 --source-address 10.0.0.0/24
  metalcloud-cli eig net acl add 12 5 --config-source rule.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateLogicalNetworkACL

			if endpointInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.CreateLogicalNetworkACL{
					RuleType:         sdk.ACLType(endpointInstanceGroupFlags.ruleType),
					Direction:        sdk.ACLDirection(endpointInstanceGroupFlags.direction),
					Sequence:         endpointInstanceGroupFlags.sequence,
					ForwardingAction: sdk.ACLForwardingAction(endpointInstanceGroupFlags.forwardingAction),
					EnforcementPoint: sdk.ACLEnforcementPoint(endpointInstanceGroupFlags.enforcementPoint),
				}
				if endpointInstanceGroupFlags.networkProtocol != "" {
					create.NetworkProtocol = sdk.PtrString(endpointInstanceGroupFlags.networkProtocol)
				}
				if endpointInstanceGroupFlags.sourceAddress != "" {
					create.SourceAddress = sdk.PtrString(endpointInstanceGroupFlags.sourceAddress)
				}
				if endpointInstanceGroupFlags.destinationAddress != "" {
					create.DestinationAddress = sdk.PtrString(endpointInstanceGroupFlags.destinationAddress)
				}
				if endpointInstanceGroupFlags.sourcePort != "" {
					create.SourcePort = sdk.PtrString(endpointInstanceGroupFlags.sourcePort)
				}
				if endpointInstanceGroupFlags.destinationPort != "" {
					create.DestinationPort = sdk.PtrString(endpointInstanceGroupFlags.destinationPort)
				}
			}

			return endpoint_instance.EndpointInstanceGroupACLAdd(cmd.Context(), args[0], args[1], create)
		},
	}

	endpointInstanceGroupACLUpdateCmd = &cobra.Command{
		Use:     "update endpoint_instance_group_id connection_id rule_id",
		Aliases: []string{"edit"},
		Short:   "Update a security rule of a network connection",
		Long: `Update a security rule of one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection
  rule_id                      The numeric ID of the security rule

Required Flags:
  --config-source string   Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli endpoint-instance-group network acl update 12 5 3 --config-source rule.json
  echo '{"forwardingAction":"deny"}' | metalcloud-cli eig net acl update 12 5 3 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(endpointInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return endpoint_instance.EndpointInstanceGroupACLUpdate(cmd.Context(), args[0], args[1], args[2], config)
		},
	}

	endpointInstanceGroupACLRemoveCmd = &cobra.Command{
		Use:     "remove endpoint_instance_group_id connection_id rule_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a security rule from a network connection",
		Long: `Remove a security rule from one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection
  rule_id                      The numeric ID of the security rule

Examples:
  metalcloud-cli endpoint-instance-group network acl remove 12 5 3
  metalcloud-cli eig net acl rm 12 5 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ENDPOINT_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return endpoint_instance.EndpointInstanceGroupACLRemove(cmd.Context(), args[0], args[1], args[2])
		},
	}
)

// eiParseOptionalBoolFlag parses a tri-state string flag: an empty value leaves
// the field unset, anything else must parse as a boolean.
func eiParseOptionalBoolFlag(value string, flagName string) (*bool, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid --%s value: '%s'", flagName, value)
	}

	return &parsed, nil
}

func init() {
	// Endpoint Instance management commands.
	rootCmd.AddCommand(endpointInstanceCmd)

	endpointInstanceCmd.AddCommand(endpointInstanceListCmd)
	endpointInstanceListCmd.Flags().StringVar(&endpointInstanceFlags.infrastructure, "infrastructure", "", "List only the endpoint instances of this infrastructure (ID or label).")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterGroupId, "filter-group-id", nil, "Filter by endpoint instance group ID.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterEndpointId, "filter-endpoint-id", nil, "Filter by endpoint ID.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterServiceStatus, "filter-service-status", nil, "Filter by service status.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterConfigEndpointId, "filter-config-endpoint-id", nil, "Filter by the endpoint ID of the pending configuration.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterConfigDeployStatus, "filter-config-deploy-status", nil, "Filter by the deploy status of the pending configuration.")
	endpointInstanceListCmd.Flags().StringSliceVar(&endpointInstanceFlags.filterConfigDeployType, "filter-config-deploy-type", nil, "Filter by the deploy type of the pending configuration.")

	endpointInstanceCmd.AddCommand(endpointInstanceGetCmd)
	endpointInstanceCmd.AddCommand(endpointInstanceConfigExampleCmd)

	endpointInstanceCmd.AddCommand(endpointInstanceCreateCmd)
	endpointInstanceCreateCmd.Flags().StringVar(&endpointInstanceFlags.configSource, "config-source", "", "Source of the new endpoint instance configuration. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceCreateCmd.Flags().StringVar(&endpointInstanceFlags.label, "label", "", "Label of the new endpoint instance.")
	endpointInstanceCreateCmd.Flags().Int64Var(&endpointInstanceFlags.groupId, "group-id", 0, "ID of the endpoint instance group the instance belongs to.")
	endpointInstanceCreateCmd.Flags().Int64Var(&endpointInstanceFlags.endpointId, "endpoint-id", 0, "ID of the endpoint deployed by this instance.")
	endpointInstanceCreateCmd.Flags().StringSliceVar(&endpointInstanceFlags.tags, "tags", nil, "Tags of the new endpoint instance.")
	endpointInstanceCreateCmd.MarkFlagsOneRequired("config-source", "endpoint-id")
	endpointInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "endpoint-id")
	endpointInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	endpointInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "group-id")
	endpointInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "tags")

	endpointInstanceCmd.AddCommand(endpointInstanceDeleteCmd)
	endpointInstanceCmd.AddCommand(endpointInstanceConfigCmd)

	endpointInstanceCmd.AddCommand(endpointInstanceUpdateConfigCmd)
	endpointInstanceUpdateConfigCmd.Flags().StringVar(&endpointInstanceFlags.configSource, "config-source", "", "Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceUpdateConfigCmd.MarkFlagRequired("config-source")

	endpointInstanceCmd.AddCommand(endpointInstanceUpdateMetaCmd)
	endpointInstanceUpdateMetaCmd.Flags().StringVar(&endpointInstanceFlags.configSource, "config-source", "", "Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceUpdateMetaCmd.MarkFlagRequired("config-source")

	// Endpoint Instance Group management commands.
	rootCmd.AddCommand(endpointInstanceGroupCmd)

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupListCmd)
	endpointInstanceGroupListCmd.Flags().StringSliceVar(&endpointInstanceGroupFlags.filterExtensionInstance, "filter-extension-instance-id", nil, "Filter by extension instance ID.")
	endpointInstanceGroupListCmd.Flags().StringSliceVar(&endpointInstanceGroupFlags.filterServiceStatus, "filter-service-status", nil, "Filter by service status.")
	endpointInstanceGroupListCmd.Flags().StringSliceVar(&endpointInstanceGroupFlags.filterConfigDeployStatus, "filter-config-deploy-status", nil, "Filter by the deploy status of the pending configuration.")
	endpointInstanceGroupListCmd.Flags().StringSliceVar(&endpointInstanceGroupFlags.filterConfigDeployType, "filter-config-deploy-type", nil, "Filter by the deploy type of the pending configuration.")

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupGetCmd)
	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupConfigExampleCmd)

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupCreateCmd)
	endpointInstanceGroupCreateCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the new group configuration. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupCreateCmd.Flags().StringVar(&endpointInstanceGroupFlags.label, "label", "", "Label of the new endpoint instance group.")
	endpointInstanceGroupCreateCmd.Flags().StringVar(&endpointInstanceGroupFlags.endpointGroupName, "endpoint-group-name", "", "Name of the endpoint group deployed by this group.")
	endpointInstanceGroupCreateCmd.Flags().Int64Var(&endpointInstanceGroupFlags.extensionInstanceId, "extension-instance-id", 0, "ID of the extension instance backing the group.")
	endpointInstanceGroupCreateCmd.Flags().StringVar(&endpointInstanceGroupFlags.hostname, "hostname", "", "Custom hostname for the DNS load balancing record.")
	endpointInstanceGroupCreateCmd.Flags().Int64Var(&endpointInstanceGroupFlags.resourcePoolId, "resource-pool-id", 0, "ID of the resource pool assigned to the group.")
	endpointInstanceGroupCreateCmd.Flags().StringSliceVar(&endpointInstanceGroupFlags.tags, "tags", nil, "Tags of the new group.")
	endpointInstanceGroupCreateCmd.MarkFlagsOneRequired("config-source", "label")
	endpointInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupDeleteCmd)
	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupConfigCmd)

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupUpdateConfigCmd)
	endpointInstanceGroupUpdateConfigCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupUpdateConfigCmd.MarkFlagRequired("config-source")

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupUpdateMetaCmd)
	endpointInstanceGroupUpdateMetaCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupUpdateMetaCmd.MarkFlagRequired("config-source")

	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupInstancesCmd)

	// Endpoint Instance Group network commands.
	endpointInstanceGroupCmd.AddCommand(endpointInstanceGroupNetworkCmd)

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkListCmd)

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkReplaceCmd)
	endpointInstanceGroupNetworkReplaceCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.")

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkConnectionsCmd)
	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkGetCmd)
	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkConfigExampleCmd)

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkConnectCmd)
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the new connection configuration. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.logicalNetworkId, "logical-network-id", "", "ID of the logical network to connect to.")
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.connectAccessMode, "access-mode", string(sdk.NETWORKENDPOINTGROUPALLOWEDACCESSMODE_L2), "Access mode of the connection.")
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.tagged, "tagged", "", "Whether the logical network is tagged (true/false).")
	endpointInstanceGroupNetworkConnectCmd.Flags().Int32Var(&endpointInstanceGroupFlags.mtu, "mtu", 0, "MTU of the connection.")
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.redundancy, "redundancy", "", "Redundancy mode of the connection.")
	endpointInstanceGroupNetworkConnectCmd.Flags().StringVar(&endpointInstanceGroupFlags.providesDefaultRoute, "provides-default-route", "", "Whether the connection provides the default route (true/false).")
	endpointInstanceGroupNetworkConnectCmd.MarkFlagsOneRequired("config-source", "logical-network-id")
	endpointInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "logical-network-id")

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkUpdateCmd)
	endpointInstanceGroupNetworkUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the updated connection. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupNetworkUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.accessMode, "access-mode", "", "New access mode of the connection.")
	endpointInstanceGroupNetworkUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.tagged, "tagged", "", "Whether the logical network is tagged (true/false).")
	endpointInstanceGroupNetworkUpdateCmd.Flags().Int32Var(&endpointInstanceGroupFlags.mtu, "mtu", 0, "New MTU of the connection.")
	endpointInstanceGroupNetworkUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.redundancy, "redundancy", "", "New redundancy mode of the connection.")
	endpointInstanceGroupNetworkUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.providesDefaultRoute, "provides-default-route", "", "Whether the connection provides the default route (true/false).")
	endpointInstanceGroupNetworkUpdateCmd.MarkFlagsOneRequired("config-source", "access-mode", "tagged", "mtu", "redundancy", "provides-default-route")
	endpointInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "access-mode")
	endpointInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "tagged")

	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupNetworkDisconnectCmd)

	// Endpoint Instance Group network security rule commands.
	endpointInstanceGroupNetworkCmd.AddCommand(endpointInstanceGroupACLCmd)

	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLListCmd)
	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLGetCmd)
	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLConfigExampleCmd)

	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLAddCmd)
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.ruleType, "rule-type", "", "Rule type: ipv4, ipv6 or mac.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.direction, "direction", string(sdk.ACLDIRECTION_IN), "Rule direction: in or out.")
	endpointInstanceGroupACLAddCmd.Flags().Int32Var(&endpointInstanceGroupFlags.sequence, "sequence", 0, "Evaluation order of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.forwardingAction, "forwarding-action", string(sdk.ACLFORWARDINGACTION_ALLOW), "Forwarding action: allow, deny, transit or discard.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.enforcementPoint, "enforcement-point", string(sdk.ACLENFORCEMENTPOINT_SVI), "Enforcement point of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.networkProtocol, "network-protocol", "", "Network protocol of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.sourceAddress, "source-address", "", "Source address of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.destinationAddress, "destination-address", "", "Destination address of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.sourcePort, "source-port", "", "Source port of the rule.")
	endpointInstanceGroupACLAddCmd.Flags().StringVar(&endpointInstanceGroupFlags.destinationPort, "destination-port", "", "Destination port of the rule.")
	endpointInstanceGroupACLAddCmd.MarkFlagsOneRequired("config-source", "rule-type")
	endpointInstanceGroupACLAddCmd.MarkFlagsMutuallyExclusive("config-source", "rule-type")

	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLUpdateCmd)
	endpointInstanceGroupACLUpdateCmd.Flags().StringVar(&endpointInstanceGroupFlags.configSource, "config-source", "", "Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.")
	endpointInstanceGroupACLUpdateCmd.MarkFlagRequired("config-source")

	endpointInstanceGroupACLCmd.AddCommand(endpointInstanceGroupACLRemoveCmd)
}
