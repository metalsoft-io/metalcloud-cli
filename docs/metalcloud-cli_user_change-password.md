## metalcloud-cli user change-password

Change the password of the current user

### Synopsis

Change the password of the user owning the API key in use.

The new password can be supplied through flags or through a JSON/YAML payload. The
current password is required unless the platform policy allows changing it without.

Required Flags (when not using --config-source):
  --new-password          The new password

Optional Flags:
  --current-password      The password currently in use
  --config-source         Source of the password change payload (JSON/YAML file or 'pipe')

Configuration File Format (JSON):
  {
    "newPassword": "newSecret123",
    "oldPassword": "oldSecret123"
  }

Examples:
  metalcloud-cli user change-password --current-password oldSecret123 --new-password newSecret123
  echo '{"newPassword":"newSecret123","oldPassword":"oldSecret123"}' | metalcloud-cli user change-password --config-source pipe

```
metalcloud-cli user change-password [flags]
```

### Options

```
      --config-source string      Source of the password change payload. Can be 'pipe' or path to a JSON/YAML file.
      --current-password string   The password currently in use.
  -h, --help                      help for change-password
      --new-password string       The new password of the current user.
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

* [metalcloud-cli user](metalcloud-cli_user.md)	 - Manage user accounts and their properties

