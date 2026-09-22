## metalcloud-cli user update-meta

Update the metadata of a user

### Synopsis

Update the metadata of a specific user account.

User metadata carries the GUI settings of the user. The payload replaces the stored
metadata, so provide the complete object.

Required Arguments:
  user_id                 The numeric ID of the user whose metadata to update

Required Flags:
  --config-source         Source of the metadata (JSON/YAML file path or 'pipe')

Configuration File Format (JSON):
  {
    "guiSettings": {
      "defaultPage": "infrastructures"
    }
  }

Examples:
  metalcloud-cli user update-meta 12345 --config-source meta.json
  echo '{"guiSettings":{}}' | metalcloud-cli user update-meta 12345 --config-source pipe

```
metalcloud-cli user update-meta user_id [flags]
```

### Options

```
      --config-source string   Source of the user metadata. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update-meta
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

* [metalcloud-cli user](metalcloud-cli_user.md)	 - Manage user accounts and their properties

