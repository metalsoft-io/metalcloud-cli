## metalcloud-cli server replacement-candidates

List the servers that can replace a server in its instance

### Synopsis

List the available servers that can replace the given server in its instance.

Eligible candidates are listed first, best first. Ineligible candidates are listed
with the reasons they were skipped. Use json or yaml output to also see the
requirements (RAM, cores, disks, GPUs, network links) the candidates were checked against.

Required Arguments:
  server_id              The ID of the server to replace

Examples:
  # List replacement candidates for server 123
  metalcloud-cli server replacement-candidates 123

  # Include the replacement requirements
  metalcloud-cli server replacement-candidates 123 -f json


```
metalcloud-cli server replacement-candidates server_id [flags]
```

### Options

```
  -h, --help   help for replacement-candidates
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

