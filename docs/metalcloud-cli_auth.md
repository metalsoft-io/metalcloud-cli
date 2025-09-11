## metalcloud-cli auth

Manage authentication settings

### Synopsis

Manage authentication settings for the MetalCloud platform.

This command provides subcommands to configure and manage various authentication
methods including LDAP integration and user role mappings.

Available authentication methods:
- LDAP: Configure LDAP server integration and group-to-role mappings

Examples:
  # List all available auth subcommands
  metalcloud-cli auth --help
  
  # Work with LDAP authentication
  metalcloud-cli auth ldap --help

### Options

```
  -h, --help   help for auth
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
* [metalcloud-cli auth ldap](metalcloud-cli_auth_ldap.md)	 - Manage LDAP authentication settings

