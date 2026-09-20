## metalcloud-cli vm-pool update

Update a VM pool from a configuration file or pipe

### Synopsis

Update an existing virtual machine pool from a JSON or YAML configuration.

Only the fields present in the configuration are changed; everything else keeps
its current value. The endpoint does not use optimistic concurrency, so no
revision needs to be fetched first.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to update

Required Flags:
  --config-source  Source of the VM pool update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields (all optional):
  description, managementHost, managementPort, certificate, privateKey,
  username, password, inMaintenance, isExperimental, tags, options,
  networkFabricId

Examples:
  # Update from a file
  metalcloud-cli vm-pool update 123 --config-source update.json

  # Put a pool into maintenance mode
  echo '{"inMaintenance": 1}' | metalcloud-cli vm-pool update 123 --config-source pipe

  # Start from the generated example
  metalcloud-cli vm-pool update-config-example > update.json

```
metalcloud-cli vm-pool update vm_pool_id [flags]
```

### Options

```
      --config-source string   Source of the VM pool update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli vm-pool](metalcloud-cli_vm-pool.md)	 - Manage virtual machine pools and their resources

