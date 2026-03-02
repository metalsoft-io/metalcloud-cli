## metalcloud-cli file-share snapshot restore

Restore a file share to a specific snapshot

### Synopsis

Restore a file share to the state captured by a specific snapshot.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  file_share_id                The unique identifier of the file share

Required Flags:
  --name string                 The name of the snapshot to restore to

Examples:
  # Restore a file share to a snapshot
  metalcloud-cli file-share snapshot restore my-infrastructure 12345 --name my-snapshot

```
metalcloud-cli file-share snapshot restore infrastructure_id_or_label file_share_id [flags]
```

### Options

```
  -h, --help          help for restore
      --name string   Name of the snapshot to restore to.
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

* [metalcloud-cli file-share snapshot](metalcloud-cli_file-share_snapshot.md)	 - Manage file share snapshots

