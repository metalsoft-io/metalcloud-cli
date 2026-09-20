## metalcloud-cli network-device run-extension

Run an extension against a network device

### Synopsis

Run an extension against a network device. The operation is asynchronous and
returns the job that carries it out.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the extension invocation. Both 'extensionId' and
                      'inputArguments' must be present; pass an empty object (or
                      null) when the extension takes no arguments.
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Run an extension described in a file
  metalcloud-cli network-device run-extension 12345 --config-source extension.json

  # Run extension 7 with no arguments
  echo '{"extensionId":7,"inputArguments":{}}' | metalcloud-cli network-device run-extension 12345 --config-source pipe

```
metalcloud-cli network-device run-extension <network_device_id> [flags]
```

### Options

```
      --config-source string   Source of the extension invocation. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for run-extension
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure

