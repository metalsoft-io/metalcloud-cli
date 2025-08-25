## metalcloud-cli site get-config

Retrieve the configuration settings for a site

### Synopsis

Retrieve the configuration settings for a specific site (datacenter) in JSON format.

This command fetches the complete configuration settings for a site including
infrastructure parameters, deployment options, and other site-specific settings.
The configuration is returned in JSON format for easy parsing and modification.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to retrieve configuration for

Required Permissions:
  sites:read - Permission to view site information

Optional Flags:
  Common output flags are available (--format, --output, etc.)

Examples:
  # Get site configuration by name
  metalcloud-cli site get-config "datacenter-01"

  # Get site configuration by ID
  metalcloud-cli site get-config 12345

  # Save configuration to file
  metalcloud-cli site get-config "datacenter-01" > site-config.json

```
metalcloud-cli site get-config site_id_or_name [flags]
```

### Options

```
  -h, --help   help for get-config
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

