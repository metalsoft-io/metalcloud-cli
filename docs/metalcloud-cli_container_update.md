## metalcloud-cli container update

Update a container

### Synopsis

Update the comments and tags of a container.

Required Arguments:
  container_id  The numeric ID of the container to update

Required Flags:
  --config-source string  Source of the container updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update container 100 from a file
  metalcloud-cli container update 100 --config-source updates.json

  # Update the tags of container 100 from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli containers edit 100 --config-source pipe

```
metalcloud-cli container update container_id [flags]
```

### Options

```
      --config-source string   Source of the container updates. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

* [metalcloud-cli container](metalcloud-cli_container.md)	 - Manage provisioned containers

