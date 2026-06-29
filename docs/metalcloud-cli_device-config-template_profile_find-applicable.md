## metalcloud-cli device-config-template profile find-applicable

Find device configuration template profiles applicable to a device or fabric

### Synopsis

Find profiles applicable to a network device or fabric.

Required Flags:
  --config-source   Source of the request (required)

The request body accepts: networkDeviceId, networkFabricId, lifecycleStage, includeDisabled.

```
metalcloud-cli device-config-template profile find-applicable [flags]
```

### Options

```
      --config-source string   Source of the request. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for find-applicable
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

