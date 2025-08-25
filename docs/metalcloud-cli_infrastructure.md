## metalcloud-cli infrastructure

Manage infrastructure resources and configurations

### Synopsis

Manage infrastructure resources including creation, deployment, monitoring, and user access control.

Infrastructure represents a collection of compute instances, storage drives, and network resources 
that can be managed as a single unit. Each infrastructure belongs to a specific site and can be 
deployed, updated, or deleted as needed.

Available Commands:
  list         List all infrastructures with filtering options
  get          Show detailed information about a specific infrastructure
  create       Create a new infrastructure in a site
  update       Update infrastructure configuration and metadata
  delete       Delete an infrastructure and all its resources
  deploy       Deploy infrastructure changes to physical resources
  revert       Revert infrastructure to the last deployed state
  users        Manage user access to infrastructures
  statistics   View infrastructure deployment and job statistics
  utilization  Generate resource utilization reports

Use "metalcloud-cli infrastructure [command] --help" for more information about a specific command.

### Options

```
  -h, --help   help for infrastructure
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli infrastructure add-user](metalcloud-cli_infrastructure_add-user.md)	 - Add a user to an infrastructure with access permissions
* [metalcloud-cli infrastructure all-statistics](metalcloud-cli_infrastructure_all-statistics.md)	 - Get deployment statistics for all infrastructures
* [metalcloud-cli infrastructure config-info](metalcloud-cli_infrastructure_config-info.md)	 - Get configuration information for an infrastructure
* [metalcloud-cli infrastructure create](metalcloud-cli_infrastructure_create.md)	 - Create a new infrastructure in a specific site
* [metalcloud-cli infrastructure delete](metalcloud-cli_infrastructure_delete.md)	 - Delete an infrastructure and all its resources
* [metalcloud-cli infrastructure deploy](metalcloud-cli_infrastructure_deploy.md)	 - Deploy infrastructure changes to physical resources
* [metalcloud-cli infrastructure get](metalcloud-cli_infrastructure_get.md)	 - Show detailed information about a specific infrastructure
* [metalcloud-cli infrastructure list](metalcloud-cli_infrastructure_list.md)	 - List infrastructures with optional filtering
* [metalcloud-cli infrastructure remove-user](metalcloud-cli_infrastructure_remove-user.md)	 - Remove a user's access from an infrastructure
* [metalcloud-cli infrastructure revert](metalcloud-cli_infrastructure_revert.md)	 - Revert infrastructure to the last deployed state
* [metalcloud-cli infrastructure statistics](metalcloud-cli_infrastructure_statistics.md)	 - Get deployment statistics for an infrastructure
* [metalcloud-cli infrastructure update](metalcloud-cli_infrastructure_update.md)	 - Update infrastructure configuration and metadata
* [metalcloud-cli infrastructure user-limits](metalcloud-cli_infrastructure_user-limits.md)	 - Display resource limits for an infrastructure
* [metalcloud-cli infrastructure users](metalcloud-cli_infrastructure_users.md)	 - List users with access to an infrastructure
* [metalcloud-cli infrastructure utilization](metalcloud-cli_infrastructure_utilization.md)	 - Get resource utilization report for infrastructures

