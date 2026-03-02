## metalcloud-cli drive snapshot

Manage drive snapshots

### Synopsis

Manage snapshots for drives within infrastructures.

Snapshots allow you to capture and restore the state of a drive at a specific point in time.

Available Commands:
  list          List all snapshots for a drive
  create        Create a new snapshot for a drive
  delete        Delete a snapshot by name
  restore       Restore a drive to a specific snapshot

Examples:
  # List snapshots for a drive
  metalcloud-cli drive snapshot list my-infrastructure 12345

  # Create a snapshot
  metalcloud-cli drive snapshot create my-infrastructure 12345

  # Delete a snapshot by name
  metalcloud-cli drive snapshot delete my-infrastructure 12345 --name my-snapshot

  # Restore a drive to a snapshot
  metalcloud-cli drive snapshot restore my-infrastructure 12345 --name my-snapshot

### Options

```
  -h, --help   help for snapshot
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

* [metalcloud-cli drive](metalcloud-cli_drive.md)	 - Manage drives within infrastructures
* [metalcloud-cli drive snapshot create](metalcloud-cli_drive_snapshot_create.md)	 - Create a new snapshot for a drive
* [metalcloud-cli drive snapshot delete](metalcloud-cli_drive_snapshot_delete.md)	 - Delete a snapshot by name
* [metalcloud-cli drive snapshot list](metalcloud-cli_drive_snapshot_list.md)	 - List all snapshots for a drive
* [metalcloud-cli drive snapshot restore](metalcloud-cli_drive_snapshot_restore.md)	 - Restore a drive to a specific snapshot

