## metalcloud-cli permission delete

Delete a permission

### Synopsis

Delete a permission.

This command permanently deletes a permission. This action cannot be undone,
so use with caution.

Required Arguments:
  permission_name       The name of the permission to delete

Examples:
  # Delete a permission
  metalcloud-cli permission delete custom_read

  # Delete a permission using alias
  metalcloud-cli permission rm custom_read


```
metalcloud-cli permission delete permission_name [flags]
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

* [metalcloud-cli permission](metalcloud-cli_permission.md)	 - Permission management

