## metalcloud-cli extension-instance

Manage extension instances within infrastructure deployments

### Synopsis

Manage extension instances within infrastructure deployments.

Extension instances are concrete deployments of application extensions within a specific infrastructure.
They represent the running or configured state of an extension with specific input variables and configurations.

Each extension instance is tied to an infrastructure and can be configured with custom input
variables that define its behavior. Extension instances go through various lifecycle states
including deployment, running, and deletion.

Available Commands:
  list     List all extension instances in an infrastructure
  get      Retrieve detailed extension instance information
  create   Deploy new extension instance in infrastructure
  update   Modify existing extension instance configuration
  delete   Remove extension instance from infrastructure

Examples:
  metalcloud extension-instance list my-infrastructure
  metalcloud extension-instance create my-infra --extension-id 123 --label "web-server"
  metalcloud extension-instance update inst456 --config-source updated-config.json
  metalcloud extension-instance delete inst456

### Options

```
  -h, --help   help for extension-instance
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
* [metalcloud-cli extension-instance create](metalcloud-cli_extension-instance_create.md)	 - Deploy new extension instance in specified infrastructure
* [metalcloud-cli extension-instance delete](metalcloud-cli_extension-instance_delete.md)	 - Remove extension instance from infrastructure
* [metalcloud-cli extension-instance get](metalcloud-cli_extension-instance_get.md)	 - Retrieve detailed information about a specific extension instance
* [metalcloud-cli extension-instance list](metalcloud-cli_extension-instance_list.md)	 - List all extension instances in an infrastructure
* [metalcloud-cli extension-instance update](metalcloud-cli_extension-instance_update.md)	 - Modify existing extension instance configuration

