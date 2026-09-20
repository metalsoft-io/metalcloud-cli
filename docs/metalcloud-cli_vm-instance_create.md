## metalcloud-cli vm-instance create

Create a new VM instance

### Synopsis

Create a new VM instance in an infrastructure.

The instance can be described either by a complete configuration document
(--config-source) or by individual flags. The new instance joins an existing
VM instance group and is materialised when the infrastructure is deployed.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string  Source of the new VM instance configuration.
                          Can be 'pipe' or path to a JSON/YAML file.
  --vm-type-id string     The VM type of the new VM instance.
  --group-id string       The VM instance group of the new VM instance.

Optional Flags:
  --disk-size-gb string  Disk size in GB of the new VM instance.
  --tags strings         Tags of the new VM instance.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --vm-type-id is required; --group-id is required when
  --vm-type-id is used.

Examples:
  # Create a VM instance from flags
  metalcloud-cli vm-instance create my-infra --vm-type-id 5 --group-id 67890 --disk-size-gb 40

  # Create a VM instance from a file
  metalcloud-cli vmi new 12345 --config-source vm-instance.json

  # Create a VM instance from stdin
  echo '{"typeId":5,"groupId":67890}' | metalcloud-cli vm create 12345 --config-source pipe

```
metalcloud-cli vm-instance create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the new VM instance configuration. Can be 'pipe' or path to a JSON/YAML file.
      --disk-size-gb string    Disk size in GB of the new VM instance.
      --group-id string        The VM instance group of the new VM instance.
  -h, --help                   help for create
      --tags strings           Tags of the new VM instance.
      --vm-type-id string      The VM type of the new VM instance.
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

* [metalcloud-cli vm-instance](metalcloud-cli_vm-instance.md)	 - Manage individual VM instances within infrastructures

