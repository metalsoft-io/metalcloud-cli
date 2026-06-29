## metalcloud-cli extension activate

Activate a suspended extension

### Synopsis

Activate an extension, returning it to active status so it can be used across the platform.

This command transitions an extension to the active status. It is typically used to
re-enable an extension that was previously suspended.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to activate

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

