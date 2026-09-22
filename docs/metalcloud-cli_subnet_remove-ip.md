## metalcloud-cli subnet remove-ip

Remove an allocated IP from a subnet

### Synopsis

Remove one allocated IP from a subnet.

The subnet's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  subnet_id    The ID of the subnet
  ip_id        The ID of the IP to remove (see 'subnet ips subnet_id')

Examples:
  metalcloud-cli subnet remove-ip 123 456

```
metalcloud-cli subnet remove-ip subnet_id ip_id [flags]
```

### Options

```
  -h, --help   help for remove-ip
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

