## metalcloud-cli job issue-command

Issue an operational command for a job

### Synopsis

Issue a command that changes the operational state of a job.

The command can be described either with individual flags (--command,
--execute-immediately) or with a JSON/YAML configuration (--config-source).
The 'job kill' command is a shortcut for issuing the "kill" command with
--execute-immediately.

Required Arguments:
  job_id    The numeric ID of the job

Required Flags (one of):
  --config-source string   Source of the command configuration. Can be 'pipe' or path to a JSON/YAML file.
  --command string         The command to issue (e.g. kill)

Optional Flags (when not using --config-source):
  --execute-immediately    Execute the command immediately instead of queueing it.

Examples:
  # Mark a job for death and execute immediately
  metalcloud-cli job issue-command 12345 --command kill --execute-immediately

  # Issue a command described in a file
  metalcloud-cli job issue-command 12345 --config-source ./command.json

  # Issue a command from piped input
  echo '{"command":"kill","executeImmediately":true}' | metalcloud-cli job issue-command 12345 --config-source pipe

Permissions:
  Requires job queue write permissions to execute this command.

```
metalcloud-cli job issue-command job_id [flags]
```

### Options

```
      --command string         The command to issue for the job (e.g. kill).
      --config-source string   Source of the job command configuration. Can be 'pipe' or path to a JSON file.
      --execute-immediately    Execute the command immediately.
  -h, --help                   help for issue-command
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

* [metalcloud-cli job](metalcloud-cli_job.md)	 - Manage MetalCloud jobs and job execution

