## metalcloud-cli extension activate

Activate a draft or suspended extension

### Synopsis

Activate an extension, making it available for use across the platform.

This command transitions an extension from draft or suspended status to active status.
Use it both to activate a newly created draft extension and to re-enable an extension
that was previously suspended. Only active extensions are available for use in
workflows, applications, and actions.

This command replaces the deprecated 'publish' command.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to activate

Requirements:
- Extension must be in draft or suspended status
- User must have write permissions for extensions

Examples:
  # Activate extension by ID
  metalcloud extension activate 12345

  # Activate extension by label
  metalcloud extension activate my-workflow-v1

```
metalcloud-cli extension activate extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for activate
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

