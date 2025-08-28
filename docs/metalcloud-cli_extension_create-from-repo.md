## metalcloud-cli extension create-from-repo

Create a new extension by cloning from a repository

### Synopsis

Create a new extension by cloning an existing extension from a repository.

This command downloads and creates a local extension based on an extension
available in a remote repository. You can optionally customize the name
and label during the creation process.

Required arguments:
  extension_path    Path to the extension within the repository
                   Use 'list-repo' command to see available extensions

Optional flags:
  --repo-url        URL of the repository to clone from
                   Defaults to the official MetalSoft extension repository
  --repo-username   Username for private repository authentication  
  --repo-password   Password for private repository authentication
  --name           Custom name for the new extension (overrides original)
  --label          Custom label for the new extension (overrides original)

Flag dependencies:
  - If accessing a private repository, both --repo-username and --repo-password
    are required together

Examples:
  # Clone extension from default public repository
  metalcloud extension create-from-repo workflows/deployment/basic-deployment
  
  # Clone with custom name and label
  metalcloud extension create-from-repo workflows/cleanup/resource-cleanup \
    --name "My Resource Cleanup" --label "my-resource-cleanup"
  
  # Clone from private repository
  metalcloud extension create-from-repo actions/monitoring/health-check \
    --repo-url https://private.com/extensions \
    --repo-username user --repo-password pass \
    --name "Custom Health Check"

```
metalcloud-cli extension create-from-repo <extension_path> [flags]
```

### Options

```
  -h, --help                   help for create-from-repo
      --label string           Label of the extension.
      --name string            Name of the extension.
      --repo-password string   Private repo password.
      --repo-url string        Private repo to use.
      --repo-username string   Private repo username.
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

