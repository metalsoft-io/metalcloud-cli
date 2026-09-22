## metalcloud-cli vm-pool update-cluster-host-config-example

Display a cluster host update configuration example

### Synopsis

Display a sample payload for 'vm-pool update-cluster-host'.

Examples:
  # Display the example
  metalcloud-cli vm-pool update-cluster-host-config-example

  # Save it for editing
  metalcloud-cli vm-pool update-cluster-host-config-example > host.json

```
metalcloud-cli vm-pool update-cluster-host-config-example [flags]
```

### Options

```
  -h, --help   help for update-cluster-host-config-example
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

