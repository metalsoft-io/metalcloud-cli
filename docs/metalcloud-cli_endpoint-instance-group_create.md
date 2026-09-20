## metalcloud-cli endpoint-instance-group create

Create an endpoint instance group in an infrastructure

### Synopsis

Create a new endpoint instance group in an infrastructure. The group starts empty.

The group can be described either with a configuration file (--config-source) or with
individual flags (--label, --endpoint-group-name, ...).

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Required Flags (one of):
  --config-source string        Source of the new group configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string                Label of the new endpoint instance group

Optional Flags (when not using --config-source):
  --endpoint-group-name string  Name of the endpoint group deployed by this group
  --extension-instance-id int   ID of the extension instance backing the group
  --hostname string             Custom hostname for the DNS load balancing record
  --resource-pool-id int        ID of the resource pool assigned to the group
  --tags strings                Tags of the new group

Examples:
  metalcloud-cli endpoint-instance-group create 1234 --label web-endpoints
  metalcloud-cli endpoint-instance-group create prod-env --config-source group.json
  cat group.yaml | metalcloud-cli eig new prod-env --config-source pipe

```
metalcloud-cli endpoint-instance-group create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string         Source of the new group configuration. Can be 'pipe' or path to a JSON/YAML file.
      --endpoint-group-name string   Name of the endpoint group deployed by this group.
      --extension-instance-id int    ID of the extension instance backing the group.
  -h, --help                         help for create
      --hostname string              Custom hostname for the DNS load balancing record.
      --label string                 Label of the new endpoint instance group.
      --resource-pool-id int         ID of the resource pool assigned to the group.
      --tags strings                 Tags of the new group.
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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management

