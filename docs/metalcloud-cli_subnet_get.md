## metalcloud-cli subnet get

Get detailed information about a specific subnet

### Synopsis

Get detailed information about a specific subnet including all its configuration properties.

The command shows comprehensive subnet details including:
- Basic subnet information (ID, name, label)
- Network configuration (address, prefix, netmask, gateway)
- IP pool configuration status
- Allocation denylist and rules
- Creation and modification timestamps
- Associated tags and annotations

Arguments:
  subnet_id    The ID of the subnet to retrieve

Examples:
  metalcloud-cli subnet get 123
  metalcloud-cli subnets show 456
  metalcloud-cli net get 789

```
metalcloud-cli subnet get subnet_id [flags]
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

* [metalcloud-cli subnet](metalcloud-cli_subnet.md)	 - Manage network subnets and IP address pools

