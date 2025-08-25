## metalcloud-cli firmware-baseline get

Get detailed information about a specific firmware baseline

### Synopsis

Get detailed information about a specific firmware baseline.

This command displays comprehensive information about a firmware baseline including:
- Basic metadata (name, description, level, level filter)
- Catalog associations
- Creation timestamp
- Unique baseline identifier

The firmware baseline contains the essential configuration for standardized firmware
deployments across compatible hardware.

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to retrieve

Examples:
  metalcloud-cli firmware-baseline get 54321
  metalcloud-cli fw-baseline show dell-r640-standard
  metalcloud-cli baseline get production-baseline-v2.1

```
metalcloud-cli firmware-baseline get firmware_baseline_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli firmware-baseline](metalcloud-cli_firmware-baseline.md)	 - Manage firmware baselines for consistent hardware configurations

