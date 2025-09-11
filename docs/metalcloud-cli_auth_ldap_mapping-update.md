## metalcloud-cli auth ldap mapping-update

Update LDAP group mapping

### Synopsis

Update an existing LDAP group-to-role mapping.

This command modifies an existing mapping between an LDAP group and a MetalCloud role.
You can update either the role name, the priority, or both. At least one of these
flags must be provided.

The priority value determines which role is assigned when a user belongs to multiple
LDAP groups with different mappings. Lower priority numbers take precedence over
higher ones.

Examples:
  # Update the role name for "Power Users" group
  metalcloud-cli auth ldap mapping-update "Power Users" --role-name senior-developer
  
  # Update the priority for "Developers" group
  metalcloud-cli auth ldap mapping-update "Developers" --priority 15
  
  # Update both role name and priority
  metalcloud-cli auth ldap mapping-update "Guests" --role-name read-only --priority 20

```
metalcloud-cli auth ldap mapping-update group_name [flags]
```

### Options

```
  -h, --help               help for mapping-update
      --priority int32     Mapping priority. (default 10)
      --role-name string   Role name to map to the LDAP group.
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

