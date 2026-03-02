## metalcloud-cli drive snapshot delete

Delete a snapshot by name

### Synopsis

Delete a specific snapshot by name for a drive within an infrastructure.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Required Flags:
  --name string                 The name of the snapshot to delete

Examples:
  # Delete a snapshot by name
  metalcloud-cli drive snapshot delete my-infrastructure 12345 --name my-snapshot

```
metalcloud-cli drive snapshot delete infrastructure_id_or_label drive_id [flags]
```

### Options

```
  -h, --help          help for delete
      --name string   Name of the snapshot to delete.
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

