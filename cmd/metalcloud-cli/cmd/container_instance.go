package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/container_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Container Instance management commands.
var (
	containerInstanceFlags = struct {
		configSource    string
		containerTypeId string
		groupId         string
		diskSizeGB      string
		tags            []string
		usage           string
		removeEmpty     bool
	}{}

	containerInstanceCmd = &cobra.Command{
		Use:     "container-instance [command]",
		Aliases: []string{"ci"},
		Short:   "Manage container instances within infrastructures",
		Long: `Manage container instances within infrastructures.

A container instance is the desired-state description of a single container in
an infrastructure. It belongs to a container instance group and is materialised
into a container when the infrastructure is deployed.

Available commands:
  list                  List the container instances of an infrastructure
  get                   Get details of a container instance
  config-example        Print an example container instance configuration
  create                Create a new container instance
  delete                Delete a container instance
  config                Show the pending configuration of a container instance
  update-config         Update the pending configuration of a container instance
  update-meta           Update the metadata (tags) of a container instance
  credentials           Show the credentials of a container instance
  variables             Show the variables of a container instance
  os-installation-data  Show the OS installation data of a container instance
  power-status          Show the power state of a container instance
  start                 Power on a container instance
  shutdown              Power off a container instance
  reboot                Reboot a container instance
  apply-type            Apply a container type on a container instance

Examples:
  metalcloud-cli container-instance list my-infra
  metalcloud-cli ci get my-infra 5678
  metalcloud-cli ci power-status my-infra 5678`,
	}

	containerInstanceListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the container instances of an infrastructure",
		Long: `List all container instances of an infrastructure.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List the container instances of infrastructure 1234
  metalcloud-cli container-instance list 1234

  # List them by infrastructure label
  metalcloud-cli ci ls my-infra`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceList(cmd.Context(), args[0])
		},
	}

	containerInstanceGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label container_instance_id",
		Aliases: []string{"show"},
		Short:   "Get container instance details",
		Long: `Get detailed information about a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Get the details of container instance 5678
  metalcloud-cli container-instance get my-infra 5678

  # Using the alias
  metalcloud-cli ci show 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGet(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print a container instance configuration example",
		Long: `Print a container instance configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli container-instance config-example

  # Save the example to a file
  metalcloud-cli ci config-example > container-instance.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceConfigExample(cmd.Context())
		},
	}

	containerInstanceCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new container instance",
		Long: `Create a new container instance in an infrastructure.

The instance can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string       Source of the new container instance configuration.
                               Can be 'pipe' or path to a JSON/YAML file.
  --container-type-id string   The container type of the new instance.
  --group-id string            The container instance group of the new instance.

Optional Flags:
  --disk-size-gb string  Disk size in GB. Defaults to the group disk size.
  --tags strings         Tags of the new container instance.

Flag Dependencies:
  --config-source is mutually exclusive with --container-type-id, --group-id,
  --disk-size-gb and --tags. One of --config-source or --container-type-id is
  required; --group-id is required when --container-type-id is used.

Examples:
  # Create a container instance from a file
  metalcloud-cli container-instance create my-infra --config-source instance.json

  # Create a container instance from flags
  metalcloud-cli ci new my-infra --container-type-id 42 --group-id 77 --disk-size-gb 40`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if containerInstanceFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(containerInstanceFlags.configSource)
				if err != nil {
					return err
				}

				return container_instance.ContainerInstanceCreate(cmd.Context(), args[0], config)
			}

			return container_instance.ContainerInstanceCreateFromFlags(cmd.Context(), args[0],
				containerInstanceFlags.containerTypeId,
				containerInstanceFlags.groupId,
				containerInstanceFlags.diskSizeGB,
				containerInstanceFlags.tags)
		},
	}

	containerInstanceDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label container_instance_id",
		Aliases: []string{"rm"},
		Short:   "Delete a container instance",
		Long: `Delete a container instance from an infrastructure.

The change takes effect when the infrastructure is deployed.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance to delete

Examples:
  # Delete container instance 5678
  metalcloud-cli container-instance delete my-infra 5678

  # Using the alias
  metalcloud-cli ci rm 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceDelete(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceConfigCmd = &cobra.Command{
		Use:     "config infrastructure_id_or_label container_instance_id",
		Aliases: []string{"get-config"},
		Short:   "Show the pending configuration of a container instance",
		Long: `Show the pending configuration of a container instance.

