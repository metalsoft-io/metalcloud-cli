## metalcloud-cli container-instance os-installation-data

Show the OS installation data of a container instance

### Synopsis

Show the OS installation data of a container instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  container_instance_id       The numeric ID of the container instance

Optional Flags:
  --usage string    Restrict the data to one usage type. One of HTTPRequest,
                    JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset.
  --remove-empty    Omit the empty entries from the response.

Examples:
  # Show the OS installation data of container instance 5678
  metalcloud-cli container-instance os-installation-data my-infra 5678

  # Show only the non-empty OS asset entries
  metalcloud-cli ci os-data 1234 5678 --usage OSAsset --remove-empty

```
metalcloud-cli container-instance os-installation-data infrastructure_id_or_label container_instance_id [flags]
```

### Options

```
  -h, --help           help for os-installation-data
      --remove-empty   Omit the empty entries from the response.
      --usage string   Restrict the OS installation data to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).
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

* [metalcloud-cli container-instance](metalcloud-cli_container-instance.md)	 - Manage container instances within infrastructures

