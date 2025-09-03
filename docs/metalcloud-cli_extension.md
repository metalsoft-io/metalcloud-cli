## metalcloud-cli extension

Manage platform extensions for workflows, applications, and actions

### Synopsis

Manage platform extensions including workflows, applications, and actions.

Extensions are modular components that extend the platform's functionality. They can be:
- workflows: Define automated sequences of operations
- applications: Provide custom application deployment logic
- actions: Implement specific operational tasks

Extension lifecycle includes draft, active, and archived states. Only published extensions
become active and available for use across the platform.

Available Commands:
  list                List and filter extensions
  get                 Retrieve detailed extension information
  create              Create new extension from definition
  update              Modify existing extension properties
  publish             Activate draft extension for platform use
  archive             Deactivate published extension
  make-public         Make extension publicly available to all users
  list-repo           List extensions available in a remote repository
  create-from-repo    Create extension by cloning from a repository

Examples:
  metalcloud extension list --filter-kind workflow --filter-status active
  metalcloud extension create my-workflow workflow "Custom deployment workflow" --definition-source definition.json
  metalcloud extension update ext123 "Updated Name" "New description"
  metalcloud extension publish ext123
  metalcloud extension make-public ext123

### Options

```
  -h, --help   help for extension
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli extension archive](metalcloud-cli_extension_archive.md)	 - Deactivate published extension and make it unavailable
* [metalcloud-cli extension create](metalcloud-cli_extension_create.md)	 - Create new extension from definition
* [metalcloud-cli extension create-from-repo](metalcloud-cli_extension_create-from-repo.md)	 - Create a new extension by cloning from a repository
* [metalcloud-cli extension get](metalcloud-cli_extension_get.md)	 - Retrieve detailed information about a specific extension
* [metalcloud-cli extension list](metalcloud-cli_extension_list.md)	 - List and filter platform extensions
* [metalcloud-cli extension list-repo](metalcloud-cli_extension_list-repo.md)	 - List available extensions from a remote repository
* [metalcloud-cli extension make-public](metalcloud-cli_extension_make-public.md)	 - Make extension publicly available to all users
* [metalcloud-cli extension publish](metalcloud-cli_extension_publish.md)	 - Activate draft extension for platform use
* [metalcloud-cli extension update](metalcloud-cli_extension_update.md)	 - Modify existing extension properties and definition

