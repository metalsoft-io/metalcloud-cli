## metalcloud-cli infrastructure update-metadata

Update infrastructure name, description and tags

### Synopsis

Update the metadata of an infrastructure: its name, its description and its tags.

The metadata is not part of the infrastructure configuration, so it is updated
independently of 'infrastructure update' and is applied immediately (no deploy is
needed). The API does not use an entity tag for this endpoint, so no revision is sent.

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Required Flags:
  --config-source   Source of the metadata configuration. Can be 'pipe' or path to a JSON/YAML file.

The configuration accepts the following properties:
  name         (required) the display name of the infrastructure
  description  (optional) a free-text description
  tags         (optional) a list of tags

Examples:
  # Update the metadata from a file
  metalcloud-cli infrastructure update-metadata my-infra --config-source ./meta.json

  # Update the metadata from piped input
  echo '{"name":"Production","description":"prod stack","tags":["prod"]}' | metalcloud-cli infrastructure update-metadata 123 --config-source pipe

```
metalcloud-cli infrastructure update-metadata infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the infrastructure metadata configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-metadata
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

* [metalcloud-cli infrastructure](metalcloud-cli_infrastructure.md)	 - Manage infrastructure resources and configurations

