## metalcloud-cli os-template list-directory

List available OS templates from a local directory

### Synopsis

List all available OS templates from a local directory.

This command retrieves and displays the templates found in a local directory,
showing their basic information and configuration. It is the offline alternative
to 'list-repo': everything it reports is read from the local filesystem, so it
needs no network access and no API endpoint or API key.

The directory must hold the templates in the same layout used by the template
repositories:

  <directory>/<vendor>/<os>/<version>/template.yaml
  <directory>/<vendor>/<os>/<version>/<asset files>

Only templates found at that depth are listed. README.md files are ignored and
templates using an older format are skipped with a warning.

Required flags:
  --dir             Path of the local directory holding the templates

Examples:
  # List templates from a local directory
  metalcloud-cli os-template list-directory --dir /var/lib/os-templates

  # List templates from a local directory on Windows
  metalcloud-cli os-template list-directory --dir C:\os-templates

```
metalcloud-cli os-template list-directory [flags]
```

### Options

```
      --dir string   Local directory holding the OS templates.
  -h, --help         help for list-directory
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

