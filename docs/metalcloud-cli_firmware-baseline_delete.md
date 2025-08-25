## metalcloud-cli firmware-baseline delete

Delete a firmware baseline permanently

### Synopsis

Delete a firmware baseline permanently.

This command removes a firmware baseline definition from the system. This action is irreversible
and will delete all associated configuration data, component specifications, and deployment policies.

Note: Deleting a firmware baseline does not affect the underlying firmware catalogs or binaries.
Only the baseline definition and its configuration are removed.

Before deletion, ensure that:
- No active deployments are using this baseline
- No automated processes reference this baseline
- You have backups of the configuration if needed for future reference

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to delete

Examples:
  metalcloud-cli firmware-baseline delete 54321
  metalcloud-cli fw-baseline rm production-baseline-v2.1
  metalcloud-cli baseline delete dell-r640-standard

```
metalcloud-cli firmware-baseline delete firmware_baseline_id [flags]
```

### Options

```
  -h, --help   help for delete
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

