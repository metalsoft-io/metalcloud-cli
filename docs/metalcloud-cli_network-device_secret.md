## metalcloud-cli network-device secret

Manage the secrets of a network device

### Synopsis

Manage the named secrets stored for a network device.

Secret values are stored encrypted and are only revealed by the
'get-credentials' sub-command.

### Options

```
  -h, --help   help for secret
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure
* [metalcloud-cli network-device secret get-credentials](metalcloud-cli_network-device_secret_get-credentials.md)	 - Reveal the value of a network device secret
* [metalcloud-cli network-device secret list](metalcloud-cli_network-device_secret_list.md)	 - List the secret names of a network device
* [metalcloud-cli network-device secret remove](metalcloud-cli_network-device_secret_remove.md)	 - Delete one secret of a network device
* [metalcloud-cli network-device secret remove-all](metalcloud-cli_network-device_secret_remove-all.md)	 - Delete every secret of a network device
* [metalcloud-cli network-device secret set](metalcloud-cli_network-device_secret_set.md)	 - Store a named secret of a network device

