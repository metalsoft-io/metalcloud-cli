## metalcloud-cli server replace

Replace a server in its instance with another server

### Synopsis

Replace a server in its instance with another server.

The replacement server is put in the instance and the infrastructure is deployed,
which moves the instance onto the new server with the same networking and OS. The
old server becomes defective and should be wiped with "server cleanup".

Use "server replacement-candidates" to find an eligible replacement server.

Required Arguments:
  server_id                  The ID of the server to replace

Required Flags:
  --replacement-server-id    The ID of the server to put in its place

Optional Flags:
  --accept-downgrade         Accept a replacement with less RAM, cores or disk

Examples:
  # Replace server 123 with server 456
  metalcloud-cli server replace 123 --replacement-server-id 456

  # Replace with a smaller server
  metalcloud-cli server replace 123 --replacement-server-id 456 --accept-downgrade


```
metalcloud-cli server replace server_id [flags]
```

### Options

```
      --accept-downgrade               Accept a replacement with less RAM, cores or disk than the instance had
  -h, --help                           help for replace
      --replacement-server-id string   The ID of the server to put in its place
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

