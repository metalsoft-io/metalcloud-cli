## metalcloud-cli extension site-config

Manage per-site configuration for extensions

### Synopsis

Manage per-site configuration variables for extensions.

Site configurations bind an extension to a specific site with a set of configuration
variables, controlling how the extension behaves on that site.

Available Commands:
  list             List the site configurations for an extension
  get              Get the configuration values for an extension on a site
  set              Set the configuration values for an extension on a site
  delete           Remove an extension's configuration for a site
  list-for-site    List the extension configurations defined for a site

### Options

```
  -h, --help   help for site-config
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

* [metalcloud-cli extension](metalcloud-cli_extension.md)	 - Manage platform extensions for workflows, applications, and actions
* [metalcloud-cli extension site-config delete](metalcloud-cli_extension_site-config_delete.md)	 - Remove an extension's configuration for a site
* [metalcloud-cli extension site-config get](metalcloud-cli_extension_site-config_get.md)	 - Get the configuration values for an extension on a site
* [metalcloud-cli extension site-config list](metalcloud-cli_extension_site-config_list.md)	 - List the site configurations for an extension
* [metalcloud-cli extension site-config list-for-site](metalcloud-cli_extension_site-config_list-for-site.md)	 - List the extension configurations defined for a site
* [metalcloud-cli extension site-config set](metalcloud-cli_extension_site-config_set.md)	 - Set the configuration values for an extension on a site

