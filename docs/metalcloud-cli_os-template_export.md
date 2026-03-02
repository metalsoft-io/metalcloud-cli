## metalcloud-cli os-template export

Export an OS template and its assets to a zip archive

### Synopsis

Export an OS template and all its assets to a zip archive file.

This command fetches a template by ID and packs its configuration and content
assets into a portable zip archive. URL-based assets (such as ISO links)
are preserved as references without downloading the actual files.

The archive contains:
  - template.yaml: Template configuration in YAML format
  - assets/: Directory containing decoded asset file contents

Required arguments:
  os_template_id    The numeric ID of the template to export

Optional flags:
  --output          Output file path (default: <template-name-slug>.zip)

Examples:
  # Export template with ID 123
  metalcloud-cli os-template export 123

  # Export to a specific file
  metalcloud-cli os-template export 123 --output my-template.zip

```
metalcloud-cli os-template export <os_template_id> [flags]
```

### Options

```
  -h, --help            help for export
      --output string   Output file path for the exported archive.
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

