## metalcloud-cli device-config-template create

Create a new device configuration template

### Synopsis

Create a new device configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration.

Examples:
  metalcloud-cli device-config-template create --config-source template.json
  cat template.json | metalcloud-cli dct create --config-source pipe

```
metalcloud-cli device-config-template create [flags]
```

### Options

```
      --config-source string   Source of the new template configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli device-config-template](metalcloud-cli_device-config-template.md)	 - Manage device configuration templates and profiles

