## metalcloud-cli infrastructure users

List users with access to an infrastructure

### Synopsis

Display all users who have access to a specific infrastructure, including their 
access levels and contact information.

This command shows the current user permissions for an infrastructure, which is useful 
for managing access control and understanding who can modify or view the infrastructure.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # List users for infrastructure by ID
  metalcloud-cli infrastructure users 123

  # List users for infrastructure by label
  metalcloud-cli infrastructure users "web-cluster"

  # Using aliases
  metalcloud-cli infrastructure list-users production-env
  metalcloud-cli infrastructure get-users my-infrastructure

```
metalcloud-cli infrastructure users infrastructure_id_or_label [flags]
```

### Options

```
  -h, --help   help for users
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

* [metalcloud-cli infrastructure](metalcloud-cli_infrastructure.md)	 - Manage infrastructure resources and configurations

