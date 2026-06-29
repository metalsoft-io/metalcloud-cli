## metalcloud-cli quota-profile get

Get detailed quota profile information

### Synopsis

Get detailed information for a specific quota profile.

This command retrieves comprehensive information about a quota profile including
its configuration, limits, and other metadata.

Required Arguments:
  profile_id            The ID of the quota profile to retrieve information for

Examples:
  # Get quota profile information
  metalcloud-cli quota-profile get example-quota-profile

  # Get quota profile information using alias
  metalcloud-cli quota-profile show example-quota-profile


```
metalcloud-cli quota-profile get profile_id [flags]
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

* [metalcloud-cli quota-profile](metalcloud-cli_quota-profile.md)	 - Quota Profile management

