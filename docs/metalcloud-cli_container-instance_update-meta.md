## metalcloud-cli container-instance update-meta

Update the metadata of a container instance

### Synopsis

Update the metadata (tags) of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Required Flags:
  --config-source string  Source of the container instance metadata updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the tags from a file
  metalcloud-cli container-instance update-meta my-infra 5678 --config-source meta.json

  # Update the tags from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli ci edit-meta 1234 5678 --config-source pipe

```
metalcloud-cli container-instance update-meta infrastructure_id_or_label container_instance_id [flags]
```

### Options

```
      --config-source string   Source of the container instance metadata updates. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-meta
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

