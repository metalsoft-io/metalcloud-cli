## metalcloud-cli network-device secret get-credentials

Reveal the value of a network device secret

### Synopsis

Reveal the unencrypted value of one named secret of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Examples:
  # Reveal the enable password of device 12345
  metalcloud-cli network-device secret get-credentials 12345 enable_password

```
metalcloud-cli network-device secret get-credentials <network_device_id> <name> [flags]
```

### Options

```
  -h, --help   help for get-credentials
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

