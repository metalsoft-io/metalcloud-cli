## metalcloud-cli server firmware update-component

Update firmware component settings

### Synopsis

Update firmware component settings.

This command updates configuration settings for a specific firmware component
using a JSON configuration file or piped JSON data.

Required Arguments:
  server_id              The ID of the server
  component_id           The ID of the firmware component to update

Required Flags:
  --config-source        Source of the component update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update component using JSON configuration file
  metalcloud-cli server firmware update-component 123 456 --config-source ./component-config.json

  # Update component using piped JSON configuration
  echo '{"setting": "value"}' | metalcloud-cli server firmware update-component 123 456 --config-source pipe


```
metalcloud-cli server firmware update-component server_id component_id [flags]
```

### Options

```
      --config-source string   Source of the component update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-component
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

* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management

