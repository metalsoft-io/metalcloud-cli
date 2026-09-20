## metalcloud-cli container-instance config

Show the pending configuration of a container instance

### Synopsis

Show the pending configuration of a container instance.

The configuration holds the changes that will be applied at the next deploy.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Examples:
  # Show the configuration of container instance 5678
  metalcloud-cli container-instance config my-infra 5678

  # Using the alias
  metalcloud-cli ci get-config 1234 5678

```
metalcloud-cli container-instance config infrastructure_id_or_label container_instance_id [flags]
```

### Options

```
  -h, --help   help for config
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

