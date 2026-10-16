## metalcloud-cli event list

List and filter system events

### Synopsis

List and filter system events with comprehensive filtering and search capabilities.

This command allows you to view events from across your MetalCloud infrastructure
with support for multiple filtering options, sorting, pagination, and text search.

FILTERING OPTIONS:
  --filter-id                  Filter by specific event IDs (supports multiple values)
  --filter-type               Filter by event type (e.g., deployment, server_provision, job_execution)
  --filter-severity           Filter by event level (info, warn, error, debug)
  --filter-visibility         Filter by visibility scope (e.g., public, private, internal)
  --filter-infrastructure-id  Filter by infrastructure ID (supports multiple values)
  --filter-user-id            Filter by user ID who triggered the event (supports multiple values)
  --filter-user-id-authenticated  Filter by the ID of the authenticated user (supports multiple values)
  --filter-server-id          Filter by server ID associated with the event (supports multiple values)
  --filter-network-device-id  Filter by network device ID (supports multiple values)
  --filter-job-id             Filter by job ID (supports multiple values)
  --filter-site-id            Filter by site/datacenter ID (supports multiple values)

OUTPUT:
  The table, CSV and markdown formats include, besides the event ID, type, level, visibility, time
  and title: the message, the user and authenticated user IDs, and the IDs of the entities the
  event belongs to (infrastructure, server, network device, network device controller, storage,
  VM pool, site and job). In the text and markdown tables, the ID columns that are empty for every
  event shown are left out; CSV always carries all of them. Use -f json or -f yaml for the full
  record, including the context, the authenticated user's email and the remote IP address.

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
  metalcloud event list --filter-type server_provision --filter-severity warn,error --filter-site-id 5678

  # Paginated results
  metalcloud event list --page 2 --limit 50

  # Filter by multiple user IDs
  metalcloud event list --filter-user-id 123,456,789

```
metalcloud-cli event list [flags]
```

### Options

```
      --filter-id strings                      Filter by event ID.
      --filter-infrastructure-id strings       Filter by infrastructure ID.
      --filter-job-id strings                  Filter by job ID.
      --filter-network-device-id strings       Filter by network device ID.
      --filter-server-id strings               Filter by server ID.
      --filter-severity strings                Filter by event level (info, warn, error, debug).
      --filter-site-id strings                 Filter by site ID.
      --filter-type strings                    Filter by event type.
      --filter-user-id strings                 Filter by user ID.
      --filter-user-id-authenticated strings   Filter by the ID of the authenticated user.
      --filter-visibility strings              Filter by event visibility.
  -h, --help                                   help for list
      --limit int                              Limit number of results.
      --page int                               Page number.
      --search string                          Search term.
      --search-by strings                      Fields to search by.
      --sort-by strings                        Sort by fields (e.g., id:ASC, createdAt:DESC).
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API key
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli event](metalcloud-cli_event.md)	 - Manage and monitor system events

