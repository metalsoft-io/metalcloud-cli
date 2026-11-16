## metalcloud-cli user api-token

Manage the current user's API tokens

### Synopsis

Manage the API tokens of the currently authenticated user.

A user can own several API tokens, each with its own label, scope, permissions
and expiry. The token value is shown only when the token is created or
regenerated; store it right away.

### Options

```
  -h, --help   help for api-token
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
* [metalcloud-cli user api-token create](metalcloud-cli_user_api-token_create.md)	 - Create an API token for the current user
* [metalcloud-cli user api-token delete](metalcloud-cli_user_api-token_delete.md)	 - Delete an API token of the current user
* [metalcloud-cli user api-token get](metalcloud-cli_user_api-token_get.md)	 - Get an API token of the current user
* [metalcloud-cli user api-token list](metalcloud-cli_user_api-token_list.md)	 - List the current user's API tokens
* [metalcloud-cli user api-token regenerate](metalcloud-cli_user_api-token_regenerate.md)	 - Regenerate an API token of the current user
* [metalcloud-cli user api-token settings](metalcloud-cli_user_api-token_settings.md)	 - Show the API token limits of the controller