The configuration holds the changes that will be applied at the next deploy.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Show the configuration of container instance 5678
  metalcloud-cli container-instance config my-infra 5678

  # Using the alias
  metalcloud-cli ci get-config 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label container_instance_id",
		Aliases: []string{"edit-config"},
		Short:   "Update the pending configuration of a container instance",
		Long: `Update the pending configuration of a container instance.

The current configuration revision is fetched automatically and sent as the
If-Match header, so concurrent modifications are rejected by the API.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Required Flags:
  --config-source string  Source of the container instance configuration updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the configuration from a file
  metalcloud-cli container-instance update-config my-infra 5678 --config-source updates.json

  # Update the label from stdin
  echo '{"label":"web-01"}' | metalcloud-cli ci edit-config 1234 5678 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return container_instance.ContainerInstanceUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	containerInstanceUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label container_instance_id",
		Aliases: []string{"edit-meta"},
		Short:   "Update the metadata of a container instance",
		Long: `Update the metadata (tags) of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Required Flags:
  --config-source string  Source of the container instance metadata updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the tags from a file
  metalcloud-cli container-instance update-meta my-infra 5678 --config-source meta.json

  # Update the tags from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli ci edit-meta 1234 5678 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return container_instance.ContainerInstanceUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	containerInstanceCredentialsCmd = &cobra.Command{
		Use:     "credentials infrastructure_id_or_label container_instance_id",
		Aliases: []string{"creds"},
		Short:   "Show the credentials of a container instance",
		Long: `Show the credentials of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Show the credentials of container instance 5678
  metalcloud-cli container-instance credentials my-infra 5678

  # Using the alias
  metalcloud-cli ci creds 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceCredentials(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceVariablesCmd = &cobra.Command{
		Use:     "variables infrastructure_id_or_label container_instance_id",
		Aliases: []string{"vars"},
		Short:   "Show the variables of a container instance",
		Long: `Show the variables available in the context of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Optional Flags:
  --usage string  Restrict the variables to one usage type. One of HTTPRequest,
                  JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset.

Examples:
  # Show all variables of container instance 5678
  metalcloud-cli container-instance variables my-infra 5678

  # Show only the variables used by OS assets
  metalcloud-cli ci vars 1234 5678 --usage OSAsset`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceVariables(cmd.Context(), args[0], args[1], containerInstanceFlags.usage)
		},
	}

	containerInstanceOsInstallationDataCmd = &cobra.Command{
		Use:     "os-installation-data infrastructure_id_or_label container_instance_id",
		Aliases: []string{"os-data"},
		Short:   "Show the OS installation data of a container instance",
		Long: `Show the OS installation data of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Optional Flags:
  --usage string    Restrict the data to one usage type. One of HTTPRequest,
                    JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset.
  --remove-empty    Omit the empty entries from the response.

