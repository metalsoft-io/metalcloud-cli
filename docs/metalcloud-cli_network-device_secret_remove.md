## metalcloud-cli network-device secret remove

Delete one secret of a network device

### Synopsis

Delete one named secret of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Examples:
  # Delete the enable password of device 12345
  metalcloud-cli network-device secret remove 12345 enable_password

```
metalcloud-cli network-device secret remove <network_device_id> <name> [flags]
```

### Options

```
  -h, --help   help for remove
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

* [metalcloud-cli network-device secret](metalcloud-cli_network-device_secret.md)	 - Manage the secrets of a network device

