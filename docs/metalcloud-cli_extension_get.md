## metalcloud-cli extension get

Retrieve detailed information about a specific extension

### Synopsis

Retrieve detailed information about a specific extension by ID or label.

This command displays comprehensive information about an extension including its
metadata, definition, current status, and configuration. The extension can be
identified by either its unique ID or label.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to retrieve

Examples:
  # Get extension by ID
  metalcloud extension get 12345
  
  # Get extension by label
  metalcloud extension get my-workflow-v1
  
  # Show extension details
  metalcloud extension show production-deployment

```
metalcloud-cli extension get extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for get
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

