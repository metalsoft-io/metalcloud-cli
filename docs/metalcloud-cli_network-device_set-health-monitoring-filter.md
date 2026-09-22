## metalcloud-cli network-device set-health-monitoring-filter

Set the network device filter of a health monitoring socket

### Synopsis

Set the network device filter applied to the health monitoring WebSocket stream
of the given socket. Only devices matching the filter are streamed to that
socket. The endpoint returns no content.

Required Flags:
  --socket-id             The id of the health monitoring WebSocket connection

Optional Flags:
  --filter-id             Restrict the stream to these network device ids
  --filter-site-id        Restrict the stream to these site ids
  --filter-status         Restrict the stream to these device statuses
  --filter-health-status  Restrict the stream to these health statuses

Examples:
  # Stream only the devices of site 1
  metalcloud-cli network-device set-health-monitoring-filter --socket-id abc123 --filter-site-id 1

  # Stream only unhealthy devices
  metalcloud-cli network-device set-health-monitoring-filter --socket-id abc123 --filter-health-status unhealthy

```
metalcloud-cli network-device set-health-monitoring-filter [flags]
```

### Options

```
      --filter-health-status strings   Restrict the stream to these health statuses.
      --filter-id strings              Restrict the stream to these network device ids.
      --filter-site-id strings         Restrict the stream to these site ids.
      --filter-status strings          Restrict the stream to these network device statuses.
  -h, --help                           help for set-health-monitoring-filter
      --socket-id string               The id of the health monitoring WebSocket connection.
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

