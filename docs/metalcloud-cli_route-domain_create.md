## metalcloud-cli route-domain create

Create a new route domain

### Synopsis

Create a new route domain (tenant VRF) from a JSON/YAML configuration.

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

The configuration must include the route domain kind (evpn_l3vpn | mpls_l3vpn |
vrf_lite) and its VRF allocation strategy; an l3evpn tenant VRF also carries an
L3VNI allocation strategy. Run 'route-domain config-example' for a template.

Examples:
  metalcloud-cli route-domain config-example > route-domain.yaml
  metalcloud-cli route-domain create --config-source route-domain.yaml
  cat route-domain.yaml | metalcloud-cli route-domain create --config-source pipe

```
metalcloud-cli route-domain create [flags]
```

### Options

```
      --config-source string   Source of the new route domain configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
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

* [metalcloud-cli route-domain](metalcloud-cli_route-domain.md)	 - Manage route domains (tenant VRFs)

