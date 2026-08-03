## metalcloud-cli os-template create-from-directory

Create a new OS template from a local directory

### Synopsis

Create a new OS template from a template stored in a local directory.

This command creates an OS template based on a template definition found in a
local directory. It is the alternative to 'create-from-repo' for air-gapped
environments: the template and its assets are read from the local filesystem
instead of being cloned from a repository, so no repository access is needed.
The template itself is still created through the MetalSoft API, so an endpoint
and an API key are required.

The directory must hold the templates in the same layout used by the template
repositories:

  <directory>/<vendor>/<os>/<version>/template.yaml
  <directory>/<vendor>/<os>/<version>/<asset files>

Required arguments:
  os_template_path  Path of the template inside the local directory
                    Use 'list-directory' command to see available templates

Required flags:
  --dir             Path of the local directory holding the templates

Optional flags:
  --name           Custom name for the new template (overrides original)
  --label          Custom label for the new template (overrides original)
  --source-iso     Custom source ISO image path (overrides original)

Examples:
  # Create a template from a local directory
  metalcloud-cli os-template create-from-directory Ubuntu/24.04/oob-u24-04-3-lts-v7 \
    --dir /var/lib/os-templates

  # Create with custom name and label
  metalcloud-cli os-template create-from-directory ubuntu/22.04/server \
    --dir /var/lib/os-templates \
    --name "My Ubuntu 22.04" --label "my-ubuntu-2204"

  # Create from a local directory on Windows, overriding the source ISO
  metalcloud-cli os-template create-from-directory Ubuntu/24.04/oob-u24-04-3-lts-v7 \
    --dir C:\os-templates \
    --source-iso http://repo.local/isos/ubuntu-24.04.iso

```
metalcloud-cli os-template create-from-directory <os_template_path> [flags]
```

### Options

```
      --dir string          Local directory holding the OS templates.
  -h, --help                help for create-from-directory
      --label string        Label of the OS template.
      --name string         Name of the OS template.
      --source-iso string   The source ISO image path.
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

