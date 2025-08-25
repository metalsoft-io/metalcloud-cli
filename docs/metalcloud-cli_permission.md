## metalcloud-cli permission

Manage system permissions and access control

### Synopsis

Manage system permissions and access control.

Permissions define what actions users and roles can perform within the MetalCloud platform.
This command group provides functionality to view and manage the available permissions
in the system.

Available Commands:
  list        List all available permissions in the system

Examples:
  # List all permissions
  metalcloud-cli permission list
  
  # List permissions with short alias
  metalcloud-cli permissions ls

Required Permissions:
  Most permission management operations require the 'roles:read' permission.

### Options

```
  -h, --help   help for permission
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
* [metalcloud-cli permission list](metalcloud-cli_permission_list.md)	 - List all available system permissions

