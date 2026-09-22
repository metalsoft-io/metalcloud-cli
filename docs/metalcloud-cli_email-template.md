## metalcloud-cli email-template

Email template management

### Synopsis

Manage the email templates the platform uses for its notifications.

Templates are addressed by NAME on every command - they have no addressable
numeric ID.

Command categories:
  Read:    list, get
  Write:   create, update, delete
  Helper:  config-example

### Options

```
  -h, --help   help for email-template
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
* [metalcloud-cli email-template config-example](metalcloud-cli_email-template_config-example.md)	 - Print an example email template configuration
* [metalcloud-cli email-template create](metalcloud-cli_email-template_create.md)	 - Create an email template
* [metalcloud-cli email-template delete](metalcloud-cli_email-template_delete.md)	 - Delete an email template
* [metalcloud-cli email-template get](metalcloud-cli_email-template_get.md)	 - Get email template details
* [metalcloud-cli email-template list](metalcloud-cli_email-template_list.md)	 - List email templates
* [metalcloud-cli email-template update](metalcloud-cli_email-template_update.md)	 - Update an email template

