## metalcloud-cli server firmware generate-audit

Generate firmware upgrade audit for servers

### Synopsis

Generate firmware upgrade audit for servers.

This command generates a comprehensive firmware upgrade audit report for servers
using a JSON configuration file or piped JSON data. The audit helps identify
which servers need firmware updates and provides detailed upgrade recommendations.

Required Flags:
  --config-source        Source of the audit configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Generate audit using JSON configuration file
  metalcloud-cli server firmware generate-audit --config-source ./audit-config.json

  # Generate audit using piped JSON configuration
  echo '{"serverIds": [123, 456]}' | metalcloud-cli server firmware generate-audit --config-source pipe

  # Also run the safety checks an upgrade would run now and list the jobs it would create
  echo '{"serverIds": [123]}' | metalcloud-cli server firmware generate-audit --config-source pipe --include-preflight

Optional Flags:
  --include-preflight    Run the upgrade safety checks (BMC checks included) and list the planned
                         jobs, in order, for each server. Nothing is changed.


```
metalcloud-cli server firmware generate-audit [flags]
```

### Options

```
      --config-source string   Source of the audit configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for generate-audit
      --include-preflight      Run the upgrade safety checks and list the planned jobs for each server.
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API token (msat_...), or a legacy API key (<user id>:<key>), which works only until it expires
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management

