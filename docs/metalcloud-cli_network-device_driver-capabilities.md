## metalcloud-cli network-device driver-capabilities

List the onboarding actions supported by each driver

### Synopsis

List, for each network device driver, the onboarding actions it supports
(for example 'start-registration' or 'mark-install-ready').

Optional Flags:
  --search   Restrict the listing to the drivers matching this text

Examples:
  # List the capabilities of every driver
  metalcloud-cli network-device driver-capabilities

  # List the SONiC driver capabilities
  metalcloud-cli network-device driver-capabilities --search sonic

```
metalcloud-cli network-device driver-capabilities [flags]
```

### Options

```
  -h, --help            help for driver-capabilities
      --search string   Restrict the listing to the drivers matching this text.
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

