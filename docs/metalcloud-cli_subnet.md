## metalcloud-cli subnet

Manage network subnets and IP address pools

### Synopsis

Manage network subnets and IP address pools in the MetalCloud infrastructure.

Subnets define network segments with specific IP address ranges and can be configured as:
- Regular subnets: Fixed network segments with defined address ranges
- IP pools: Dynamic address pools for automatic IP allocation

Available commands allow you to list, create, update, delete subnets and view configuration examples.

Available Commands:
  list, get, create, update, delete, config-example
  ips              List the allocated IPs of a subnet
  ip-ranges        List the IP ranges of a subnet
  capacity         Show how much of a subnet is still available
  remove-ip        Remove an allocated IP from a subnet
  remove-ip-range  Remove an IP range from a subnet

### Options

```
  -h, --help   help for subnet
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli subnet capacity](metalcloud-cli_subnet_capacity.md)	 - Show how much of a subnet is still available
* [metalcloud-cli subnet config-example](metalcloud-cli_subnet_config-example.md)	 - Display subnet configuration example
* [metalcloud-cli subnet create](metalcloud-cli_subnet_create.md)	 - Create a new subnet
* [metalcloud-cli subnet delete](metalcloud-cli_subnet_delete.md)	 - Delete a subnet
* [metalcloud-cli subnet get](metalcloud-cli_subnet_get.md)	 - Get detailed information about a specific subnet
* [metalcloud-cli subnet ip-ranges](metalcloud-cli_subnet_ip-ranges.md)	 - List IP ranges in a subnet
* [metalcloud-cli subnet ips](metalcloud-cli_subnet_ips.md)	 - List allocated IPs in a subnet
* [metalcloud-cli subnet list](metalcloud-cli_subnet_list.md)	 - List all subnets and IP pools
* [metalcloud-cli subnet remove-ip](metalcloud-cli_subnet_remove-ip.md)	 - Remove an allocated IP from a subnet
* [metalcloud-cli subnet remove-ip-range](metalcloud-cli_subnet_remove-ip-range.md)	 - Remove an IP range from a subnet
* [metalcloud-cli subnet update](metalcloud-cli_subnet_update.md)	 - Update an existing subnet

