## metalcloud-cli auth ldap mapping-add

Add LDAP group mapping

### Synopsis

Add a new LDAP group-to-role mapping for authentication.

This command creates a new mapping between an LDAP group and a MetalCloud role.
Users who belong to the specified LDAP group will be assigned the specified role
when they authenticate via LDAP.

The priority value determines which role is assigned when a user belongs to multiple
LDAP groups with different mappings. Lower priority numbers take precedence over
higher ones. If two mappings have the same priority, the behavior is undefined.

Examples:
  # Map "Domain Admins" group to admin role with highest priority
  metalcloud-cli auth ldap mapping-add "Domain Admins" --role-name admin --priority 1
  
  # Map "Power Users" group to power-user role  
  metalcloud-cli auth ldap mapping-add "Power Users" --role-name power-user --priority 5
  
  # Map "Developers" group to developer role
  metalcloud-cli auth ldap mapping-add "Developers" --role-name developer --priority 10

```
metalcloud-cli auth ldap mapping-add group_name [flags]
```

### Options

```
  -h, --help               help for mapping-add
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

