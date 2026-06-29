## metalcloud-cli extension site-config set

Set the configuration values for an extension on a site

### Synopsis

Set the configuration variable values for an extension on a specific site.

The configuration values must be provided through the --config-source flag, which
accepts either 'pipe' for stdin input or a path to a JSON file.

Arguments:
  extension_id_or_label    The unique ID or label of the extension
  site_id_or_label         The unique ID or label of the site

Required Flags:
  --config-source string   Source of the configuration values (pipe or JSON file path)

JSON Configuration Format:
  [
    {"label": "variable1", "value": "value1"},
    {"label": "variable2", "value": true}
  ]

Examples:
  # Set from JSON file
  metalcloud extension site-config set 12345 1 --config-source ./config.json

  # Set from pipe
  echo '[{"label":"env","value":"prod"}]' | metalcloud extension site-config set my-workflow-v1 my-site --config-source pipe

```
metalcloud-cli extension site-config set extension_id_or_label site_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the site configuration values. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for set
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

