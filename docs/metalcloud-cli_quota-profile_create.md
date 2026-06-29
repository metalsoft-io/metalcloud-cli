## metalcloud-cli quota-profile create

Create a new quota profile

### Synopsis

Create a new quota profile in MetalSoft.

The quota profile configuration must be provided using the --config-source flag.
The configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags:
  --config-source       Source of the quota profile configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Create using JSON configuration file
  metalcloud-cli quota-profile create --config-source ./quota-profile.json

  # Create using piped JSON configuration
  cat quota-profile.json | metalcloud-cli quota-profile create --config-source pipe


```
metalcloud-cli quota-profile create [flags]
```

### Options

```
      --config-source string   Source of the new quota profile configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli quota-profile](metalcloud-cli_quota-profile.md)	 - Quota Profile management

