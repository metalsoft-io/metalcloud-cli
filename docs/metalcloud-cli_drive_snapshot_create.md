## metalcloud-cli drive snapshot create

Create a new snapshot for a drive

### Synopsis

Create a new snapshot for a specific drive within an infrastructure.

This captures the current state of the drive as a point-in-time snapshot.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # Create a snapshot for a drive
  metalcloud-cli drive snapshot create my-infrastructure 12345

  # Create a snapshot using infrastructure ID
  metalcloud-cli drive snapshot create 1001 67890

```
metalcloud-cli drive snapshot create infrastructure_id_or_label drive_id [flags]
```

### Options

```
  -h, --help   help for create
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

