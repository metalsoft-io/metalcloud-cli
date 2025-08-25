## metalcloud-cli role update

Update an existing role's permissions or metadata

### Synopsis

Update an existing role's permissions or metadata in the MetalCloud platform.

This command modifies an existing role based on configuration provided through a JSON file
or piped input. You can update the role's label, description, and permissions.
System roles cannot be updated.

Arguments:
  role_name    Name of the role to update

Configuration fields (all optional):
  label         New human-readable name for the role
  description   New description of the role's purpose
  permissions   New array of permission strings (replaces existing permissions)

Flags:
  --config-source string   Required. Source of the role configuration.
                          Can be 'pipe' for stdin input or path to a JSON file.

Configuration format (JSON):
{
  "label": "Updated Admin Role",
  "description": "Updated administrative role with modified access",
  "permissions": ["roles:read", "roles:write", "users:read"]
}

Note: When updating permissions, the provided array completely replaces the existing
permissions. To add permissions, include both existing and new permissions in the array.

Examples:
  # Update role from JSON file
  metalcloud-cli role update custom-admin --config-source role-update.json

  # Update role using piped JSON
  echo '{"description": "Updated role description"}' | metalcloud-cli role update my-role --config-source pipe

  # Update role permissions and label
  metalcloud-cli roles edit editor-role --config-source /path/to/updated-role.json

```
metalcloud-cli role update <role_name> [flags]
```

### Options

```
      --config-source string   Source of the role configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