Examples:
  # Show the OS installation data of container instance 5678
  metalcloud-cli container-instance os-installation-data my-infra 5678

  # Show only the non-empty OS asset entries
  metalcloud-cli ci os-data 1234 5678 --usage OSAsset --remove-empty`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceOSInstallationData(cmd.Context(), args[0], args[1],
				containerInstanceFlags.usage, containerInstanceFlags.removeEmpty)
		},
	}

	containerInstancePowerStatusCmd = &cobra.Command{
		Use:     "power-status infrastructure_id_or_label container_instance_id",
		Aliases: []string{"power-state"},
		Short:   "Get the power status of a container instance",
		Long: `Get the current power status of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Show the power status of container instance 5678
  metalcloud-cli container-instance power-status my-infra 5678

  # Using the alias
  metalcloud-cli ci power-state 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstancePowerStatus(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceStartCmd = &cobra.Command{
		Use:     "start infrastructure_id_or_label container_instance_id",
		Aliases: []string{"power-on"},
		Short:   "Power on a container instance",
		Long: `Power on a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Power on container instance 5678
  metalcloud-cli container-instance start my-infra 5678

  # Using the alias
  metalcloud-cli ci power-on 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstancePowerControl(cmd.Context(), args[0], args[1], "start")
		},
	}

	containerInstanceShutdownCmd = &cobra.Command{
		Use:     "shutdown infrastructure_id_or_label container_instance_id",
		Aliases: []string{"power-off", "stop"},
		Short:   "Power off a container instance",
		Long: `Power off a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Power off container instance 5678
  metalcloud-cli container-instance shutdown my-infra 5678

  # Using an alias
  metalcloud-cli ci stop 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstancePowerControl(cmd.Context(), args[0], args[1], "shutdown")
		},
	}

	containerInstanceRebootCmd = &cobra.Command{
		Use:     "reboot infrastructure_id_or_label container_instance_id",
		Aliases: []string{"restart"},
		Short:   "Reboot a container instance",
		Long: `Reboot a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Reboot container instance 5678
  metalcloud-cli container-instance reboot my-infra 5678

  # Using the alias
  metalcloud-cli ci restart 1234 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstancePowerControl(cmd.Context(), args[0], args[1], "reboot")
		},
	}

	containerInstanceApplyTypeCmd = &cobra.Command{
		Use:     "apply-type infrastructure_id_or_label container_instance_id container_type_id",
		Aliases: []string{"set-type"},
		Short:   "Apply a container type on a container instance",
		Long: `Apply a container type on a container instance.

The current container instance revision is fetched automatically and sent as the
If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance
  container_type_id           The numeric ID of the container type to apply

Examples:
  # Apply container type 42 on container instance 5678
  metalcloud-cli container-instance apply-type my-infra 5678 42

  # Using the alias
  metalcloud-cli ci set-type 1234 5678 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceApplyType(cmd.Context(), args[0], args[1], args[2])
		},
	}
)

