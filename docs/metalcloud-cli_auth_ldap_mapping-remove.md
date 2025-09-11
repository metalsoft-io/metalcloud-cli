## metalcloud-cli auth ldap mapping-remove

Delete LDAP group mapping

### Synopsis

Delete an existing LDAP group-to-role mapping.

This command removes a mapping between an LDAP group and a MetalCloud role.
Once removed, users who belong only to the specified LDAP group will no longer
receive the associated role when authenticating via LDAP.

This operation is irreversible. If you need the mapping again, you will need
to recreate it using the mapping-add command.

Examples:
  # Remove mapping for "Guests" group
  metalcloud-cli auth ldap mapping-remove "Guests"
  
  # Remove mapping for "Contractors" group using short alias
  metalcloud-cli auth ldap map-rm "Contractors"
  
  # Remove mapping for group with spaces in name
  metalcloud-cli auth ldap mapping-remove "External Users"

```
metalcloud-cli auth ldap mapping-remove group_name [flags]
```

### Options

```
  -h, --help   help for mapping-remove
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

