## metalcloud-cli os-template import

Import an OS template from a zip archive

### Synopsis

Import an OS template from a zip archive file.

This command reads a previously exported zip archive and creates a new
OS template with all its assets. The new template is always created with
private visibility.

The archive should contain:
  - template.yaml: Template configuration in YAML format
  - assets/: Directory containing asset file contents

Required arguments:
  archive_path      Path to the zip archive file

Required flags:
  --name            Name for the new template

Optional flags:
  --label           Label for the new template (default: slug of name)

Examples:
  # Import a template
  metalcloud-cli os-template import my-template.zip --name "My Imported Template"

  # Import with custom label
  metalcloud-cli os-template import my-template.zip --name "My Template" --label "my-template-v2"

```
metalcloud-cli os-template import <archive_path> [flags]
```

### Options

```
  -h, --help           help for import
      --label string   Label of the new OS template.
      --name string    Name of the new OS template.
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

