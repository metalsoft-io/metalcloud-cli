## metalcloud-cli server-instance-group list

List all server instance groups in an infrastructure

### Synopsis

List all server instance groups in an infrastructure.

This command displays all server instance groups within a specified infrastructure,
showing their configuration details including ID, label, status, and timestamps.

Arguments:
  infrastructure_id_or_label  The infrastructure ID (numeric) or label (string) to list groups from

Examples:
  # List all instance groups in infrastructure with ID 1234
  metalcloud-cli server-instance-group list 1234

  # List all instance groups in infrastructure with label "prod-env"
  metalcloud-cli server-instance-group list prod-env

  # Using alias
  metalcloud-cli ig ls 1234

```
metalcloud-cli server-instance-group list infrastructure_id_or_label [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures

