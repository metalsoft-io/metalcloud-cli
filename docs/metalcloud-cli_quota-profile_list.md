## metalcloud-cli quota-profile list

List quota profiles

### Synopsis

List all quota profiles in the MetalSoft infrastructure.

This command displays information about all quota profiles including their IDs,
names, and descriptions.

Examples:
  # List all quota profiles
  metalcloud-cli quota-profile list


```
metalcloud-cli quota-profile list [flags]
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

* [metalcloud-cli quota-profile](metalcloud-cli_quota-profile.md)	 - Quota Profile management

