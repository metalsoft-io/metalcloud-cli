## metalcloud-cli scheduled-job

Manage scheduled jobs

### Synopsis

Manage scheduled jobs in MetalCloud.

Scheduled jobs allow you to schedule recurring operations that are executed
automatically on a defined schedule.

Available Commands:
  list    List all scheduled jobs
  get     Get detailed information about a specific scheduled job
  create  Create a new scheduled job
  update  Update an existing scheduled job
  delete  Delete a scheduled job

Use "metalcloud-cli scheduled-job [command] --help" for more information about a command.

### Options

```
  -h, --help   help for scheduled-job
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
* [metalcloud-cli scheduled-job create](metalcloud-cli_scheduled-job_create.md)	 - Create a new scheduled job from configuration
* [metalcloud-cli scheduled-job delete](metalcloud-cli_scheduled-job_delete.md)	 - Delete a scheduled job
* [metalcloud-cli scheduled-job get](metalcloud-cli_scheduled-job_get.md)	 - Get detailed information about a specific scheduled job
* [metalcloud-cli scheduled-job list](metalcloud-cli_scheduled-job_list.md)	 - List all scheduled jobs
* [metalcloud-cli scheduled-job update](metalcloud-cli_scheduled-job_update.md)	 - Update an existing scheduled job

