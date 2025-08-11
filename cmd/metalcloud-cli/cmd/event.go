package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/event"
	"github.com/spf13/cobra"
)

var (
	eventFlags = struct {
		filterId               []string
		filterType             []string
		filterSeverity         []string
		filterVisibility       []string
		filterInfrastructureId []string
		filterUserId           []string
		filterServerId         []string
		filterJobId            []string
		filterSiteId           []string
		sortBy                 []string
		page                   int
		limit                  int
		search                 string
		searchBy               []string
	}{}

	eventCmd = &cobra.Command{
		Use:     "event [command]",
		Aliases: []string{"evt"},
		Short:   "Manage and monitor system events",
		Long: `Manage and monitor system events in MetalCloud.

Events represent important system activities such as infrastructure changes,
server deployments, job executions, and user actions. This command provides
tools to list, filter, search, and retrieve detailed information about events.

Available Commands:
  list    List events with filtering and search capabilities
  get     Retrieve detailed information about a specific event

Examples:
  # List all events
  metalcloud event list

  # Get details of a specific event
  metalcloud event get 12345

  # List events with filters
  metalcloud event list --filter-type deployment --filter-severity error`,
	}

	eventListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List and filter system events",
		Long: `List and filter system events with comprehensive filtering and search capabilities.

This command allows you to view events from across your MetalCloud infrastructure
with support for multiple filtering options, sorting, pagination, and text search.

FILTERING OPTIONS:
  --filter-id                  Filter by specific event IDs (supports multiple values)
  --filter-type               Filter by event type (e.g., deployment, server_provision, job_execution)
  --filter-severity           Filter by severity level (e.g., info, warning, error, critical)
  --filter-visibility         Filter by visibility scope (e.g., public, private, internal)
  --filter-infrastructure-id  Filter by infrastructure ID (supports multiple values)
  --filter-user-id            Filter by user ID who triggered the event (supports multiple values)
  --filter-server-id          Filter by server ID associated with the event (supports multiple values)
  --filter-job-id             Filter by job ID (supports multiple values)
  --filter-site-id            Filter by site/datacenter ID (supports multiple values)

SEARCH AND SORTING:
  --search                    Free-text search across event content
  --search-by                 Specify fields to search in (e.g., title, description, metadata)
  --sort-by                   Sort results by field with direction (e.g., "id:ASC", "createdAt:DESC")

PAGINATION:
  --page                      Page number for pagination (default: 0 for first page)
  --limit                     Maximum number of results per page (default: system default)

Examples:
  # List all events
  metalcloud event list

  # Filter by event type and severity
  metalcloud event list --filter-type deployment --filter-severity error

  # Filter by infrastructure and sort by creation time (newest first)
  metalcloud event list --filter-infrastructure-id 1234 --sort-by createdAt:DESC

  # Search for specific text in events
  metalcloud event list --search "failed deployment" --search-by title,description

  # Combine multiple filters
  metalcloud event list --filter-type server_provision --filter-severity warning,error --filter-site-id 5678

  # Paginated results
  metalcloud event list --page 2 --limit 50

  # Filter by multiple user IDs
  metalcloud event list --filter-user-id 123,456,789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EVENTS_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return event.EventList(cmd.Context(), event.ListFlags{
				FilterId:               eventFlags.filterId,
				FilterType:             eventFlags.filterType,
				FilterSeverity:         eventFlags.filterSeverity,
				FilterVisibility:       eventFlags.filterVisibility,
				FilterInfrastructureId: eventFlags.filterInfrastructureId,
				FilterUserId:           eventFlags.filterUserId,
				FilterServerId:         eventFlags.filterServerId,
				FilterJobId:            eventFlags.filterJobId,
				FilterSiteId:           eventFlags.filterSiteId,
				SortBy:                 eventFlags.sortBy,
				Page:                   eventFlags.page,
				Limit:                  eventFlags.limit,
				Search:                 eventFlags.search,
				SearchBy:               eventFlags.searchBy,
			})
		},
	}

	eventGetCmd = &cobra.Command{
		Use:     "get event_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific event",
		Long: `Retrieve detailed information about a specific event by its ID.

This command displays comprehensive information about a single event, including
its metadata, timestamp, severity, type, associated resources, and full description.

REQUIRED ARGUMENTS:
  event_id                    The unique identifier of the event to retrieve

Examples:
  # Get details of event with ID 12345
  metalcloud event get 12345

  # Get event details using the 'show' alias
  metalcloud event show 67890

The output includes:
- Event ID and timestamp
- Event type and severity level
- Associated infrastructure, server, job, or site information
- Full event description and metadata
- User who triggered the event (if applicable)`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EVENTS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return event.EventGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(eventCmd)

	eventCmd.AddCommand(eventListCmd)
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterId, "filter-id", nil, "Filter by event ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterType, "filter-type", nil, "Filter by event type.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterSeverity, "filter-severity", nil, "Filter by event severity.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterVisibility, "filter-visibility", nil, "Filter by event visibility.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterUserId, "filter-user-id", nil, "Filter by user ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterServerId, "filter-server-id", nil, "Filter by server ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterJobId, "filter-job-id", nil, "Filter by job ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.filterSiteId, "filter-site-id", nil, "Filter by site ID.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., id:ASC, createdAt:DESC).")
	eventListCmd.Flags().IntVar(&eventFlags.page, "page", 0, "Page number.")
	eventListCmd.Flags().IntVar(&eventFlags.limit, "limit", 0, "Limit number of results.")
	eventListCmd.Flags().StringVar(&eventFlags.search, "search", "", "Search term.")
	eventListCmd.Flags().StringSliceVar(&eventFlags.searchBy, "search-by", nil, "Fields to search by.")

	eventCmd.AddCommand(eventGetCmd)
}
