## metalcloud-cli custom-iso delete

Delete a custom ISO permanently

### Synopsis

Delete a custom ISO image permanently from your account.

This action cannot be undone. The custom ISO will be removed from all servers
where it might be mounted and will no longer be available for provisioning.

Arguments:
  custom_iso_id   ID of the custom ISO to delete (required)

Required permissions:
  - custom_iso:write

Dependencies:
  - Custom ISO must exist and be accessible
  - Custom ISO should not be actively used by running servers

Examples:
  # Delete custom ISO with ID 12345
  metalcloud-cli custom-iso delete 12345
  
  # Delete using shorter alias
  metalcloud-cli iso rm 12345

```
metalcloud-cli custom-iso delete custom_iso_id [flags]
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

* [metalcloud-cli custom-iso](metalcloud-cli_custom-iso.md)	 - Manage custom ISO images for server provisioning

