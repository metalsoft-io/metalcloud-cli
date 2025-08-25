## metalcloud-cli fabric remove-device

Remove network device from a fabric

### Synopsis

Remove a network device from an existing fabric.

This command disassociates a network device from a fabric, removing it from the fabric's
network topology. The device will no longer be managed by the fabric configuration.

Arguments:
  fabric_id    The ID or label of the fabric to remove the device from
  device_id    The ID or label of the device to remove from the fabric

Examples:
  # Remove device from fabric by IDs
  metalcloud fabric remove-device 12345 device123
  
  # Remove device using labels
  metalcloud fabric remove-device my-fabric switch-01
  
  # Using alias
  metalcloud fabric delete-device my-fabric device123

```
metalcloud-cli fabric remove-device fabric_id device_id [flags]
```

### Options

```
  -h, --help   help for remove-device
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

