## metalcloud-cli server-instance os-installation-data

Get the OS installation data of a server instance

### Synopsis

Get the OS installation data of a server instance: the site, server, instance,
group, infrastructure, drive and network values the OS installer is rendered with.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --usage string   Restrict the returned values to one usage type (HTTPRequest, JavaScript,
                   APICall, AnsibleBundle, SSHExec, Copy, OSAsset).

Examples:
  metalcloud-cli server-instance os-installation-data 5678
  metalcloud-cli inst os-data 5678 --usage AnsibleBundle

```
metalcloud-cli server-instance os-installation-data server_instance_id [flags]
```

### Options

```
  -h, --help           help for os-installation-data
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

