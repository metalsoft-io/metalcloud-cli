## metalcloud-cli network-device disable-syslog

Unsubscribe a network device from remote syslog

### Synopsis

Unsubscribe a network device from remote syslog, the counterpart of
'enable-syslog'. The operation is asynchronous and returns the job that carries
it out.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Disable syslog for device 12345
  metalcloud-cli network-device disable-syslog 12345

```
metalcloud-cli network-device disable-syslog <network_device_id> [flags]
```

### Options

```
  -h, --help   help for disable-syslog
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

