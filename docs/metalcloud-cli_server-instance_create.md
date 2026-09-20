## metalcloud-cli server-instance create

Create a server instance in an infrastructure

### Synopsis

Create a new server instance in an infrastructure.

The instance can be described either with a configuration file (--config-source) or
with individual flags (--label, --group-id, ...).

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure

Required Flags (one of):
  --config-source string   Source of the new server instance configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Label of the new server instance

Optional Flags (when not using --config-source):
  --group-id int           ID of the server instance group the instance belongs to
  --server-type-id int     ID of the server type to allocate
  --hostname string        Custom hostname (subdomain) of the instance
  --os-template-id int     ID of the OS template to deploy
  --tags strings           Tags of the new server instance

Examples:
  metalcloud-cli server-instance create 1234 --label web-01 --server-type-id 100
  metalcloud-cli server-instance create prod-env --config-source server-instance.json
  cat server-instance.yaml | metalcloud-cli inst new prod-env --config-source pipe

```
metalcloud-cli server-instance create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the new server instance configuration. Can be 'pipe' or path to a JSON/YAML file.
      --group-id int           ID of the server instance group the instance belongs to.
  -h, --help                   help for create
      --hostname string        Custom hostname (subdomain) of the instance.
      --label string           Label of the new server instance.
      --os-template-id int     ID of the OS template to deploy.
      --server-type-id int     ID of the server type to allocate.
      --tags strings           Tags of the new server instance.
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

