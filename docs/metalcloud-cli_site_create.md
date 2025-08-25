## metalcloud-cli site create

Create a new site with the specified name

### Synopsis

Create a new site (datacenter) with the specified name in the system.

This command creates a new site that can be used to host infrastructure components.
The site name must be unique within the system and will serve as the identifier
for the new datacenter location.

Required Arguments:
  name    The name for the new site (must be unique)

Required Permissions:
  sites:write - Permission to create sites

Examples:
  # Create a new site
  metalcloud-cli site create "datacenter-west"

  # Create a site with a descriptive name
  metalcloud-cli site create "production-datacenter-01"

```
metalcloud-cli site create name [flags]
```

### Options

```
  -h, --help   help for create
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

