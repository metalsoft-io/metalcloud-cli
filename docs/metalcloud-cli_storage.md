## metalcloud-cli storage

Manage storage pools and related resources

### Synopsis

Manage storage pools and their associated resources in the MetalCloud infrastructure.

Storage pools are external storage systems that provide block, file, or object storage
to instances. This command group allows you to create, configure, and manage storage
pools, as well as access their drives, file shares, buckets, and network configurations.

Available commands:
  list             List all storage pools
  get              Get detailed information about a specific storage pool
  create           Create a new storage pool
  delete           Delete an existing storage pool
  config-example   Display a configuration template for creating storage pools
  credentials      Retrieve credentials for a storage pool
  drives           List drives available in a storage pool
  file-shares      List file shares in a storage pool
  buckets          List object storage buckets in a storage pool
  network-configs  List network device configurations for a storage pool

Use "metalcloud storage [command] --help" for more information about a command.

### Options

```
  -h, --help   help for storage
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
* [metalcloud-cli storage buckets](metalcloud-cli_storage_buckets.md)	 - List object storage buckets in a storage pool
* [metalcloud-cli storage config-example](metalcloud-cli_storage_config-example.md)	 - Display a configuration template for creating storage pools
* [metalcloud-cli storage create](metalcloud-cli_storage_create.md)	 - Create a new storage pool
* [metalcloud-cli storage credentials](metalcloud-cli_storage_credentials.md)	 - Retrieve credentials for a storage pool
* [metalcloud-cli storage delete](metalcloud-cli_storage_delete.md)	 - Delete an existing storage pool
* [metalcloud-cli storage drives](metalcloud-cli_storage_drives.md)	 - List drives available in a storage pool
* [metalcloud-cli storage file-shares](metalcloud-cli_storage_file-shares.md)	 - List file shares in a storage pool
* [metalcloud-cli storage get](metalcloud-cli_storage_get.md)	 - Get detailed information about a specific storage pool
* [metalcloud-cli storage list](metalcloud-cli_storage_list.md)	 - List all storage pools

