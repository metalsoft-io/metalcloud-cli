## metalcloud-cli device-config-template profile render

Render a device configuration template profile for a device

### Synopsis

Render a profile for a given network device.

Required Flags:
  --config-source   Source of the render request (required)

The request body accepts: networkDeviceId (required), extraVariables, debug.

```
metalcloud-cli device-config-template profile render <profile_id> [flags]
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

* [metalcloud-cli device-config-template profile](metalcloud-cli_device-config-template_profile.md)	 - Manage device configuration template profiles

