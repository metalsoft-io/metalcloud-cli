## metalcloud-cli license add

Install a license on the system

### Synopsis

Install a license on the system.

The license must be provided as a Base64-encoded signed document, exactly as
returned by MetalSoft. It is forwarded verbatim to preserve the original signed
bytes.

Required Flags:
  --source    Source of the license document. Can be 'pipe' to read from stdin
              or a path to a file containing the Base64-encoded license.

Examples:
  # Install a license from a file
  metalcloud-cli license add --source ./license.txt

  # Install a license from stdin
  cat license.txt | metalcloud-cli license add --source pipe

```
metalcloud-cli license add [flags]
```

### Options

```
  -h, --help            help for add
      --source string   Source of the license document. Can be 'pipe' or a path to a file containing the Base64-encoded license.
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

