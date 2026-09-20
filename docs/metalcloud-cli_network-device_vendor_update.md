## metalcloud-cli network-device vendor update

Update one network device vendor profile

### Synopsis

Update the SNMP OID groups, health check rules and backup file list of one
vendor profile.

Required Arguments:
  vendor_id         The numeric id of the vendor profile

Required Flags:
  --config-source   Source of the vendor configuration
                    Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Update vendor profile 3 from a file
  metalcloud-cli network-device vendor update 3 --config-source vendor.json

```
metalcloud-cli network-device vendor update <vendor_id> [flags]
```

### Options

```
      --config-source string   Source of the vendor configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

