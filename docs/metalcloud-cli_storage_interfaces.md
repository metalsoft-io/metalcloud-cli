## metalcloud-cli storage interfaces

List the interfaces of a storage pool

### Synopsis

List all interfaces of a storage pool.

The output shows each interface ID, its name, the protocols it serves, the
storage nodes it belongs to, whether it is an uplink, whether it is used for
deploys and the network device interface it is linked to.

Required Arguments:
  storage_id    The numeric ID of the storage pool

Examples:
  # List the interfaces of storage pool 123
  metalcloud-cli storage interfaces 123

```
metalcloud-cli storage interfaces storage_id [flags]
```

### Options

```
  -h, --help   help for interfaces
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

* [metalcloud-cli storage](metalcloud-cli_storage.md)	 - Manage storage pools and related resources

