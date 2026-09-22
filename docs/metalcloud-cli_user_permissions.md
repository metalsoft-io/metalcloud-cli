## metalcloud-cli user permissions

List the permissions of the current user

### Synopsis

List the permissions of the user owning the API key in use.

The permissions come from the roles assigned to the user and determine which API
operations, and therefore which CLI commands, are available.

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

