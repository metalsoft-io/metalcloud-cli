## metalcloud-cli fabric remove-link

Remove a network fabric link

### Synopsis

Remove a network fabric link from an existing fabric.

This command removes a link from the fabric, disconnecting the associated network devices.

Arguments:
  fabric_id    The ID or label of the fabric to remove the link from
  link_id      The ID of the link to remove from the fabric

Examples:
  # Remove link from fabric by IDs
  metalcloud fabric remove-link 12345 67890
  
  # Remove link using fabric label
  metalcloud fabric remove-link my-fabric 67890
  
  # Using alias
  metalcloud fabric delete-link my-fabric 67890

```
metalcloud-cli fabric remove-link fabric_id link_id [flags]
```

### Options

```
  -h, --help   help for remove-link
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

