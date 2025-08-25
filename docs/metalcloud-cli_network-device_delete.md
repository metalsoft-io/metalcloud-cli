## metalcloud-cli network-device delete

Delete a network device from the infrastructure

### Synopsis

Delete a network device from the infrastructure. This operation will remove
the device from management and monitoring. The physical device will no longer
be controlled by the system.

WARNING: This operation is irreversible. Ensure the device is not in use
before deletion.

Arguments:
  network_device_id   The unique identifier of the network device to delete

Examples:
  # Delete network device
  metalcloud-cli network-device delete 12345

  # Using alias
  metalcloud-cli switch rm 12345

```
metalcloud-cli network-device delete <network_device_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

