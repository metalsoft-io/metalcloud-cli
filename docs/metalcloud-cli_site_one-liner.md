## metalcloud-cli site one-liner

Get the site controller agent install script

```
metalcloud-cli site one-liner site_id_or_label [flags]
```

### Options

```
      --github-tag string      GitHub tag for deploy-agents.sh script [33m(default: main)[0m
  -h, --help                   help for one-liner
      --images-tag string      Docker images tag version [33m(default: auto-detected from system version)[0m
      --inband-mode            Install in inband mode
      --registry string        Container registry URL [33m(default: registry.metalsoft.dev)[0m
      --ssl-hostname string    SSL hostname [33m(default: eveng-qa02.metalcloud.io)[0m
      --tunnel-secret string   MS Tunnel secret for secure connections (required)
      --use-podman             Use Podman instead of Docker
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

