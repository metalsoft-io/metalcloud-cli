## metalcloud-cli device-config-template profile bulk-assign

Bulk-assign a device configuration template to multiple devices

### Synopsis

Bulk-assign a device configuration template to multiple network devices as profiles.

Required Flags:
  --config-source   Source of the request (required)

The request body accepts: deviceConfigurationTemplateId (required), networkFabricId,
networkDeviceIds, lifecycleStage, variables, isEnabled, priority, applyMode, annotations, tags.

```
metalcloud-cli device-config-template profile bulk-assign [flags]
```

### Options

```
      --config-source string   Source of the request. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for bulk-assign
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

