## metalcloud-cli custom-iso list

List all available custom ISO images

### Synopsis

List all custom ISO images available in your account.

This command displays a table of all custom ISOs showing their ID, name, 
description, size, creation date, and availability status.

Required permissions:
  - custom_iso:read

Examples:
  # List all custom ISOs
  metalcloud-cli custom-iso list
  
  # List custom ISOs with shorter alias
  metalcloud-cli iso ls

```
metalcloud-cli custom-iso list [flags]
```

### Options

```
  -h, --help   help for list
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

