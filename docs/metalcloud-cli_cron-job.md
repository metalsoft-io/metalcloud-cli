## metalcloud-cli cron-job

Manage scheduled cron jobs

### Synopsis

Manage scheduled cron jobs in MetalCloud.

Cron jobs allow you to schedule recurring operations that are executed
automatically on a defined schedule.

Available Commands:
  list    List all cron jobs
  get     Get detailed information about a specific cron job
  create  Create a new cron job
  update  Update an existing cron job
  delete  Delete a cron job

Use "metalcloud-cli cron-job [command] --help" for more information about a command.

### Options

```
  -h, --help   help for cron-job
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
* [metalcloud-cli cron-job create](metalcloud-cli_cron-job_create.md)	 - Create a new cron job from configuration
* [metalcloud-cli cron-job delete](metalcloud-cli_cron-job_delete.md)	 - Delete a cron job
* [metalcloud-cli cron-job get](metalcloud-cli_cron-job_get.md)	 - Get detailed information about a specific cron job
* [metalcloud-cli cron-job list](metalcloud-cli_cron-job_list.md)	 - List all cron jobs
* [metalcloud-cli cron-job update](metalcloud-cli_cron-job_update.md)	 - Update an existing cron job

