## metalcloud-cli endpoint create-bulk

Create multiple endpoints in one call

### Synopsis

Create multiple endpoints in a single bulk call.

The configuration must be a JSON/YAML list of endpoint definitions. Each entry
requires at least siteId, name, and label; endpointInterfaces and other fields
are optional.

Each endpoint interface can be specified in one of two ways:
  - by numeric id:    { "networkDeviceInterfaceId": 12345 }
  - by label:         { "networkDevice": "leaf-01", "interface": "swp9s0" }

When specified by label, the network device is resolved by numeric id or by its
identifierString (switch hostname), and the interface by its name (e.g.
"swp9s0"); each device's ports are fetched once and cached. An optional
macAddress may be set on either form.

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file
                    containing a list of endpoints.

Examples:
  metalcloud-cli endpoint create-bulk --config-source endpoints.yaml
  cat endpoints.json | metalcloud-cli endpoint create-bulk --config-source pipe

```
metalcloud-cli endpoint create-bulk [flags]
```

### Options

```
      --config-source string   Source of the endpoints list. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create-bulk
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

* [metalcloud-cli endpoint](metalcloud-cli_endpoint.md)	 - Endpoint management

