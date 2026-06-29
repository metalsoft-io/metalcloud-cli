## metalcloud-cli fabric import-devices

Bulk-import switches and attach them to a fabric (idempotent)

### Synopsis

Bulk-import network devices (switches) from a declarative YAML/JSON file and
attach them to a fabric, in one idempotent step.

For each switch the command finds-or-creates the device (matched against the
existing devices in the fabric's site by management address, identifier string,
or serial number) and attaches any that are not yet attached to the fabric.
Re-running it is safe: existing devices are not recreated and already-attached
ones are left alone. The switches' site is derived from the fabric, so do not
set siteId on the switches. No deploy is triggered.

Arguments:
  fabric_id    The ID or label of the (pre-existing) fabric to import into

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a YAML/JSON file.

Optional Flags:
  --dry-run         Report the plan (creates + attaches) without writing.

Configuration format:
  defaults:   optional map deep-merged into every switch (per-switch keys win)
  switches:   non-empty list of switch entries. Per switch (after defaults):
              driver, position, username, managementPassword, managementAddress,
              managementPort, identifierString are required; tagsMap and the rest
              are optional. tagsMap scalar values are coerced to strings.

Examples:
  metalcloud-cli fabric import-devices 5 --config-source switches.yaml --dry-run
  metalcloud-cli fabric import-devices my-fabric --config-source switches.yaml
  cat switches.yaml | metalcloud-cli fabric import-switches 5 --config-source pipe

```
metalcloud-cli fabric import-devices fabric_id [flags]
```

### Options

```
      --config-source string   Source of the switch import file. Can be 'pipe' or path to a YAML/JSON file.
      --dry-run                Report the plan without creating or attaching anything.
  -h, --help                   help for import-devices
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

