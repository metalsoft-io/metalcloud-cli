## metalcloud-cli extension delete

Permanently delete an extension

### Synopsis

Permanently delete an extension from the platform.

This command permanently removes an extension identified by its ID or label.
Unlike archiving, which preserves the extension's definition and history in an
inactive state, deletion is irreversible and removes the extension entirely.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to delete

Requirements:
- Extension must exist and be accessible
- User must have write permissions for extensions
- Extension should not be in use by any active extension instances

Examples:
  # Delete extension by ID
  metalcloud extension delete 12345

  # Delete extension by label
  metalcloud extension rm deprecated-workflow-v1

```
metalcloud-cli extension delete extension_id_or_label [flags]
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

* [metalcloud-cli extension](metalcloud-cli_extension.md)	 - Manage platform extensions for workflows, applications, and actions

