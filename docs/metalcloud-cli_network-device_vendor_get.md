## metalcloud-cli network-device vendor get

Get one network device vendor profile

### Synopsis

Display one vendor profile: its SNMP OID groups, health check rules and the
optional files included in configuration backups.

Required Arguments:
  vendor_id   The numeric id of the vendor profile

Examples:
  # Show vendor profile 3
  metalcloud-cli network-device vendor get 3

```
metalcloud-cli network-device vendor get <vendor_id> [flags]
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

* [metalcloud-cli network-device vendor](metalcloud-cli_network-device_vendor.md)	 - Manage the network device vendor profiles

