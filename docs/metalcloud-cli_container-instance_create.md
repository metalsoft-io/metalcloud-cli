## metalcloud-cli container-instance create

Create a new container instance

### Synopsis

Create a new container instance in an infrastructure.

The instance can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string       Source of the new container instance configuration.
                               Can be 'pipe' or path to a JSON/YAML file.
  --container-type-id string   The container type of the new instance.
  --group-id string            The container instance group of the new instance.

Optional Flags:
  --disk-size-gb string  Disk size in GB. Defaults to the group disk size.
  --tags strings         Tags of the new container instance.

Flag Dependencies:
  --config-source is mutually exclusive with --container-type-id, --group-id,
  --disk-size-gb and --tags. One of --config-source or --container-type-id is
  required; --group-id is required when --container-type-id is used.

Examples:
  # Create a container instance from a file
  metalcloud-cli container-instance create my-infra --config-source instance.json

  # Create a container instance from flags
  metalcloud-cli ci new my-infra --container-type-id 42 --group-id 77 --disk-size-gb 40

```
metalcloud-cli container-instance create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string       Source of the new container instance configuration. Can be 'pipe' or path to a JSON/YAML file.
      --container-type-id string   The container type of the new container instance.
      --disk-size-gb string        Disk size in GB of the new container instance.
      --group-id string            The container instance group of the new container instance.
  -h, --help                       help for create
      --tags strings               Tags of the new container instance.
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

