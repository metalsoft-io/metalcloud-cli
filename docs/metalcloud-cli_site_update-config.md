## metalcloud-cli site update-config

Update site configuration using JSON input

### Synopsis

Update the configuration settings for a specific site (datacenter) using JSON input.

This command allows you to modify the configuration settings of an existing site.
The configuration can be provided through a file or piped from standard input.
The configuration must be in valid JSON format and contain the appropriate
site configuration parameters.

Required Arguments:
  site_id_or_name    Site identifier (ID or name) to update configuration for

Required Flags:
  --config-source    Source of the site configuration. Can be 'pipe' for stdin input
                     or path to a JSON file containing the configuration

Required Permissions:
  sites:write - Permission to modify sites

Dependencies:
  The --config-source flag is mandatory and must specify either:
  - 'pipe' to read JSON configuration from standard input
  - Path to a valid JSON file containing site configuration

Examples:
  # Update site configuration from a file
  metalcloud-cli site update-config "datacenter-01" --config-source config.json

  # Update site configuration from standard input
  cat config.json | metalcloud-cli site update-config 12345 --config-source pipe

  # Update site configuration with inline JSON
  echo '{"key": "value"}' | metalcloud-cli site update-config "datacenter-01" --config-source pipe

```
metalcloud-cli site update-config site_id_or_name [flags]
```

### Options

```
      --config-source string   Source of the site configuration. Can be 'pipe' for stdin or path to a JSON file (required).
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

* [metalcloud-cli site](metalcloud-cli_site.md)	 - Manage sites (datacenters) and their configurations

