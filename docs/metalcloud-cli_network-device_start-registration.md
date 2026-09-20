## metalcloud-cli network-device start-registration

Start the onboarding registration of a network device

### Synopsis

Start the onboarding registration of a network device. Supported only by the
drivers that report the 'start-registration' activation action; see
'network-device driver-capabilities'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Start the registration of device 12345
  metalcloud-cli network-device start-registration 12345

```
metalcloud-cli network-device start-registration <network_device_id> [flags]
```

### Options

```
  -h, --help   help for start-registration
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

