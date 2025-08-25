## metalcloud-cli file-share config-info

Get configuration information for a file share

### Synopsis

Get configuration information for a specific file share including technical details,
settings, and current configuration state.

This command displays the complete configuration profile of a file share including
storage configuration, networking settings, access control, and other technical
parameters that may be needed for troubleshooting or integration purposes.

Required Arguments:
  infrastructure_id_or_label    The infrastructure ID (numeric) or label containing the file share
  file_share_id                 The ID of the file share to get configuration information for

Examples:
  # Get configuration information for a file share
  metalcloud-cli file-share config-info my-infrastructure 12345

  # Get configuration info using infrastructure ID
  metalcloud-cli file-share config-info 100 12345

```
metalcloud-cli file-share config-info infrastructure_id_or_label file_share_id [flags]
```

### Options

```
  -h, --help   help for config-info
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

