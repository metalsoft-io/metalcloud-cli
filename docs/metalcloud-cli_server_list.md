## metalcloud-cli server list

List servers

### Synopsis

List all servers in the MetalSoft infrastructure.

This command displays information about all servers including their IDs, site locations, 
types, UUIDs, serial numbers, management addresses, vendors, models, and current status.

The table output also shows the operational fields: health state, details and last check, firmware
baseline compliance and target/current baselines, drift detection status, backup, drift detection and
health monitoring settings, and tags. These are independent of the lifecycle status.

Filters are applied server-side. Every --filter-* flag accepts a comma-separated list or can be
repeated; multiple values are combined with OR. A single value is sent as-is, so it may also use the
API operator syntax (e.g. '$not:$eq:value' or '$ilike:value').
Without --page/--limit all matching servers are returned.

Optional Flags:
  --show-credentials                      Display server IPMI credentials (username and password)
  --filter-status                         Filter by server status (e.g., active, registered, provisioning)
  --filter-type                           Filter by server type ID
  --filter-id                             Filter by server ID
  --filter-site-id                        Filter by site ID
  --filter-tags                           Filter by tag (e.g., a rack or room label)
  --filter-vendor                         Filter by vendor
  --filter-model                          Filter by model
  --filter-serial-number                  Filter by serial number
  --filter-management-address             Filter by BMC management address
  --filter-administration-state           Filter by administration state
  --filter-power-status                   Filter by power status
  --filter-dhcp-status                    Filter by DHCP status
  --filter-class                          Filter by server class
  --filter-disk-wipe                      Filter by disk wipe setting
  --filter-requires-manual-cleaning       Filter by manual cleaning requirement
  --filter-drift-detection-sync-status    Filter by drift detection sync status
  --filter-firmware-baseline-compliance   Filter by firmware baseline compliance status (unknown, compliant, non_compliant)
  --filter-firmware-backup-enabled        Filter by firmware backup setting
  --filter-firmware-drift-detection       Filter by firmware drift detection setting
  --filter-infrastructure-id              Filter by the infrastructure the server is allocated to
  --filter-owner-id                       Filter by the owner (user ID) of the infrastructure
  --page                                  Page number to retrieve (default: return all records)
  --limit                                 Number of records per page, or maximum records when used without --page

Examples:
  # List all servers
  metalcloud-cli server list

  # List servers with IPMI credentials
  metalcloud-cli server list --show-credentials

  # Filter servers by status
  metalcloud-cli server list --filter-status active,registered

  # Every server in one site carrying a given tag
  metalcloud-cli server list --filter-site-id 2 --filter-tags rack-12

  # Every server that is not compliant with its firmware baseline
  metalcloud-cli server list --filter-firmware-baseline-compliance non_compliant

  # Second page of 50 servers
  metalcloud-cli server list --page 2 --limit 50

  # Combine filters
  metalcloud-cli server list --filter-status active --filter-type 1


```
metalcloud-cli server list [flags]
```

### Options

```
      --filter-administration-state strings           Filter the result by administration state.
      --filter-class strings                          Filter the result by server class.
      --filter-dhcp-status strings                    Filter the result by DHCP status.
      --filter-disk-wipe strings                      Filter the result by disk wipe setting.
      --filter-drift-detection-sync-status strings    Filter the result by drift detection sync status.
      --filter-firmware-backup-enabled strings        Filter the result by firmware backup setting.
      --filter-firmware-baseline-compliance strings   Filter the result by firmware baseline compliance status (unknown, compliant, non_compliant).
      --filter-firmware-drift-detection strings       Filter the result by firmware drift detection setting.
      --filter-id strings                             Filter the result by server ID.
      --filter-infrastructure-id strings              Filter the result by the infrastructure the server is allocated to.
      --filter-management-address strings             Filter the result by management address.
      --filter-model strings                          Filter the result by model.
      --filter-owner-id strings                       Filter the result by the owner (user ID) of the infrastructure.
      --filter-power-status strings                   Filter the result by power status.
      --filter-requires-manual-cleaning strings       Filter the result by manual cleaning requirement.
      --filter-serial-number strings                  Filter the result by serial number.
      --filter-site-id strings                        Filter the result by site ID.
      --filter-status strings                         Filter the result by server status.
      --filter-tags strings                           Filter the result by server tag.
      --filter-type strings                           Filter the result by server type.
      --filter-vendor strings                         Filter the result by vendor.
  -h, --help                                          help for list
      --limit int                                     Number of records per page, or maximum records without --page (default: return all records).
      --page int                                      Page number to retrieve (default: return all records).
      --show-credentials                              If set returns the server IPMI credentials.
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

