## metalcloud-cli firmware-catalog delete

Delete a firmware catalog

### Synopsis

Delete a firmware catalog permanently.

This command removes a firmware catalog from the system. This action is irreversible
and will delete all associated firmware package information.

Arguments:
  firmware_catalog_id    The ID of the firmware catalog to delete

Examples:
  metalcloud-cli firmware-catalog delete 12345
  metalcloud-cli fw-catalog rm dell-r640-catalog

```
metalcloud-cli firmware-catalog delete firmware_catalog_id [flags]
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

* [metalcloud-cli firmware-catalog](metalcloud-cli_firmware-catalog.md)	 - Manage firmware catalogs for server hardware updates