// Container Instance Group management commands.
var (
	containerInstanceGroupFlags = struct {
		configSource    string
		containerTypeId string
		osTemplateId    string
		vmPoolId        string
		diskSizeGB      string
		instanceCount   string
		tags            []string
		accessMode      string
		tagged          string
		redundancy      string
	}{}

	containerInstanceGroupCmd = &cobra.Command{
		Use:     "container-instance-group [command]",
		Aliases: []string{"cig"},
		Short:   "Manage container instance groups within infrastructures",
		Long: `Manage container instance groups within infrastructures.

A container instance group is a collection of identically configured container
instances. The group owns the shared sizing, the OS template and the network
connections of its instances.

Available commands:
  list            List the container instance groups of an infrastructure
  get             Get details of a container instance group
  config-example  Print an example container instance group configuration
  create          Create a new container instance group
  delete          Delete a container instance group
  config          Show the pending configuration of a group
  update-config   Update the pending configuration of a group
  update-meta     Update the metadata (tags) of a group
  instances       List the container instances of a group
  interfaces      List the network interfaces of a group
  interface       Get one network interface of a group
  apply-type      Apply a container type on a group
  network         Manage the network connections of a group

Examples:
  metalcloud-cli container-instance-group list my-infra
  metalcloud-cli cig get my-infra 77
  metalcloud-cli cig network list my-infra 77`,
	}

	containerInstanceGroupListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List the container instance groups of an infrastructure",
		Long: `List all container instance groups of an infrastructure.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List the container instance groups of infrastructure 1234
  metalcloud-cli container-instance-group list 1234

  # List them by infrastructure label
  metalcloud-cli cig ls my-infra`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupList(cmd.Context(), args[0])
		},
	}

	containerInstanceGroupGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"show"},
		Short:   "Get container instance group details",
		Long: `Get detailed information about a container instance group.

Required Arguments:
  infrastructure_id_or_label     The ID or the label of the infrastructure
  container_instance_group_id    The numeric ID of the container instance group

Examples:
  # Get the details of container instance group 77
  metalcloud-cli container-instance-group get my-infra 77

  # Using the alias
  metalcloud-cli cig show 1234 77`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupGet(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print a container instance group configuration example",
		Long: `Print a container instance group configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli container-instance-group config-example

  # Save the example to a file
  metalcloud-cli cig config-example > group.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupConfigExample(cmd.Context())
		},
	}

	containerInstanceGroupCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new container instance group",
		Long: `Create a new container instance group in an infrastructure.

The group can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string      Source of the new group configuration.
                              Can be 'pipe' or path to a JSON/YAML file.
  --container-type-id string  The container type of the group instances.
  --os-template-id string     The OS template of the group instances.
  --vm-pool-id string         The VM pool the group is provisioned on.
  --disk-size-gb string       Disk size in GB of each group instance.

Optional Flags:
  --instance-count string  Number of instances in the group.
  --tags strings           Tags of the new group.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --container-type-id is required; --os-template-id,
  --vm-pool-id and --disk-size-gb are required when --container-type-id is used.

Examples:
  # Create a group from a file
  metalcloud-cli container-instance-group create my-infra --config-source group.json

  # Create a group from flags
  metalcloud-cli cig new my-infra --container-type-id 42 --os-template-id 7 --vm-pool-id 3 --disk-size-gb 40 --instance-count 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if containerInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(containerInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}

				return container_instance.ContainerInstanceGroupCreate(cmd.Context(), args[0], config)
			}

			return container_instance.ContainerInstanceGroupCreateFromFlags(cmd.Context(), args[0],
				containerInstanceGroupFlags.containerTypeId,
				containerInstanceGroupFlags.osTemplateId,
				containerInstanceGroupFlags.vmPoolId,
				containerInstanceGroupFlags.diskSizeGB,
				containerInstanceGroupFlags.instanceCount,
				containerInstanceGroupFlags.tags)
		},
	}

	containerInstanceGroupDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"rm"},
		Short:   "Delete a container instance group",
		Long: `Delete a container instance group from an infrastructure.

All the container instances of the group are removed as well. The change takes
effect when the infrastructure is deployed.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the group to delete

Examples:
  # Delete container instance group 77
  metalcloud-cli container-instance-group delete my-infra 77

  # Using the alias
  metalcloud-cli cig rm 1234 77`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupDelete(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupConfigCmd = &cobra.Command{
		Use:     "config infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"get-config"},
		Short:   "Show the pending configuration of a container instance group",
		Long: `Show the pending configuration of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Examples:
  # Show the configuration of group 77
  metalcloud-cli container-instance-group config my-infra 77

  # Using the alias
  metalcloud-cli cig get-config 1234 77`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"edit-config"},
		Short:   "Update the pending configuration of a container instance group",
		Long: `Update the pending configuration of a container instance group.

The current configuration revision is fetched automatically and sent as the
If-Match header, so concurrent modifications are rejected by the API.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Required Flags:
  --config-source string  Source of the group configuration updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the configuration from a file
  metalcloud-cli container-instance-group update-config my-infra 77 --config-source updates.json

  # Update the label from stdin
  echo '{"label":"web"}' | metalcloud-cli cig edit-config 1234 77 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return container_instance.ContainerInstanceGroupUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	containerInstanceGroupUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"edit-meta"},
		Short:   "Update the metadata of a container instance group",
		Long: `Update the metadata (tags) of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Required Flags:
  --config-source string  Source of the group metadata updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the tags from a file
  metalcloud-cli container-instance-group update-meta my-infra 77 --config-source meta.json

  # Update the tags from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli cig edit-meta 1234 77 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return container_instance.ContainerInstanceGroupUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	containerInstanceGroupInstancesCmd = &cobra.Command{
		Use:     "instances infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"list-instances"},
		Short:   "List the container instances of a group",
		Long: `List the container instances that belong to a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Examples:
  # List the instances of group 77
  metalcloud-cli container-instance-group instances my-infra 77

  # Using the alias
  metalcloud-cli cig list-instances 1234 77`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupInstances(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupInterfacesCmd = &cobra.Command{
		Use:     "interfaces infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"list-interfaces"},
		Short:   "List the network interfaces of a group",
		Long: `List the network interfaces of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Examples:
  # List the interfaces of group 77
  metalcloud-cli container-instance-group interfaces my-infra 77

  # Using the alias
  metalcloud-cli cig list-interfaces 1234 77`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupInterfaces(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupInterfaceCmd = &cobra.Command{
		Use:     "interface infrastructure_id_or_label container_instance_group_id interface_id",
		Aliases: []string{"get-interface"},
		Short:   "Get one network interface of a group",
		Long: `Get detailed information about one network interface of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  interface_id                 The numeric ID of the interface

Examples:
  # Get interface 3 of group 77
  metalcloud-cli container-instance-group interface my-infra 77 3

  # Using the alias
  metalcloud-cli cig get-interface 1234 77 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupInterfaceGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	containerInstanceGroupApplyTypeCmd = &cobra.Command{
		Use:     "apply-type infrastructure_id_or_label container_instance_group_id container_type_id",
		Aliases: []string{"set-type"},
		Short:   "Apply a container type on a container instance group",
		Long: `Apply a container type on a container instance group.

The current group revision is fetched automatically and sent as the If-Match
header.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  container_type_id            The numeric ID of the container type to apply

Examples:
  # Apply container type 42 on group 77
  metalcloud-cli container-instance-group apply-type my-infra 77 42

  # Using the alias
  metalcloud-cli cig set-type 1234 77 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupApplyType(cmd.Context(), args[0], args[1], args[2])
		},
	}

	containerInstanceGroupNetworkCmd = &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"net"},
		Short:   "Manage the network connections of a container instance group",
		Long: `Manage the network connections of a container instance group.

Available commands:
  list        Show the network configuration of a group
  get         Get one network connection of a group
  connect     Connect a group to a logical network
  update      Update one network connection of a group
  disconnect  Remove one network connection of a group

Examples:
  metalcloud-cli container-instance-group network list my-infra 77
  metalcloud-cli cig net connect my-infra 77 5 trunk true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
	}

	containerInstanceGroupNetworkListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label container_instance_group_id",
		Aliases: []string{"ls"},
		Short:   "List the network connections of a container instance group",
		Long: `List all network connections of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Optional Flags:
  --configuration  Show the network endpoint group configuration instead of the
                   list of connections.

Examples:
  # List the network connections of group 77
  metalcloud-cli container-instance-group network list my-infra 77

  # Show the network endpoint group configuration
  metalcloud-cli cig net ls 1234 77 --configuration`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if containerInstanceGroupNetworkFlags.configuration {
				return container_instance.ContainerInstanceGroupNetworkConfiguration(cmd.Context(), args[0], args[1])
			}

			return container_instance.ContainerInstanceGroupNetworkList(cmd.Context(), args[0], args[1])
		},
	}

	containerInstanceGroupNetworkGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label container_instance_group_id connection_id",
		Aliases: []string{"show"},
		Short:   "Get one network connection of a container instance group",
		Long: `Get detailed information about one network connection of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  connection_id                The numeric ID of the network connection

Examples:
  # Get network connection 5 of group 77
  metalcloud-cli container-instance-group network get my-infra 77 5

  # Using the alias
  metalcloud-cli cig net show 1234 77 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupNetworkGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	containerInstanceGroupNetworkConnectCmd = &cobra.Command{
		Use:     "connect infrastructure_id_or_label container_instance_group_id logical_network_id access_mode tagged [redundancy]",
		Aliases: []string{"new", "add"},
		Short:   "Connect a container instance group to a logical network",
		Long: `Connect a container instance group to a logical network.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  logical_network_id           The ID of the logical network to connect to
  access_mode                  The network access mode (e.g. 'trunk', 'access')
  tagged                       Whether VLAN tagging is enabled (true/false)
  redundancy                   Optional. The redundancy mode of the connection

Examples:
  # Connect group 77 to logical network 5 in trunk mode
  metalcloud-cli container-instance-group network connect my-infra 77 5 trunk true

  # Connect with a redundancy mode
  metalcloud-cli cig net add 1234 77 5 trunk true active-backup`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.RangeArgs(5, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			redundancy := ""
			if len(args) == 6 {
				redundancy = args[5]
			}

			return container_instance.ContainerInstanceGroupNetworkConnect(cmd.Context(), args[0], args[1], args[2], args[3], args[4], redundancy)
		},
	}

	containerInstanceGroupNetworkUpdateCmd = &cobra.Command{
		Use:     "update infrastructure_id_or_label container_instance_group_id connection_id",
		Aliases: []string{"edit"},
		Short:   "Update one network connection of a container instance group",
		Long: `Update one network connection of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  connection_id                The numeric ID of the network connection

Required Flags:
  --access-mode string  The network connection access mode (e.g. 'trunk', 'access').
  --tagged string       Whether VLAN tagging is enabled (true/false).
  --redundancy string   The network connection redundancy mode.

Flag Dependencies:
  At least one of --access-mode, --tagged or --redundancy must be provided.

Examples:
  # Change the access mode of connection 5
  metalcloud-cli container-instance-group network update my-infra 77 5 --access-mode trunk

  # Enable VLAN tagging
  metalcloud-cli cig net edit 1234 77 5 --tagged true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupNetworkUpdate(cmd.Context(), args[0], args[1], args[2],
				containerInstanceGroupFlags.accessMode,
				containerInstanceGroupFlags.tagged,
				containerInstanceGroupFlags.redundancy)
		},
	}

	containerInstanceGroupNetworkDisconnectCmd = &cobra.Command{
		Use:     "disconnect infrastructure_id_or_label container_instance_group_id connection_id",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove one network connection of a container instance group",
		Long: `Remove one network connection of a container instance group.

This disconnects every instance of the group from the logical network.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  connection_id                The numeric ID of the network connection to remove

Examples:
  # Remove network connection 5 of group 77
  metalcloud-cli container-instance-group network disconnect my-infra 77 5

  # Using an alias
  metalcloud-cli cig net rm 1234 77 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_instance.ContainerInstanceGroupNetworkDisconnect(cmd.Context(), args[0], args[1], args[2])
		},
	}

	containerInstanceGroupNetworkFlags = struct {
		configuration bool
	}{}
)

