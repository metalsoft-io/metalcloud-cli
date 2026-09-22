## metalcloud-cli network-device secret list

List the secret names of a network device

### Synopsis

List the names of the secrets stored for a network device. Values are not
returned by this command.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the secret names of device 12345
  metalcloud-cli network-device secret list 12345

```
metalcloud-cli network-device secret list <network_device_id> [flags]
```

### Options

```
  -h, --help   help for list
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

