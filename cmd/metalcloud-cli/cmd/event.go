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
		Short:   "Event management",
		Long:    `Event management commands.`,
	}

	eventListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List events.",
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
		Use:          "get event_id",
		Aliases:      []string{"show"},
		Short:        "Get event details.",
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
