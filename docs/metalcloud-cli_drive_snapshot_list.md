## metalcloud-cli drive snapshot list

List all snapshots for a drive

### Synopsis

List all snapshots for a specific drive within an infrastructure.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # List snapshots for a drive
  metalcloud-cli drive snapshot list my-infrastructure 12345

  # List snapshots using infrastructure ID
  metalcloud-cli drive snapshot list 1001 67890

```
metalcloud-cli drive snapshot list infrastructure_id_or_label drive_id [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli drive snapshot](metalcloud-cli_drive_snapshot.md)	 - Manage drive snapshots

