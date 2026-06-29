## metalcloud-cli network-device update-port-config

Update the staged config (enable/description) of a port

### Synopsis

Update the staged configuration of a single network device port, addressed
by its numeric interface id.

This patches the port's persistent config (applied on the next fabric deploy),
not its live administrative status (use 'set-port-status' for that). You can set
the enabled flag and/or the interface description.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags:
  --port-id           Numeric interface id of the port to configure

Optional Flags (at least one must be specified):
  --enabled           Whether the port should be enabled (true/false)
  --description       Interface description text

Examples:
  # Enable a port and set its description
  metalcloud-cli network-device update-port-config 12345 --port-id 67890 --enabled --description "to_spine-s00_swp1s0"

  # Only set a description
  metalcloud-cli nd update-port-config 12345 --port-id 67890 --description "to_hgx-su00-h00_enp26s0f0np0"

```
metalcloud-cli network-device update-port-config <network_device_id> [flags]
```

### Options

```
      --description string   Interface description text.
      --enabled              Whether the port should be enabled.
  -h, --help                 help for update-port-config
      --port-id string       Numeric interface id of the port to configure.
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

