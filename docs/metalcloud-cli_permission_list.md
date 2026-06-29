## metalcloud-cli permission list

List permissions

### Synopsis

List all permissions.

This command displays information about all permissions including their IDs,
names, labels, types, and descriptions.

Examples:
  # List all permissions
  metalcloud-cli permission list


```
metalcloud-cli permission list [flags]
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

* [metalcloud-cli permission](metalcloud-cli_permission.md)	 - Permission management

