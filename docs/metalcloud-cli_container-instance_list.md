## metalcloud-cli container-instance list

List the container instances of an infrastructure

### Synopsis

List all container instances of an infrastructure.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List the container instances of infrastructure 1234
  metalcloud-cli container-instance list 1234

  # List them by infrastructure label
  metalcloud-cli ci ls my-infra

```
metalcloud-cli container-instance list infrastructure_id_or_label [flags]
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

* [metalcloud-cli container-instance](metalcloud-cli_container-instance.md)	 - Manage container instances within infrastructures

