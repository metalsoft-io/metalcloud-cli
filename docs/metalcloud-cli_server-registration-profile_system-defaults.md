## metalcloud-cli server-registration-profile system-defaults

Show the built-in server registration settings

### Synopsis

Show the built-in server registration settings that are used when no server
registration profile applies.

These settings are a useful starting point for the "settings" section of a new
server registration profile.

Examples:
  # Show the system default registration settings
  metalcloud-cli server-registration-profile system-defaults


```
metalcloud-cli server-registration-profile system-defaults [flags]
```

### Options

```
  -h, --help   help for system-defaults
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

