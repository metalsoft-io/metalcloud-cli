## metalcloud-cli container-instance-group list

List the container instance groups of an infrastructure

### Synopsis

List all container instance groups of an infrastructure.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List the container instance groups of infrastructure 1234
  metalcloud-cli container-instance-group list 1234

  # List them by infrastructure label
  metalcloud-cli cig ls my-infra

```
metalcloud-cli container-instance-group list infrastructure_id_or_label [flags]
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

* [metalcloud-cli container-instance-group](metalcloud-cli_container-instance-group.md)	 - Manage container instance groups within infrastructures

