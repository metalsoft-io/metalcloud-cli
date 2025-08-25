## metalcloud-cli extension publish

Activate draft extension for platform use

### Synopsis

Activate a draft extension making it available for use across the platform.

This command publishes a draft extension, changing its status from draft to active.
Only published extensions are available for use in workflows, applications, and
actions. Once published, an extension cannot be modified directly - you must
create a new version or archive and recreate it.

Publishing validates the extension definition and ensures it meets all platform
requirements before making it available to users.

Arguments:
  extension_id_or_label    The unique ID or label of the draft extension to publish

Requirements:
- Extension must be in draft status
- Extension definition must be valid
- User must have write permissions for extensions

Examples:
  # Publish extension by ID
  metalcloud extension publish 12345
  
  # Publish extension by label
  metalcloud extension publish my-workflow-v1

```
metalcloud-cli extension publish extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for publish
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

