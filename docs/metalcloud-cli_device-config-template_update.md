## metalcloud-cli device-config-template update

Update an existing device configuration template

### Synopsis

Update an existing device configuration template using JSON configuration provided via
file or pipe. Only the specified fields will be updated; other configuration remains unchanged.

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  metalcloud-cli device-config-template update 12345 --config-source updates.json
  cat updates.json | metalcloud-cli dct update 12345 --config-source pipe

```
metalcloud-cli device-config-template update <device_configuration_template_id> [flags]
```

### Options

```
      --config-source string   Source of the template updates. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli device-config-template](metalcloud-cli_device-config-template.md)	 - Manage device configuration templates and profiles

