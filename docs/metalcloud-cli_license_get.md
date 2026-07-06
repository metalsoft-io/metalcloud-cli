## metalcloud-cli license get

Show the installed license document

### Synopsis

Show the license document currently installed on the system as a Base64-encoded signed blob. In text output the raw document is printed so it can be piped or saved to a file.

```
metalcloud-cli license get [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli license](metalcloud-cli_license.md)	 - Manage the system license

