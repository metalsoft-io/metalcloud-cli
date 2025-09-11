## metalcloud-cli auth ldap

Manage LDAP authentication settings

### Synopsis

Manage LDAP authentication settings for the MetalCloud platform.

Configure LDAP server integration and manage group-to-role mappings that determine
user permissions based on their LDAP group memberships.

Available operations:
- List existing LDAP group mappings
- Add new LDAP group mappings
- Update existing LDAP group mappings  
- Remove LDAP group mappings

Examples:
  # List all LDAP group mappings
  metalcloud-cli auth ldap mapping-list
  
  # Add a new LDAP group mapping
  metalcloud-cli auth ldap mapping-add "Domain Admins" --role-name admin --priority 1
  
  # Update an existing mapping
  metalcloud-cli auth ldap mapping-update "Power Users" --role-name power-user --priority 5
  
  # Remove a mapping
  metalcloud-cli auth ldap mapping-remove "Guests"

### Options

```
  -h, --help   help for ldap
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

* [metalcloud-cli auth](metalcloud-cli_auth.md)	 - Manage authentication settings
* [metalcloud-cli auth ldap mapping-add](metalcloud-cli_auth_ldap_mapping-add.md)	 - Add LDAP group mapping
* [metalcloud-cli auth ldap mapping-list](metalcloud-cli_auth_ldap_mapping-list.md)	 - List all LDAP group mappings
* [metalcloud-cli auth ldap mapping-remove](metalcloud-cli_auth_ldap_mapping-remove.md)	 - Delete LDAP group mapping
* [metalcloud-cli auth ldap mapping-update](metalcloud-cli_auth_ldap_mapping-update.md)	 - Update LDAP group mapping

