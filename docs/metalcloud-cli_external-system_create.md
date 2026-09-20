## metalcloud-cli external-system create

Create an external system

### Synopsis

Create a new external system.

The external system can be described either with a configuration document
(--config-source) or with individual flags (--label, --name, --annotations).

Required Flags (one of):
  --config-source string   Source of the new external system configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Unique label of the new external system (lowercase letters, digits and dashes)

Optional Flags (when not using --config-source):
  --name string            Display name. Defaults to the label.
  --annotations string     Annotations as a JSON object, e.g. '{"owner":"platform-team"}'

Examples:
  metalcloud external-system create --config-source external-system.json
  cat external-system.yaml | metalcloud ext-system create --config-source pipe
  metalcloud external-system create --label my-system --name "My System" --annotations '{"owner":"platform-team"}'

```
metalcloud-cli external-system create [flags]
```

### Options

```
      --annotations string     Annotations of the new external system, as a JSON object.
      --config-source string   Source of the new external system configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
      --label string           Unique label of the new external system.
      --name string            Display name of the new external system.
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

* [metalcloud-cli external-system](metalcloud-cli_external-system.md)	 - External system management

