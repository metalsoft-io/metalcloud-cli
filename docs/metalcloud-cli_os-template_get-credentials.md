## metalcloud-cli os-template get-credentials

Show default credentials for an OS template

### Synopsis

Display the default credentials for an OS template.

This command retrieves and displays the default username and password
that are configured for a specific OS template. These credentials are
used for initial access to servers deployed with this template.

Required arguments:
  os_template_id    The numeric ID of the template

Examples:
  # Get credentials for template with ID 123
  metalcloud-cli os-template get-credentials 123
  
  # Get credentials using alias
  metalcloud-cli templates creds 456

```
metalcloud-cli os-template get-credentials <os_template_id> [flags]
```

### Options

```
  -h, --help   help for get-credentials
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

