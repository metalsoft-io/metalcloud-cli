## metalcloud-cli vm-type update

Update an existing VM type configuration

### Synopsis

Update an existing VM type in the MetalCloud platform using a configuration file or piped input.

The configuration must be provided in JSON format. Only the fields you want to update need to be included.
You can modify any of the following fields:
- name: VM type name
- displayName: Display name for the VM type
- label: Label for the VM type
- cpuCores: Number of CPU cores
- ramGB: Amount of RAM in gigabytes
- isExperimental: Whether the VM type is experimental (0 or 1)
- forUnmanagedVMsOnly: Whether restricted to unmanaged VMs (0 or 1)
- tags: Array of tags

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to update

REQUIRED FLAGS:
  --config-source string    Source of the VM type configuration (required)
                           Can be 'pipe' for stdin or path to a JSON file

EXAMPLES:
  # Update VM type from file
  metalcloud vm-type update 123 --config-source vm-type-update.json
  
  # Update VM type from stdin
  echo '{"cpuCores":8,"ramGB":16}' | metalcloud vm-type update 123 --config-source pipe
  
  # Generate example config and update VM type
  metalcloud vm-type config-example > config.json
  # Edit config.json as needed
  metalcloud vm-type update 123 --config-source config.json

```
metalcloud-cli vm-type update vm_type_id [flags]
```

### Options

```
      --config-source string   Source of the VM type update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli vm-type](metalcloud-cli_vm-type.md)	 - Manage VM types and configurations

