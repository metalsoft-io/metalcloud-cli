## metalcloud-cli os-template delete

Delete an OS template

### Synopsis

Delete an OS template from the system.

This command permanently removes an OS template from the system. The template
must not be in use by any active deployments before it can be deleted.

Required arguments:
  os_template_id    The numeric ID of the template to delete

Examples:
  # Delete template with ID 123
  metalcloud-cli os-template delete 123
  
  # Delete template using alias
  metalcloud-cli templates rm 456

```
metalcloud-cli os-template delete <os_template_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

