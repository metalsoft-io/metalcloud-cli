## metalcloud-cli vm-pool cluster-host-interface-network-device config-example

Display a network device assignment configuration example

### Synopsis

Display a sample payload for
'vm-pool cluster-host-interface-network-device add --config-source'.

Examples:
  # Display the example
  metalcloud-cli vm-pool chind config-example

  # Save it for editing
  metalcloud-cli vm-pool chind config-example > assignment.json

```
metalcloud-cli vm-pool cluster-host-interface-network-device config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

* [metalcloud-cli vm-pool cluster-host-interface-network-device](metalcloud-cli_vm-pool_cluster-host-interface-network-device.md)	 - Manage the network device assignments of a cluster host interface

