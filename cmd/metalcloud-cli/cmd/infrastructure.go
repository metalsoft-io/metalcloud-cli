package cmd

import (
	"time"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/spf13/cobra"
)

var (
	infrastructureFlags = struct {
		showAll             bool
		showOrdered         bool
		showDeleted         bool
		customVariables     string
		allowDataLoss       bool
		attemptSoftShutdown bool
		attemptHardShutdown bool
		softShutdownTimeout int
		forceShutdown       bool
	}{}

	infrastructureUtilFlags = struct {
		userId            int
		startTime         time.Time
		endTime           time.Time
		siteIds           []int
		infrastructureIds []int
		showAll           bool
		showInstances     bool
		showSubnets       bool
		showDrives        bool
	}{}

	infrastructureCmd = &cobra.Command{
		Use:     "infrastructure [command]",
		Aliases: []string{"infra"},
		Short:   "Manage infrastructure resources and configurations",
		Long: `Manage infrastructure resources including creation, deployment, monitoring, and user access control.

Infrastructure represents a collection of compute instances, storage drives, and network resources 
that can be managed as a single unit. Each infrastructure belongs to a specific site and can be 
deployed, updated, or deleted as needed.

Available Commands:
  list         List all infrastructures with filtering options
  get          Show detailed information about a specific infrastructure
  create       Create a new infrastructure in a site
  update       Update infrastructure configuration and metadata
  delete       Delete an infrastructure and all its resources
  deploy       Deploy infrastructure changes to physical resources
  revert       Revert infrastructure to the last deployed state
  users        Manage user access to infrastructures
  statistics   View infrastructure deployment and job statistics
  utilization  Generate resource utilization reports

Use "metalcloud-cli infrastructure [command] --help" for more information about a specific command.`,
	}

	infrastructureListCmd = &cobra.Command{
		Use:     "list [flags...]",
		Aliases: []string{"ls"},
		Short:   "List infrastructures with optional filtering",
		Long: `List all infrastructures with various filtering options to control visibility and scope.

By default, this command shows only active infrastructures owned by the current user. 
Use the filtering flags to customize the output based on your needs.

Flags:
  --show-all      Show infrastructures from all users (requires admin privileges)
  --show-ordered  Include ordered (created but not yet deployed) infrastructures
  --show-deleted  Include deleted infrastructures in the output

Examples:
  # List your active infrastructures
  metalcloud-cli infrastructure list

  # List all infrastructures including ordered and deleted ones
  metalcloud-cli infrastructure list --show-all --show-ordered --show-deleted

  # List only your ordered infrastructures
  metalcloud-cli infrastructure list --show-ordered

  # List all infrastructures from all users (admin only)
  metalcloud-cli infrastructure list --show-all`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureList(cmd.Context(),
				infrastructureFlags.showAll,
				infrastructureFlags.showOrdered,
				infrastructureFlags.showDeleted)
		},
	}

	infrastructureGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label",
		Aliases: []string{"show"},
		Short:   "Show detailed information about a specific infrastructure",
		Long: `Display comprehensive details about a specific infrastructure including its configuration, 
status, resources, and deployment information.

The infrastructure can be identified by either its numeric ID or its label (name).

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Get infrastructure details by ID
  metalcloud-cli infrastructure get 123

  # Get infrastructure details by label
  metalcloud-cli infrastructure get "my-web-cluster"

  # Show infrastructure details (alias)
  metalcloud-cli infrastructure show production-env`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGet(cmd.Context(), args[0])
		},
	}

	infrastructureCreateCmd = &cobra.Command{
		Use:     "create site_id label",
		Aliases: []string{"new"},
		Short:   "Create a new infrastructure in a specific site",
		Long: `Create a new infrastructure with the specified label in the given site.

The infrastructure will be created in an "ordered" state and must be deployed to provision
actual resources. After creation, you can add compute instances, drives, and networks
to the infrastructure before deploying it.

Arguments:
  site_id  The numeric ID of the site where the infrastructure will be created
  label    A unique label (name) for the infrastructure

Examples:
  # Create a new infrastructure in site 1
  metalcloud-cli infrastructure create 1 "web-cluster"

  # Create infrastructure with a descriptive name
  metalcloud-cli infrastructure create 2 "production-database-cluster"

  # Using the alias
  metalcloud-cli infrastructure new 1 "test-environment"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureCreate(cmd.Context(), args[0], args[1])
		},
	}

	infrastructureUpdateCmd = &cobra.Command{
		Use:     "update infrastructure_id_or_label [new_label]",
		Aliases: []string{"edit"},
		Short:   "Update infrastructure configuration and metadata",
		Long: `Update various properties of an infrastructure including its label and custom variables.

This command allows you to modify infrastructure metadata without affecting the deployed
resources. Changes to the infrastructure configuration require a subsequent deploy to
take effect on the actual infrastructure.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure to update
  new_label                   Optional new label for the infrastructure

Flags:
  --custom-variables  JSON string containing custom variables to set on the infrastructure

Examples:
  # Update infrastructure label
  metalcloud-cli infrastructure update 123 "new-cluster-name"

  # Update only custom variables
  metalcloud-cli infrastructure update web-cluster --custom-variables '{"env":"production","version":"1.2.3"}'

  # Update both label and custom variables
  metalcloud-cli infrastructure update 123 "prod-cluster" --custom-variables '{"tier":"production"}'

  # Using the alias
  metalcloud-cli infrastructure edit my-infrastructure new-name`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label := ""
			if len(args) > 1 {
				label = args[1]
			}

			return infrastructure.InfrastructureUpdate(cmd.Context(), args[0], label,
				infrastructureFlags.customVariables)
		},
	}

	infrastructureDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label",
		Aliases: []string{"rm"},
		Short:   "Delete an infrastructure and all its resources",
		Long: `Delete an infrastructure and all associated resources including compute instances, 
drives, and network configurations.

WARNING: This operation is irreversible and will permanently destroy all data 
and resources associated with the infrastructure. Make sure to backup any 
important data before proceeding.

The infrastructure must be in a non-deployed state before it can be deleted.
If the infrastructure is currently deployed, you must first revert it to
remove all active resources.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure to delete

Examples:
  # Delete infrastructure by ID
  metalcloud-cli infrastructure delete 123

  # Delete infrastructure by label
  metalcloud-cli infrastructure delete "test-cluster"

  # Using the alias
  metalcloud-cli infrastructure rm old-infrastructure`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureDelete(cmd.Context(), args[0])
		},
	}

	infrastructureDeployCmd = &cobra.Command{
		Use:     "deploy infrastructure_id_or_label",
		Aliases: []string{"apply"},
		Short:   "Deploy infrastructure changes to physical resources",
		Long: `Deploy an infrastructure configuration to provision and configure physical resources.

This command applies all pending changes to the infrastructure including creating, 
modifying, or destroying compute instances, drives, and network configurations.
The deployment process may take several minutes depending on the complexity of changes.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure to deploy

Flags:
  --allow-data-loss           Allow deployment even if data loss is expected (default: false)
  --attempt-soft-shutdown     Attempt ACPI power off before deployment (default: true)
  --attempt-hard-shutdown     Force hard power off after timeout if soft shutdown fails (default: true)
  --soft-shutdown-timeout     Timeout in seconds for soft shutdown before forcing hard shutdown (default: 180)
  --force-shutdown           Force immediate shutdown of all servers (default: false)

Examples:
  # Deploy infrastructure with default settings
  metalcloud-cli infrastructure deploy 123

  # Deploy with custom shutdown settings
  metalcloud-cli infrastructure deploy web-cluster --soft-shutdown-timeout 300

  # Deploy allowing data loss (dangerous)
  metalcloud-cli infrastructure deploy test-env --allow-data-loss

  # Force deployment with immediate shutdown
  metalcloud-cli infrastructure deploy emergency-fix --force-shutdown

  # Using the alias
  metalcloud-cli infrastructure apply production-cluster`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureDeploy(cmd.Context(), args[0],
				infrastructureFlags.allowDataLoss,
				infrastructureFlags.attemptSoftShutdown,
				infrastructureFlags.attemptHardShutdown,
				infrastructureFlags.softShutdownTimeout,
				infrastructureFlags.forceShutdown)
		},
	}

	infrastructureRevertCmd = &cobra.Command{
		Use:     "revert infrastructure_id_or_label",
		Aliases: []string{"undo"},
		Short:   "Revert infrastructure to the last deployed state",
		Long: `Revert an infrastructure to its last successfully deployed state, discarding all 
pending changes.

This command is useful when you want to undo configuration changes that haven't been 
deployed yet, or when you need to roll back to a known working state. All pending 
modifications to compute instances, drives, and networks will be reverted.

Note: This operation only affects the infrastructure configuration and does not 
modify any deployed physical resources. To apply the reverted configuration to 
physical resources, you need to deploy the infrastructure after reverting.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure to revert

Examples:
  # Revert infrastructure by ID
  metalcloud-cli infrastructure revert 123

  # Revert infrastructure by label
  metalcloud-cli infrastructure revert "web-cluster"

  # Using the alias
  metalcloud-cli infrastructure undo production-env`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureRevert(cmd.Context(), args[0])
		},
	}

	infrastructureGetUsersCmd = &cobra.Command{
		Use:     "users infrastructure_id_or_label",
		Aliases: []string{"list-users", "get-users"},
		Short:   "List users with access to an infrastructure",
		Long: `Display all users who have access to a specific infrastructure, including their 
access levels and contact information.

This command shows the current user permissions for an infrastructure, which is useful 
for managing access control and understanding who can modify or view the infrastructure.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # List users for infrastructure by ID
  metalcloud-cli infrastructure users 123

  # List users for infrastructure by label
  metalcloud-cli infrastructure users "web-cluster"

  # Using aliases
  metalcloud-cli infrastructure list-users production-env
  metalcloud-cli infrastructure get-users my-infrastructure`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUsers(cmd.Context(), args[0])
		},
	}

	infrastructureAddUserCmd = &cobra.Command{
		Use:   "add-user infrastructure_id_or_label user_email [create_if_not_exists]",
		Short: "Add a user to an infrastructure with access permissions",
		Long: `Grant a user access to a specific infrastructure by their email address.

This command adds the specified user to the infrastructure's access control list,
allowing them to view and potentially modify the infrastructure based on their
permissions level.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure
  user_email                  Email address of the user to add
  create_if_not_exists       Optional: "true" to create user if they don't exist (default: "false")

Examples:
  # Add existing user to infrastructure
  metalcloud-cli infrastructure add-user 123 "user@example.com"

  # Add user by infrastructure label
  metalcloud-cli infrastructure add-user "web-cluster" "admin@company.com"

  # Add user and create if they don't exist
  metalcloud-cli infrastructure add-user my-infra "newuser@company.com" true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			createMissing := "false"
			if len(args) == 3 {
				createMissing = args[2]
			}
			return infrastructure.InfrastructureAddUser(cmd.Context(), args[0], args[1], createMissing)
		},
	}

	infrastructureRemoveUserCmd = &cobra.Command{
		Use:     "remove-user infrastructure_id_or_label user_id",
		Aliases: []string{"delete-user"},
		Short:   "Remove a user's access from an infrastructure",
		Long: `Remove a user's access permissions from a specific infrastructure.

This command revokes the specified user's access to the infrastructure, preventing them 
from viewing or modifying it. The user is identified by their numeric user ID.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure
  user_id                     The numeric ID of the user to remove

Examples:
  # Remove user by infrastructure ID and user ID
  metalcloud-cli infrastructure remove-user 123 456

  # Remove user by infrastructure label
  metalcloud-cli infrastructure remove-user "web-cluster" 789

  # Using the alias
  metalcloud-cli infrastructure delete-user my-infrastructure 101`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureRemoveUser(cmd.Context(), args[0], args[1])
		},
	}

	infrastructureGetUserLimitsCmd = &cobra.Command{
		Use:     "user-limits infrastructure_id_or_label",
		Aliases: []string{"get-user-limits"},
		Short:   "Display resource limits for an infrastructure",
		Long: `Show the resource limits configured for a specific infrastructure including compute nodes,
drives, and infrastructure count limits.

This information helps understand the maximum resources that can be provisioned within
the infrastructure and plan capacity accordingly.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Get user limits for infrastructure by ID
  metalcloud-cli infrastructure user-limits 123

  # Get user limits for infrastructure by label
  metalcloud-cli infrastructure user-limits "web-cluster"

  # Using the alias
  metalcloud-cli infrastructure get-user-limits production-env`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUserLimits(cmd.Context(), args[0])
		},
	}

	infrastructureGetStatisticsCmd = &cobra.Command{
		Use:     "statistics infrastructure_id_or_label",
		Aliases: []string{"stats", "get-statistics"},
		Short:   "Get deployment statistics for an infrastructure",
		Long: `Display deployment and job execution statistics for a specific infrastructure.

This command shows information about deployment groups, job completion rates, error counts,
and timing information which helps monitor infrastructure deployment health and performance.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Get statistics for infrastructure by ID
  metalcloud-cli infrastructure statistics 123

  # Get statistics for infrastructure by label
  metalcloud-cli infrastructure statistics "web-cluster"

  # Using aliases
  metalcloud-cli infrastructure stats production-env
  metalcloud-cli infrastructure get-statistics my-infrastructure`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetStatistics(cmd.Context(), args[0])
		},
	}

	infrastructureGetConfigInfoCmd = &cobra.Command{
		Use:     "config-info infrastructure_id_or_label",
		Aliases: []string{"get-config-info"},
		Short:   "Get configuration information for an infrastructure",
		Long: `Display detailed configuration information for a specific infrastructure including
deployment status, configuration revision, and update timestamps.

This command provides insights into the current configuration state of the infrastructure,
helping to understand deployment progress and configuration changes.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Get config info for infrastructure by ID
  metalcloud-cli infrastructure config-info 123

  # Get config info for infrastructure by label
  metalcloud-cli infrastructure config-info "web-cluster"

  # Using the alias
  metalcloud-cli infrastructure get-config-info production-env`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetConfigInfo(cmd.Context(), args[0])
		},
	}

	infrastructureGetAllStatisticsCmd = &cobra.Command{
		Use:     "all-statistics",
		Aliases: []string{"all-stats", "get-all-statistics"},
		Short:   "Get deployment statistics for all infrastructures",
		Long: `Display aggregated deployment and job execution statistics for all infrastructures.

This command provides a global overview of infrastructure deployment health, including 
total infrastructure counts, active deployments, error rates, and ongoing operations
across the entire system.

The statistics include:
- Total infrastructure count and service status breakdown
- Number of ongoing deployments and their status
- Error counts and retry statistics for failed deployments

Examples:
  # Get statistics for all infrastructures
  metalcloud-cli infrastructure all-statistics

  # Using aliases
  metalcloud-cli infrastructure all-stats
  metalcloud-cli infrastructure get-all-statistics`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetAllStatistics(cmd.Context())
		},
	}

	infrastructureUtilizationCmd = &cobra.Command{
		Use:     "utilization",
		Aliases: []string{"get-utilization"},
		Short:   "Get resource utilization report for infrastructures",
		Long: `Get detailed utilization report for infrastructure resources within a specified time range. 
The report provides insights into resource usage patterns and capacity planning for infrastructures.

Required flags:
  --user-id       ID of the user to include in the report
  --start-time    Start time for the report (RFC3339 or date format)
  --end-time      End time for the report (RFC3339 or date format)

Optional flags:
  --site-id            Site IDs to include in the report (can be specified multiple times)
  --infrastructure-id  Infrastructure IDs to include in the report (can be specified multiple times)
  --show-all           Show all utilizations
  --show-instances     Show instance utilizations
  --show-drives        Show drive utilizations
  --show-subnets       Show subnet utilizations

Examples:
  # Get utilization report for user 123 for the last 7 days
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08

  # Get utilization for specific sites and infrastructures
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01T00:00:00Z --end-time 2025-08-08T23:59:59Z --site-id 1 --site-id 2 --infrastructure-id 100 --infrastructure-id 101

  # Get utilization for all infrastructures of a user in specific sites
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08 --site-id 1 --site-id 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_UTILIZATION_REPORTS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUtilization(
				cmd.Context(),
				infrastructureUtilFlags.userId,
				infrastructureUtilFlags.startTime,
				infrastructureUtilFlags.endTime,
				infrastructureUtilFlags.siteIds,
				infrastructureUtilFlags.infrastructureIds,
				infrastructureUtilFlags.showAll || infrastructureUtilFlags.showInstances,
				infrastructureUtilFlags.showAll || infrastructureUtilFlags.showDrives,
				infrastructureUtilFlags.showAll || infrastructureUtilFlags.showSubnets)
		},
	}
)

func init() {
	rootCmd.AddCommand(infrastructureCmd)

	infrastructureCmd.AddCommand(infrastructureListCmd)
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showAll, "show-all", false, "If set will return all infrastructures.")
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showOrdered, "show-ordered", false, "If set will also return ordered (created but not deployed) infrastructures.")
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showDeleted, "show-deleted", false, "If set will also return deleted infrastructures.")

	infrastructureCmd.AddCommand(infrastructureGetCmd)

	infrastructureCmd.AddCommand(infrastructureCreateCmd)

	infrastructureCmd.AddCommand(infrastructureUpdateCmd)
	infrastructureUpdateCmd.Flags().StringVar(&infrastructureFlags.customVariables, "custom-variables", "", "Set of infrastructure custom variables.")

	infrastructureCmd.AddCommand(infrastructureDeleteCmd)

	infrastructureCmd.AddCommand(infrastructureDeployCmd)
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.allowDataLoss, "allow-data-loss", false, "If set, deploy will not throw error if data loss is expected.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.attemptSoftShutdown, "attempt-soft-shutdown", true, "If set, attempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.attemptHardShutdown, "attempt-hard-shutdown", true, "If set, force a hard power off after timeout expired and the server is not powered off")
	infrastructureDeployCmd.Flags().IntVar(&infrastructureFlags.softShutdownTimeout, "soft-shutdown-timeout", 180, "Timeout to wait for soft shutdown before forcing hard shutdown.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.forceShutdown, "force-shutdown", false, "If set, deploy will force shutdown of all servers in the infrastructure.")

	infrastructureCmd.AddCommand(infrastructureRevertCmd)

	infrastructureCmd.AddCommand(infrastructureGetUsersCmd)

	infrastructureCmd.AddCommand(infrastructureAddUserCmd)

	infrastructureCmd.AddCommand(infrastructureRemoveUserCmd)

	infrastructureCmd.AddCommand(infrastructureGetUserLimitsCmd)

	infrastructureCmd.AddCommand(infrastructureGetStatisticsCmd)

	infrastructureCmd.AddCommand(infrastructureGetConfigInfoCmd)

	infrastructureCmd.AddCommand(infrastructureGetAllStatisticsCmd)

	infrastructureCmd.AddCommand(infrastructureUtilizationCmd)
	infrastructureUtilizationCmd.Flags().IntVar(&infrastructureUtilFlags.userId, "user-id", 0, "ID of the user to include in the report.")
	infrastructureUtilizationCmd.Flags().TimeVar(&infrastructureUtilFlags.startTime, "start-time", time.Now().Add(-time.Duration(time.Now().Day())), []string{time.RFC3339, time.DateOnly}, "Start time for the report.")
	infrastructureUtilizationCmd.Flags().TimeVar(&infrastructureUtilFlags.endTime, "end-time", time.Now(), []string{time.RFC3339, time.DateOnly}, "End time for the report.")
	infrastructureUtilizationCmd.Flags().IntSliceVar(&infrastructureUtilFlags.siteIds, "site-id", []int{}, "Site IDs to include in the report.")
	infrastructureUtilizationCmd.Flags().IntSliceVar(&infrastructureUtilFlags.infrastructureIds, "infrastructure-id", []int{}, "Infrastructure IDs to include in the report.")
	infrastructureUtilizationCmd.Flags().BoolVar(&infrastructureUtilFlags.showAll, "show-all", false, "If set, will display all utilizations.")
	infrastructureUtilizationCmd.Flags().BoolVar(&infrastructureUtilFlags.showInstances, "show-instances", false, "If set, will display instance utilizations.")
	infrastructureUtilizationCmd.Flags().BoolVar(&infrastructureUtilFlags.showDrives, "show-drives", false, "If set, will display drive utilizations.")
	infrastructureUtilizationCmd.Flags().BoolVar(&infrastructureUtilFlags.showSubnets, "show-subnets", false, "If set, will display subnet utilizations.")
	infrastructureUtilizationCmd.MarkFlagRequired("user-id")
	infrastructureUtilizationCmd.MarkFlagRequired("start-time")
	infrastructureUtilizationCmd.MarkFlagRequired("end-time")
}
