## metalcloud-cli endpoint create

Create a new endpoint

### Synopsis

Create a new endpoint in MetalSoft.

You can specify the endpoint configuration either by providing individual flags (--site-id, --name, --label, --external-id) 
or by supplying a configuration file or piped JSON/YAML using --config-source. 
When using --config-source, the file or piped content must contain a valid endpoint configuration in JSON or YAML format.

Required flags (when not using --config-source):
  --site-id     Site ID where the endpoint will be created
  --name        Name of the endpoint
  --label       Label of the endpoint

Optional flags:
  --external-id string       External ID of the endpoint
  --config-source string     Source of configuration (file path or 'pipe')

Flag dependencies:
  - When using --config-source, all other flags are ignored
  - When not using --config-source, --site-id, --name, and --label are required together

Examples:
  metalcloud-cli endpoint create --site-id 1 --name "api-endpoint" --label "API Endpoint"
  metalcloud-cli endpoint create --site-id 1 --name "api-endpoint" --label "API Endpoint" --external-id "ext-001"
  metalcloud-cli endpoint create --config-source ./endpoint.json
  cat endpoint.yaml | metalcloud-cli endpoint create --config-source pipe

```
metalcloud-cli endpoint create [flags]
```

### Options

```
      --config-source string   Source of the new endpoint configuration. Can be 'pipe' or path to a JSON file.
      --external-id string     The external ID of the endpoint.
  -h, --help                   help for create
      --label string           The label of the endpoint.
      --name string            The name of the endpoint.
      --site-id int            The site ID where the endpoint will be created.
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

