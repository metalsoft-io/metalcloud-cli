## metalcloud-cli configuration

Global platform configuration management

### Synopsis

Read and change the GLOBAL configuration of the MetalSoft platform.

This is not the configuration of the CLI (see the --config flag for that): these
commands read and write the platform-wide settings served by /api/v2/config, which
affect every user and every infrastructure of the installation.

The configuration is split into service sections, addressed by their name:
  auth, platform, notification, tunnel, gateway-api, image-builder

Command categories:
  Read:    get
  Write:   replace (full section), update (partial section)

### Options

```
  -h, --help   help for configuration
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli configuration get](metalcloud-cli_configuration_get.md)	 - Get the global platform configuration
* [metalcloud-cli configuration replace](metalcloud-cli_configuration_replace.md)	 - Replace a section of the global platform configuration
* [metalcloud-cli configuration update](metalcloud-cli_configuration_update.md)	 - Update part of a section of the global platform configuration

