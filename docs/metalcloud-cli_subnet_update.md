## metalcloud-cli subnet update

Update an existing subnet

### Synopsis

Update an existing subnet in the MetalCloud infrastructure.

This command updates a subnet with the configuration provided through JSON input.
Only the fields included in the configuration will be updated, other fields remain unchanged.

Arguments:
  subnet_id    The ID of the subnet to update

Required Flags:
  --config-source    Source of the subnet configuration updates. Can be 'pipe' to read
                     from stdin, or a path to a JSON file containing the configuration.

Configuration Format:
The JSON configuration can contain the following fields (all optional for updates):
- label: Human-readable label for the subnet
- name: Subnet name
- defaultGatewayAddress: Gateway IP address
- isPool: Whether this subnet is an IP pool (true/false)
- allocationDenylist: List of IP ranges to exclude from allocation
- childOverlapAllowRules: Rules for allowing child subnet overlaps
- tags: Key-value pairs for tagging
- annotations: Additional metadata

Note: Core network settings (networkAddress, prefixLength) typically cannot be modified
after subnet creation due to infrastructure constraints.

Examples:
  # Update subnet from stdin
  echo '{"label":"updated-subnet","isPool":true}' | metalcloud-cli subnet update 123 --config-source pipe
  
  # Update subnet from file
  metalcloud-cli subnet update 456 --config-source updates.json
  
  # Update only tags
  echo '{"tags":{"environment":"production","team":"networking"}}' | metalcloud-cli subnet update 789 --config-source pipe

```
metalcloud-cli subnet update subnet_id [flags]
```

### Options

```
      --config-source string   Source of the subnet configuration updates. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli subnet](metalcloud-cli_subnet.md)	 - Manage network subnets and IP address pools

