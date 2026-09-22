## metalcloud-cli network-device-controller get-credentials

Get the management credentials of a controller

### Synopsis

Get the management credentials used to connect to a network device controller.

Required Arguments:
  controller_id_or_identifier   The ID or identifier (hostname) of the controller

Examples:
  metalcloud network-device-controller get-credentials 7

```
metalcloud-cli network-device-controller get-credentials controller_id_or_identifier [flags]
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

* [metalcloud-cli network-device-controller](metalcloud-cli_network-device-controller.md)	 - Network device controller management

