## metalcloud-cli network-device-configuration-template config-example

Generate example configuration template for network device configuration template

### Synopsis

Generate an example JSON configuration template that can be used to create
or update network device configuration templates. This template includes all available configuration
options with example values and documentation.

Preparation and configuration fields need to be base64 encoded when submitted.

The generated template can be saved to a file and modified as needed for actual
template configuration.

Examples:
  # Display example configuration
  metalcloud-cli network-device-configuration-template config-example -f json

  # Save example to file
  metalcloud-cli network-device-configuration-template config-example -f json > template.json

```
metalcloud-cli network-device-configuration-template config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

* [metalcloud-cli network-device-configuration-template](metalcloud-cli_network-device-configuration-template.md)	 - Manage network devices configuration templates

