## metalcloud-cli os-template list-repo

List available OS templates from a remote repository

### Synopsis

List all available OS templates from a remote repository.

This command retrieves and displays templates available in a remote repository,
showing their basic information and configuration.

Optional flags:
  --repo-url        URL of the repository to list templates from
                   Defaults to the official MetalSoft template repository
  --repo-username   Username for private repository authentication
  --repo-password   Password for private repository authentication

Flag dependencies:
  - If accessing a private repository, both --repo-username and --repo-password
    are required together

Examples:
  # List templates from default public repository
  metalcloud-cli os-template list-repo
  
  # List templates from a custom repository
  metalcloud-cli os-template list-repo --repo-url https://example.com/templates
  
  # List templates from private repository
  metalcloud-cli os-template list-repo --repo-url https://private.com/templates \
    --repo-username user --repo-password pass

```
metalcloud-cli os-template list-repo [flags]
```

### Options

```
  -h, --help                   help for list-repo
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

* [metalcloud-cli os-template](metalcloud-cli_os-template.md)	 - Manage OS templates for server deployments

