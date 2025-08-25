## metalcloud-cli server list

List servers

### Synopsis

List all servers in the MetalSoft infrastructure.

This command displays information about all servers including their IDs, site locations, 
types, UUIDs, serial numbers, management addresses, vendors, models, and current status.

Optional Flags:
  --show-credentials     Display server IPMI credentials (username and password)
  --filter-status        Filter servers by status (e.g., active, registered, provisioning)
  --filter-type          Filter servers by type ID

Examples:
  # List all servers
  metalcloud-cli server list

  # List servers with IPMI credentials
  metalcloud-cli server list --show-credentials

  # Filter servers by status
  metalcloud-cli server list --filter-status active,registered

  # Filter servers by type
  metalcloud-cli server list --filter-type 1,2,3

  # Combine filters
  metalcloud-cli server list --filter-status active --filter-type 1


```
metalcloud-cli server list [flags]
```

### Options

```
      --filter-status strings   Filter the result by server status.
      --filter-type strings     Filter the result by server type.
  -h, --help                    help for list
      --show-credentials        If set returns the server IPMI credentials.
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

