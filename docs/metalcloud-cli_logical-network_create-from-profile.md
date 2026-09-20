## metalcloud-cli logical-network create-from-profile

Create a logical network from a logical network profile

### Synopsis

Create a logical network whose configuration is taken from a logical network profile.

The network can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

Required Flags (one of):
  --profile-id       The ID of the logical network profile to instantiate
  --config-source    'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --label              The label of the new logical network
  --name               The name of the new logical network
  --infrastructure-id  The infrastructure the network belongs to
  --mtu                Maximum Transmission Unit in bytes

Examples:
  metalcloud-cli logical-network create-from-profile --profile-id 3 --label my-network
  metalcloud-cli logical-network create-from-profile --profile-id 3 --infrastructure-id 7 --mtu 9000
  metalcloud-cli logical-network create-from-profile --config-source network.json
  echo '{"logicalNetworkProfileId":3,"label":"my-network"}' | \
    metalcloud-cli logical-network create-from-profile --config-source pipe

```
metalcloud-cli logical-network create-from-profile [flags]
```

### Options

```
      --config-source string    Source of the create-from-profile configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                    help for create-from-profile
      --infrastructure-id int   The infrastructure the new logical network belongs to.
      --label string            The label of the new logical network.
      --mtu int                 Maximum Transmission Unit (MTU) in bytes.
      --name string             The name of the new logical network.
      --profile-id int          The ID of the logical network profile to instantiate.
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

