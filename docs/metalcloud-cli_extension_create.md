## metalcloud-cli extension create

Create new extension from definition

### Synopsis

Create a new extension with the specified name, kind, and description.

This command creates a new extension in draft status. The extension definition must be
provided through the --definition-source flag, which accepts either 'pipe' for stdin
input or a path to a JSON file containing the extension definition.

Extension kinds:
- workflow: Automated sequences of operations
- application: Custom application deployment logic
- action: Specific operational tasks

The newly created extension will be in draft status and must be published before
it becomes available for use on the platform.

Arguments:
  name          The name of the extension to create
  kind          The extension type (workflow, application, action)
  description   Description of the extension's purpose and functionality

Required Flags:
  --definition-source string   Source of the extension definition (required)
                              Can be 'pipe' for stdin or path to a JSON file

Examples:
  # Create extension from JSON file
  metalcloud extension create my-workflow workflow "Custom deployment workflow" --definition-source workflow.json
  
  # Create extension from stdin
  cat definition.json | metalcloud extension create my-app application "Custom app logic" --definition-source pipe
  
  # Create action extension
  metalcloud extension create cleanup-action action "Cleanup resources" --definition-source ./actions/cleanup.json

```
metalcloud-cli extension create name kind description [flags]
```

### Options

```
      --definition-source string   Source of the extension definition. Can be 'pipe' or path to a JSON file.
  -h, --help                       help for create
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

