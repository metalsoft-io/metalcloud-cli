## metalcloud-cli extension update

Modify existing extension properties and definition

### Synopsis

Modify existing extension properties including name, description, and definition.

This command allows you to update various properties of an existing extension.
You can update the name, description, and/or the extension definition. All
parameters are optional, allowing you to update only specific properties.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to update
  name                     New name for the extension (optional)
  description              New description for the extension (optional)

Optional Flags:
  --definition-source string   Source of the updated extension definition
                              Can be 'pipe' for stdin or path to a JSON file

Flag Dependencies:
- --definition-source is independent of other parameters
- name and description are positional arguments

Examples:
  # Update only the name
  metalcloud extension update ext123 "New Extension Name"
  
  # Update name and description
  metalcloud extension update ext123 "New Name" "Updated description"
  
  # Update only the definition
  metalcloud extension update ext123 --definition-source updated-definition.json
  
  # Update name, description, and definition
  metalcloud extension update ext123 "New Name" "New description" --definition-source definition.json
  
  # Update definition from stdin
  cat new-definition.json | metalcloud extension update ext123 --definition-source pipe

```
metalcloud-cli extension update extension_id_or_label [name [description]] [flags]
```

### Options

```
      --definition-source string   Source of the updated extension definition. Can be 'pipe' or path to a JSON file.
  -h, --help                       help for update
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

