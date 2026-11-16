## metalcloud-cli user api-token create

Create an API token for the current user

### Synopsis

Create an API token for the currently authenticated user.

The token value is printed once; store it right away.

Flags:
  --name                  Display name of the token (required unless --config-source is used)
  --label                 Unique (per user) label; derived from the name when omitted
  --description           Free-text description
  --scope                 account (default) or system; system needs the global_token_write permission
  --expires               Expiry as an RFC 3339 timestamp or a YYYY-MM-DD date (UTC);
                          defaults to the controller's default token validity
  --permission            Permission label granted to the token (repeatable or comma-separated)
  --config-source         Full CreateApiToken body (JSON/YAML file path or 'pipe'),
                          needed to grant policy entries

Examples:
  metalcloud-cli user api-token create --name "CI pipeline" --expires 2027-06-30 --permission infrastructures_read
  metalcloud-cli user api-token create --config-source token.json

```
metalcloud-cli user api-token create [flags]
```

### Options

```
      --config-source string   Source of the CreateApiToken body. Can be 'pipe' or path to a JSON/YAML file.
      --description string     Free-text description of the token.
      --expires string         Expiry as an RFC 3339 timestamp or YYYY-MM-DD date (UTC).
  -h, --help                   help for create
      --label string           Unique (per user) label of the token.
      --name string            Display name of the token.
      --permission strings     Permission label granted to the token (repeatable).
      --scope string           Token scope: account (default) or system.
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

* [metalcloud-cli user api-token](metalcloud-cli_user_api-token.md)	 - Manage the current user's API tokens

