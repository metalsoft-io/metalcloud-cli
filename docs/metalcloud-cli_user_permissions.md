## metalcloud-cli user permissions

List the permissions of the current user

### Synopsis

List the permissions of the roles assigned to the user the command runs as: the owner
of the API token in use, or the user signed in with 'login'.

The list is the same whatever the credential allows, and the credential can allow less:
an API token with the Account scope or other restrictions is refused operations that these
permissions cover (403 Forbidden), and a session signed in with 'login' is refused a few
account-security operations, such as managing API tokens.

Examples:
  metalcloud-cli user permissions
  metalcloud-cli user permissions -f json

```
metalcloud-cli user permissions [flags]
```

### Options

```
  -h, --help   help for permissions
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API token (msat_...), or a legacy API key (<user id>:<key>), which works only until it expires
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

