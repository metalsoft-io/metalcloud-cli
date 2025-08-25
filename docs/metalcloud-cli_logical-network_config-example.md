## metalcloud-cli logical-network config-example

Generate example configuration for a logical network kind

### Synopsis

Generate example configuration templates for different logical network kinds.

This command provides sample JSON configurations that can be used as templates when
creating logical networks. The configuration examples show the structure and required
fields for each network kind.

Arguments:
  kind  The type of logical network for which to generate example configuration
        Supported kinds include: vlan, vxlan, flat, and others

The generated configuration can be used with the 'create' command by saving it to a file
and using the --config-source flag, or by piping it directly.

Examples:
  # Get example configuration for a VLAN network
  metalcloud-cli logical-network config-example vlan

  # Save example to file for editing
  metalcloud-cli logical-network config-example vxlan > network-config.json

  # Use with create command via pipe
  metalcloud-cli logical-network config-example vlan | metalcloud-cli logical-network create vlan --config-source pipe

```
metalcloud-cli logical-network config-example kind [flags]
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

