## metalcloud-cli server firmware schedule-upgrade-batch

Schedule a firmware upgrade for several servers at once

### Synopsis

Schedule a firmware upgrade for several servers at once.

The API models this as an action on a single server, so the first server ID given
is used to address the endpoint while the full list of servers is sent in the
request body. The result lists the error reported for every server that could not
be scheduled.

Required Arguments:
  server_id...           The IDs of the servers to schedule the upgrade for

Optional Flags:
  --schedule-timestamp       When the firmware upgrade should run (RFC 3339 timestamp)
  --confirmation-required    Require a confirmation before the scheduled upgrade runs

Examples:
  # Schedule an upgrade for servers 123 and 124 as soon as possible
  metalcloud-cli server firmware schedule-upgrade-batch 123 124

  # Schedule an upgrade for a maintenance window, with confirmation
  metalcloud-cli server firmware schedule-upgrade-batch 123 124 \
    --schedule-timestamp 2024-01-01T10:00:00Z --confirmation-required


```
metalcloud-cli server firmware schedule-upgrade-batch server_id... [flags]
```

### Options

```
      --confirmation-required       Require a confirmation before the scheduled upgrade runs.
  -h, --help                        help for schedule-upgrade-batch
      --schedule-timestamp string   When the firmware upgrade should run (RFC 3339 timestamp).
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

* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management

