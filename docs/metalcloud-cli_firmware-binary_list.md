## metalcloud-cli firmware-binary list

List all firmware binaries

### Synopsis

List all firmware binaries in the system.

This command displays all available firmware binaries with their basic information including:
- Binary ID and name
- Catalog ID and package information
- Target hardware components and models
- Version information and update severity
- Reboot requirements and vendor release details
- Release timestamps and download URLs
- External ID and vendor info URLs

The output includes both catalog-managed binaries and individually registered ones.

Filters are applied server-side. Every --filter-* flag accepts a comma-separated list or can be
repeated; multiple values are combined with OR. A single value is sent as-is, so it may also use the
API operator syntax (e.g. '$not:$eq:value', '$ilike:value' or '$gt:2025-01-01').
Without --page/--limit all matching binaries are returned.

Optional Flags:
  --filter-id                        Filter by binary ID
  --filter-catalog-id                Filter by catalog ID
  --filter-update-severity           Filter by update severity (e.g. critical)
  --filter-reboot-required           Filter by reboot requirement (true or false)
  --filter-package-version           Filter by package version
  --filter-external-id               Filter by the vendor's artefact ID
  --filter-package-id                Filter by package ID
  --filter-created-timestamp         Filter by creation timestamp
  --filter-vendor-release-timestamp  Filter by vendor release timestamp
  --page                             Page number to retrieve (default: return all records)
  --limit                            Number of records per page, or maximum records when used without --page

Examples:
  # List all firmware binaries
  metalcloud-cli firmware-binary list

  # Critical binaries in catalog 8813
  metalcloud-cli firmware-binary list --filter-catalog-id 8813 --filter-update-severity critical

  # Binaries that do not require a reboot
  metalcloud-cli fw-binary ls --filter-reboot-required false

  # Second page of 50 binaries
  metalcloud-cli firmware-binary list --page 2 --limit 50

```
metalcloud-cli firmware-binary list [flags]
```

### Options

```
      --filter-catalog-id strings                 Filter the result by catalog ID.
      --filter-created-timestamp strings          Filter the result by creation timestamp.
      --filter-external-id strings                Filter the result by the vendor's artefact ID.
      --filter-id strings                         Filter the result by binary ID.
      --filter-package-id strings                 Filter the result by package ID.
      --filter-package-version strings            Filter the result by package version.
      --filter-reboot-required strings            Filter the result by reboot requirement (true or false).
      --filter-update-severity strings            Filter the result by update severity.
      --filter-vendor-release-timestamp strings   Filter the result by vendor release timestamp.
  -h, --help                                      help for list
      --limit int                                 Number of records per page, or maximum records without --page (default: return all records).
      --page int                                  Page number to retrieve (default: return all records).
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

* [metalcloud-cli firmware-binary](metalcloud-cli_firmware-binary.md)	 - Manage individual firmware binary files and packages

