## metalcloud-cli fabric get-links

List links in a fabric

### Synopsis

List all network fabric links in a specific fabric.

This command displays a table showing all links that are part of the specified fabric,
including their link information, status, and connection details.

Arguments:
  fabric_id    The ID or label of the fabric to list links for

Examples:
  # List links in fabric by ID
  metalcloud fabric get-links 12345
  
  # List links in fabric by label
  metalcloud fabric get-links my-fabric-label
  
  # Using alias
  metalcloud fabric list-links 12345

```
metalcloud-cli fabric get-links fabric_id [flags]
```

### Options

```
  -h, --help   help for get-links
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

