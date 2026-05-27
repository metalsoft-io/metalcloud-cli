## metalcloud-cli os-template create-from-repo

Create a new OS template by cloning from a repository

### Synopsis

Create a new OS template by cloning an existing template from a repository.

This command downloads and creates a local OS template based on a template
available in a remote repository. You can optionally customize the name,
label, and source ISO image during the creation process.

Required arguments:
  os_template_path  Path to the template within the repository
                   Use 'list-repo' command to see available templates

Optional flags:
  --repo-url        URL of the repository to clone from
                   Defaults to the official MetalSoft template repository
  --repo-username   Username for private repository authentication  
  --repo-password   Password for private repository authentication
  --name           Custom name for the new template (overrides original)
  --label          Custom label for the new template (overrides original)
  --source-iso     Custom source ISO image path (overrides original)

Flag dependencies:
  - If accessing a private repository, both --repo-username and --repo-password
    are required together

Examples:
  # Clone template from default public repository
  metalcloud-cli os-template create-from-repo ubuntu/22.04/server
  
  # Clone with custom name and label
  metalcloud-cli os-template create-from-repo ubuntu/22.04/server \
    --name "My Ubuntu 22.04" --label "my-ubuntu-2204"
  
  # Clone from private repository with custom ISO
  metalcloud-cli os-template create-from-repo centos/7/server \
    --repo-url https://private.com/templates \
    --repo-username user --repo-password pass \
    --source-iso /path/to/custom.iso

  # Clone from private repository on Windows OS (folder C:\os-templates)
  # (on Windows, please use / to replace \ and do not include the drive letter (for
  # example, if the os-template folder is in c:\os-templates, then the repo-url
  # would be /os-templates and the command would need to be issued from that drive))
  metalcloud-cli os-template create-from-repo Ubuntu/24.04/oob-u24-04-3-lts-v7
    --repo-url /os-templates --name "Ubuntu2404" --label "Ubuntu2404"
```
metalcloud-cli os-template create-from-repo <os_template_path> [flags]
```

### Options

```
  -h, --help                   help for create-from-repo
      --label string           Label of the OS template.
      --name string            Name of the OS template.
      --repo-password string   Private repo password.
      --repo-url string        Private repo to use.
      --repo-username string   Private repo username.
      --source-iso string      The source ISO image path.
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

