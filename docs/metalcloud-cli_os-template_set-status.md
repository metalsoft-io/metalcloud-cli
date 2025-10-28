## metalcloud-cli os-template set-status

Set the status of an OS template

### Synopsis

Set the status of an OS template.

This command updates the status of an existing OS template. Valid status values
include: ready, active, used, archived.

Required arguments:
  os_template_id    The numeric ID of the template to update
  status           The new status value (ready, active, used, archived)

Examples:
  # Set template status to active
  metalcloud-cli os-template set-status 123 active
  
  # Archive a template
  metalcloud-cli os-template set-status 456 archived
  
  # Set template to ready status using alias
  metalcloud-cli templates status 789 ready

```
metalcloud-cli os-template set-status <os_template_id> <status> [flags]
```

### Options

```
  -h, --help   help for set-status
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

