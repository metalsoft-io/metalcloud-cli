## metalcloud-cli subnet list

List all subnets and IP pools

### Synopsis

List all subnets and IP address pools in the MetalCloud infrastructure.

This command displays a tabular view of all subnets with key information including:
- Subnet ID and name
- IP version (IPv4/IPv6)  
- Network address and prefix length
- Netmask
- Pool status (whether it's configured as an IP pool)
- Creation timestamp

Examples:
  metalcloud-cli subnet list
  metalcloud-cli subnets ls
  metalcloud-cli net list

```
metalcloud-cli subnet list [flags]
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

* [metalcloud-cli subnet](metalcloud-cli_subnet.md)	 - Manage network subnets and IP address pools