func init() {
	// Container Instance management commands.
	rootCmd.AddCommand(containerInstanceCmd)

	containerInstanceCmd.AddCommand(containerInstanceListCmd)

	containerInstanceCmd.AddCommand(containerInstanceGetCmd)

	containerInstanceCmd.AddCommand(containerInstanceConfigExampleCmd)

	containerInstanceCmd.AddCommand(containerInstanceCreateCmd)
	containerInstanceCreateCmd.Flags().StringVar(&containerInstanceFlags.configSource, "config-source", "", "Source of the new container instance configuration. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceCreateCmd.Flags().StringVar(&containerInstanceFlags.containerTypeId, "container-type-id", "", "The container type of the new container instance.")
	containerInstanceCreateCmd.Flags().StringVar(&containerInstanceFlags.groupId, "group-id", "", "The container instance group of the new container instance.")
	containerInstanceCreateCmd.Flags().StringVar(&containerInstanceFlags.diskSizeGB, "disk-size-gb", "", "Disk size in GB of the new container instance.")
	containerInstanceCreateCmd.Flags().StringSliceVar(&containerInstanceFlags.tags, "tags", nil, "Tags of the new container instance.")
	containerInstanceCreateCmd.MarkFlagsOneRequired("config-source", "container-type-id")
	containerInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "container-type-id")
	containerInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "group-id")
	containerInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "disk-size-gb")
	containerInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "tags")
	containerInstanceCreateCmd.MarkFlagsRequiredTogether("container-type-id", "group-id")

	containerInstanceCmd.AddCommand(containerInstanceDeleteCmd)

	containerInstanceCmd.AddCommand(containerInstanceConfigCmd)

	containerInstanceCmd.AddCommand(containerInstanceUpdateConfigCmd)
	containerInstanceUpdateConfigCmd.Flags().StringVar(&containerInstanceFlags.configSource, "config-source", "", "Source of the container instance configuration updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	containerInstanceCmd.AddCommand(containerInstanceUpdateMetaCmd)
	containerInstanceUpdateMetaCmd.Flags().StringVar(&containerInstanceFlags.configSource, "config-source", "", "Source of the container instance metadata updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	containerInstanceCmd.AddCommand(containerInstanceCredentialsCmd)

	containerInstanceCmd.AddCommand(containerInstanceVariablesCmd)
	containerInstanceVariablesCmd.Flags().StringVar(&containerInstanceFlags.usage, "usage", "", "Restrict the variables to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).")

	containerInstanceCmd.AddCommand(containerInstanceOsInstallationDataCmd)
	containerInstanceOsInstallationDataCmd.Flags().StringVar(&containerInstanceFlags.usage, "usage", "", "Restrict the OS installation data to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).")
	containerInstanceOsInstallationDataCmd.Flags().BoolVar(&containerInstanceFlags.removeEmpty, "remove-empty", false, "Omit the empty entries from the response.")

	containerInstanceCmd.AddCommand(containerInstancePowerStatusCmd)
	containerInstanceCmd.AddCommand(containerInstanceStartCmd)
	containerInstanceCmd.AddCommand(containerInstanceShutdownCmd)
	containerInstanceCmd.AddCommand(containerInstanceRebootCmd)
	containerInstanceCmd.AddCommand(containerInstanceApplyTypeCmd)

	// Container Instance Group management commands.
	rootCmd.AddCommand(containerInstanceGroupCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupListCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupGetCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupConfigExampleCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupCreateCmd)
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.configSource, "config-source", "", "Source of the new container instance group configuration. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.containerTypeId, "container-type-id", "", "The container type of the group instances.")
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.osTemplateId, "os-template-id", "", "The OS template of the group instances.")
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.vmPoolId, "vm-pool-id", "", "The VM pool the group is provisioned on.")
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.diskSizeGB, "disk-size-gb", "", "Disk size in GB of each group instance.")
	containerInstanceGroupCreateCmd.Flags().StringVar(&containerInstanceGroupFlags.instanceCount, "instance-count", "", "Number of instances in the group.")
	containerInstanceGroupCreateCmd.Flags().StringSliceVar(&containerInstanceGroupFlags.tags, "tags", nil, "Tags of the new container instance group.")
	containerInstanceGroupCreateCmd.MarkFlagsOneRequired("config-source", "container-type-id")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "container-type-id")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "os-template-id")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "vm-pool-id")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "disk-size-gb")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "instance-count")
	containerInstanceGroupCreateCmd.MarkFlagsMutuallyExclusive("config-source", "tags")
	containerInstanceGroupCreateCmd.MarkFlagsRequiredTogether("container-type-id", "os-template-id", "vm-pool-id", "disk-size-gb")

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupDeleteCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupConfigCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupUpdateConfigCmd)
	containerInstanceGroupUpdateConfigCmd.Flags().StringVar(&containerInstanceGroupFlags.configSource, "config-source", "", "Source of the container instance group configuration updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceGroupUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupUpdateMetaCmd)
	containerInstanceGroupUpdateMetaCmd.Flags().StringVar(&containerInstanceGroupFlags.configSource, "config-source", "", "Source of the container instance group metadata updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerInstanceGroupUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupInstancesCmd)
	containerInstanceGroupCmd.AddCommand(containerInstanceGroupInterfacesCmd)
	containerInstanceGroupCmd.AddCommand(containerInstanceGroupInterfaceCmd)
	containerInstanceGroupCmd.AddCommand(containerInstanceGroupApplyTypeCmd)

	containerInstanceGroupCmd.AddCommand(containerInstanceGroupNetworkCmd)

	containerInstanceGroupNetworkCmd.AddCommand(containerInstanceGroupNetworkListCmd)
	containerInstanceGroupNetworkListCmd.Flags().BoolVar(&containerInstanceGroupNetworkFlags.configuration, "configuration", false, "Show the network endpoint group configuration instead of the list of connections.")

	containerInstanceGroupNetworkCmd.AddCommand(containerInstanceGroupNetworkGetCmd)

	containerInstanceGroupNetworkCmd.AddCommand(containerInstanceGroupNetworkConnectCmd)

	containerInstanceGroupNetworkCmd.AddCommand(containerInstanceGroupNetworkUpdateCmd)
	containerInstanceGroupNetworkUpdateCmd.Flags().StringVar(&containerInstanceGroupFlags.accessMode, "access-mode", "", "Network connection access mode.")
	containerInstanceGroupNetworkUpdateCmd.Flags().StringVar(&containerInstanceGroupFlags.tagged, "tagged", "", "Network connection VLAN tagging (true/false).")
	containerInstanceGroupNetworkUpdateCmd.Flags().StringVar(&containerInstanceGroupFlags.redundancy, "redundancy", "", "Network connection redundancy mode.")
	containerInstanceGroupNetworkUpdateCmd.MarkFlagsOneRequired("access-mode", "tagged", "redundancy")

	containerInstanceGroupNetworkCmd.AddCommand(containerInstanceGroupNetworkDisconnectCmd)
}
