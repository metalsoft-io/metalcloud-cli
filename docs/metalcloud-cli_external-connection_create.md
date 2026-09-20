## metalcloud-cli external-connection create

Create an external connection

### Synopsis

Create a new external connection.

The external connection can be described either with a configuration file
(--config-source) or with individual flags (--label, --name, --fabric).

Required Flags (one of):
  --config-source string   Source of the external connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Label of the new external connection (used together with the flags below)

Optional Flags (when not using --config-source):
  --name string            Name of the new external connection (required with --label)
  --fabric string          ID or name of the fabric the connection belongs to (required with --label)
  --interface-ids strings  Network device interface IDs to attach. Repeatable or comma-separated.

Examples:
  metalcloud external-connection create --config-source external-connection.json
  cat ext-conn.yaml | metalcloud ext-conn create --config-source pipe
  metalcloud ext-conn create --label dc1-ext --name "DC1 external" --fabric dc1-fabric --interface-ids 101,102

```
metalcloud-cli external-connection create [flags]
```

### Options

```
      --config-source string    Source of the new external connection configuration. Can be 'pipe' or path to a JSON/YAML file.
      --fabric string           ID or name of the fabric the external connection belongs to.
  -h, --help                    help for create
      --interface-ids strings   Network device interface IDs to attach to the new external connection.
      --label string            Label of the new external connection.
      --name string             Name of the new external connection.
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

* [metalcloud-cli external-connection](metalcloud-cli_external-connection.md)	 - External connection management

