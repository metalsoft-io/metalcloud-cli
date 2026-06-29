## metalcloud-cli device-config-template render

Render arbitrary device configuration template content

### Synopsis

Render device configuration template content provided inline (not saved) with the given variables.

Required Flags:
  --config-source   Source of the render request (required)
                   Values: 'pipe' for stdin input, or path to JSON file

The request body accepts: templateContent (required), variables, debug.

```
metalcloud-cli device-config-template render [flags]
```

### Options

```
      --config-source string   Source of the render request. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for render
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

* [metalcloud-cli device-config-template](metalcloud-cli_device-config-template.md)	 - Manage device configuration templates and profiles

