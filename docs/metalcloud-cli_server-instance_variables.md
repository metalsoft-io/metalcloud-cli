## metalcloud-cli server-instance variables

Get the variables of a server instance

### Synopsis

Get the variables of a server instance: the site, server, instance, group,
infrastructure, drive and network values available to extensions and templates.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --usage string   Restrict the returned values to one usage type (HTTPRequest, JavaScript,
                   APICall, AnsibleBundle, SSHExec, Copy, OSAsset).

Examples:
  metalcloud-cli server-instance variables 5678
  metalcloud-cli inst vars 5678 --usage APICall

```
metalcloud-cli server-instance variables server_instance_id [flags]
```

### Options

```
  -h, --help           help for variables
      --usage string   Restrict the returned values to one usage type.
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

