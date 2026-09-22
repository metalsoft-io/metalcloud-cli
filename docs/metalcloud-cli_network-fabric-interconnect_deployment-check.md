## metalcloud-cli network-fabric-interconnect deployment-check

Validate whether the interconnect links can be activated

### Synopsis

Validate the links of a network fabric interconnect, reporting for each link
whether it can be activated, the resolved template variables and any errors.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Optional Flags:
  --link-ids strings         Restrict the check to these link IDs. Repeatable or comma-separated.

Examples:
  metalcloud network-fabric-interconnect deployment-check 12
  metalcloud nfi check dc1-dc2 --link-ids 3,4

```
metalcloud-cli network-fabric-interconnect deployment-check interconnect_id_or_label [flags]
```

### Options

```
  -h, --help               help for deployment-check
      --link-ids strings   Restrict the check to these link IDs.
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

