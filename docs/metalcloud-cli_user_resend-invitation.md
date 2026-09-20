## metalcloud-cli user resend-invitation

Resend the platform invitation to a user

### Synopsis

Resend the platform invitation message to a specific user.

The user receives a new invitation link. Use --redirect-url to control where the user
lands after accepting the invitation.

Required Arguments:
  user_id                 The numeric ID of the user to invite again

Optional Flags:
  --redirect-url          URL the user is redirected to after accepting the invitation

Examples:
  metalcloud-cli user resend-invitation 12345
  metalcloud-cli user resend-invitation 12345 --redirect-url https://metalsoft.io/welcome

```
metalcloud-cli user resend-invitation user_id [flags]
```

### Options

```
  -h, --help                  help for resend-invitation
      --redirect-url string   URL the user is redirected to after accepting the invitation.
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

