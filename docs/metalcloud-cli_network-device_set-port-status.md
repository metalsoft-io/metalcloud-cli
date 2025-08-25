## metalcloud-cli network-device set-port-status

Enable or disable a specific port on the network device

### Synopsis

Set the administrative status of a specific port on the network device.
This allows you to enable (bring up) or disable (bring down) individual
ports for maintenance or troubleshooting purposes.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags (both must be specified):
  --port-id    ID or name of the port to modify
  --action     Action to perform on the port
               Values: 'up' (enable port), 'down' (disable port)

Examples:
  # Bring port down for maintenance
  metalcloud-cli network-device set-port-status 12345 --port-id eth0/1 --action down

  # Bring port back up
  metalcloud-cli network-device set-port-status 12345 --port-id eth0/1 --action up

  # Using port number
  metalcloud-cli nd set-port-status 12345 --port-id 24 --action up

```
metalcloud-cli network-device set-port-status <network_device_id> [flags]
```

### Options

```
      --action string    Action to perform on the port (up/down).
  -h, --help             help for set-port-status
      --port-id string   ID of the port to change status.
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

