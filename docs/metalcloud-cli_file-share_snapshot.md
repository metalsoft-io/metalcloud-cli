## metalcloud-cli file-share snapshot

Manage file share snapshots

### Synopsis

Manage snapshots for file shares within infrastructures.

Snapshots allow you to capture and restore the state of a file share at a specific point in time.

Available Commands:
  list          List all snapshots for a file share
  create        Create a new snapshot for a file share
  delete        Delete a snapshot by name
  restore       Restore a file share to a specific snapshot

Examples:
  # List snapshots for a file share
  metalcloud-cli file-share snapshot list my-infrastructure 12345

  # Create a snapshot
  metalcloud-cli file-share snapshot create my-infrastructure 12345

  # Delete a snapshot by name
  metalcloud-cli file-share snapshot delete my-infrastructure 12345 --name my-snapshot

  # Restore a file share to a snapshot
  metalcloud-cli file-share snapshot restore my-infrastructure 12345 --name my-snapshot

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

* [metalcloud-cli file-share](metalcloud-cli_file-share.md)	 - Manage file shares for infrastructure resources
* [metalcloud-cli file-share snapshot create](metalcloud-cli_file-share_snapshot_create.md)	 - Create a new snapshot for a file share
* [metalcloud-cli file-share snapshot delete](metalcloud-cli_file-share_snapshot_delete.md)	 - Delete a snapshot by name
* [metalcloud-cli file-share snapshot list](metalcloud-cli_file-share_snapshot_list.md)	 - List all snapshots for a file share
* [metalcloud-cli file-share snapshot restore](metalcloud-cli_file-share_snapshot_restore.md)	 - Restore a file share to a specific snapshot

