## metalcloud-cli route-domain update-config

Update the global settings of a route domain config

### Synopsis

Update the global settings of a route domain's config.

Only the config's global settings are updated here (autoRouteDistinguisher and
autoRouteTarget); the allocation strategies are managed with the 'route-domain
allocation-strategy' commands. The config's own revision is sent as the
If-Match entity tag, so a concurrent change is rejected instead of being
overwritten.

Required Arguments:
  route_domain_id  The ID of the route domain

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli route-domain update-config 2 --config-source settings.json
  echo '{"autoRouteTarget":true}' | metalcloud-cli route-domain update-config 2 --config-source pipe

```
metalcloud-cli route-domain update-config route_domain_id [flags]
```

### Options

```
      --config-source string   Source of the route domain config updates. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-config
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

