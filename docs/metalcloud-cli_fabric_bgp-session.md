## metalcloud-cli fabric bgp-session

Manage the BGP sessions of a fabric

### Synopsis

Manage the BGP sessions of a network fabric.

A BGP session attaches a BGP numbering and link configuration to one fabric
link or link aggregation, and carries the custom variables exposed to the
device configuration templates as bgp_session_custom_variables.

Commands:
  list, get, create, update, delete, config-example

### Options

```
  -h, --help   help for bgp-session
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics
* [metalcloud-cli fabric bgp-session config-example](metalcloud-cli_fabric_bgp-session_config-example.md)	 - Show a BGP session configuration example
* [metalcloud-cli fabric bgp-session create](metalcloud-cli_fabric_bgp-session_create.md)	 - Create a BGP session on a fabric
* [metalcloud-cli fabric bgp-session delete](metalcloud-cli_fabric_bgp-session_delete.md)	 - Delete a BGP session of a fabric
* [metalcloud-cli fabric bgp-session get](metalcloud-cli_fabric_bgp-session_get.md)	 - Get one BGP session of a fabric
* [metalcloud-cli fabric bgp-session list](metalcloud-cli_fabric_bgp-session_list.md)	 - List the BGP sessions of a fabric
* [metalcloud-cli fabric bgp-session update](metalcloud-cli_fabric_bgp-session_update.md)	 - Update a BGP session of a fabric

