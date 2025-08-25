## metalcloud-cli role delete

Delete a role from the system

### Synopsis

Delete a role from the MetalCloud platform.

This command permanently removes a role from the system. Once deleted, the role
cannot be recovered and any users assigned to this role will lose those permissions.
System roles cannot be deleted.

Arguments:
  role_name    Name of the role to delete

Warning:
  - This operation is irreversible
  - Users assigned to this role will lose the associated permissions
  - System roles cannot be deleted

Examples:
  # Delete a custom role
  metalcloud-cli role delete custom-editor

  # Using the alias
  metalcloud-cli role rm temp-role

  # Using the remove alias
  metalcloud-cli roles remove old-role

```
metalcloud-cli role delete <role_name> [flags]
```

### Options

```
  -h, --help   help for delete
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

* [metalcloud-cli role](metalcloud-cli_role.md)	 - Manage user roles and permissions

