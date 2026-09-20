## metalcloud-cli user initiate-email-change

Start changing the e-mail address of the current user

### Synopsis

Start the e-mail address change flow for the user owning the API key in use.

A verification message is sent to the new address; the change takes effect only after
the link in that message is followed.

Required Flags:
  --email                 The new e-mail address

Optional Flags:
  --redirect-url          URL the user is redirected to after verifying the address

Examples:
  metalcloud-cli user initiate-email-change --email new.address@company.com

```
metalcloud-cli user initiate-email-change [flags]
```

### Options

```
      --email string          The new e-mail address of the current user.
  -h, --help                  help for initiate-email-change
      --redirect-url string   URL the user is redirected to after verifying the e-mail address.
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

