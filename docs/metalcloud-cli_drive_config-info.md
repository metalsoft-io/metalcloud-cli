## metalcloud-cli drive config-info

Get configuration information for a drive

### Synopsis

Get configuration information for a drive within an infrastructure.

This command retrieves the current configuration information for a specified drive,
including all configuration parameters, settings, and current values.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # Get configuration information for a drive
  metalcloud-cli drive config-info my-infrastructure 12345

  # Get config info using infrastructure ID
  metalcloud-cli drive get-config-info 1001 67890

```
metalcloud-cli drive config-info infrastructure_id_or_label drive_id [flags]
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

* [metalcloud-cli drive](metalcloud-cli_drive.md)	 - Manage drives within infrastructures

