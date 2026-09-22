## metalcloud-cli network-device drift

Inspect the configuration drift of a network device

### Synopsis

Inspect and acknowledge the configuration drift detected between the running
configuration of a network device and its target snapshot.

### Options

```
  -h, --help   help for drift
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
* [metalcloud-cli network-device drift acknowledge](metalcloud-cli_network-device_drift_acknowledge.md)	 - Acknowledge a configuration drift entry
* [metalcloud-cli network-device drift get](metalcloud-cli_network-device_drift_get.md)	 - Get one configuration drift entry
* [metalcloud-cli network-device drift list](metalcloud-cli_network-device_drift_list.md)	 - List the configuration drift of a network device

