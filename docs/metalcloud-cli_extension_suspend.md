## metalcloud-cli extension suspend

Suspend an active extension

### Synopsis

Suspend an extension, temporarily disabling it without archiving or deleting it.

This command transitions an active extension to the suspended status. Suspended
extensions are unavailable for use but can be reactivated later with the activate
command.

Arguments:
  extension_id_or_label    The unique ID or label of the extension to suspend

Examples:
  # Suspend extension by ID
  metalcloud extension suspend 12345

  # Suspend extension by label
  metalcloud extension suspend my-workflow-v1

```
metalcloud-cli extension suspend extension_id_or_label [flags]
```

### Options

```
  -h, --help   help for suspend
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

* [metalcloud-cli extension](metalcloud-cli_extension.md)	 - Manage platform extensions for workflows, applications, and actions

