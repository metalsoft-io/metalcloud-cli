## metalcloud-cli extension site-config list

List the site configurations for an extension

### Synopsis

List all per-site configurations defined for an extension.

Arguments:
  extension_id_or_label    The unique ID or label of the extension

Examples:
  metalcloud extension site-config list 12345
  metalcloud extension site-config list my-workflow-v1

```
metalcloud-cli extension site-config list extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for list
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

