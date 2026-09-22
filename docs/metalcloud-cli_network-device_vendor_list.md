## metalcloud-cli network-device vendor list

List the network device vendor profiles

### Synopsis

List the vendor profiles, one per network device driver.

Optional Flags:
  --filter-kind   Restrict the listing to these drivers

Examples:
  # List all vendor profiles
  metalcloud-cli network-device vendor list

  # List the SONiC profile only
  metalcloud-cli network-device vendor list --filter-kind sonic_enterprise

```
metalcloud-cli network-device vendor list [flags]
```

### Options

```
      --filter-kind strings   Filter the vendor profiles by network device driver.
  -h, --help                  help for list
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

* [metalcloud-cli network-device vendor](metalcloud-cli_network-device_vendor.md)	 - Manage the network device vendor profiles

