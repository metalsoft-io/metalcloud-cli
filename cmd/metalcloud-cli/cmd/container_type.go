package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/container_type"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	containerTypeFlags = struct {
		configSource      string
		filterId          []string
		filterLabel       []string
		filterName        []string
		filterDisplayName []string
	}{}

	containerTypeCmd = &cobra.Command{
		Use:     "container-type [command]",
		Aliases: []string{"ct", "container-types"},
		Short:   "Manage container types",
		Long: `Manage container types.

Container types describe the hardware envelope (CPU cores, RAM, GPUs) that a
container instance receives when it is provisioned. They are global objects and
are referenced by container instances and container instance groups.

Available commands:
  list            List all container types
  get             Get details of a container type
  config-example  Print an example container type configuration
  create          Create a new container type
  update          Update an existing container type
  delete          Delete a container type
  containers      List the containers provisioned with a container type

Examples:
  metalcloud-cli container-type list
  metalcloud-cli ct get 42
  metalcloud-cli ct containers 42`,
	}

	containerTypeListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all container types",
		Long: `List all container types.

Retrieves every container type visible to the current user. The output can be
narrowed with the filter flags below; each filter accepts multiple values.

Optional Flags:
  --filter-id strings            Filter by container type ID
  --filter-label strings         Filter by container type label
  --filter-name strings          Filter by container type name
  --filter-display-name strings  Filter by container type display name

Examples:
  # List all container types
  metalcloud-cli container-type list

  # List only the container types with a given label
  metalcloud-cli ct ls --filter-label small-container

  # List several container types by ID
  metalcloud-cli ct ls --filter-id 10 --filter-id 11`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_type.ContainerTypeList(cmd.Context(), container_type.ContainerTypeFilters{
				Id:          containerTypeFlags.filterId,
				Label:       containerTypeFlags.filterLabel,
				Name:        containerTypeFlags.filterName,
				DisplayName: containerTypeFlags.filterDisplayName,
			})
		},
	}

	containerTypeGetCmd = &cobra.Command{
		Use:     "get container_type_id",
		Aliases: []string{"show"},
		Short:   "Get container type details",
		Long: `Get detailed information about a container type.

Required Arguments:
  container_type_id  The numeric ID of the container type

Examples:
  # Get the details of container type 42
  metalcloud-cli container-type get 42

  # Using the alias
  metalcloud-cli ct show 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_type.ContainerTypeGet(cmd.Context(), args[0])
		},
	}

	containerTypeConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print a container type configuration example",
		Long: `Print a container type configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli container-type config-example

  # Save the example to a file
  metalcloud-cli ct config-example > container-type.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_type.ContainerTypeConfigExample(cmd.Context())
		},
	}

	containerTypeCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new container type",
		Long: `Create a new container type.

The new container type is described by a JSON or YAML document. Use the
config-example command to obtain a template.

Required Flags:
  --config-source string  Source of the container type configuration.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Create a container type from a file
  metalcloud-cli container-type create --config-source container-type.json

  # Create a container type from stdin
  cat container-type.json | metalcloud-cli ct new --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerTypeFlags.configSource)
			if err != nil {
				return err
			}

			return container_type.ContainerTypeCreate(cmd.Context(), config)
		},
	}

	containerTypeUpdateCmd = &cobra.Command{
		Use:     "update container_type_id",
		Aliases: []string{"edit"},
		Short:   "Update a container type",
		Long: `Update an existing container type.

Only the display name, label, tags and the experimental / unmanaged-only flags
can be changed; CPU and RAM sizing is immutable.

Required Arguments:
  container_type_id  The numeric ID of the container type to update

Required Flags:
  --config-source string  Source of the container type updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update container type 42 from a file
  metalcloud-cli container-type update 42 --config-source updates.json

  # Update container type 42 from stdin
  echo '{"label":"new-label"}' | metalcloud-cli ct edit 42 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerTypeFlags.configSource)
			if err != nil {
				return err
			}

			return container_type.ContainerTypeUpdate(cmd.Context(), args[0], config)
		},
	}

	containerTypeDeleteCmd = &cobra.Command{
		Use:     "delete container_type_id",
		Aliases: []string{"rm"},
		Short:   "Delete a container type",
		Long: `Delete a container type.

The container type must no longer be referenced by any container instance.

Required Arguments:
  container_type_id  The numeric ID of the container type to delete

Examples:
  # Delete container type 42
  metalcloud-cli container-type delete 42

  # Using the alias
  metalcloud-cli ct rm 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_type.ContainerTypeDelete(cmd.Context(), args[0])
		},
	}

	containerTypeContainersCmd = &cobra.Command{
		Use:     "containers container_type_id",
		Aliases: []string{"list-containers"},
		Short:   "List the containers of a container type",
		Long: `List every container that has been provisioned with a container type.

Required Arguments:
  container_type_id  The numeric ID of the container type

Examples:
  # List the containers of container type 42
  metalcloud-cli container-type containers 42

  # Using the alias
  metalcloud-cli ct list-containers 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINER_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container_type.ContainerTypeContainers(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(containerTypeCmd)

	containerTypeCmd.AddCommand(containerTypeListCmd)
	containerTypeListCmd.Flags().StringSliceVar(&containerTypeFlags.filterId, "filter-id", nil, "Filter by container type ID.")
	containerTypeListCmd.Flags().StringSliceVar(&containerTypeFlags.filterLabel, "filter-label", nil, "Filter by container type label.")
	containerTypeListCmd.Flags().StringSliceVar(&containerTypeFlags.filterName, "filter-name", nil, "Filter by container type name.")
	containerTypeListCmd.Flags().StringSliceVar(&containerTypeFlags.filterDisplayName, "filter-display-name", nil, "Filter by container type display name.")

	containerTypeCmd.AddCommand(containerTypeGetCmd)

	containerTypeCmd.AddCommand(containerTypeConfigExampleCmd)

	containerTypeCmd.AddCommand(containerTypeCreateCmd)
	containerTypeCreateCmd.Flags().StringVar(&containerTypeFlags.configSource, "config-source", "", "Source of the new container type configuration. Can be 'pipe' or path to a JSON/YAML file.")
	containerTypeCreateCmd.MarkFlagsOneRequired("config-source")

	containerTypeCmd.AddCommand(containerTypeUpdateCmd)
	containerTypeUpdateCmd.Flags().StringVar(&containerTypeFlags.configSource, "config-source", "", "Source of the container type updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerTypeUpdateCmd.MarkFlagsOneRequired("config-source")

	containerTypeCmd.AddCommand(containerTypeDeleteCmd)

	containerTypeCmd.AddCommand(containerTypeContainersCmd)
}
