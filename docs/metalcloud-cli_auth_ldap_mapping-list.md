## metalcloud-cli auth ldap mapping-list

List all LDAP group mappings

### Synopsis

List all configured LDAP group mappings and their associated roles.

This command displays all LDAP group-to-role mappings that are currently configured
in the system. Each mapping shows the LDAP group name, the MetalCloud role it maps to,
and the priority of the mapping.

The priority determines which role is assigned when a user belongs to multiple LDAP
groups with different mappings. Lower priority numbers take precedence.

Examples:
  # List all LDAP group mappings
  metalcloud-cli auth ldap mapping-list
  
  # List mappings with short alias
  metalcloud-cli auth ldap map-ls

```
metalcloud-cli auth ldap mapping-list [flags]
```

### Options

```
  -h, --help   help for mapping-list
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

* [metalcloud-cli auth ldap](metalcloud-cli_auth_ldap.md)	 - Manage LDAP authentication settings

