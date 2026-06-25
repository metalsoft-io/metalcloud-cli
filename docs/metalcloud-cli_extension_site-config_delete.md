## metalcloud-cli extension site-config delete

Remove an extension's configuration for a site

### Synopsis

Remove the per-site configuration of an extension for a specific site.

Arguments:
  extension_id_or_label    The unique ID or label of the extension
  site_id_or_label         The unique ID or label of the site

Examples:
  metalcloud extension site-config delete 12345 1
  metalcloud extension site-config rm my-workflow-v1 my-site

```
metalcloud-cli extension site-config delete extension_id_or_label site_id_or_label [flags]
```

### Options

```
  -h, --help   help for delete
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

* [metalcloud-cli extension site-config](metalcloud-cli_extension_site-config.md)	 - Manage per-site configuration for extensions

