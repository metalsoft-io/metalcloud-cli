## metalcloud-cli server firmware upgrade-batch

Upgrade the firmware of several servers at once

### Synopsis

Upgrade the firmware of several servers at once.

A firmware upgrade job is started for each of the given servers. The result lists
the job started for every server that could be upgraded and the error reported for
every server that could not.

Required Arguments:
  server_id...           The IDs of the servers to upgrade

Examples:
  # Upgrade the firmware of servers 123, 124 and 125
  metalcloud-cli server firmware upgrade-batch 123 124 125


```
metalcloud-cli server firmware upgrade-batch server_id... [flags]
```

### Options

```
  -h, --help   help for upgrade-batch
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

