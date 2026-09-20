package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/container"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	containerFlags = struct {
		configSource              string
		filterId                  []string
		filterSiteId              []string
		filterName                []string
		filterHost                []string
		filterTypeId              []string
		filterPoolId              []string
		filterAdministrationState []string
		filterInfrastructureId    []string
	}{}

	containerCmd = &cobra.Command{
		Use:     "container [command]",
		Aliases: []string{"containers"},
		Short:   "Manage provisioned containers",
		Long: `Manage provisioned containers.

Containers are the runtime objects allocated for container instances. They are
created by the platform when an infrastructure is deployed; this command group
inspects them and controls their power state.

Available commands:
  list                 List all containers
  get                  Get details of a container
  update               Update container comments and tags
  power-status         Show the current power state of a container
  start                Power on a container
  shutdown             Power off a container
  reboot               Reboot a container
  remote-console-info  Show remote console information for a container

Examples:
  metalcloud-cli container list
  metalcloud-cli containers get 100
  metalcloud-cli container power-status 100`,
	}

	containerListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all containers",
		Long: `List all containers.

Retrieves every container visible to the current user. The output can be
narrowed with the filter flags below; each filter accepts multiple values.

Optional Flags:
  --filter-id strings                   Filter by container ID
  --filter-site-id strings              Filter by site ID
  --filter-name strings                 Filter by container name
  --filter-host strings                 Filter by host name
  --filter-type-id strings              Filter by container type ID
  --filter-pool-id strings              Filter by VM pool ID
  --filter-administration-state strings Filter by administration state
  --filter-infrastructure-id strings    Filter by infrastructure ID

Examples:
  # List all containers
  metalcloud-cli container list

  # List the containers of one infrastructure
  metalcloud-cli containers ls --filter-infrastructure-id 1234

  # List the containers of a given type
  metalcloud-cli container ls --filter-type-id 42`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerList(cmd.Context(), container.ContainerFilters{
				Id:                  containerFlags.filterId,
				SiteId:              containerFlags.filterSiteId,
				Name:                containerFlags.filterName,
				Host:                containerFlags.filterHost,
				TypeId:              containerFlags.filterTypeId,
				PoolId:              containerFlags.filterPoolId,
				AdministrationState: containerFlags.filterAdministrationState,
				InfrastructureId:    containerFlags.filterInfrastructureId,
			})
		},
	}

	containerGetCmd = &cobra.Command{
		Use:     "get container_id",
		Aliases: []string{"show"},
		Short:   "Get container details",
		Long: `Get detailed information about a container.

Required Arguments:
  container_id  The numeric ID of the container

Examples:
  # Get the details of container 100
  metalcloud-cli container get 100

  # Using the alias
  metalcloud-cli containers show 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerGet(cmd.Context(), args[0])
		},
	}

	containerUpdateCmd = &cobra.Command{
		Use:     "update container_id",
		Aliases: []string{"edit"},
		Short:   "Update a container",
		Long: `Update the comments and tags of a container.

Required Arguments:
  container_id  The numeric ID of the container to update

Required Flags:
  --config-source string  Source of the container updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update container 100 from a file
  metalcloud-cli container update 100 --config-source updates.json

  # Update the tags of container 100 from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli containers edit 100 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(containerFlags.configSource)
			if err != nil {
				return err
			}

			return container.ContainerUpdate(cmd.Context(), args[0], config)
		},
	}

	containerPowerStatusCmd = &cobra.Command{
		Use:     "power-status container_id",
		Aliases: []string{"power-state"},
		Short:   "Get the power status of a container",
		Long: `Get the current power status of a container.

Required Arguments:
  container_id  The numeric ID of the container

Examples:
  # Show the power status of container 100
  metalcloud-cli container power-status 100

  # Using the alias
  metalcloud-cli containers power-state 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerPowerStatus(cmd.Context(), args[0])
		},
	}

	containerStartCmd = &cobra.Command{
		Use:     "start container_id",
		Aliases: []string{"power-on"},
		Short:   "Power on a container",
		Long: `Power on a container.

Required Arguments:
  container_id  The numeric ID of the container to power on

Examples:
  # Power on container 100
  metalcloud-cli container start 100

  # Using the alias
  metalcloud-cli containers power-on 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerPowerControl(cmd.Context(), args[0], "start")
		},
	}

	containerShutdownCmd = &cobra.Command{
		Use:     "shutdown container_id",
		Aliases: []string{"power-off", "stop"},
		Short:   "Power off a container",
		Long: `Power off a container.

Required Arguments:
  container_id  The numeric ID of the container to power off

Examples:
  # Power off container 100
  metalcloud-cli container shutdown 100

  # Using an alias
  metalcloud-cli containers stop 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerPowerControl(cmd.Context(), args[0], "shutdown")
		},
	}

	containerRebootCmd = &cobra.Command{
		Use:     "reboot container_id",
		Aliases: []string{"restart"},
		Short:   "Reboot a container",
		Long: `Reboot a container.

Required Arguments:
  container_id  The numeric ID of the container to reboot

Examples:
  # Reboot container 100
  metalcloud-cli container reboot 100

  # Using the alias
  metalcloud-cli containers restart 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerPowerControl(cmd.Context(), args[0], "reboot")
		},
	}

	containerRemoteConsoleInfoCmd = &cobra.Command{
		Use:     "remote-console-info container_id",
		Aliases: []string{"console-info"},
		Short:   "Get remote console information for a container",
		Long: `Get remote console information for a container.

Shows how many remote console connections are currently active for the container.

Required Arguments:
  container_id  The numeric ID of the container

Examples:
  # Show the remote console info of container 100
  metalcloud-cli container remote-console-info 100

  # Using the alias
  metalcloud-cli containers console-info 100`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_CONTAINERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return container.ContainerRemoteConsoleInfo(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(containerCmd)

	containerCmd.AddCommand(containerListCmd)
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterId, "filter-id", nil, "Filter by container ID.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterSiteId, "filter-site-id", nil, "Filter by site ID.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterName, "filter-name", nil, "Filter by container name.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterHost, "filter-host", nil, "Filter by host name.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterTypeId, "filter-type-id", nil, "Filter by container type ID.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterPoolId, "filter-pool-id", nil, "Filter by VM pool ID.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterAdministrationState, "filter-administration-state", nil, "Filter by administration state.")
	containerListCmd.Flags().StringSliceVar(&containerFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")

	containerCmd.AddCommand(containerGetCmd)

	containerCmd.AddCommand(containerUpdateCmd)
	containerUpdateCmd.Flags().StringVar(&containerFlags.configSource, "config-source", "", "Source of the container updates. Can be 'pipe' or path to a JSON/YAML file.")
	containerUpdateCmd.MarkFlagsOneRequired("config-source")

	containerCmd.AddCommand(containerPowerStatusCmd)
	containerCmd.AddCommand(containerStartCmd)
	containerCmd.AddCommand(containerShutdownCmd)
	containerCmd.AddCommand(containerRebootCmd)
	containerCmd.AddCommand(containerRemoteConsoleInfoCmd)
}
