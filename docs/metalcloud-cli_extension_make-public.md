## metalcloud-cli extension make-public

Make extension publicly available to all users

### Synopsis

Make an extension publicly available to all users in the organization.

This command changes the visibility of an extension from private (accessible only
to the owner) to public (accessible to all users with appropriate permissions).
Public extensions can be discovered and used by other users within the organization.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to make public

Required permissions:
  - extensions:write

Dependencies:
  - Extension must exist and be accessible
  - User must be the owner of the extension or have admin privileges

Examples:
  # Make extension public by ID
  metalcloud extension make-public 12345
  
  # Make extension public by label
  metalcloud extension make-public my-workflow-v1

```
metalcloud-cli extension make-public extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for make-public
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
