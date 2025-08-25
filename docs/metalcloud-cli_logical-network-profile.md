## metalcloud-cli logical-network-profile

Manage logical network profiles for network configuration templates

### Synopsis

Manage logical network profiles which define network configuration templates
that can be applied to infrastructure deployments. These profiles contain
network settings, routing rules, and connectivity configurations.

Available commands:
  list          List all logical network profiles with filtering options
  get           Get detailed information about a specific profile
  create        Create a new logical network profile from configuration
  update        Update an existing logical network profile
  delete        Delete a logical network profile
  config-example Get example configuration for a specific profile kind

Examples:
  metalcloud-cli logical-network-profile list
  metalcloud-cli lnp get 12345
  metalcloud-cli network-profile create --config-source profile.json

### Options

```
  -h, --help   help for logical-network-profile
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli logical-network-profile config-example](metalcloud-cli_logical-network-profile_config-example.md)	 - Get example configuration for a specific profile kind
* [metalcloud-cli logical-network-profile create](metalcloud-cli_logical-network-profile_create.md)	 - Create a new logical network profile from configuration
* [metalcloud-cli logical-network-profile delete](metalcloud-cli_logical-network-profile_delete.md)	 - Delete a logical network profile
* [metalcloud-cli logical-network-profile get](metalcloud-cli_logical-network-profile_get.md)	 - Get detailed information about a logical network profile
* [metalcloud-cli logical-network-profile list](metalcloud-cli_logical-network-profile_list.md)	 - List logical network profiles with optional filtering
* [metalcloud-cli logical-network-profile update](metalcloud-cli_logical-network-profile_update.md)	 - Update an existing logical network profile

