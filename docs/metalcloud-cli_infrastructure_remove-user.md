## metalcloud-cli infrastructure remove-user

Remove a user's access from an infrastructure

### Synopsis

Remove a user's access permissions from a specific infrastructure.

This command revokes the specified user's access to the infrastructure, preventing them 
from viewing or modifying it. The user is identified by their numeric user ID.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure
  user_id                     The numeric ID of the user to remove

Examples:
  # Remove user by infrastructure ID and user ID
  metalcloud-cli infrastructure remove-user 123 456

  # Remove user by infrastructure label
  metalcloud-cli infrastructure remove-user "web-cluster" 789

  # Using the alias
  metalcloud-cli infrastructure delete-user my-infrastructure 101

```
metalcloud-cli infrastructure remove-user infrastructure_id_or_label user_id [flags]
```

### Options

```
  -h, --help   help for remove-user
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

