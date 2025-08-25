## metalcloud-cli extension-instance create

Deploy new extension instance in specified infrastructure

### Synopsis

Deploy a new extension instance in the specified infrastructure.

This command creates and deploys a new extension instance within an infrastructure.
Extension instances are concrete deployments of extensions (workflows, applications,
or actions) with specific configurations and input variables.

You can provide the configuration using two methods:

Method 1: Configuration file or pipe (--config-source)
  Use a JSON file or pipe the configuration via stdin. This method allows for
  complete configuration including complex input variables.

Method 2: Individual flags (--extension-id with optional flags)
  Specify the extension ID directly with optional label and input variables.
  This method is suitable for simple configurations.

Arguments:
  infrastructure_id_or_label    The unique ID or label of the target infrastructure

Required Flags (mutually exclusive):
  --config-source string        Source of configuration (pipe or JSON file path)
  --extension-id int           ID of the extension to instantiate

Optional Flags (only with --extension-id):
  --label string               Custom label for the extension instance
  --input-variable strings     Input variables in 'label=value' format (repeatable)

Flag Dependencies:
- --config-source and --extension-id are mutually exclusive
- One of --config-source or --extension-id is required
- --label and --input-variable only work with --extension-id

JSON Configuration Format:
  {
    "extensionId": 123,
    "label": "optional-instance-label",
    "inputVariables": [
      {"label": "variable1", "value": "value1"},
      {"label": "variable2", "value": "value2"}
    ]
  }

Examples:
  # Create from JSON file
  metalcloud extension-instance create my-infra --config-source ./config.json

  # Create from pipe
  echo '{"extensionId": 123, "label": "web-app"}' | metalcloud extension-instance create my-infra --config-source pipe

  # Create using individual flags
  metalcloud extension-instance create my-infra --extension-id 123 --label "database-server"

  # Create with input variables
  metalcloud extension-instance create my-infra --extension-id 123 --input-variable "env=production" --input-variable "replicas=3"

  # Create minimal instance (auto-generated label)
  metalcloud ext-inst create prod-infra --extension-id 456

```
metalcloud-cli extension-instance create infrastructure_id_or_label [flags]
```

### Options

```
      --config-source string         Source of the new extension instance configuration. Can be 'pipe' or path to a JSON file.
      --extension-id int             The extension ID to create an instance of.
  -h, --help                         help for create
      --input-variable stringArray   Input variables in format 'label=value'. Can be specified multiple times.
      --label string                 The extension instance label (optional, will be auto-generated if not provided).
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

* [metalcloud-cli extension-instance](metalcloud-cli_extension-instance.md)	 - Manage extension instances within infrastructure deployments

