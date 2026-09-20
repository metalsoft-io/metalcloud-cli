## metalcloud-cli container-instance-group create

Create a new container instance group

### Synopsis

Create a new container instance group in an infrastructure.

The group can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string      Source of the new group configuration.
                              Can be 'pipe' or path to a JSON/YAML file.
  --container-type-id string  The container type of the group instances.
  --os-template-id string     The OS template of the group instances.
  --vm-pool-id string         The VM pool the group is provisioned on.
  --disk-size-gb string       Disk size in GB of each group instance.

Optional Flags:
  --instance-count string  Number of instances in the group.
  --tags strings           Tags of the new group.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --container-type-id is required; --os-template-id,
  --vm-pool-id and --disk-size-gb are required when --container-type-id is used.

Examples:
  # Create a group from a file
  metalcloud-cli container-instance-group create my-infra --config-source group.json

  # Create a group from flags
  metalcloud-cli cig new my-infra --container-type-id 42 --os-template-id 7 --vm-pool-id 3 --disk-size-gb 40 --instance-count 2

```
metalcloud-cli container-instance-group create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string       Source of the new container instance group configuration. Can be 'pipe' or path to a JSON/YAML file.
      --container-type-id string   The container type of the group instances.
      --disk-size-gb string        Disk size in GB of each group instance.
  -h, --help                       help for create
      --instance-count string      Number of instances in the group.
      --os-template-id string      The OS template of the group instances.
      --tags strings               Tags of the new container instance group.
      --vm-pool-id string          The VM pool the group is provisioned on.
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

