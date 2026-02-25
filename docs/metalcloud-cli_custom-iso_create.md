## metalcloud-cli custom-iso create

Create a new custom ISO from configuration

### Synopsis

Create a new custom ISO image using command-line flags or a JSON configuration file.

You can provide the custom ISO configuration either via command-line flags or by
specifying a configuration source using the --config-source flag. The configuration
source can be a path to a JSON file or 'pipe' to read from standard input.

If --config-source is not provided, you must specify at least --label and --access-url.

Required Flags (when not using --config-source):
  --label             Label for the custom ISO
  --access-url        URL where the ISO file can be accessed

Optional Flags:
  --config-source     Source of configuration (JSON file path or 'pipe')
  --name              Display name for the custom ISO

Required permissions:
  - custom_iso:write

Examples:
  # Create using command line flags
  metalcloud-cli custom-iso create --label my-iso --access-url http://example.com/my.iso

  # Create with optional name
  metalcloud-cli custom-iso create --label my-iso --access-url http://example.com/my.iso --name "My Custom ISO"

  # Create custom ISO from a JSON file
  metalcloud-cli custom-iso create --config-source config.json

  # Create custom ISO from piped JSON
  cat config.json | metalcloud-cli custom-iso create --config-source pipe

```
metalcloud-cli custom-iso create [flags]
```

### Options

```
      --access-url string      URL where the ISO file can be accessed
      --config-source string   Source of the new custom ISO configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
      --label string           Label for the custom ISO
      --name string            Display name for the custom ISO
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

* [metalcloud-cli custom-iso](metalcloud-cli_custom-iso.md)	 - Manage custom ISO images for server provisioning

