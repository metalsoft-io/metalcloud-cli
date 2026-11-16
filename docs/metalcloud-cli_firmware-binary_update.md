## metalcloud-cli firmware-binary update

Update a firmware binary from a configuration file

### Synopsis

Update a firmware binary from a configuration file.

Only the fields in the configuration change; the others keep their current values. A field set to
null is cleared. Use it to add prerequisites to a binary, or to correct the steps of a bundle.

Required Flags:
  --config-source    Source of the changes (JSON/YAML file path or 'pipe')

Examples:
  # Require a minimum host BIOS and BMC before the bundle is scheduled
  echo '{"prerequisites": [{"platform": "NVIDIA DGXH100", "componentClass": "bios", "minimumVersion": "01.02.05"},
                           {"componentClass": "bmc", "minimumVersion": "23.06.00"}]}' |
    metalcloud-cli firmware-binary update 4796 --config-source pipe

  # Change the severity
  metalcloud-cli fw-binary update 4796 --config-source ./severity.json

```
metalcloud-cli firmware-binary update firmware_binary_id [flags]
```

### Options

```
      --config-source string   Source of the changes. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli firmware-binary](metalcloud-cli_firmware-binary.md)	 - Manage individual firmware binary files and packages

