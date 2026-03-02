## metalcloud-cli os-template clone

Clone an existing OS template

### Synopsis

Clone an existing OS template to create a new copy.

This command creates a new OS template that is an exact copy of an existing one,
including all its assets. The cloned template is always created with private
visibility. Name and label can be optionally overridden.

Required arguments:
  os_template_id    The numeric ID of the template to clone

Optional flags:
  --name            Name for the cloned template (default: "<original-name> (clone)")
  --label           Label for the cloned template (default: slug of name)

Examples:
  # Clone template with ID 123
  metalcloud-cli os-template clone 123

  # Clone with custom name
  metalcloud-cli os-template clone 123 --name "My Custom Ubuntu"

  # Clone with custom name and label
  metalcloud-cli os-template clone 123 --name "My Custom Ubuntu" --label "my-custom-ubuntu"

```
metalcloud-cli os-template clone <os_template_id> [flags]
```

### Options

```
  -h, --help           help for clone
      --label string   Label of the cloned OS template.
      --name string    Name of the cloned OS template.
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

* [metalcloud-cli os-template](metalcloud-cli_os-template.md)	 - Manage OS templates for server deployments

