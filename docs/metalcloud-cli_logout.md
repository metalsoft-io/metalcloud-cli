## metalcloud-cli logout

End the saved session and forget it

### Synopsis

End the session saved by 'metalcloud-cli login' on the server and remove it from
$HOME/.metalcloud/sessions.json. The saved session is removed even if the server can't be reached.

Examples:
  metalcloud-cli logout

```
metalcloud-cli logout [flags]
```

### Options

```
  -h, --help   help for logout
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI

