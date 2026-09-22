## metalcloud-cli network-device secret set

Store a named secret of a network device

### Synopsis

Store (or replace) one named secret of a network device. The value is stored
encrypted.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  name                The name of the secret

Required Flags:
  --value             The secret value to store

Examples:
  # Store the enable password of device 12345
  metalcloud-cli network-device secret set 12345 enable_password --value s3cr3t

```
metalcloud-cli network-device secret set <network_device_id> <name> [flags]
```

### Options

```
  -h, --help           help for set
      --value string   The secret value to store.
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

