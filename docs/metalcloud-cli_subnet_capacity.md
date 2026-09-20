## metalcloud-cli subnet capacity

Show how much of a subnet is still available

### Synopsis

Report how much of a subnet is still available.

A subnet that is not a pool reports IP counts (total, free, used); a pool
reports its free blocks instead, in largest-block CIDR form. Counts are decimal
strings because a large IPv6 subnet holds more addresses than a JSON number can
carry exactly.

Required Arguments:
  subnet_id    The ID of the subnet

Optional Flags:
  --prefix-length  Pools only: also report how many blocks of this prefix length
                   still fit. Must not be shorter than the pool prefix itself.

Examples:
  metalcloud-cli subnet capacity 123
  metalcloud-cli subnet capacity 123 --prefix-length 26

```
metalcloud-cli subnet capacity subnet_id [flags]
```

### Options

```
  -h, --help                help for capacity
      --prefix-length int   Pools only: also report how many blocks of this prefix length still fit.
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

