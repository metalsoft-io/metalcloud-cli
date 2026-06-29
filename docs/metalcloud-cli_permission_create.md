## metalcloud-cli permission create

Create a new permission

### Synopsis

Create a new permission.

You must provide the permission configuration using the --config-source flag.
The configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags:
  --config-source       Source of the permission configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Create using JSON configuration file
  metalcloud-cli permission create --config-source ./permission.json

  # Create using piped JSON configuration
  cat permission.json | metalcloud-cli permission create --config-source pipe


```
metalcloud-cli permission create [flags]
```

### Options

```
      --config-source string   Source of the new permission configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli permission](metalcloud-cli_permission.md)	 - Permission management

