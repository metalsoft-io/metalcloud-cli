## metalcloud-cli device-config-template profile list

List device configuration template profiles with optional filtering

```
metalcloud-cli device-config-template profile list [flags]
```

### Options

```
      --filter-id strings                  Filter by profile ID.
      --filter-network-device-id strings   Filter by network device ID.
      --filter-network-fabric-id strings   Filter by network fabric ID.
      --filter-template-id strings         Filter by device configuration template ID.
  -h, --help                               help for list
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

* [metalcloud-cli device-config-template profile](metalcloud-cli_device-config-template_profile.md)	 - Manage device configuration template profiles

