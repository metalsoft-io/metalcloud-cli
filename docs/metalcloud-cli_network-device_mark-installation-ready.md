## metalcloud-cli network-device mark-installation-ready

Mark the physical installation of a network device as done

### Synopsis

Mark the physical installation of a network device as complete so that
onboarding can continue. Supported only by the drivers that report the
'mark-install-ready' activation action; see 'network-device driver-capabilities'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Mark device 12345 as installed
  metalcloud-cli network-device mark-installation-ready 12345

```
metalcloud-cli network-device mark-installation-ready <network_device_id> [flags]
```

### Options

```
  -h, --help   help for mark-installation-ready
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

