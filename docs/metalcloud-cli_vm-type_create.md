## metalcloud-cli vm-type create

Create a new VM type from configuration

### Synopsis

Create a new VM type in the MetalCloud platform using a configuration file or piped input.

The configuration must be provided in JSON format and include all required fields:
- name: Unique name for the VM type
- cpuCores: Number of CPU cores
- ramGB: Amount of RAM in gigabytes

Optional fields include displayName, label, isExperimental, forUnmanagedVMsOnly, and tags.

REQUIRED FLAGS:
  --config-source string    Source of the VM type configuration (required)
                           Can be 'pipe' for stdin or path to a JSON file

EXAMPLES:
  # Create VM type from file
  metalcloud vm-type create --config-source vm-type-config.json
  
  # Create VM type from stdin
  echo '{"name":"test-vm","cpuCores":2,"ramGB":4}' | metalcloud vm-type create --config-source pipe
  
  # Generate example config and create VM type
  metalcloud vm-type config-example > config.json
  # Edit config.json as needed
  metalcloud vm-type create --config-source config.json

```
metalcloud-cli vm-type create [flags]
```

### Options

```
      --config-source string   Source of the new VM type configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli vm-type](metalcloud-cli_vm-type.md)	 - Manage VM types and configurations

