## metalcloud-cli network-device snapshot list

List the configuration snapshots of a network device

### Synopsis

List the configuration snapshots stored for a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --kind   Restrict the listing to one snapshot class

Examples:
  # List all snapshots of device 12345
  metalcloud-cli network-device snapshot list 12345

  # List only the backup snapshots
  metalcloud-cli network-device snapshot list 12345 --kind backup

```
metalcloud-cli network-device snapshot list <network_device_id> [flags]
```

### Options

```
  -h, --help          help for list
      --kind string   Restrict the listing to one snapshot class.
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

* [metalcloud-cli network-device snapshot](metalcloud-cli_network-device_snapshot.md)	 - Inspect the configuration snapshots of a network device

