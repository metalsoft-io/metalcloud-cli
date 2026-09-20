## metalcloud-cli site registry-urls

List the container registry URLs available for site agents

### Synopsis

List the container registry URLs that site agents can be deployed from.

The returned values are the ones accepted by the --registry flag of
'metalcloud-cli site one-liner'.

Required Permissions:
  sites:read - Permission to view site information

Examples:
  # List the available registry URLs
  metalcloud-cli site registry-urls

  # List them as JSON
  metalcloud-cli site registry-urls -f json

```
metalcloud-cli site registry-urls [flags]
```

### Options

```
  -h, --help   help for registry-urls
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

* [metalcloud-cli site](metalcloud-cli_site.md)	 - Manage sites (datacenters) and their configurations

